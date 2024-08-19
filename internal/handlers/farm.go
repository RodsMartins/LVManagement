package handlers

import (
	pages "lvm/internal/templates/pages/farm"
	"net/http"

	"github.com/a-h/templ"
)

type FarmHandLer struct {
	BaseHandler
}

func (h FarmHandLer) ViewCrops(w http.ResponseWriter, r *http.Request) {
	var template templ.Component
	if h.BaseHandler.UsesHtmx(r) {
		template = pages.Crops()
	} else {
		template = pages.CropsPage(r)
	}

	err := template.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h FarmHandLer) NewCrop(w http.ResponseWriter, r *http.Request) {
	var template templ.Component
	if h.BaseHandler.UsesHtmx(r) {
		template = pages.NewCrop()
	} else {
		template = pages.NewCropPage(r)
	}

	err := template.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
