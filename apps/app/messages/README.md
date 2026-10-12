<!-- SPDX-License-Identifier: CC0-1.0 -->

# Messages

The app speaks English, Spanish and Portuguese. A visitor's language is the `locale` cookie, set by the language
switch, or else the best match of the browser's `Accept-Language`; the URL does not change.

## Layout

`messages/<language>/<namespace>.json`, one file per area of the app, with the same keys in every language.
`messages/en` is the source: its keys are the type of every message key in the code (`MessageKey`), so a key that
does not exist in English does not compile. `catalogs.test.ts` fails when a language lacks a key or has an extra one,
when a message is not valid ICU MessageFormat, or when a translation does not use the same placeholders.

## Writing a message

- Components write their own strings with `useTranslations("<namespace>")`.
- Functions that decide what a visitor is told (`lib/*.ts`) return a `Msg`, `{ key, values }`, built with `msg()`. A
  component writes it with `useMsg()`, a test with `say(msg, locale)`. What is decided is tested in every language
  alike, and the language is the screen's.
- A value is text the language does not change (an asset, a figure already formatted) or another `Msg`, which is
  written first.
- Counts use ICU plurals: `{days, plural, one {# day} other {# days}}`. Text with markup uses `t.rich`.

## Numbers, dates and typed amounts

Every figure goes through `useFormat()`, or `formatters(locale)` outside a component. The arithmetic is exact, on
bigints; only the separators change. A typed amount is read with the language's decimal separator alone
(`parseAmount`), so `1.500` is never read as `1,5` or the reverse.

## Glossary

| English | Spanish | Portuguese |
| --- | --- | --- |
| band | banda | banda |
| margin | margen | margem |
| requirement | requerimiento | requisito |
| equity | capital | patrimônio |
| loan | préstamo | empréstimo |
| basket | cesta | cesta |
| gap cover | cobertura de brecha | cobertura de gap |
| vault | bóveda | cofre |
| issuer | emisor | emissor |
| guardian | guardián | guardião |
| wallet | wallet | carteira |
| smart account | cuenta inteligente | conta inteligente |
| short | corto | venda a descoberto |

Stock Token, tickers, USDG, WETH and contract or error names stay as they are. Spanish speaks to the visitor as
`tú`, Portuguese is Brazilian.
