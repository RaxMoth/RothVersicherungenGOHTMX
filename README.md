# Eumel — Go + HTMX Website Template

A batteries-included base template for server-rendered websites. Copy it, rename it, build on it.

**Stack:** Go stdlib (router + `html/template`) · [HTMX](https://htmx.org) (vendored) · [Tailwind CSS v4](https://tailwindcss.com) (standalone CLI, no Node.js) · SQLite (pure-Go driver, no CGO) · JSON-based i18n

## Quick start

```sh
make dev        # downloads the Tailwind binary on first run, builds CSS, starts the server
```

Open http://localhost:8080. In dev mode, templates and static files are re-read from disk on every request — edit and reload, no restart needed. Run `make css-watch` in a second terminal if you're changing Tailwind classes.

## Project layout

```
cmd/server/            Entrypoint: config, DB, i18n, renderer, graceful shutdown
internal/
  config/              Env-var configuration (ADDR, ENV, DB_PATH, DEFAULT_LANG)
  server/              Routes + middleware (language, logging, panic recovery)
  handlers/            HTTP handlers — pages render full templates, HTMX endpoints render partials
  db/                  SQLite setup, embedded migrations, query functions
  db/migrations/       0001_*.sql, 0002_*.sql, ... applied in order at startup
  i18n/                String loading + per-request language resolution
  view/                Template renderer (dev: live reload, prod: parsed once, embedded)
locales/               All UI strings: en.json, de.json, ... (flat JSON, dot keys)
web/
  templates/layouts/   base.html — the page shell
  templates/pages/     One file per page, fills the "content" block
  templates/partials/  Shared snippets, also rendered standalone for HTMX responses
  static/              css/ (Tailwind in+out), js/ (vendored htmx), img/
```

## How to…

### Add a page

1. Create `web/templates/pages/contact.html` with a `{{define "content"}}` block.
2. Add a handler in `internal/handlers/` calling `h.View.Render(w, r, http.StatusOK, "contact.html", data)`.
3. Register the route in `internal/server/server.go`.
4. Add its strings to every file in `locales/`.

### Add an HTMX endpoint

Handler renders a partial instead of a page: `h.View.RenderPartial(w, r, "my-partial", data)`. In the template, point at it with `hx-post`/`hx-get` + `hx-target`. See the todo demo (`internal/handlers/todos.go`, `web/templates/partials/todo-list.html`) for the pattern.

### Add a language

Drop `locales/fr.json` with the same keys as `en.json`. That's it — it's picked up at startup and appears in the nav switcher. Strings are `fmt.Sprintf` format strings when you pass arguments: `{{t "todos.count" 5}}`.

### Add a migration

Create `internal/db/migrations/0002_something.sql`. Applied automatically at startup, tracked in `schema_migrations`.

### Deploy

```sh
make build      # builds CSS, compiles bin/server with templates/static/locales/migrations embedded
ENV=prod ./bin/server
```

The binary is fully self-contained — copy it to the server and run it. Configuration via env vars (see `.env.example`).

## Starting a new project from this template

1. Copy the repo, then rename the module: `go mod edit -module github.com/you/newproject && grep -rl maxroth/eumel --include='*.go' . | xargs sed -i '' 's|github.com/maxroth/eumel|github.com/you/newproject|g'`
2. Replace the todo demo: delete `internal/handlers/todos.go`, `internal/db/todos.go`, the demo section in `home.html`, `partials/todo-list.html`, and write your own `0001_*.sql`.
3. Update `locales/*.json` with your site name and strings.

## Make targets

Run `make help` to list them: `dev`, `run`, `build`, `css`, `css-watch`, `tailwind`, `htmx` (update vendored htmx), `test`, `tidy`, `clean`.
