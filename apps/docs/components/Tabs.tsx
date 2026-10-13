// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { Children, isValidElement, useId, useState, useSyncExternalStore, type ReactElement, type ReactNode } from "react";

const listeners = new Set<() => void>();
const memory = new Map<string, string>();

function key(group: string) {
  return `tapehouse-docs-tab:${group}`;
}

function read(group: string): string | null {
  if (memory.has(group)) return memory.get(group) ?? null;
  try {
    return localStorage.getItem(key(group));
  } catch {
    return null;
  }
}

function write(group: string, title: string) {
  memory.set(group, title);
  try {
    localStorage.setItem(key(group), title);
  } catch {}
  listeners.forEach((l) => l());
}

function subscribe(listener: () => void) {
  listeners.add(listener);
  return () => void listeners.delete(listener);
}

export function Tab({ children }: { title: string; children: ReactNode }) {
  return <>{children}</>;
}

/**
 * Tabs whose choice follows the reader: every group of the same name on the site
 * shows the title the reader picked last, where it has one.
 */
export function Tabs({ group, children }: { group?: string; children: ReactNode }) {
  const id = useId();
  const tabs = Children.toArray(children).filter(
    (c): c is ReactElement<{ title: string; children: ReactNode }> => isValidElement(c) && typeof (c.props as { title?: unknown }).title === "string",
  );
  const stored = useSyncExternalStore(
    subscribe,
    () => (group ? read(group) : null),
    () => null,
  );
  const [local, setLocal] = useState(tabs[0]?.props.title ?? "");
  const wanted = group ? stored : local;
  const active = tabs.find((t) => t.props.title === wanted) ?? tabs[0];
  if (!active) return null;
  return (
    <div className="tabs">
      <div className="tablist" role="tablist">
        {tabs.map((t, i) => {
          const selected = t === active;
          return (
            <button
              key={t.props.title}
              id={`${id}-tab-${i}`}
              role="tab"
              type="button"
              className="tab"
              aria-selected={selected}
              aria-controls={`${id}-panel`}
              tabIndex={selected ? 0 : -1}
              onClick={() => (group ? write(group, t.props.title) : setLocal(t.props.title))}
              onKeyDown={(e) => {
                const step = e.key === "ArrowRight" ? 1 : e.key === "ArrowLeft" ? -1 : 0;
                if (!step) return;
                const next = tabs[(i + step + tabs.length) % tabs.length].props.title;
                if (group) write(group, next);
                else setLocal(next);
                document.getElementById(`${id}-tab-${(i + step + tabs.length) % tabs.length}`)?.focus();
              }}
            >
              {t.props.title}
            </button>
          );
        })}
      </div>
      <div className="tabpanel" role="tabpanel" id={`${id}-panel`}>
        {active.props.children}
      </div>
    </div>
  );
}
