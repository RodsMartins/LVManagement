package routes

import (
	handlers "lvm/internal/handlers/admin"
	"lvm/database"

	"github.com/go-chi/chi/v5"
)

func AdminRoutes(repository *database.Queries) chi.Router {
	r := chi.NewRouter()
	seedHandler := handlers.AdminSeedHandLer{
		Repository: repository,
	}

	r.Get("/seeds", seedHandler.ViewSeeds)
	r.Get("/seeds/new", seedHandler.NewSeedForm)
	r.Post("/seeds/new", seedHandler.NewSeed)

	return r
}
