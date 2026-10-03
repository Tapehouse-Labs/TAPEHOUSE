// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

/** Switches between the light and dark themes and remembers the choice. */
export function ThemeToggle() {
  const toggle = () => {
    const root = document.documentElement;
    const dark = root.dataset.theme
      ? root.dataset.theme === "dark"
      : window.matchMedia("(prefers-color-scheme: dark)").matches;
    const next = dark ? "light" : "dark";
    root.dataset.theme = next;
    try {
      localStorage.setItem("theme", next);
    } catch {}
  };
  return (
    <button type="button" onClick={toggle} aria-label="Switch between light and dark" className="pill !p-2.5">
      <svg viewBox="0 0 16 16" aria-hidden="true" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="1.5">
        <circle cx="8" cy="8" r="5.25" />
        <path d="M8 2.75 A5.25 5.25 0 0 1 8 13.25 Z" fill="currentColor" stroke="none" />
      </svg>
    </button>
  );
}
