// SPDX-License-Identifier: MIT OR Apache-2.0
import type { ComponentProps } from "react";
import { LANG_LABEL } from "@/lib/highlight";
import { CopyButton } from "./CopyButton";

type Props = ComponentProps<"pre"> & { "data-lang"?: string; "data-title"?: string; "data-raw"?: string };

/** A highlighted code block with its title or language and a copy button. */
export function Pre({ children, className, style, ...rest }: Props) {
  const lang = rest["data-lang"] ?? "text";
  const raw = rest["data-raw"] ?? "";
  const title = rest["data-title"];
  return (
    <div className="code">
      <div className="code-head">
        <span className="code-label">{title ?? LANG_LABEL[lang] ?? lang}</span>
        <CopyButton text={raw} />
      </div>
      <pre className={className} style={style} tabIndex={0}>
        {children}
      </pre>
    </div>
  );
}
