// Package server wires routes, middleware and handlers together.
package server

import (
	"database/sql"
	"io/fs"
	"net/http"

	"github.com/maxroth/eumel/internal/config"
	"github.com/maxroth/eumel/internal/handlers"
	"github.com/maxroth/eumel/internal/i18n"
	"github.com/maxroth/eumel/internal/view"
	"github.com/maxroth/eumel/web"
)

func New(cfg *config.Config, database *sql.DB, v *view.View, tr *i18n.Translator) http.Handler {
	h := &handlers.Handler{DB: database, View: v}

	mux := http.NewServeMux()

	// Pages
	mux.HandleFunc("GET /{$}", h.Home)
	mux.HandleFunc("GET /about", h.About)

	// HTMX demo endpoints (todos)
	mux.HandleFunc("POST /todos", h.CreateTodo)
	mux.HandleFunc("POST /todos/{id}/toggle", h.ToggleTodo)
	mux.HandleFunc("DELETE /todos/{id}", h.DeleteTodo)

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
