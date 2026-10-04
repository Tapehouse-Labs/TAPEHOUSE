---
name: Tapehouse App
description: Portfolio margin for Stock Tokens. The account drawn as the mark, one column, live bands off a navy stem.
colors:
  registry-navy: "#16326b"
  registry-navy-deep: "#0e2149"
  navy-mist: "#e3e7f0"
  navy-haze: "#c2cce2"
  navy-slate: "#7c8fbf"
  navy-ink: "#3d5795"
  bone-paper: "#f5f2ec"
  bone-paper-cool: "#faf9f6"
  panel-white: "#ffffff"
  sunk-linen: "#ebe7de"
  hairline-rule: "#dcd6ca"
  warm-ink: "#1a1713"
  ledger-green: "#1c6b4a"
  brick-loss: "#b0503f"
  liquidation-red: "#9e1409"
  amber-warning: "#8a5a00"
  night-paper: "#14120e"
  night-panel: "#1b1813"
  night-sunk: "#221e18"
  night-rule: "#2f2a22"
  night-ink: "#f2efe8"
  night-navy: "#8fa9e8"
  night-text-on-accent: "#0e1116"
  night-green: "#5fbf92"
  night-loss: "#e59484"
  night-liquidation: "#ff7a66"
  night-amber: "#e0a93c"
typography:
  figure:
    fontFamily: "Inter Tight, system-ui, -apple-system, Helvetica Neue, Arial, sans-serif"
    fontSize: "clamp(4rem, 13vw, 6.5rem)"
    fontWeight: 600
    lineHeight: 1
    letterSpacing: "-0.035em"
    fontFeature: "\"tnum\" 1, \"lnum\" 1"
  figure-empty:
    fontFamily: "Inter Tight, system-ui, -apple-system, Helvetica Neue, Arial, sans-serif"
    fontSize: "clamp(2.75rem, 8vw, 4.5rem)"
    fontWeight: 500
    lineHeight: 1
    letterSpacing: "-0.03em"
  lede:
    fontFamily: "Inter Tight, system-ui, -apple-system, Helvetica Neue, Arial, sans-serif"
    fontSize: "19px"
    fontWeight: 400
    lineHeight: 1.375
  body:
    fontFamily: "Inter Tight, system-ui, -apple-system, Helvetica Neue, Arial, sans-serif"
    fontSize: "14px"
    fontWeight: 400
    lineHeight: 1.5
  control:
    fontFamily: "Inter Tight, system-ui, -apple-system, Helvetica Neue, Arial, sans-serif"
    fontSize: "15px"
    fontWeight: 500
    lineHeight: 1.2
  label:
    fontFamily: "Spline Sans Mono, ui-monospace, SFMono-Regular, Menlo, monospace"
    fontSize: "10.5px"
    fontWeight: 500
    lineHeight: 1.3
    letterSpacing: "0.1em"
  seal:
    fontFamily: "Spline Sans Mono, ui-monospace, SFMono-Regular, Menlo, monospace"
    fontSize: "11px"
    fontWeight: 400
    letterSpacing: "0.09em"
    fontFeature: "\"tnum\" 1, \"zero\" 0"
  wordmark:
    fontFamily: "Archivo, Helvetica Neue, Arial, sans-serif"
    fontSize: "13px"
    fontWeight: 600
    letterSpacing: "0.18em"
rounded:
  hairline: "2px"
  icon: "5px"
  option: "6px"
  menu: "10px"
  pill: "999px"
spacing:
  xs: "12px"
  sm: "16px"
  md: "24px"
  lg: "32px"
  xl: "64px"
  stem-gap: "2.65rem"
  stem-gap-wide: "2.75rem"
  column-max: "780px"
  header-height: "64px"
  row-height: "60px"
