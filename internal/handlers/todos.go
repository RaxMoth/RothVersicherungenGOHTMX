package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/maxroth/eumel/internal/db"
)

// The todo endpoints are the HTMX demo: every mutation re-renders the
// "todo-list" partial, which the client swaps in via hx-target.

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.FormValue("title"))
	if title != "" {
		if err := db.CreateTodo(h.DB, title); err != nil {
			h.Error(w, r, err)
			return
		}
	}
	h.renderTodoList(w, r)
}

func (h *Handler) ToggleTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := db.ToggleTodo(h.DB, id); err != nil {
		h.Error(w, r, err)
		return
	}
	h.renderTodoList(w, r)
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := db.DeleteTodo(h.DB, id); err != nil {
		h.Error(w, r, err)
		return
	}
	h.renderTodoList(w, r)
}

func (h *Handler) renderTodoList(w http.ResponseWriter, r *http.Request) {
	todos, err := db.ListTodos(h.DB)
	if err != nil {
		h.Error(w, r, err)
		return
	}
	h.View.RenderPartial(w, r, "todo-list", map[string]any{
		"Todos": todos,
	})
}
