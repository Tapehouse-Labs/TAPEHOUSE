// SPDX-License-Identifier: MIT OR Apache-2.0
import Link from "next/link";
import type { ComponentProps } from "react";
import type { MDXComponents } from "mdx/types";
import { Callout, Card, Cards, Step, Steps } from "./blocks";
import { Deployments } from "./Deployments";
import { Pre } from "./Pre";
import { Tab, Tabs } from "./Tabs";

function heading(Tag: "h2" | "h3" | "h4") {
  return function Heading({ id, children, ...rest }: ComponentProps<"h2">) {
    return (
      <Tag id={id} {...rest}>
        {children}
        {id ? (
          <a className="anchor" href={`#${id}`} aria-label="Link to this section">
            #
          </a>
        ) : null}
      </Tag>
    );
  };
}

function A({ href = "", children, ...rest }: ComponentProps<"a">) {
  if (/^https?:/.test(href)) {
    return (
      <a href={href} rel="noreferrer" {...rest}>
        {children}
      </a>
    );
  }
  if (href.startsWith("#") || href.startsWith("mailto:")) {
    return (
      <a href={href} {...rest}>
        {children}
      </a>
    );
  }
  return (
    <Link href={href} {...rest}>
      {children}
    </Link>
  );
}

export const components: MDXComponents = {
  h2: heading("h2"),
  h3: heading("h3"),
  h4: heading("h4"),
  a: A,
  pre: Pre as MDXComponents["pre"],
  table: (props) => (
    <div className="table-wrap">
      <table {...props} />
    </div>
  ),
  Callout,
  Note: (p: { title?: string; children: React.ReactNode }) => <Callout kind="note" {...p} />,
  Warning: (p: { title?: string; children: React.ReactNode }) => <Callout kind="warning" {...p} />,
  Tip: (p: { title?: string; children: React.ReactNode }) => <Callout kind="tip" {...p} />,
  Steps,
  Step,
  Cards,
  Card,
  Tabs,
  Tab,
  Deployments,
};
