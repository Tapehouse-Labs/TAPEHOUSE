// SPDX-License-Identifier: MIT OR Apache-2.0
//! Reads SPY's band and its BandFeed on a dev node, sends RedStone packages from a file through a `PackageSource` and
//! decodes the band's refusal of their age, then authorizes a fresh address in the margin accounts, which sells 1 SPY
//! short for the account at a QuoterV2 quote and buys it back. Then it mints 1 share of the basket PAIR for the fresh
//! address, which deposits it into its cross position, reads it there as NVDA and SPY, unwraps it and withdraws them.
//!
//! Usage: `PRIVATE_KEY=0x… cargo run --example devnode -- RPC_URL DEPLOYMENTS_JSON PAYLOAD_FILE`

use alloy::network::TransactionBuilder;
use alloy::primitives::{Bytes, I256, U256, utils::parse_ether};
use alloy::providers::{Provider, ProviderBuilder};
use alloy::rpc::types::{BlockId, TransactionRequest};
use alloy::signers::local::PrivateKeySigner;
use alloy::sol_types::SolEvent;
use tapehouse_sdk::bindings::{MarginAccounts, ShortPositions};
use tapehouse_sdk::{CROSS, Deployments, Error, PackageSource, Revert, Tapehouse, to_bytes32};

struct FileSource(String);

impl PackageSource for FileSource {
    async fn payload(&self, _: &[&str]) -> Result<Bytes, Box<dyn std::error::Error + Send + Sync>> {
        Ok(std::fs::read_to_string(&self.0)?.trim().parse()?)
    }
}

fn check(ok: bool, message: impl std::fmt::Display) -> Result<(), Box<dyn std::error::Error>> {
    if ok {
        Ok(())
    } else {
        Err(format!("FAIL: {message}").into())
    }
}

fn refusal<T>(
    result: Result<T, alloy::contract::Error>,
) -> Result<Revert, Box<dyn std::error::Error>> {
    match result.map_err(Error::from) {
        Err(Error::Revert(revert)) => Ok(revert),
        Err(error) => {
            Err(format!("FAIL: the call did not revert with a decodable error: {error}").into())
        }
        Ok(_) => Err("FAIL: the call did not revert".into()),
    }
}