components:
  button-primary:
    backgroundColor: "{colors.registry-navy}"
    textColor: "{colors.bone-paper}"
    typography: "{typography.control}"
    rounded: "{rounded.pill}"
    padding: "0.7em 1.4em"
    height: "44px"
  button-quiet:
    backgroundColor: "transparent"
    textColor: "{colors.warm-ink}"
    typography: "{typography.control}"
    rounded: "{rounded.pill}"
    padding: "0.7em 1.4em"
    height: "44px"
  pill:
    backgroundColor: "{colors.panel-white}"
    textColor: "{colors.warm-ink}"
    rounded: "{rounded.pill}"
    padding: "0.5em 1em"
  nav-tab-current:
    backgroundColor: "{colors.panel-white}"
    textColor: "{colors.warm-ink}"
    rounded: "{rounded.pill}"
    padding: "8px 14px"
  seal:
    backgroundColor: "{colors.panel-white}"
    textColor: "{colors.warm-ink}"
    typography: "{typography.seal}"
    rounded: "{rounded.hairline}"
    padding: "0.3em 0.75em"
  seal-unknown:
    textColor: "{colors.amber-warning}"
  menu:
    backgroundColor: "{colors.panel-white}"
    rounded: "{rounded.menu}"
    padding: "16px"
    width: "256px"
  menu-option:
    textColor: "{colors.warm-ink}"
    rounded: "{rounded.option}"
    padding: "10px 12px"
  menu-option-hover:
    backgroundColor: "{colors.sunk-linen}"
  band-row:
    height: "{spacing.row-height}"
  band-track:
    height: "18px"
  skeleton:
    backgroundColor: "{colors.sunk-linen}"
    rounded: "{rounded.hairline}"
---

# Design System: Tapehouse App

## Overview

**Creative North Star: "The Thread"**

The account is the brand mark drawn at working scale. A 2px registry-navy stem runs down the left edge of one centred column; a 2px bar rises across it at 8°, and the liquidation price sits directly under that crossing at display size. Everything else hangs off the stem: a short hairline reaches from the stem to each live band row. There is one column, one accent and one line language; no tiles, no side rail, no cards.

The ground is warm bone paper in light and a warm near-black in dark. Text hierarchy is made by the ink's alpha over that ground, never by separate greys. Navy is used as line and as the single filled primary action. Semantic colour is reserved for market events: gain, loss, warning and liquidation. The interface is dense and exact where it speaks numbers (tabular figures, monospaced labels, the bracketed session seal) and warm where it speaks to the person (one plain sentence in second person under the figure).

Motion is drawn once and then holds still: the stem and the bar draw on load, the sentence under the figure swaps with a short rise when the wallet state changes, band spans glide when the band moves. Nothing loops except the loading skeleton.

**Key Characteristics:**
- One centred column (780px max) with a navy stem on its left edge and the 8° bar crossing above the main figure.
- Warm ink at four alpha steps (100 / 80 / 62 / 42%) over bone paper; no grey palette.
- Navy accent as 2px lines and one filled pill button; never as a tint wash outside the band area.
- Inter Tight for interface and figures, Spline Sans Mono for labels, addresses and the seal, Archivo only for the wordmark.
- Pill controls with inset 1px hairline rings; the seal is the one squared object.
- Every market state is written as a word; colour only reinforces it.

## Colors

Warm bone and warm ink with a single deep registry navy; semantic colours are named for the market event they report.

### Primary
- **Registry Navy** (`--navy-80`, aliased as `--accent`): the stem, the bar, the band-span edges, the mark in the header, focus outlines, text selection, and the fill of the primary button. In dark theme it lifts to **Night Navy** so the 2px lines keep their weight on the dark ground.
- **Navy stops** (Navy Mist, Navy Haze, Navy Slate, Navy Ink, Registry Navy Deep): the six-stop accent scale shared with the landing. Navy Mist is `--accent-surface`. The band area is the accent mixed at 12% into transparent (18% in dark), the only place navy appears as a fill other than the button.

