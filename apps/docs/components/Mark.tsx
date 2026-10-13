// SPDX-License-Identifier: MIT OR Apache-2.0
/** The continuity mark: a support, and a line that crosses it and keeps rising. */
export function Mark({ className }: { className?: string }) {
  return (
    <svg viewBox="10 26 100 82" className={className} aria-hidden="true" focusable="false" fill="currentColor">
      <path d="M14 42.94 L106 30.01 L106 41.12 L14 54.05 Z" />
      <path d="M40.5 42 H51.5 V104 H40.5 Z" />
    </svg>
  );
}
