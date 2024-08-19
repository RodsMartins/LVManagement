package routes

import (
	"github.com/go-chi/chi/v5"
	"lvm/internal/handlers"
)

func MyDayRoutes() chi.Router {
    r := chi.NewRouter()

	myDayHandler := handlers.MyDayHandler{}

    r.Get("/", myDayHandler.Index);

    return r
}