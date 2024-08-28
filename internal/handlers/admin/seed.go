package handlers

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"lvm/database"
	"lvm/internal/handlers"
	pages "lvm/internal/templates/pages/admin/seed"
	"net/http"
	"os"
)

type AdminSeedHandLer struct {
	handlers.BaseHandler
	Repository *database.Queries
}

func (h AdminSeedHandLer) ViewSeeds(w http.ResponseWriter, r *http.Request) {
	var template templ.Component

	ctx := context.Background()
	seeds, err := h.Repository.ListSeeds(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to query: %v\n", err)
		os.Exit(1)
	}

	if h.BaseHandler.UsesHtmx(r) {
		template = pages.Seeds(seeds)
	} else {
		template = pages.SeedsPage(seeds, r)
	}

	err = template.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h AdminSeedHandLer) NewSeedForm(w http.ResponseWriter, r *http.Request) {
	var template templ.Component

	if h.BaseHandler.UsesHtmx(r) {
		template = pages.NewSeed()
	} else {
		template = pages.NewSeedPage(r)
	}

	err := template.Render(r.Context(), w)

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

	// Extract data from the form
	seedName := r.FormValue("name")
	seedType := r.FormValue("type")

	// Validate that the required fields are not empty
	if seedName == "" || seedType == "" {
		http.Error(w, "Missing form fields", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	h.Repository.NewSeed(ctx, database.NewSeedParams{
		SeedID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Name:   pgtype.Text{String: seedName, Valid: true},
		Type:   pgtype.Text{String: seedType, Valid: true},
	})

	h.ViewSeeds(w, r)
}
