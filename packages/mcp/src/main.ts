#!/usr/bin/env node
// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from 'node:fs'
import { createServer as createHttpServer } from 'node:http'
import type { AddressInfo } from 'node:net'
import { parseArgs } from 'node:util'
import { localhostHostValidation, localhostOriginValidation, toNodeHandler } from '@modelcontextprotocol/node'
import { createMcpHandler } from '@modelcontextprotocol/server'
import { serveStdio } from '@modelcontextprotocol/server/stdio'
import { parseDeployments } from '@tapehouse/sdk'
import { BaseError, createPublicClient, http } from 'viem'
import { getChainId } from 'viem/actions'
import { configFromEnv } from './config.js'
import { createServer, RateLimit } from './server.js'

try {
  const { values } = parseArgs({ options: { http: { type: 'string' } } })
  const config = configFromEnv(process.env)
  const deployments = parseDeployments(JSON.parse(readFileSync(config.deployments, 'utf8')))
  const client = createPublicClient({ transport: http(config.rpcUrl) })
  const chainId = await getChainId(client)
  if (chainId !== deployments.chainId)
    throw new Error(`TAPEHOUSE_RPC_URL serves chain ${chainId}, TAPEHOUSE_DEPLOYMENTS chain ${deployments.chainId}.`)
  const context = {
    client,
    deployments,
    maxSlippageBps: config.maxSlippageBps,
    rateLimit: new RateLimit(config.maxCallsPerMinute),
  }
  if (values.http === undefined) {
    serveStdio(() => createServer(context))
    console.error(`tapehouse-mcp serves chain ${chainId} on stdio`)
  } else {
    const port = Number(values.http)
    if (!/^\d{1,5}$/.test(values.http) || port > 65_535) throw new Error(`--http ${values.http} is not a port.`)
    const handler = createMcpHandler(() => createServer(context))
    const serve = toNodeHandler(handler)
    const host = localhostHostValidation()
    const origin = localhostOriginValidation()
    const listener = createHttpServer((req, res) => {
      if (host(req, res) && origin(req, res)) void serve(req, res)
    })
    listener.listen(port, '127.0.0.1', () => {
      const { port: bound } = listener.address() as AddressInfo
      console.error(`tapehouse-mcp serves chain ${chainId} on http://127.0.0.1:${bound}/mcp`)
    })
    const stop = async () => {
      await handler.close()
      listener.close()
    }
    process.once('SIGINT', stop)
    process.once('SIGTERM', stop)
  }
} catch (error) {
  console.error(error instanceof BaseError ? error.shortMessage : error instanceof Error ? error.message : error)
  process.exitCode = 1
}
