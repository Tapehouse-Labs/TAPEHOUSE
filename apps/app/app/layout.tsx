// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata, Viewport } from "next";
import { Archivo, Inter_Tight, Spline_Sans_Mono } from "next/font/google";
import { loadRegistry } from "@/lib/registry";
import { Providers } from "./providers";
import "./globals.css";

const archivo = Archivo({ variable: "--font-archivo", subsets: ["latin"], display: "swap" });
const interTight = Inter_Tight({ variable: "--font-inter-tight", subsets: ["latin"], display: "swap" });
const splineMono = Spline_Sans_Mono({ variable: "--font-spline-mono", subsets: ["latin"], display: "swap" });

export const metadata: Metadata = {
  title: { default: "Tapehouse", template: "%s — Tapehouse" },
  description: "Portfolio margin for Stock Tokens.",
};

export const viewport: Viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "#f5f2ec" },
    { media: "(prefers-color-scheme: dark)", color: "#14120e" },
  ],
};

const THEME = `(function(){try{var t=localStorage.getItem("theme");if(t==="light"||t==="dark")document.documentElement.setAttribute("data-theme",t)}catch(e){}})()`;

export default async function RootLayout({ children }: LayoutProps<"/">) {
  const { deployments, rpcUrl, bundlerUrl, sponsorUrl } = await loadRegistry(process.env);
  return (
    <html
      lang="en"
      className={`${archivo.variable} ${interTight.variable} ${splineMono.variable} antialiased`}
      suppressHydrationWarning
    >
      <head>
        <script dangerouslySetInnerHTML={{ __html: THEME }} />
      </head>
      <body>
        <Providers deployments={deployments} rpcUrl={rpcUrl} bundlerUrl={bundlerUrl} sponsorUrl={sponsorUrl}>
          {children}
        </Providers>
      </body>
    </html>
  );
}
