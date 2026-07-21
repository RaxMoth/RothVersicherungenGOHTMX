# Eumel — Go + HTMX template

Base template for server-rendered Go websites: stdlib router + `html/template`, HTMX (vendored in `web/static/js/`), Tailwind v4 standalone CLI, SQLite (modernc, no CGO), JSON i18n.

## Commands

- `make dev` — build CSS + run server in dev mode (templates/static re-read from disk per request)
- `make test` — run tests; `make build` — production binary with everything embedded
- `make css` — rebuild Tailwind output (needed after adding new utility classes)

## Conventions

- **Never hardcode UI text** in templates or handlers. Every string goes into ALL `locales/*.json` files (flat JSON, dot keys like `nav.home`) and is used via `{{t "key"}}` in templates. Args make it a Sprintf: `{{t "todos.count" 5}}`.
- Pages live in `web/templates/pages/` and define a `content` block rendered inside `layouts/base.html` via `View.Render(w, r, status, "name.html", data)`. Page data is accessed as `.Data.X` in templates.
- HTMX endpoints render partials from `web/templates/partials/` via `View.RenderPartial`; mutations re-render the affected partial and the client swaps it with `hx-target`/`hx-swap="outerHTML"`.
- Routes are registered in `internal/server/server.go` using Go 1.22+ patterns (`GET /path/{id}`).
- Migrations: sequential `internal/db/migrations/NNNN_name.sql`, auto-applied at startup. Query functions live in `internal/db/` as plain functions taking `*sql.DB`.
- Config is env-vars only (`internal/config`); defaults must keep the server runnable with zero setup.
- The todo feature (handlers/todos.go, db/todos.go, partials/todo-list.html, migration 0001) is demo code meant to be replaced in real projects.
