package handlers

import (
	templates "lvm/internal/templates"
	"net/http"
)

type NotFoundHandler struct{}

func (h *NotFoundHandler) NotFound(w http.ResponseWriter, r *http.Request) {

	err := templates.NotFound(r).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
