package routes

import (
	"lvm/database"
	"lvm/internal/db/repositories"
	"lvm/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func MyDayRoutes(db *database.Queries) chi.Router {
    r := chi.NewRouter()

	cropRepo := repositories.NewCropRepository(db)


	myDayHandler := handlers.NewMyDayHandler(
		cropRepo,
	)

    r.Get("/", myDayHandler.Index);

    return r
}