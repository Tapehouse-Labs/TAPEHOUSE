// SPDX-License-Identifier: MIT OR Apache-2.0
export type Group = { title: string; pages: string[] };
export type Tab = { id: string; title: string; home: string; groups: Group[] };

/** The sidebar of each tab, in reading order. A page is a file under `content/`, named by its slug. */
export const tabs: Tab[] = [
  {
    id: "docs",
    title: "Documentation",
    home: "/",
    groups: [
      { title: "Get started", pages: ["index", "quickstart", "networks"] },
      {
        title: "Concepts",
        pages: [
          "concepts/band",
          "concepts/session",
          "concepts/portfolio-margin",
          "concepts/liquidation",
          "concepts/weekend-gap-cover",
          "concepts/trust",
        ],
      },
      {
        title: "Guides",
        pages: [
          "guides/read-a-price",
          "guides/open-an-account",
          "guides/borrow",
          "guides/earn",
          "guides/short-and-baskets",
          "guides/morpho-oracle",
          "guides/auction",
        ],
      },
      {
        title: "Operate",
        pages: ["operate/indexer", "operate/keepers", "operate/sponsor", "operate/bidder", "operate/alerts"],
      },
    ],
  },
  {
    id: "sdks",
    title: "SDKs",
    home: "/sdks/overview",
    groups: [
      { title: "SDKs", pages: ["sdks/overview", "sdks/typescript", "sdks/go", "sdks/rust", "sdks/errors"] },
    ],
  },
  {
    id: "agents",
    title: "Agents",
    home: "/agents/overview",
    groups: [
      {
        title: "Agents",
        pages: ["agents/overview", "agents/mcp", "agents/claude-code", "agents/mandates", "agents/x402"],
      },
    ],
  },
  {
    id: "api",
    title: "API reference",
    home: "/api/overview",
    groups: [{ title: "Indexer API", pages: ["api/overview"] }],
  },
  {
    id: "contracts",
    title: "Contracts",
    home: "/contracts/overview",
    groups: [
      { title: "Contracts", pages: ["contracts/overview", "contracts/addresses"] },
      {
        title: "Pricing",
        pages: ["contracts/band", "contracts/band-feeds", "contracts/margin-engine", "contracts/morpho-oracle"],
      },
      {
        title: "Lending",
        pages: [
          "contracts/supply-vault",
          "contracts/margin-accounts",
          "contracts/liquidator",
          "contracts/gap-backstop",
          "contracts/reopening-auction",
          "contracts/stock-lending",
        ],
      },
      {
        title: "Products",
        pages: ["contracts/short-positions", "contracts/baskets", "contracts/gap-cover", "contracts/sponsor-paymaster"],
      },
    ],
  },
  {
    id: "resources",
    title: "Resources",
    home: "/resources/backtest",
    groups: [
      {
        title: "Resources",
        pages: [
          "resources/backtest",
          "resources/calibration",
          "resources/stylus-vs-solidity",
          "resources/verification",
          "resources/glossary",
        ],
      },
    ],
  },
];

export function slugHref(slug: string): string {
  return slug === "index" ? "/" : `/${slug}`;
}

export function hrefSlug(segments: string[] | undefined): string {
  return segments && segments.length > 0 ? segments.join("/") : "index";
}

export function tabOf(slug: string): Tab {
  const found = tabs.find((t) => t.groups.some((g) => g.pages.includes(slug)));
  if (found) return found;
  const first = slug.split("/")[0];
  return tabs.find((t) => t.id === first) ?? tabs[0];
}