### Neutral
- **Bone Paper** (`--paper`): page ground and the 90%-opaque header behind its blur. Also the browser theme colour and the text colour on navy.
- **Panel White** (`--panel`): the fill of pills, the seal, the current nav tab and popover menus. It sits one step above the paper without a shadow.
- **Sunk Linen** (`--sunk`): skeleton blocks, menu option hover, placeholder wallet icon.
- **Hairline Rule** (`--rule`): every 1px line: header bottom border, table row borders, the stem-to-row rules, the inset ring on pills and the seal, the scrollbar thumb.
- **Warm Ink** (`--ink`): the source of all text colour. Strong is ink at 100% (tickers, figures, control labels), body at 80% (sentences), muted at 62% (labels, secondary figures, empty figure), faint at 42% (the resting figure, pill hover ring).
- **Night set** (Night Paper, Night Panel, Night Sunk, Night Rule, Night Ink): the same roles in dark theme, all warm, never blue-black.

### Semantic
- **Ledger Green** (`--gain`): the live-network dot in the network pill and the connected-wallet dot. Reserved for gain.
- **Brick Loss** (`--loss`): reserved for loss.
- **Error** (`--error`, `#9e1409` light, `#ff7a66` dark): the inline wallet error sentence.
- **Liquidation Red** (`--liquidation`): only the main figure when it holds a live liquidation price.
- **Amber Warning** (`--warning`): band state words HALTED and DEGRADED, the seal when the session is UNKNOWN, the wrong-network wallet dot.

### Named Rules
**The Alpha Ink Rule.** Text hierarchy is ink over paper at 100, 80, 62 or 42% via `color-mix`. Never introduce a grey hex.

**The One Red Rule.** Liquidation red is spent on exactly one element: the live liquidation-price figure. Errors and losses use their own tokens.

**The Line Accent Rule.** Navy appears as 2px lines and one filled button per view. It is not used for backgrounds, cards or decorative tints.

## Typography

**Display Font:** Archivo (with Helvetica Neue, Arial), wordmark only
**Body Font:** Inter Tight (with system-ui, -apple-system, Helvetica Neue, Arial)
**Label/Mono Font:** Spline Sans Mono (with ui-monospace, SFMono-Regular, Menlo)

**Character:** A tight grotesque carries both prose and the large figure, so the number reads as part of the sentence rather than a dashboard readout; the mono gives labels, addresses and the seal the voice of a printed ledger.

### Hierarchy
- **Figure** (600, `clamp(4rem, 13vw, 6.5rem)`, line-height 1, -0.035em): the liquidation price under the crossing. Faint ink at rest, liquidation red when live. Empty state ("None yet") drops to 500 at `clamp(2.75rem, 8vw, 4.5rem)`, -0.03em, muted ink.
- **Lede** (400, 19px, 1.375, max 34ch): the one sentence that explains the figure and the wallet state.
- **Body** (400, 14px): band table rows, menus, error text (max 52ch). Tickers are 600 strong ink.
- **Control** (500, 15px): primary and quiet buttons. Pills use 13px at 400; nav tabs 14px.
- **Label** (500, 10.5px, 1.3, 0.1em, uppercase, mono, muted ink): field names above figures, the band table heading, band state words.
- **Seal** (mono, 11px, 0.09em, uppercase by content): the session seal.
- **Secondary figure** (12px, muted, tabular): band half-width under the centre price.
- **Wordmark** (Archivo 600, 13px, 0.18em, uppercase): header only.

### Named Rules
**The Tabular Money Rule.** Every figure that decides money carries the `num` treatment: tabular, lining numerals, ligatures and contextual alternates off, left-to-right. Addresses add mono with slashed zero.

**The Wordmark-Only Display Rule.** Archivo appears in the wordmark and nowhere else in the app; headings and figures are Inter Tight.

## Layout

A sticky 64px header spans a 1320px max container (16px gutter on phones, 32px from 640px). Below it, a single column, 780px max, centred. The column is padded asymmetrically to leave room for the stem: 3.5rem top, 1.25rem right, 7rem bottom, 4.25rem left on phones; 4rem, 2rem, 8rem, 6rem from 640px. The stem sits 1.6rem from the column's left edge (3.25rem from 640px) and runs its full height. The stem-gap (2.65rem, 2.75rem from 640px) is the distance from the content edge back to the stem and drives both the bar's overhang and the row rules.

