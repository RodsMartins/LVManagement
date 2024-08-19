package routes

import (
	"lvm/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func FarmRoutes() chi.Router {
	r := chi.NewRouter()
	farmHandler := handlers.FarmHandLer{}

	r.Get("/", farmHandler.ViewCrops)
	r.Get("/new/crop", farmHandler.NewCrop)

	return r
}
