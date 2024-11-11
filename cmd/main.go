package main

import (
	"context"
	"errors"
	"fmt"
	//"log"
	"log/slog"
	"lvm/database"
	"lvm/internal/config"
	"lvm/internal/handlers"
	"lvm/internal/routes"

	//database "lvm/internal/store/db"
	//"lvm/internal/store/dbstore"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	//"github.com/jackc/pgx/v5/pgtype"
)

/*
* Set to production at build time
* used to determine what assets to load
 */
var Environment = "development"

func init() {
	os.Setenv("env", Environment)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r := chi.NewRouter()

	cfg := config.MustLoadConfig()

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	repository := database.New(conn)

	//insertedAuthor, err := queries.CreateTmp(ctx, database.CreateTmpParams{
	//	FertilizerID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
	//	Name:         pgtype.Text{String: "Hello Helene", Valid: true},
	//})
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//log.Println(insertedAuthor)
	//fmt.Println(insertedAuthor)

	//var name string
	//var weight int64
	//err = conn.QueryRow(context.Background(), "select * from widgets where id=$1", 42).Scan(&name, &weight)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	//	os.Exit(1)
	//}

	//fmt.Println(name, weight)

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
		)

		r.Handle("/", http.RedirectHandler("/my-day", http.StatusPermanentRedirect))

		notFoundHandler := handlers.NotFoundHandler{}
		r.NotFound(notFoundHandler.NotFound)

		r.Mount("/my-day", routes.MyDayRoutes())
		r.Mount("/farm", routes.FarmRoutes(repository))
		r.Mount("/admin", routes.AdminRoutes(repository))
	})

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Info("Server started", slog.String("port", cfg.Port), slog.String("env", Environment))
	<-killSig

	logger.Info("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
