// Package handlers contains all HTTP handlers. Full pages render
// through View.Render, HTMX endpoints through View.RenderPartial.
package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/maxroth/eumel/internal/db"
	"github.com/maxroth/eumel/internal/view"
)

type Handler struct {
	DB   *sql.DB
	View *view.View
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	todos, err := db.ListTodos(h.DB)
	if err != nil {
		h.Error(w, r, err)
		return
	}
	h.View.Render(w, r, http.StatusOK, "home.html", map[string]any{
		"Todos": todos,
	})
}

func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	h.View.Render(w, r, http.StatusOK, "about.html", nil)
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	h.View.Render(w, r, http.StatusNotFound, "error.html", map[string]any{
		"Status": http.StatusNotFound,
	})
}

// Error logs err and renders the 500 page.
func (h *Handler) Error(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("handler", "path", r.URL.Path, "err", err)
	h.View.Render(w, r, http.StatusInternalServerError, "error.html", map[string]any{
		"Status": http.StatusInternalServerError,
	})
}
