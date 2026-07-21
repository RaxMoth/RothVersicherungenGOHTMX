// Package handler is the Vercel serverless entrypoint. The whole site
// runs as one function: vercel.json rewrites every route here, and the
// request is served by the same mux as the regular binary
// (cmd/server), just without ListenAndServe.
package handler

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/maxroth/eumel/internal/config"
	"github.com/maxroth/eumel/internal/i18n"
	"github.com/maxroth/eumel/internal/server"
	"github.com/maxroth/eumel/internal/view"
)

var (
	once    sync.Once
	app     http.Handler
	initErr error
)

func setup() {
	cfg := config.Load()
	// The serverless filesystem is read-only: force embedded assets and
	// skip the SQLite database (the site is pure content; nothing reads it).
	cfg.Env = "prod"
	cfg.Dev = false

	tr, err := i18n.Load(cfg.DefaultLang)
	if err != nil {
		initErr = err
		return
	}
	v, err := view.New(cfg.Dev, tr)
	if err != nil {
		initErr = err
		return
	}
	app = server.New(cfg, nil, v, tr)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(setup)
	if initErr != nil {
		slog.Error("vercel init", "err", initErr)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	app.ServeHTTP(w, r)
}
