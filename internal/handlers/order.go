package handlers

import (
	"fmt"
	"lvm/internal/db/repositories"
	"lvm/internal/dtos"
	formDtos "lvm/internal/dtos/form"
	"lvm/internal/services"
	"lvm/internal/templates/components/form"
	pages "lvm/internal/templates/pages/farm"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type FarmHandler struct {
	BaseHandler
	cropRepository repositories.CropRepository
	seedRepository repositories.SeedRepository
	orderService   services.OrderService
}

func NewFarmHandler(
	cropRepository repositories.CropRepository,
	seedRepository repositories.SeedRepository,
	orderService services.OrderService,
) *FarmHandler {
	return &FarmHandler{
		cropRepository: cropRepository,
		seedRepository: seedRepository,
		orderService:   orderService,
	}
}

func (h FarmHandler) ViewCrops(w http.ResponseWriter, r *http.Request) {
	template := pages.Crops(r)

	err := template.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h FarmHandler) CropForm(w http.ResponseWriter, r *http.Request) {
	var template templ.Component

	seeds, err := h.seedRepository.ListSeeds(r.Context())
	if err != nil {
		http.Error(w, "Error fetching seeds", http.StatusInternalServerError)
		return
	}

	var seedList []dtos.Seed

	options := make([]form.SelectOption, len(seeds))

	for _, seed := range seeds {
		seedModel := dtos.Seed{}.FromDatabaseModel(seed)
		seedList = append(seedList, seedModel)
	}

	for i, seed := range seedList {
		options[i] = form.SelectOption{Value: seed.SeedID.String(), Label: seed.Name}
	}

	template = pages.CropForm(r, options)

	err = template.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h FarmHandler) NewCrop(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	order, err := getAndValidateUpsertFields(nil, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate crops from order using the service
	crops, err := h.orderService.GenerateCropsFromOrder(*order)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate crops: %v", err), http.StatusInternalServerError)
		return
	}

	// Create context for database operations
	ctx := r.Context()

	// Insert each generated crop into the database
	for _, crop := range crops {
		_, err = h.cropRepository.NewCrop(ctx, crop.GetNewCropParams())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create crop: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Show the updated list of crops
	h.ViewCrops(w, r)
}

func getAndValidateUpsertFields(crop *formDtos.Order, r *http.Request) (*formDtos.Order, error) {
	cropID := r.FormValue("crop")

	seedID := r.FormValue("seed")
	if seedID == "" {
		return nil, fmt.Errorf("seed ID is required")
	}

	dateStr := r.FormValue("datetime-toggle-picker-date")
	if dateStr == "" {
		return nil, fmt.Errorf("soaking/harvest date is required")
	}

	parsedDate, err := time.Parse("01/02/2006", dateStr) // Adjust the layout to match your date format
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	timeStr := r.FormValue("datetime-toggle-picker-time")
	if timeStr == "" {
		return nil, fmt.Errorf("soaking/harvest time is required")
	}

	parsedTime, err := time.Parse("15:04", timeStr) // Adjust the layout to match your time format
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %v", err)
	}

	quantityStr := r.FormValue("quantity")
	if quantityStr == "" {
		return nil, fmt.Errorf("quantity is required")
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return nil, fmt.Errorf("invalid quantity: %v", err)
	}

	isHarvest := r.FormValue("datetime-toggle-picker-toggle") == "1"

	if crop == nil {
		crop = &formDtos.Order{}
	}

	if cropID != "" {
		cropUUID, err := uuid.Parse(cropID)
		if err != nil {
			return nil, fmt.Errorf("invalid crop ID: %v", err)
		}

		crop.CropID = cropUUID
	}

	SeedUUID, err := uuid.Parse(seedID)
	if err != nil {
		return nil, fmt.Errorf("invalid crop ID: %v", err)
	}
	crop.SeedID = SeedUUID
	crop.Date = parsedDate
	crop.Time = parsedTime
	crop.IsHarvest = isHarvest
	crop.Yield = quantity

	return crop, nil
}
