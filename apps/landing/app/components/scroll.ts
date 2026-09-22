/** Run `fn` once per animation frame while the page scrolls or resizes. */
export function onScrollFrame(fn: () => void) {
  let frame = 0;
  const request = () => {
    if (!frame)
      frame = requestAnimationFrame(() => {
        frame = 0;
        fn();
      });
  };
  fn();
  window.addEventListener("scroll", request, { passive: true });
  window.addEventListener("resize", request);
  return () => {
    window.removeEventListener("scroll", request);
    window.removeEventListener("resize", request);
    if (frame) cancelAnimationFrame(frame);
  };
}

export const clamp01 = (v: number) => Math.min(1, Math.max(0, v));
export const prefersStill = () => window.matchMedia("(prefers-reduced-motion: reduce)").matches;
