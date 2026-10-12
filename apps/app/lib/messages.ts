// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Messages } from "@/messages";

type Paths<T, Prefix extends string = ""> = {
  [K in keyof T & string]: T[K] extends string ? `${Prefix}${K}` : Paths<T[K], `${Prefix}${K}.`>;
}[keyof T & string];

/** The key of a message, namespace first: `account.deposit.title`. Only a key the English catalog has is one. */
export type MessageKey = Paths<Messages>;

/**
 * The values a message is filled with: text the language does not change, such as an asset, numbers already formatted,
 * or another message, which is written first, as a duration is inside the sentence that states it.
 */
export type Values = Readonly<Record<string, string | number | Msg>>;

/**
 * A message not yet in a language: the key and the values that fill it. The functions that decide what a visitor is
 * told return these, so that what they decide is tested whatever the language, and the language is the screen's.
 */
export type Msg = { readonly key: MessageKey; readonly values?: Values };

export function msg(key: MessageKey, values?: Values): Msg {
  return values ? { key, values } : { key };
}

/** Text that is not a message: a figure, a name, a time, as it is. */
export function plain(text: string): Msg {
  return msg("common.plain", { text });
}

/** Writes a message with `translate`, the messages inside its values first. */
export function render(translate: (key: string, values?: Record<string, string | number>) => string, message: Msg): string {
  const values = message.values
    ? Object.fromEntries(
        Object.entries(message.values).map(([name, value]) => [name, typeof value === "object" ? render(translate, value) : value]),
      )
    : undefined;
  return translate(message.key, values);
}

/** An error whose words are a message, so that the screen that shows it writes it in the visitor's language. */
export class MessageError extends Error {
  constructor(readonly msg: Msg) {
    super(msg.key);
    this.name = "MessageError";
  }
}
