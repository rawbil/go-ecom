package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	Config Config
	db     *sql.DB
}

type Config struct {
	Addr string
	DB   DBConfig
}

type DBConfig struct {
	DSN string
}

// mount
func (app *Application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// test database health
		err := app.db.Ping()
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
			return
		}
		fmt.Fprintln(w, "Server and Database OK")
	})

	return r
}

// run
func (app *Application) run(m http.Handler) error {
	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      m,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	if app.Config.Addr != "" {
		fmt.Printf("Server is running on http://localhost%s", app.Config.Addr)
	}

	return srv.ListenAndServe()
}
