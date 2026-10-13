// SPDX-License-Identifier: MIT OR Apache-2.0
import Link from "next/link";
import type { ReactNode } from "react";

const TITLES = { note: "Note", warning: "Warning", danger: "Danger", tip: "Tip" } as const;

export function Callout({ kind = "note", title, children }: { kind?: keyof typeof TITLES; title?: string; children: ReactNode }) {
  return (
    <aside className="callout" data-kind={kind}>
      <div className="callout-title">{title ?? TITLES[kind]}</div>
      {children}
    </aside>
  );
}

export function Steps({ children }: { children: ReactNode }) {
  return <div className="steps">{children}</div>;
}

export function Step({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="step">
      <h3 className="step-title" style={{ margin: 0 }}>{title}</h3>
      {children}
    </div>
  );
}

export function Cards({ children }: { children: ReactNode }) {
  return <div className="cards">{children}</div>;
}

export function Card({ title, href, children }: { title: string; href: string; children: ReactNode }) {
  const external = /^https?:/.test(href);
  const body = (
    <>
      <div className="card-title">{title}</div>
      <div className="card-body">{children}</div>
    </>
  );
  return external ? (
    <a className="card" href={href} rel="noreferrer">
      {body}
    </a>
  ) : (
    <Link className="card" href={href}>
      {body}
    </Link>
  );
}
