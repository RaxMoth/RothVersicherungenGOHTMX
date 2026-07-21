// Package server wires routes, middleware and handlers together.
package server

import (
	"database/sql"
	"io/fs"
	"net/http"

	"handler/internal/config"
	"handler/internal/handlers"
	"handler/internal/i18n"
	"handler/internal/view"
	"handler/web"
)

func New(cfg *config.Config, database *sql.DB, v *view.View, tr *i18n.Translator) http.Handler {
	h := &handlers.Handler{DB: database, View: v}

	mux := http.NewServeMux()

	// Pages — every route renders a static page template; all content
	// comes from locales/de.json.
	mux.HandleFunc("GET /{$}", h.Page("home.html"))

	mux.HandleFunc("GET /roth-versicherungen", h.Page("versicherungen.html"))
	mux.HandleFunc("GET /roth-versicherungen/firmenkunden", h.Page("firmenkunden.html"))
	mux.HandleFunc("GET /roth-versicherungen/firmenkunden/cyber-police", h.Page("cyber.html"))
	mux.HandleFunc("GET /roth-versicherungen/privatkunden", h.Page("privatkunden.html"))
	mux.HandleFunc("GET /roth-versicherungen/privatkunden/tierkrankenversicherung", h.Page("tier.html"))
	mux.HandleFunc("GET /roth-versicherungen/wichtige-hinweise", h.Page("hinweise.html"))
	mux.HandleFunc("GET /roth-versicherungen/jobs", h.Page("jobs.html"))
	mux.HandleFunc("GET /roth-versicherungen/erstinformation", h.Page("vers-erstinformation.html"))
	mux.HandleFunc("GET /roth-versicherungen/datenschutz", h.Page("vers-datenschutz.html"))
	mux.HandleFunc("GET /roth-versicherungen/impressum", h.Page("vers-impressum.html"))

	mux.HandleFunc("GET /roth-finanz", h.Page("finanz.html"))
	mux.HandleFunc("GET /roth-finanz/altersversorgung", h.Page("altersversorgung.html"))
	mux.HandleFunc("GET /roth-finanz/sterbegeldversicherung", h.Page("sterbegeld.html"))
	mux.HandleFunc("GET /roth-finanz/erstinformation", h.Page("finanz-erstinformation.html"))
	mux.HandleFunc("GET /roth-finanz/datenschutz", h.Page("finanz-datenschutz.html"))
	mux.HandleFunc("GET /roth-finanz/impressum", h.Page("finanz-impressum.html"))

	mux.HandleFunc("GET /team", h.Page("team.html"))
	mux.HandleFunc("GET /kontakt-anfahrt", h.Page("kontakt.html"))
	mux.HandleFunc("GET /sitemap", h.Page("sitemap.html"))

	// Static assets: from disk in dev (live reload), embedded in prod.
	var staticFS http.FileSystem
	if cfg.Dev {
		staticFS = http.Dir("web/static")
	} else {
		sub, _ := fs.Sub(web.Static, "static")
		staticFS = http.FS(sub)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", cacheStatic(cfg, http.FileServer(staticFS))))

	// Operational
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Everything else is a 404 rendered through the error page.
	mux.HandleFunc("/", h.NotFound)

	var handler http.Handler = mux
	handler = Language(tr)(handler)
	handler = Logger(handler)
	handler = Recover(v)(handler)
	return handler
}

func cacheStatic(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Dev {
			w.Header().Set("Cache-Control", "public, max-age=86400")
		}
		next.ServeHTTP(w, r)
	})
}
