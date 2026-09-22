"use client";

import { useEffect, useRef } from "react";
import { clamp01, onScrollFrame, prefersStill } from "./scroll";

/**
 * The page as the mark at full scale: a bar rising at eight degrees crosses the
 * head of a stem, and the stem runs the length of the story as it is read.
 */
export function Thread() {
  const el = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const node = el.current;
    const host = node?.parentElement;
    if (!node || !host) return;
    if (prefersStill()) {
      node.style.setProperty("--thread", "1");
      node.style.setProperty("--thread-bar", "1");
      return;
    }
    return onScrollFrame(() => {
      const r = host.getBoundingClientRect();
      const vh = window.innerHeight;
      const read = vh * 0.62 - r.top;
      node.style.setProperty("--thread", clamp01(read / r.height).toFixed(4));
      node.style.setProperty("--thread-bar", clamp01(read / (vh * 0.55)).toFixed(4));
    });
  }, []);

  return (
    <div ref={el} className="thread" aria-hidden="true">
      <div className="thread-bar" />
      <div className="thread-stem" />
    </div>
  );
}
