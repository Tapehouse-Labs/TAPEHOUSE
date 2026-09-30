// SPDX-License-Identifier: MIT OR Apache-2.0
/**
 * Public addresses. Each can be overridden with an environment variable;
 * a link left empty is not rendered, so the footer never ships a dead end.
 */
export const SITE_URL = process.env.NEXT_PUBLIC_SITE_URL ?? "https://tapehouse.xyz";
export const APP_URL = process.env.NEXT_PUBLIC_APP_URL ?? "#open-app";

export const X_HANDLE = "tapehouseHQ";

export const EXTERNAL: { label: string; href?: string }[] = [
  { label: "X", href: process.env.NEXT_PUBLIC_X_URL ?? `https://x.com/${X_HANDLE}` },
  { label: "GitHub", href: process.env.NEXT_PUBLIC_GITHUB_URL ?? "https://github.com/Tapehouse-Labs/TAPEHOUSE" },
  { label: "Telegram", href: process.env.NEXT_PUBLIC_TELEGRAM_URL },
];

export const CONTACT_EMAIL = process.env.NEXT_PUBLIC_CONTACT_EMAIL ?? "hello@tapehouse.xyz";
export const SECURITY_EMAIL = "security@tapehouse.xyz";
