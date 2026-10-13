// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata, Viewport } from "next";
import { Archivo, Inter_Tight, Spline_Sans_Mono } from "next/font/google";
import { SITE_URL } from "@/lib/site";
import "./globals.css";

const archivo = Archivo({ variable: "--font-archivo", subsets: ["latin"], display: "swap" });
const interTight = Inter_Tight({ variable: "--font-inter-tight", subsets: ["latin"], display: "swap" });
const splineMono = Spline_Sans_Mono({ variable: "--font-spline-mono", subsets: ["latin"], display: "swap" });

export const metadata: Metadata = {
  metadataBase: new URL(SITE_URL),
  title: { default: "Tapehouse documentation", template: "%s · Tapehouse docs" },
  description: "Portfolio margin for Stock Tokens: the band, the margin accounts, the SDKs, the API and the contracts.",
  openGraph: { type: "website", siteName: "Tapehouse docs", url: SITE_URL },
  twitter: { card: "summary", site: "@tapehouseHQ" },
};

export const viewport: Viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "#f5f2ec" },
    { media: "(prefers-color-scheme: dark)", color: "#14120e" },
  ],
};

const THEME = `try{var t=localStorage.getItem("theme");if(t==="light"||t==="dark")document.documentElement.dataset.theme=t}catch(e){}`;

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={`${archivo.variable} ${interTight.variable} ${splineMono.variable} antialiased`} suppressHydrationWarning>
      <head>
        <script dangerouslySetInnerHTML={{ __html: THEME }} />
      </head>
      <body>
        <a href="#main" className="sr-only focus:not-sr-only focus:absolute focus:left-3 focus:top-3 focus:z-50 focus:rounded focus:bg-paper focus:px-3 focus:py-2">
          Skip to the content
        </a>
        {children}
      </body>
    </html>
  );
}
