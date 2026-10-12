// SPDX-License-Identifier: MIT OR Apache-2.0
import en from "./en";
import es from "./es";
import pt from "./pt";

/** Every message of every language, by language. English is the reference the others are held to. */
export const catalogs = { en, es, pt };

export type Messages = typeof en;
