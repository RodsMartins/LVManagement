package handlers

import (
	"context"
	"fmt"

	"lvm/database"
	"lvm/internal/dtos"
	"lvm/internal/handlers"
	pages "lvm/internal/templates/pages/admin/seed"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AdminSeedHandLer struct {
	handlers.BaseHandler
	Repository *database.Queries
}

func (h AdminSeedHandLer) ViewSeeds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seeds, err := h.Repository.ListSeeds(ctx)
	if err != nil {
		http.Error(w, "Failed to retrieve seeds", http.StatusInternalServerError)
		return
	}

	seedInstructions, err := h.Repository.ListSeedInstructions(ctx)
	if err != nil {
		http.Error(w, "Failed to retrieve seed instructions", http.StatusInternalServerError)
		return
	}

	var seedDTOs []dtos.Seed
	var seedInstructionDTOs []dtos.SeedInstruction

	for _, seed := range seeds {
		seedDTOs = append(seedDTOs, dtos.Seed{}.FromDatabaseModel(seed))
	}

	for _, instruction := range seedInstructions {
		seedInstructionDTOs = append(seedInstructionDTOs, dtos.SeedInstruction{}.FromDatabaseModel(instruction))
	}

	var template templ.Component
	if h.BaseHandler.UsesHtmx(r) {
		template = pages.SeedsContainer(seedDTOs, seedInstructionDTOs)
	} else {
		template = pages.SeedsPage(seedDTOs, seedInstructionDTOs, r)
	}

	err = template.Render(ctx, w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h AdminSeedHandLer) SeedForm(w http.ResponseWriter, r *http.Request) {
	var template templ.Component

	seedId, err := h.BaseHandler.GetUrlUuidOrEmpty("seedId", r)
	if err != nil {
		http.Error(w, "Seed ID should be valid or empty", http.StatusBadRequest)
		return
	}

	seedModel := dtos.Seed{}
	seedInstructionModel := dtos.SeedInstruction{}

	if seedId != "" {
		ctx := context.Background()
		seed, err := h.Repository.GetSeed(ctx, pgtype.UUID{Bytes: uuid.MustParse(seedId), Valid: true})
		if err != nil {
			http.Error(w, "Seed not found", http.StatusNotFound)
			return
		}

		seedModel = seedModel.FromDatabaseModel(seed)

		seedInstruction, err := h.Repository.GetSeedInstructionsBySeedId(ctx, pgtype.UUID{Bytes: uuid.MustParse(seedId), Valid: true})
		if err != nil {
			http.Error(w, "Seed instruction not found", http.StatusNotFound)
			return
		}

		seedInstructionModel = seedInstructionModel.FromDatabaseModel(seedInstruction)
	}

	if h.BaseHandler.UsesHtmx(r) {
		template = pages.UpsertSeed(seedModel, seedInstructionModel)
	} else {
		template = pages.UpsertSeedPage(seedModel, seedInstructionModel, r)
	}

	err = template.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h AdminSeedHandLer) NewSeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	seed, seedInstruction, err := getAndValidateUpsertFields(nil, nil, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	seedModel := seed.ToDatabaseModel()

	ctx := context.Background()
	h.Repository.NewSeed(ctx, database.NewSeedParams(seedModel))

	seedInstructionModel := seedInstruction.ToDatabaseModel()

	h.Repository.NewSeedInstruction(ctx, database.NewSeedInstructionParams(seedInstructionModel))

	h.ViewSeeds(w, r)
}

func (h AdminSeedHandLer) UpdateSeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Extract data from the form
	seedId, err := h.BaseHandler.GetUrlUuid("seedId", r)
	if err != nil {
		http.Error(w, "Invalid seed ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	seedInstruction, err := h.Repository.GetSeedInstructionsBySeedId(ctx, pgtype.UUID{Bytes: uuid.MustParse(seedId), Valid: true})
	if err != nil {
		http.Error(w, "Seed not found", http.StatusNotFound)
		return
	}

	seedUUID := uuid.MustParse(seedId)
	seedInstructionUUID, _ := uuid.FromBytes(seedInstruction.SeedID.Bytes[:])
	seed, updatedSeedInstruction, err := getAndValidateUpsertFields(&seedUUID, &seedInstructionUUID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	seedModel := seed.ToDatabaseModel()
	updatedSeedInstructionModel := updatedSeedInstruction.ToDatabaseModel()

	h.Repository.UpdateSeed(ctx, database.UpdateSeedParams{
		SeedID: seedModel.SeedID,
		Name:   seedModel.Name,
		Type:   seedModel.Type,
	})

	h.Repository.UpdateSeedInstruction(ctx, database.UpdateSeedInstructionParams{
		SeedID:           updatedSeedInstructionModel.SeedID,
		SeedGrams:        updatedSeedInstructionModel.SeedGrams,
		SoakingHours:     updatedSeedInstructionModel.SoakingHours,
		StackingHours:    updatedSeedInstructionModel.StackingHours,
		BlackoutHours:    updatedSeedInstructionModel.BlackoutHours,
		LightsHours:      updatedSeedInstructionModel.LightsHours,
		YieldGrams:       updatedSeedInstructionModel.YieldGrams,
		SpecialTreatment: updatedSeedInstructionModel.SpecialTreatment,
	})

	h.ViewSeeds(w, r)
}

func (h AdminSeedHandLer) DeleteSeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Extract data from the form$
	seedId, err := h.BaseHandler.GetUrlUuid("seedId", r)
	if err != nil {
		http.Error(w, "Invalid seed ID", http.StatusBadRequest)
		return
	}

	// Validate that the required fields are not empty
	if seedId == "" {
		http.Error(w, "Missing form fields", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	nbrRows, err := h.Repository.DeleteSeed(ctx, pgtype.UUID{Bytes: uuid.MustParse(seedId), Valid: true})
	if nbrRows == 0 {
		http.Error(w, "Seed not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Error deleting seed", http.StatusInternalServerError)
		return
	}

	h.ViewSeeds(w, r)
}

func getAndValidateUpsertFields(SeedUUID *uuid.UUID, SeedInstructionUUID *uuid.UUID, r *http.Request) (*dtos.Seed, *dtos.SeedInstruction, error) {
	seedName := r.FormValue("name")
	seedType := r.FormValue("type")
	seedGrams := r.FormValue("seedGrams")
	soakingHours := r.FormValue("soakingHours")
	stackingHours := r.FormValue("stackingHours")
	blackoutHours := r.FormValue("blackoutHours")
	lightsHours := r.FormValue("lightsHours")
	yieldGrams := r.FormValue("yieldGrams")
	specialTreatment := r.FormValue("specialTreatment")

	// Validate that the required fields are not empty
	if seedName == "" || seedType == "" || seedGrams == "" {
		return nil, nil, fmt.Errorf("missing form fields")
	}

	// Convert the numeric fields to the appropriate types
	seedGramsInt, err := strconv.Atoi(seedGrams)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid seed grams")
	}

	soakingHoursInt, err := strconv.Atoi(soakingHours)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid soaking hours")
	}

	stackingHoursInt, err := strconv.Atoi(stackingHours)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid stacking hours")
	}

	blackoutHoursInt, err := strconv.Atoi(blackoutHours)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid blackout hours")
	}

	lightsHoursInt, err := strconv.Atoi(lightsHours)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid lights hours")
	}

	yieldGramsInt, err := strconv.Atoi(yieldGrams)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid yield grams")
	}

	seed := &dtos.Seed{}
	if SeedUUID != nil {
		seed.SeedID = *SeedUUID
	} else {
		seed.SeedID = uuid.New()
	}
	seed.Name = seedName
	seed.Type = seedType

	seedInstruction := &dtos.SeedInstruction{}
	if SeedInstructionUUID != nil {
		seedInstruction.SeedInstructionID = *SeedInstructionUUID
	} else {
		seedInstruction.SeedInstructionID = uuid.New()
	}
	seedInstruction.SeedID = seed.SeedID
	seedInstruction.SeedGrams = seedGramsInt
	seedInstruction.SoakingHours = soakingHoursInt
	seedInstruction.StackingHours = stackingHoursInt
	seedInstruction.BlackoutHours = blackoutHoursInt
	seedInstruction.LightsHours = lightsHoursInt
	seedInstruction.YieldGrams = yieldGramsInt
	seedInstruction.SpecialTreatment = specialTreatment

	return seed, seedInstruction, nil
}
