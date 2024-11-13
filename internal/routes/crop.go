package routes

import (
	"lvm/database"
	"lvm/internal/db/repositories"
	"lvm/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func CropRoutes(db *database.Queries) chi.Router {
	r := chi.NewRouter()

	cropHandler := handlers.NewCropHandler(repositories.NewCropRepository(db))

	r.Get("/", cropHandler.ViewCrops)

	return r
}
