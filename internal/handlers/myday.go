package handlers

import (
	//"lvm/internal/middleware"
	//"lvm/internal/store"

	"fmt"
	"lvm/internal/db/repositories"
	"lvm/internal/dtos"
	pages "lvm/internal/templates/pages/myday"
	"net/http"
	"time"
)

type MyDayHandler struct {
	BaseHandler
	cropRepository repositories.CropRepository
}

func NewMyDayHandler(
	cropRepository repositories.CropRepository,
) *MyDayHandler {
	return &MyDayHandler{
		cropRepository: cropRepository,
	}
}
func (h MyDayHandler) Index(w http.ResponseWriter, r *http.Request) {
	crops, err := h.cropRepository.ListCropsByDate(r.Context(), time.Now())

	if err != nil {
		http.Error(w, "Failed to retrieve crops", http.StatusInternalServerError)
		return
	}

    cropsActiveStages := make(map[dtos.CropStage][]dtos.Crop)
    for _, crop := range crops {
        stages := crop.GetActiveStages(time.Now())
		fmt.Println(stages)
        for _, stage := range stages {
            cropsActiveStages[stage] = append(cropsActiveStages[stage], crop)
        }
    }

    cropsNewStages := make(map[dtos.CropStage][]dtos.Crop)
    for _, crop := range crops {
        stages := crop.GetNewStages(time.Now())
		fmt.Println(stages)
        for _, stage := range stages {
            cropsNewStages[stage] = append(cropsNewStages[stage], crop)
        }
    }

	template := pages.MyDay(r, cropsActiveStages, cropsNewStages)
	err = template.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
