package handlers

import (
	//"lvm/internal/middleware"
	//"lvm/internal/store"

	pages "lvm/internal/templates/pages/myday"
	"net/http"

	"github.com/a-h/templ"
)

type MyDayHandler struct {
	BaseHandler
}

func (h MyDayHandler) Index(w http.ResponseWriter, r *http.Request) {
	var template templ.Component
	if h.BaseHandler.UsesHtmx(r) {
		template = pages.MyDay()
	} else {
		template = pages.MyDayPage(r)
	}
	
	err := template.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