Vertical rhythm: label to figure 12px, figure to lede 24px, lede to actions 32px (actions reserve a 44px row even when empty, so the column does not jump between wallet states), actions to bands 64px. Band rows are 60px tall with a 1px rule on top; columns are ticker (72px), track (fluid), state (118px, hidden below 640px) and centre/half-width (112px, right-aligned).

Responsive changes: below 640px the state column collapses into the centre cell when the asset is unpriced, the chain name leaves the bands heading, and the seal moves from the header to the bands heading. Below 768px the section tabs and the network pill leave the header. The bar's length is `min(640px, 82vw)`, so on phones it runs to the viewport edge.

## Elevation & Depth

Flat by default. Depth comes from tonal steps (sunk, paper, panel) and inset 1px hairline rings, not shadows. The header uses paper at 90% opacity with a medium backdrop blur so the stem reads as passing beneath it. The only cast shadow belongs to popover menus.

### Shadow Vocabulary
- **Hairline ring** (`box-shadow: inset 0 0 0 1px var(--rule)`): pills, the seal, the current nav tab, the quiet button. Hover on pills strengthens the ring to faint ink.
- **Popover lift** (`box-shadow: 0 18px 50px -12px rgb(0 0 0 / 0.35), 0 0 0 1px var(--rule)`): the wallet list and the connected-account menu.
- **Band edges** (`box-shadow: inset 1px 0 0 var(--accent), inset -1px 0 0 var(--accent)`): the left and right bounds of a band span.

### Named Rules
**The Ring Not Border Rule.** Controls are outlined with an inset 1px ring so they keep their box size; true borders are used only for full-width rules (header bottom, table rows).

## Shapes

Two silhouettes. Interactive controls are full pills (999px). The seal and loading blocks are near-square (2px), echoing printed tape. Popover menus are softly rounded (10px), their options 6px, wallet icons 5px. Lines are the main geometry: 2px navy for the stem and bar, 1px rule for everything structural, and the bar is always at -8°, matching the mark's rise.

## Components

### Buttons
Confident and few: one filled pill per view.
- **Shape:** full pill (999px), minimum 44px tall.
- **Primary:** navy fill, text on accent (bone in light, near-black in dark), 15px/500, padding 0.7em 1.4em. Used for "Get a wallet", "Connect wallet", "Switch to <chain>".
- **Hover / Active:** lifts 1px over 0.2s on the app ease; returns flat on press.
- **Disabled / Pending:** 55% opacity, no lift, progress cursor; the label becomes "Waiting for your wallet…".
- **Quiet:** transparent with a hairline ring and strong ink; used for "Disconnect" in the account menu.
- **Compact:** in the header the same actions render as a pill instead of a button.

### Pills
- **Style:** panel fill, inset hairline ring, 13px strong ink, padding 0.5em 1em, gap 0.55em. Carries a 6px status dot (network) or 8px dot (wallet) before its text.
- **Hover:** ring darkens to faint ink over 0.2s.
- **Icon pill:** the theme toggle, 10px padding around a 16px half-filled circle.

### Navigation
Section tabs centred in the header from 768px: 14px muted ink pills (8px by 14px padding). Hover lifts to strong ink; the current page takes the panel fill and hairline ring and is marked with `aria-current="page"`.

### The Seal (signature)
The session's state in brackets: `[ CLOSED ] · OPENS 26H`. Mono 11px at 0.09em, panel fill, hairline ring, 2px corners, tabular figures, never wraps. While the session loads it shows `[ ······ ]`. An UNKNOWN session turns the seal amber. It carries `role="status"` with a full spoken label; the bracketed text is hidden from assistive technology.

