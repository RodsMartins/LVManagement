package handlers

import (
	"lvm/internal/db/repositories"
	pages "lvm/internal/templates/pages/farm"
	"net/http"
)

type CropHandler struct {
	BaseHandler
	cropRepository repositories.CropRepository
}

func NewCropHandler(cropRepository repositories.CropRepository) *CropHandler {
	return &CropHandler{
		cropRepository: cropRepository,
	}
}

func (c CropHandler) ViewCrops(w http.ResponseWriter, r *http.Request) {
	crops, err := c.cropRepository.ListCrop(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve crops", http.StatusInternalServerError)
		return
	}

	template := pages.Crops(r, crops)

	err = template.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