fn event<E: SolEvent>(
    receipt: &alloy::rpc::types::TransactionReceipt,
) -> Result<E, Box<dyn std::error::Error>> {
    receipt
        .decoded_log::<E>()
        .map(|log| log.data)
        .ok_or_else(|| format!("FAIL: the receipt has no {}", E::SIGNATURE).into())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let [_, rpc, registry, payload] = std::env::args()
        .collect::<Vec<_>>()
        .try_into()
        .map_err(|_| "usage: PRIVATE_KEY=0x… devnode RPC_URL DEPLOYMENTS_JSON PAYLOAD_FILE")?;
    let key =
        std::env::var("PRIVATE_KEY").map_err(|_| "set PRIVATE_KEY to the account owner's key")?;
    let deployments = Deployments::load(&registry)?;
    let owner_key: PrivateKeySigner = key.parse()?;
    let account = owner_key.address();
    let operator_key = PrivateKeySigner::random();
    let operator_address = operator_key.address();
    let owner_provider = ProviderBuilder::new()
        .wallet(owner_key)
        .connect_http(rpc.parse()?);
    let operator_provider = ProviderBuilder::new()
        .wallet(operator_key)
        .connect_http(rpc.parse()?);
    check(
        owner_provider.get_chain_id().await? == deployments.chain_id,
        "the registry is for another chain",
    )?;
    let owner = Tapehouse::new(owner_provider.clone(), deployments.clone());
    let operator = Tapehouse::new(operator_provider.clone(), deployments.clone());
    let spy = to_bytes32("SPY")?;

    let block = owner_provider
        .get_block(BlockId::latest())
        .await?
        .ok_or("FAIL: no latest block")?
        .header;
    let at = BlockId::number(block.number);
    let band = owner.band()?;
    let quote = band.quote(spy).block(at).call().await?;
    let session = band.session().block(at).call().await?;
    let halt = band.halt(spy).block(at).call().await?;
    check(quote.state != 0 && !halt.signedHalt, "SPY's band is halted")?;
    let feed = owner.band_feed("SPY")?;
    let round = feed.latestRoundData().block(at).call().await?;
    check(
        round._1 == I256::try_from(quote.low)?,
        "the feed's answer is not the band's low edge",
    )?;
    check(
        round._0.to::<u64>() == block.timestamp && round._3 == U256::from(block.timestamp),
        "the round is not the block time",
    )?;
    let whole = feed.latestBand().block(at).call().await?;
    check(
        whole.state == quote.state && whole.mid == quote.mid,
        "latestBand is not the band",
    )?;
    let sealed = feed.seals(session.boundaryMs).block(at).call().await?;
    let sealing = match feed.seal().from(account).call().await {
        Ok(_) => "within its window".to_string(),
        Err(error) => {
            let revert = refusal::<()>(Err(error))?;
            check(
                revert.name == "NotSealWindow"
                    && revert.to_string()
                        == format!("NotSealWindow({}, {})", session.state, session.boundaryMs),
                format!("seal() reverted with {revert}"),
            )?;
            revert.to_string()
        }
    };

    let write = owner
        .write_prices(&FileSource(payload), &["NVDA---24_7"])
        .await?;
    let stale = refusal(write.from(account).call().await)?;
    check(
        stale
            .to_string()
            .starts_with("TimestampIsTooOld(1790119750, "),
        format!("the band took the file's packages: {stale}"),
    )?;

    let funding = TransactionRequest::default()
        .with_to(operator_address)
        .with_value(parse_ether("0.01")?);
    owner_provider
        .send_transaction(funding)
        .await?
        .get_receipt()
        .await?;
    let margin = U256::from(1_000_000_000u64);
    let refused = refusal(
        operator
            .shorts()?
            .deposit(spy, margin, account)
            .call()
            .await,
    )?;
    check(
        refused.to_string() == format!("Unauthorized({operator_address}, {account})"),
        format!("the shorts did not refuse: {refused}"),
    )?;
    let usdg = deployments.tokens["USDG"];
    let withdrawal = operator
        .accounts()?
        .withdraw(CROSS, usdg, U256::from(1), account, account)
        .call()
        .await;
    check(
        refusal(withdrawal)?.to_string() == refused.to_string(),
        "the margin accounts did not refuse the operator",
    )?;

    let authorized = owner
        .accounts()?
        .setAuthorization(operator_address, true)
        .send()
        .await?
        .get_receipt()
        .await?;
    let set = event::<MarginAccounts::AuthorizationSet>(&authorized)?;
    check(
        set.authorized == operator_address && set.allowed,
        "no AuthorizationSet",
    )?;
    check(
        owner
            .accounts()?
            .isAuthorized(account, operator_address)
            .call()
            .await?,
        "the operator is not authorized",
    )?;

    let shorts_address = deployments.tapehouse["ShortPositions"];
    owner
        .usdg()?
        .approve(shorts_address, margin)
        .send()
        .await?
        .get_receipt()
        .await?;
    owner
        .shorts()?
        .deposit(spy, margin, account)
        .send()
        .await?
        .get_receipt()
        .await?;
    let one = U256::from(10u64.pow(18));
    let sale = owner.quote_sale("SPY", one, 50).await?;
    let sold = operator
        .shorts()?
        .sell(spy, one, sale.min_proceeds, account)
        .send()
        .await?
        .get_receipt()
        .await?;
    check(sold.status(), "the sale reverted")?;
    let sell = event::<ShortPositions::Sell>(&sold)?;
    check(
        sell.proceeds == sale.proceeds,
        "the sale did not pay QuoterV2's quote",
    )?;
    let open = owner.shorts()?.position(account, spy).call().await?;
    check(open.debt >= one, "the short does not owe 1 SPY and its fee")?;
    let health = owner.shorts()?.health(account, spy).call().await?;
    check(
        health.equity > I256::try_from(health.requirement)?,
        "the short falls short",
    )?;
    let buy_back = owner.quote_cover("SPY", open.debt, 50).await?;
    let covered = operator
        .shorts()?
        .cover(spy, U256::MAX, buy_back.max_cost, account)
        .send()
        .await?
        .get_receipt()
        .await?;
    check(covered.status(), "the buy-back reverted")?;
    let cover = event::<ShortPositions::Cover>(&covered)?;
    check(
        cover.cost >= buy_back.cost && cover.cost <= buy_back.max_cost,
        "the buy-back did not cost QuoterV2's quote, and the fee accrued since, within the slippage",
    )?;
    let closed = owner.shorts()?.position(account, spy).call().await?;
    check(
        closed.shares.is_zero() && closed.debt.is_zero(),
        "the short is still open",
    )?;
    operator
        .shorts()?
        .withdraw(spy, closed.usdgHeld.try_into()?, account, account)
        .send()
        .await?
        .get_receipt()
        .await?;
    owner
        .accounts()?
        .setAuthorization(operator_address, false)
        .send()
        .await?
        .get_receipt()
        .await?;
    check(
        !owner
            .accounts()?
            .isAuthorized(account, operator_address)
            .call()
            .await?,
        "the operator is still authorized",
    )?;

    let pair = owner.basket("PAIR")?;
    let components = owner.basket_components("PAIR").await?;
    check(
        components.assets == ["NVDA", "SPY"],
        format!("PAIR holds {:?}", components.assets),
    )?;
    let paid = owner.preview_mint("PAIR", one).await?;
    for (asset, amount) in components.assets.iter().zip(&paid) {
        owner
            .stock_token(asset)?
            .approve(*pair.address(), *amount)
            .send()
            .await?
            .get_receipt()
            .await?;
    }
    let minted = pair
        .mint(one, operator_address, paid.clone())
        .send()
        .await?
        .get_receipt()
        .await?;
    check(minted.status(), "the mint reverted")?;
    let margin_accounts = *operator.accounts()?.address();
    operator
        .basket("PAIR")?
        .approve(margin_accounts, one)
        .send()
        .await?
        .get_receipt()
        .await?;
    let deposited = operator
        .accounts()?
        .deposit(CROSS, *pair.address(), one, operator_address)
        .send()
        .await?
        .get_receipt()
        .await?;
    check(deposited.status(), "the deposit of the share reverted")?;
    let through = owner.in_baskets(operator_address, CROSS).await?;
    let margined = owner
        .accounts()?
        .health(operator_address, CROSS)
        .call()
        .await?;
    check(
        through["NVDA"] == paid[0] && through["SPY"] == paid[1],
        "the share does not hold what it was minted for",
    )?;
    let shares = owner
        .collateral(operator_address, CROSS, *pair.address())
        .await?;
    let unwrapped_receipt = operator
        .unwrap(operator_address, "PAIR", one)?
        .send()
        .await?
        .get_receipt()
        .await?;
    check(unwrapped_receipt.status(), "the unwrap reverted")?;
    let mut unwrapped = Vec::new();
    for token in &components.tokens {
        unwrapped.push(owner.collateral(operator_address, CROSS, *token).await?);
    }
    check(
        unwrapped == [through["NVDA"], through["SPY"]],
        "the unwrap did not put the share's tokens in the position",
    )?;
    for (token, amount) in components.tokens.iter().zip(&unwrapped) {
        operator
            .accounts()?
            .withdraw(CROSS, *token, *amount, operator_address, operator_address)
            .send()
            .await?
            .get_receipt()
            .await?;
    }

    println!(
        "rust sdk: SPY band state {} at {}, feed round {} answers {}, seal at {} {sealing} (sealed {}); writePrices {stale}; \
         before authorization {refused}; sold 1 SPY for {} and bought it back for {}",
        quote.state,
        quote.mid,
        round._0,
        round._1,
        session.boundaryMs,
        sealed.sealedAt,
        sell.proceeds,
        cover.cost
    );
    println!(
        "basket PAIR: {shares} shares in the cross position hold {} NVDA and {} SPY; equity {} requirement {}",
        through["NVDA"], through["SPY"], margined.equity, margined.requirement
    );
    println!(
        "basket PAIR: unwrapped into {} NVDA and {} SPY in the cross position, then withdrawn",
        unwrapped[0], unwrapped[1]
    );
    println!("PASS");
    Ok(())
}