### The Thread (signature)
- **Stem:** 2px navy, full column height, draws top-down over 1.1s (scaleY from 0).
- **Bar:** 2px navy, rotated -8°, starting one stem-gap plus an overhang (1.25rem, 2.25rem from 640px) left of the content, so it crosses the stem just above the figure. Draws left to right over 1s after a 0.35s delay.
- Both are decorative and hidden from assistive technology; both hold still after drawing.

### Band Rows (signature)
- **Row:** 60px, top rule, a 1px rule reaching from the stem to just short of the ticker.
- **Track:** 18px tall. The band span is the navy band-area tint bounded by 1px navy edges, inset symmetrically by the half-width against a 400 bps full scale; a 2px strong-ink mid tick extends 3px past the track on both sides. Span edges glide over 0.6s when the band moves.
- **Unpriced:** the track becomes a 1px rule and the centre shows an em dash.
- **State:** label type, `STATE · n live`; HALTED and DEGRADED turn amber.
- **Figures:** centre price in strong 600 tabular; half-width as `±0.55%` in 12px muted tabular beneath.
- **Loading:** sunk-linen skeleton blocks pulsing to 55% opacity every 1.6s. **Error:** one muted 13px sentence spanning the row.
- Built as a real table with a screen-reader caption and column headers.

### The Notches (signature)
The smart account's free actions, drawn into the mark. Three 2px navy strokes, 1.7rem long, cross the stem at the bar's -8° beside the account, 9px apart; each one goes out to the rule colour as its action is spent, over 0.5s. They draw once on arrival, left to right, staggered by 80ms, and hold still; while the account loads they rest in sunk linen. Decorative and hidden from assistive technology: the line under the seal says the same in words.
- **Seal:** the account's address in the session seal's shape, `[ 0x42CE…E870 ] · NOT CREATED` in muted ink, or `· LIVE` in ledger green once it exists.
- **Line:** one 14px body sentence: "Creating it costs you nothing, and so do your first 3 actions.", then "2 free actions left", then "Free actions used. You pay gas from here."
- **Action:** "Create account" as the view's one filled button, only while the account does not exist; it reads "Waiting for your wallet…" while the owner signs and "Creating your account…" while the bundler lands it.
- The block loads only once a wallet connects, so it costs the first load nothing.

### Menus
Native popovers anchored below the header at the right: 256px wide, panel fill, 10px corners, popover lift. The wallet list shows each connector with its 22px icon; options are 14px strong ink, 6px corners, sunk-linen on hover. The account menu shows a label, the full address in 12.5px mono, and a full-width quiet Disconnect button.

### Theming
Light and dark ship together. The theme follows `prefers-color-scheme` until the toggle sets `data-theme` on the root, which is remembered in local storage and applied by an inline script before first paint. Browser chrome follows with theme colours matching bone and night paper.

## Do's and Don'ts

### Do:
- **Do** hang every new account element off the stem: a single column, a 1px rule from the stem where a row needs anchoring.
- **Do** keep the bar at -8° and the stem and bar at 2px navy; let them draw once on load and then stay still. The notches share that angle and weight.
- **Do** set every money figure with the tabular `num` treatment and every address in mono with slashed zero.
- **Do** state market and wallet state in words (HALTED, CLOSED, "Your wallet is on another network"); colour only reinforces.
- **Do** build text hierarchy from ink alpha (100 / 80 / 62 / 42%) and depth from paper, panel and sunk steps.
- **Do** keep primary actions at least 44px tall and reserve their row so state changes don't shift the column.
- **Do** disable the draw, swap, pulse, glide and notch motions under `prefers-reduced-motion: reduce`.

### Don't:
- **Don't** lay the account out as a grid of equal tiles or add a side rail of menus.
- **Don't** use liquidation red for anything but the live liquidation-price figure.
- **Don't** introduce greys, cool neutrals or a second accent hue.
- **Don't** add drop shadows to controls or rows; only popover menus lift.
- **Don't** use Archivo for headings or figures; it belongs to the wordmark.
- **Don't** show illustrative or placeholder figures; an unpriced asset shows an em dash and an empty account says "None yet".
