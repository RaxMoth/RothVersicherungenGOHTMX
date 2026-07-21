# Roth Versicherungen — Go + HTMX

Server-rendered website for Roth Versicherungen / Roth Finanz (Langen), built on the Eumel template: stdlib router + `html/template`, HTMX (vendored in `web/static/js/`), Tailwind v4 standalone CLI, SQLite (modernc, no CGO), JSON i18n.

## Commands

- `make dev` — build CSS + run server in dev mode (templates/static re-read from disk per request)
- `make test` — run tests; `make build` — production binary with everything embedded
- `make css` — rebuild Tailwind output (needed after adding new utility classes)

## Conventions

- **Never hardcode UI text** in templates or handlers. Every string goes into `locales/de.json` (flat JSON, dot keys like `nav.team`) and is used via `{{t "key"}}` in templates. Args make it a Sprintf: `{{t "footer.copyright" year}}`. Lists use numbered keys (`x.items.1`, `x.items.2`, …) iterated with `{{range tlist "x.items"}}`.
- The site is German-only; `DEFAULT_LANG` defaults to `de`.
- Pages live in `web/templates/pages/` and define a `content` block (plus optional `title`/`description` blocks) rendered inside `layouts/base.html`. All pages are static content: routes register with `h.Page("name.html")` in `internal/server/server.go` (Go 1.22+ patterns).
- Shared building blocks are partials in `web/templates/partials/`: `page-hero`, `section-head`, `cta`, `link-card`, `nav`, `footer`, `legal-address`, `legal-text-section`. They take named parameters via the `dict` template function. HTMX endpoints (none yet) would render partials via `View.RenderPartial`.
- Brand theme (colors `brand-red`/`brand-page`/…, `shadow-card`, `rounded-4xl`, Inter font) is defined in `@theme` in `web/static/css/input.css`. Header dropdowns/mobile menu are progressive enhancement in `web/static/js/nav.js`.
- Migrations: sequential `internal/db/migrations/NNNN_name.sql`, auto-applied at startup. The DB is currently unused (0001 is a placeholder); query functions go in `internal/db/` as plain functions taking `*sql.DB`.
- Config is env-vars only (`internal/config`); defaults must keep the server runnable with zero setup.
