// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { isLocale, LANGUAGE_NAMES, LOCALES, negotiate } from "./config";

describe("negotiate", () => {
  test("takes the language a browser prefers among the app's, by its primary subtag", () => {
    expect(negotiate("es-CO,es;q=0.9,en;q=0.8")).toBe("es");
    expect(negotiate("pt-BR,pt;q=0.9")).toBe("pt");
    expect(negotiate("EN-us")).toBe("en");
  });

  test("goes by quality, then by order, and skips what the app lacks", () => {
    expect(negotiate("fr;q=0.9,es;q=0.8,pt;q=0.7")).toBe("es");
    expect(negotiate("pt;q=0.5,es;q=0.5")).toBe("pt");
    expect(negotiate("de,fr,es;q=0.1")).toBe("es");
    expect(negotiate("es;q=0,pt;q=0.1")).toBe("pt");
  });

  test("falls back to English where none matches, or the header is absent or malformed", () => {
    for (const header of ["fr,de;q=0.9", "", "*", ";;;", "es;q=abc", null, undefined]) {
      expect(negotiate(header)).toBe("en");
    }
  });
});

describe("locales", () => {
  test("each has a name in itself, and only they are locales", () => {
    for (const locale of LOCALES) {
      expect(isLocale(locale)).toBe(true);
      expect(LANGUAGE_NAMES[locale].length).toBeGreaterThan(1);
    }
    expect(isLocale("fr")).toBe(false);
    expect(isLocale(undefined)).toBe(false);
    expect(isLocale("")).toBe(false);
  });
});
