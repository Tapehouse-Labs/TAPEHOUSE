import type { Metadata, Viewport } from "next";
import { Archivo, Inter_Tight, Spline_Sans_Mono } from "next/font/google";
import "./globals.css";
import { SITE_URL, X_HANDLE } from "./site";

const archivo = Archivo({ variable: "--font-archivo", subsets: ["latin"], display: "swap" });
const interTight = Inter_Tight({ variable: "--font-inter-tight", subsets: ["latin"], display: "swap" });
const splineMono = Spline_Sans_Mono({ variable: "--font-spline-mono", subsets: ["latin"], display: "swap" });

const DESCRIPTION =
  "Borrow against every Stock Token you hold. Your whole portfolio backs one loan, the liquidation price is stated before you sign, and the price keeps moving when the market doesn’t.";

export const metadata: Metadata = {
  metadataBase: new URL(SITE_URL),
  title: { default: "Tapehouse — Portfolio margin for Stock Tokens", template: "%s — Tapehouse" },
  description: DESCRIPTION,
  alternates: { canonical: "/" },
  openGraph: {
    type: "website",
    siteName: "Tapehouse",
    url: "/",
    title: "Tapehouse — Portfolio margin for Stock Tokens",
    description: DESCRIPTION,
  },
  twitter: {
    card: "summary_large_image",
    site: `@${X_HANDLE}`,
    title: "Tapehouse — Portfolio margin for Stock Tokens",
    description: DESCRIPTION,
  },
};

export const viewport: Viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "#f5f2ec" },
    { media: "(prefers-color-scheme: dark)", color: "#14120e" },
  ],
};

/**
 * Direction contract. Emitted as an HTML comment, first child of <body>, so it
 * survives the production build and can be audited against the render.
 */
const CONTRACT = `<!--
TAPEHOUSE · LANDING · direction contract (seed 36b26739)
THESIS: The market is closed and one desk is still working. The page opens on that scene and pushes the camera into its lit screen until the product fills the frame. Refuses the category default of headline-left, dashboard-screenshot-right, three stat tiles.
OWN-WORLD: A night trading hall in registry navy and warm bone light; one lit screen as the only bright surface. Archivo display set light and tight, Inter Tight for interface and tabular figures, Spline Sans Mono for the session seal. Pill controls; a bone primary and a glass secondary. The lit screen surface never inverts with the theme.
STORY: The visitor understands that her whole portfolio backs one loan, believes it because the liquidation price is on screen before she signs, and opens the app on testnet.
FIRST VIEWPORT: Full-bleed filmed scene. Wordmark and session seal top, centred three-sentence headline in the dark upper third, footnoted sub-line, primary action "Open the app" beside a glass secondary, legal confession pinned to the bottom edge. Scroll drives a push-in anchored on the lit screen; the account panel grows out of it.
FORM: The opening shot of a film, a scene and a push-in. Candidate 1 of 7 on the grounded list, chosen by the user over the rolled assignment.
FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
-->`;

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html
      lang="en"
      data-theme="light"
      className={`${archivo.variable} ${interTight.variable} ${splineMono.variable} antialiased`}
    >
      <body>
        <div hidden dangerouslySetInnerHTML={{ __html: CONTRACT }} />
        {children}
      </body>
    </html>
  );
}
