package server

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	"github.com/rawbil/ecom2/internal/auth"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/orders"
	"github.com/rawbil/ecom2/internal/products"
	"github.com/rawbil/ecom2/internal/users"
)

type Application struct {
	Config Config
	DB     *sql.DB
}

type Config struct {
	Addr string
	DB   DBConfig
}

type DBConfig struct {
	DSN string
}

// mount
func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	repo := repository.New(app.DB)
	authService := auth.NewService(*repo, app.DB)
	authHandler := auth.NewHandler(authService)

	productsService := products.NewService(*repo)
	productsHandler := products.NewHandler(productsService)

	usersService := users.NewService(*repo)
	usersHandler := users.NewHandler(usersService)

	orderService := orders.NewService(*repo, app.DB)
	orderHandler := orders.NewHandler(orderService)

	//* Root route
	r.Route("/api/v1", func(r chi.Router) {
		// ! /api/v1/auth
		r.Route("/auth", func(r chi.Router) {
			//?POST /auth/register
			r.Post("/register", authHandler.UserRegister)
			//? POST /auth/login
			r.Post("/login", authHandler.UserLogin)
			//? POST /auth/logout
			r.With(authutils.AuthMiddleware(*repo)).Post("/logout", authHandler.UserLogout)
			//?POST /auth/password-reset
			r.With(authutils.AuthMiddleware(*repo)).Post("/password-reset", authHandler.PasswordReset)
		})

		//& Group protected routes to apply auth middleware
		r.Group(func(r chi.Router) {
			// Consume auth middleware for all protected routes
			r.Use(authutils.AuthMiddleware(*repo))
			// ! /api/v1/users
			r.Route("/users", func(r chi.Router) {
				//? GET /users/find-one
				r.Get("/id", usersHandler.ListUser)
				//? GET /users/find-all
				r.Get("/", usersHandler.ListAllUsers)
				//? POST /users/create
				r.Post("/create", usersHandler.CreateUser)
				//?DELETE /users/delete
				r.Delete("/delete", usersHandler.DeleteUser)
			})

			// ! /api/v1/products
			r.Route("/products", func(r chi.Router) {
				//? GET /products/list
				r.Get("/", productsHandler.ListProducts)
				//? GET /products/one
				r.Get("/id", productsHandler.ListProduct)
				//? POST /products
				r.Post("/create", productsHandler.CreateProduct)
				//? DELETE /products
				r.Delete("/delete", productsHandler.DeleteProduct)
			})

			r.Route("/orders", func(r chi.Router) {
				r.Post("/", orderHandler.CreateOrder)
			})
		})

		// health route GET /api/v1/health
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			// test database health
			err := app.DB.Ping()
			if err != nil {
				http.Error(w, "Database connection error", http.StatusInternalServerError)
				fmt.Fprintf(w, "%s", err)
				return
			}
			fmt.Fprintln(w, "Server and Database OK")
		})
	})

	return r
}

// run
func (app *Application) Run(m http.Handler) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      m,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	if app.Config.Addr != "" {
		message := fmt.Sprintf("Server is running on http://localhost%s", app.Config.Addr)
		slog.Info(message, "Status", http.StatusOK)
	}

	return srv.ListenAndServe()
}
