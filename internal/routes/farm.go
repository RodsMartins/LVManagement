package routes

import (
	"lvm/database"
	"lvm/internal/db/repositories"
	"lvm/internal/handlers"
	"lvm/internal/services"

	"github.com/go-chi/chi/v5"
)

func FarmRoutes(db *database.Queries) chi.Router {
	r := chi.NewRouter()

	cropRepository := repositories.NewCropRepository(db)
	seedRepository := repositories.NewSeedRepository(db)
	farmHandler := handlers.NewFarmHandler(cropRepository, seedRepository, *services.NewOrderService(db))

	r.Get("/", farmHandler.ViewCrops)
	r.Get("/crops/upsert", farmHandler.CropForm)
	r.Get("/crops/upsert/{seedId}", farmHandler.CropForm)
	r.Post("/crops/new", farmHandler.NewCrop)

	return r
}
