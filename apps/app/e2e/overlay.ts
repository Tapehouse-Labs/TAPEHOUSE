// SPDX-License-Identifier: MIT OR Apache-2.0
import { reopeningAuctionAbi } from "@tapehouse/sdk";
import {
  type Address,
  decodeFunctionData,
  encodeFunctionResult,
  erc20Abi,
  type Hex,
  keccak256,
  multicall3Abi,
  parseTransaction,
  recoverTransactionAddress,
  toHex,
} from "viem";

type Round = {
  sealMs: bigint;
  floor: bigint;
  supply: bigint;
  deposits: bigint;
  pool: bigint;
  price: bigint;
  above: bigint;
  atPrice: bigint;
  sold: bigint;
  paid: bigint;
  taken: bigint;
  cleared: boolean;
};
type Bid = { bidder: Address; quantity: bigint; price: bigint; escrow: bigint; claimed: boolean };

/** What the overlay says the reopening auction holds, and the visitor's USDG: a test double of a chain mid-reopening. */
export type World = {
  phase: { openMs: bigint; revealing: boolean };
  lastOpenMs: bigint;
  rounds: Map<string, { round: Round; bids: Bid[] }>;
  commitments: Map<string, { bidder: Address; deposit: bigint }>;
  usdg: { balance: bigint; allowance: bigint };
  /** Seconds the clock the app reads is ahead of the node's. */
  shift: bigint;
};

export const emptyRound = (): Round => ({
  sealMs: 0n,
  floor: 0n,
  supply: 0n,
  deposits: 0n,
  pool: 0n,
  price: 0n,
  above: 0n,
  atPrice: 0n,
  sold: 0n,
  paid: 0n,
  taken: 0n,
  cleared: false,
});

const AGGREGATE3 = "0x82ad56cb";

type Rpc = { id: number; jsonrpc: "2.0"; method: string; params?: unknown[] };

/**
 * A JSON-RPC front for a real node that answers the reopening auction's reads, the visitor's USDG and the sending of
 * transactions from a `World`, and passes everything else, the bands and the clock included, to the node. A
 * transaction it is sent is recorded, decoded, applied to the world and answered as mined.
 */
export class Overlay {
  readonly sent: { to: Address; name: string; args: readonly unknown[] }[] = [];
  readonly world: World = {
    phase: { openMs: 0n, revealing: false },
    lastOpenMs: 0n,
    rounds: new Map(),
    commitments: new Map(),
    usdg: { balance: 0n, allowance: 0n },
    shift: 0n,
  };
  private readonly mined = new Map<Hex, { from: Address; to: Address }>();

  constructor(
    private readonly upstream: string,
    private readonly auction: Address,
    private readonly usdg: Address,
  ) {}

  async handle(body: Rpc | Rpc[]): Promise<unknown> {
    if (Array.isArray(body)) return Promise.all(body.map((request) => this.handle(request)));
    try {
      const result = await this.answer(body);
      return { jsonrpc: "2.0", id: body.id, result };
    } catch (error) {
      return { jsonrpc: "2.0", id: body.id, error: { code: -32000, message: String(error) } };
    }
  }

  private round(symbol: Hex, openMs: bigint) {
    const key = `${symbol}:${openMs}`;
    if (!this.world.rounds.has(key)) this.world.rounds.set(key, { round: emptyRound(), bids: [] });
    return this.world.rounds.get(key)!;
  }

  private async answer(request: Rpc): Promise<unknown> {
    const { method, params = [] } = request;
    if (method === "eth_call" || method === "eth_estimateGas") {
      const { to, data } = params[0] as { to?: Address; data?: Hex };
      const target = to?.toLowerCase();
      if (method === "eth_call" && data?.startsWith(AGGREGATE3)) return this.aggregate(data, params[1]);
      if (target === this.auction.toLowerCase() && data) {
        const answered = this.callAuction(method, data);
        if (answered !== undefined) return answered;
      }
      if (target === this.usdg.toLowerCase() && data) {
        const call = decodeFunctionData({ abi: erc20Abi, data });
        if (call.functionName === "balanceOf") return encodeFunctionResult({ abi: erc20Abi, functionName: "balanceOf", result: this.world.usdg.balance });
        if (call.functionName === "allowance") return encodeFunctionResult({ abi: erc20Abi, functionName: "allowance", result: this.world.usdg.allowance });
        if (call.functionName === "approve" && method === "eth_call") return encodeFunctionResult({ abi: erc20Abi, functionName: "approve", result: true });
        if (method === "eth_estimateGas") return toHex(100_000n);
      }
    }
    if (method === "eth_getBlockByNumber") {
      const block = (await this.forward(request)) as { timestamp: Hex } | null;
      return block && { ...block, timestamp: toHex(BigInt(block.timestamp) + this.world.shift) };
    }
    if (method === "eth_sendRawTransaction") return await this.mine(params[0] as Hex);
    if (method === "eth_getTransactionReceipt") {
      const mined = this.mined.get(params[0] as Hex);
      if (mined) return this.receipt(params[0] as Hex, mined);
    }
    return this.forward(request);
  }

