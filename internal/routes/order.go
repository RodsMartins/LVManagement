package routes

import (
	"lvm/database"
	"lvm/internal/db/repositories"
	"lvm/internal/handlers"
	"lvm/internal/services/order"

	"github.com/go-chi/chi/v5"
)

func OrderRoutes(db *database.Queries) chi.Router {
	r := chi.NewRouter()

	seedRepo := repositories.NewSeedRepository(db)
	seedInstructionRepo := repositories.NewSeedInstructionRepository(db)
	cropRepo := repositories.NewCropRepository(db)
	orderHandler := handlers.NewOrderHandler(
		cropRepo,
		seedRepo,
		*order.NewOrderService(seedRepo, seedInstructionRepo, cropRepo),
		*handlers.NewCropHandler(cropRepo),
	)

	r.Get("/new", orderHandler.OrderForm)
	r.Post("/new", orderHandler.NewOrder)

	return r
}
