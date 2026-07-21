# Roth Versicherungen — Go + HTMX

Server-rendered website for Roth Versicherungen Maklergesellschaft m.b.H. and Roth Finanz Maklergesellschaft m.b.H. in Langen, converted from the original React SPA.

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
  handlers/            HTTP handlers — every page renders a static template
  db/                  SQLite setup, embedded migrations (currently unused, kept for future features)
  i18n/                String loading + per-request language resolution
  view/                Template renderer (dev: live reload, prod: parsed once, embedded)
locales/               All site content and UI strings: de.json (flat JSON, dot keys)
web/
  templates/layouts/   base.html — the page shell (header, footer, meta)
  templates/pages/     One file per page, fills the "content" block
  templates/partials/  nav, footer, page-hero, section-head, cta, link-card, legal helpers
  static/              css/ (Tailwind in+out), js/ (vendored htmx + nav.js), img/ (site images)
```

## Pages

| Route | Template |
| --- | --- |
| `/` | home.html |
| `/roth-versicherungen` (+ firmenkunden, cyber-police, privatkunden, tierkrankenversicherung, wichtige-hinweise, jobs, erstinformation, datenschutz, impressum) | versicherungen.html … |
| `/roth-finanz` (+ altersversorgung, sterbegeldversicherung, erstinformation, datenschutz, impressum) | finanz.html … |
| `/team`, `/kontakt-anfahrt`, `/sitemap` | team.html, kontakt.html, sitemap.html |

## How to…

### Add a page

1. Create `web/templates/pages/name.html` with `{{define "content"}}` (plus optional `title`/`description` blocks).
2. Register the route in `internal/server/server.go`: `mux.HandleFunc("GET /path", h.Page("name.html"))`.
3. Add its strings to `locales/de.json`.

### Edit content

All visible text lives in `locales/de.json`. Lists use numbered keys (`x.items.1`, `x.items.2`, …) and are rendered with the `tlist` template function; changing text needs no template edits.

### Add a migration

Create `internal/db/migrations/0002_something.sql`. Applied automatically at startup, tracked in `schema_migrations`.

### Deploy

```sh
make build      # builds CSS, compiles bin/server with templates/static/locales/migrations embedded
ENV=prod ./bin/server
```

The binary is fully self-contained — copy it to the server and run it. Configuration via env vars (see `.env.example`).

### Deploy to Vercel

The repo also works on Vercel as a single serverless function: `vercel.json` rewrites every route to `api/index.go`, which serves the same mux as the binary (embedded assets, no SQLite — the serverless filesystem is read-only). Just push; no build command or env vars needed. Remember to commit `web/static/css/output.css` after running `make css`, since Vercel does not run the Tailwind build.

## Make targets

Run `make help` to list them: `dev`, `run`, `build`, `css`, `css-watch`, `tailwind`, `htmx` (update vendored htmx), `test`, `tidy`, `clean`.