  /** A Multicall3 `aggregate3`, each call answered by the world where it can be and by the node where not. */
  private async aggregate(data: Hex, block: unknown): Promise<Hex> {
    const { args } = decodeFunctionData({ abi: multicall3Abi, data });
    const calls = (args as [readonly { target: Address; allowFailure: boolean; callData: Hex }[]])[0];
    const results = await Promise.all(
      calls.map(async ({ target, callData }) => {
        try {
          const returnData = (await this.answer({ id: 0, jsonrpc: "2.0", method: "eth_call", params: [{ to: target, data: callData }, block] })) as Hex;
          return { success: true, returnData };
        } catch {
          return { success: false, returnData: "0x" as Hex };
        }
      }),
    );
    return encodeFunctionResult({ abi: multicall3Abi, functionName: "aggregate3", result: results });
  }

  private callAuction(method: string, data: Hex): Hex | undefined {
    const call = decodeFunctionData({ abi: reopeningAuctionAbi, data });
    const encode = (result: unknown) =>
      encodeFunctionResult({ abi: reopeningAuctionAbi, functionName: call.functionName, result } as never);
    if (method === "eth_estimateGas") return toHex(300_000n);
    const args = (call.args ?? []) as readonly unknown[];
    switch (call.functionName) {
      case "phase":
        return encode([this.world.phase.openMs, this.world.phase.revealing]);
      case "lastOpenMs":
        return encode(this.world.lastOpenMs);
      case "round":
        return encode(this.round(args[0] as Hex, args[1] as bigint).round);
      case "lots":
        return encode([]);
      case "bids":
        return encode(this.round(args[0] as Hex, args[1] as bigint).bids);
      case "commitments": {
        const entry = this.world.commitments.get(args[2] as Hex);
        return encode([entry?.bidder ?? "0x0000000000000000000000000000000000000000", entry?.deposit ?? 0n]);
      }
      case "commit":
      case "reveal":
      case "claim":
      case "clear":
      case "forfeit":
        return "0x";
      default:
        return undefined;
    }
  }

  private async mine(raw: Hex): Promise<Hex> {
    const transaction = parseTransaction(raw);
    const hash = keccak256(raw);
    const to = transaction.to as Address;
    const from = await recoverTransactionAddress({ serializedTransaction: raw as never });
    this.mined.set(hash, { from, to });
    if (to.toLowerCase() === this.usdg.toLowerCase() && transaction.data) {
      const call = decodeFunctionData({ abi: erc20Abi, data: transaction.data });
      this.sent.push({ to, name: call.functionName, args: call.args ?? [] });
      if (call.functionName === "approve") this.world.usdg.allowance = (call.args as [Address, bigint])[1];
    } else if (transaction.data) {
      const call = decodeFunctionData({ abi: reopeningAuctionAbi, data: transaction.data });
      const args = (call.args ?? []) as readonly unknown[];
      this.sent.push({ to, name: call.functionName, args });
      this.apply(call.functionName, args, from);
    }
    return hash;
  }

  private apply(name: string, args: readonly unknown[], from: Address) {
    if (name === "commit") {
      const [, commitment, deposit] = args as [Hex, Hex, bigint];
      this.world.commitments.set(commitment, { bidder: from, deposit });
      this.world.usdg.balance -= deposit;
    }
  }

  private receipt(hash: Hex, { from, to }: { from: Address; to: Address }) {
    return {
      transactionHash: hash,
      transactionIndex: "0x0",
      blockHash: keccak256(hash),
      blockNumber: "0x1",
      from,
      to,
      cumulativeGasUsed: "0x5208",
      gasUsed: "0x5208",
      effectiveGasPrice: "0x1",
      contractAddress: null,
      logs: [],
      logsBloom: `0x${"0".repeat(512)}`,
      status: "0x1",
      type: "0x2",
    };
  }

  private async forward(request: Rpc): Promise<unknown> {
    const response = await fetch(this.upstream, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify(request),
    });
    const body = (await response.json()) as { result?: unknown; error?: unknown };
    if (body.error) throw new Error(JSON.stringify(body.error));
    return body.result;
  }
}
