# Tapehouse documentation

The site at docs.tapehouse.xyz, built with Next.js. Pages are MDX files under `content/`, listed in `lib/nav.ts`. The indexer API reference is generated from `services/internal/api/openapi.yaml`, and `contracts/addresses` from `deployments/46630.json`, both read when the site is built.

```bash
pnpm --filter docs dev
pnpm --filter docs test
```

Every page is also served as Markdown at its own address with `.md` appended, and `/llms.txt` and `/llms-full.txt` list them for agents. `lib/content.test.ts` fails on a page missing from the navigation, a broken internal link, a dash in prose or a banned word.
