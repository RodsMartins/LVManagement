package handlers

import (
	//"lvm/internal/middleware"
	//"lvm/internal/store"

	pages "lvm/internal/templates/pages/myday"
	"net/http"
)

type MyDayHandler struct {
	BaseHandler
}

func (h MyDayHandler) Index(w http.ResponseWriter, r *http.Request) {
	template := pages.MyDay(r)
	
	err := template.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
