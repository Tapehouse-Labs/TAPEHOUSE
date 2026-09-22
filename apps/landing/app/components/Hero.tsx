"use client";

import Link from "next/link";
import { useEffect, useRef } from "react";
import { APP_URL } from "../site";
import { Mark } from "./Mark";
import { Seal } from "./Seal";
import { PortfolioScreen } from "./PortfolioScreen";


/** Where the lit screen sits inside the scene plate, as fractions of the plate. */
const PLATE = { aspect: 1536 / 1024, screenX: 0.469, screenY: 0.658, screenW: 0.079 };

const clamp = (v: number) => Math.min(1, Math.max(0, v));
const span = (p: number, from: number, to: number) => clamp((p - from) / (to - from));
const easeInOut = (t: number) => (t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2);
const easeOut = (t: number) => 1 - Math.pow(1 - t, 3);

export function Hero() {
  const section = useRef<HTMLElement>(null);
  const stage = useRef<HTMLDivElement>(null);
  const screen = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const root = section.current;
    const st = stage.current;
    const sc = screen.current;
    if (!root || !st || !sc) return;
    const still = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    let frame = 0;
    const update = () => {
      frame = 0;
      const vw = st.clientWidth;
      const vh = still ? window.innerHeight : st.clientHeight;
      const travel = root.offsetHeight - vh;
      const p = !still && travel > 0 ? clamp((window.scrollY - root.offsetTop) / travel) : 0;

      // Lay the plate out by hand so the lit screen lands where the copy leaves room for it.
      // Portrait viewports zoom the plate and anchor it to the top, which drops the screen
      // below the actions and centres it; landscape is a plain centred cover.
      const portrait = vw / vh < 0.8;
      const zoom = portrait ? 1.2 : 1;
      const plateW = Math.max(vw, vh * PLATE.aspect) * zoom;
      const plateH = plateW / PLATE.aspect;
      const plateX = portrait
        ? Math.min(0, Math.max(vw - plateW, vw / 2 - PLATE.screenX * plateW))
        : (vw - plateW) / 2;
      const plateY = portrait ? 0 : (vh - plateH) / 2;
      const sx = plateX + PLATE.screenX * plateW;
      const sy = plateY + PLATE.screenY * plateH;
      const screenW = PLATE.screenW * plateW;

      st.style.setProperty("--plate-w", `${plateW.toFixed(1)}px`);
      st.style.setProperty("--plate-h", `${plateH.toFixed(1)}px`);
      st.style.setProperty("--plate-x", `${plateX.toFixed(1)}px`);
      st.style.setProperty("--plate-y", `${plateY.toFixed(1)}px`);
      if (still) return;

      // A night patch over the monitor's glow, so the push-in ends on a clean field.
      const gw = screenW * 1.3;
      const gh = (screenW / 1.6) * 1.45;
      st.style.setProperty("--gw", `${gw.toFixed(1)}px`);
      st.style.setProperty("--gh", `${gh.toFixed(1)}px`);
      st.style.setProperty("--gx", `${(sx - gw / 2).toFixed(1)}px`);
      st.style.setProperty("--gy", `${(sy - gh / 2).toFixed(1)}px`);
      st.style.setProperty("--glow-off", span(p, 0.52, 0.8).toFixed(3));

      const pw = sc.offsetWidth;
      const ph = sc.offsetHeight;
      const fit = Math.min(1, (vh - 112) / ph);
      const maxScale = (pw * fit) / screenW;

      const push = easeInOut(span(p, 0.04, 0.82));
      const scale = 1 + (maxScale - 1) * push;
      const tx = (vw / 2 - sx) * push;
      const ty = (vh / 2 + 16 - sy) * push;

      const ps = Math.min(fit, (screenW * scale) / pw);
      const cx = sx + tx;
      const cy = sy + ty;
      const fade = easeOut(span(p, 0.0, 0.2));

      st.style.setProperty("--ox", `${sx}px`);
      st.style.setProperty("--oy", `${sy}px`);
      st.style.setProperty("--tx", `${tx}px`);
      st.style.setProperty("--ty", `${ty}px`);
      st.style.setProperty("--scene-scale", scale.toFixed(4));
      st.style.setProperty("--scene-light", (1 - 0.88 * span(p, 0.4, 0.85)).toFixed(3));
      st.style.setProperty("--copy-alpha", (1 - fade).toFixed(3));
      st.style.setProperty("--copy-lift", (fade * 36).toFixed(1));
      st.style.setProperty("--ps", ps.toFixed(4));
      st.style.setProperty("--px", `${(cx - (pw * ps) / 2).toFixed(1)}px`);
      st.style.setProperty("--py", `${(cy - (ph * ps) / 2).toFixed(1)}px`);
      st.style.setProperty("--pa", easeOut(span(p, 0.3, 0.56)).toFixed(3));
      st.toggleAttribute("data-inside", p > 0.5);
    };
    const request = () => {
      if (!frame) frame = requestAnimationFrame(update);
    };

    update();
    window.addEventListener("scroll", request, { passive: true });
    window.addEventListener("resize", request);
    const observer = new ResizeObserver(request);
    observer.observe(sc);
    return () => {
      window.removeEventListener("scroll", request);
      window.removeEventListener("resize", request);
      observer.disconnect();
      if (frame) cancelAnimationFrame(frame);
    };
  }, []);

  return (
    <section ref={section} className="hero" aria-labelledby="hero-title">
      <div ref={stage} className="hero-stage">
        <div className="hero-scene" aria-hidden="true">
          <div className="hero-settle absolute inset-0">
            {/* Placeholder plate. The final asset is a filmed shot; swap for <video> with this as its poster. */}
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img src="/hero/floor-dawn.jpg" alt="" fetchPriority="high" decoding="async" />
            <div className="hero-tint" />
          </div>
          <div className="hero-glow-off" />
        </div>
        <div className="hero-scrim" />

        <header className="hero-copy relative z-10 mx-auto flex w-full max-w-[1320px] items-center justify-between gap-4 px-5 pt-5 sm:px-8 sm:pt-6">
          <Link href="/" className="flex items-center gap-3" aria-label="Tapehouse, home">
            <Mark className="h-[22px] w-auto" />
            <span className="font-display text-[15px] font-semibold tracking-[0.18em]">TAPEHOUSE</span>
          </Link>
          <div className="flex items-center gap-3 sm:gap-4">
            <span className="hidden sm:inline-flex">
              <Seal />
            </span>
            <a href={APP_URL} className="btn-ghost !py-[0.55em] text-[14px]">
              Open the app
            </a>
          </div>
        </header>

        <div className="hero-copy relative z-10 mx-auto flex max-w-[1280px] flex-col items-center px-5 pt-[clamp(3rem,11vh,7.5rem)] text-center sm:px-8">
          <h1
            id="hero-title"
            className="arrive text-balance font-display text-[clamp(2.4rem,4.9vw,4.25rem)] font-normal leading-[1.04] tracking-[-0.032em]"
            style={{ "--i": 0 } as React.CSSProperties}
          >
            <span className="block">Your whole portfolio.</span>{" "}
            <span className="block">One limit. Published in advance.</span>
          </h1>
          <p
            className="arrive mt-6 max-w-[58ch] text-pretty text-[clamp(1.02rem,1.35vw,1.22rem)] leading-[1.5] text-[rgb(245_242_236/0.84)]"
            style={{ "--i": 1 } as React.CSSProperties}
          >
            Borrow against every Stock Token you hold, with the liquidation price stated
            before you sign and a price that keeps moving when the market doesn’t.
          </p>
          <div
            className="arrive mt-9 flex flex-wrap items-center justify-center gap-3 text-[15px]"
            style={{ "--i": 2 } as React.CSSProperties}
          >
            <a href={APP_URL} className="btn-primary">
              Open the app
              <span className="font-mono text-[10.5px] uppercase tracking-[0.12em] opacity-60">Testnet</span>
            </a>
            <a href="#measured" className="btn-ghost">
              <span className="sm:hidden">What we measured</span>
              <span className="hidden sm:inline">See what we measured</span>
            </a>
          </div>
        </div>

        <p
          className="hero-copy arrive absolute inset-x-0 bottom-5 z-10 mx-auto max-w-[760px] px-5 text-center text-[12.5px] leading-snug text-[rgb(245_242_236/0.72)] sm:bottom-7"
          style={{ "--i": 4 } as React.CSSProperties}
        >
          Tapehouse is not a broker. A Stock Token is a debt security of its issuer, not a share: it follows
          the price and gives you no shareholder rights.
        </p>

        <div ref={screen} className="hero-screen z-20">
          <PortfolioScreen />
        </div>
      </div>
    </section>
  );
}
