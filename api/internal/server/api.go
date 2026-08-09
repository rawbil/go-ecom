package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	"github.com/rawbil/ecom2/internal/auth"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/authorization"
	"github.com/rawbil/ecom2/internal/middlewarefns"
	"github.com/rawbil/ecom2/internal/orders"
	"github.com/rawbil/ecom2/internal/products"
	"github.com/rawbil/ecom2/internal/users"
	"github.com/rawbil/ecom2/internal/utils"
	"golang.org/x/time/rate"
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
	utils.Slogger()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	// r.Use(middleware.Logger)
	r.Use(middlewarefns.Logger())
	r.Use(middleware.Recoverer)

	repo := repository.New(app.DB)
	authService := auth.NewService(repo, app.DB)
	authHandler := auth.NewHandler(authService)

	productsService := products.NewService(*repo, app.DB)
	productsHandler := products.NewHandler(productsService)

	usersService := users.NewService(*repo)
	usersHandler := users.NewHandler(usersService)

	orderService := orders.NewService(*repo, app.DB)
	orderHandler := orders.NewHandler(orderService)

	rl := middlewarefns.NewRateLimiter()
	// run rate limiter cleanup every 10 minutes
	go rl.StartCleanup(10*time.Minute, 30*time.Minute)

	//* Root route
	r.Route("/api/v1", func(r chi.Router) {
		// ! /api/v1/auth
		r.Route("/auth", func(r chi.Router) {
			//?POST /auth/register
			r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).Post("/register", authHandler.UserRegister)
			//? POST /auth/login
			r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).Post("/login", authHandler.UserLogin)
			//? POST /auth/logout
			r.With(authutils.AuthMiddleware(*repo)).With(rl.Limit(rate.Limit(5.0/60.0), 5)).Post("/logout", authHandler.UserLogout)
			//?POST /auth/password-reset
			r.With(authutils.AuthMiddleware(*repo)).With(rl.Limit(rate.Limit(5.0/60.0), 5)).Post("/password-reset", authHandler.PasswordReset)
			//? POST /auth/refresh-tokens
			r.With(authutils.AuthMiddleware(*repo)).With(rl.Limit(rate.Limit(5.0/60.0), 5)).Post("/refresh-tokens", authHandler.RefreshTokens)
			//? POST /auth/update-my-email
			r.With(authutils.AuthMiddleware(*repo)).With(rl.Limit(rate.Limit(5.0/60.0), 5)).Post("/update-my-email", authHandler.UpdateUserEmail)
			//? POST /auth/update-my-username
			r.With(authutils.AuthMiddleware(*repo)).With(rl.Limit(rate.Limit(5.0/60.0), 5)).Post("/update-my-username", authHandler.UpdateUsername)
			r.With(rl.Limit(rate.Limit(100.0/60.0), 5)).Post("/forgot-password", authHandler.ForgotPassword)
		})

		//& Group protected routes to apply auth middleware
		r.Group(func(r chi.Router) {
			// Consume auth middleware for all protected routes
			r.Use(authutils.AuthMiddleware(*repo))
			// ! /api/v1/users
			r.Route("/users", func(r chi.Router) {
				//? GET /users/find-one
				r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).Get("/one", usersHandler.ListUser)
				//? GET /users/find-all
				r.With(authorization.PermissionMiddleware(*repo, authorization.PermissionViewUsers)).With(rl.Limit(rate.Limit(5.0/60.0), 5)).Get("/find-all", usersHandler.ListAllUsers)
				//? POST /users/create
				r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).With(authorization.PermissionMiddleware(*repo, authorization.PermissionCreateUser)).Post("/create", usersHandler.CreateUser)
				//?DELETE /users/delete
				r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).With(authorization.PermissionMiddleware(*repo, authorization.PermissionDeleteUser)).Delete("/delete", usersHandler.DeleteUser)
			})

			// ! /api/v1/products
			r.Route("/products", func(r chi.Router) {
				//? GET /products
				r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).Get("/", productsHandler.ListProducts)
				//? GET /products/id
				r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).Get("/id", productsHandler.ListProduct)
				//? POST /products
				r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).With(authorization.PermissionMiddleware(*repo, authorization.PermissionCreateProduct)).Post("/", productsHandler.CreateProduct)
				//? DELETE /products
				r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).With(authorization.PermissionMiddleware(*repo, authorization.PermissionDeleteProduct)).Delete("/delete", productsHandler.DeleteProduct)
				//? PATCH /products
				r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).With(authorization.PermissionMiddleware(*repo, authorization.PermissionUpdateProduct)).Patch("/", productsHandler.UpdateProduct)
			})

			//! /api/v1/orders
			r.Route("/orders", func(r chi.Router) {
				//? POST /orders
				r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).With(authorization.PermissionMiddleware(*repo, authorization.PermissionCreateOrder)).Post("/", orderHandler.CreateOrder)
				//? GET /orders/my-orders
				r.With(authorization.PermissionMiddleware(*repo, authorization.PermissionViewOrder)).With(rl.Limit(rate.Limit(5.0/60.0), 5)).Get("/my-orders", orderHandler.GetMyOrder)
				//? GET /orders/all
				r.With(authorization.PermissionMiddleware(*repo, authorization.PermissionViewAllOrders)).With(rl.Limit(rate.Limit(5.0/60.0), 5)).Get("/all", orderHandler.GetAllOrders)
				//? POST /orders/id
				r.With(rl.Limit(rate.Limit(5.0/60.0), 5)).With(authorization.PermissionMiddleware(*repo, authorization.PermissionCancelOrder)).Post("/cancel", orderHandler.CancleOrder)
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

	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      m,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	if app.Config.Addr != "" {
		message := fmt.Sprintf("Server is running on http://localhost%s", app.Config.Addr)
		utils.Log.Info(message, "Status", http.StatusOK)
	}

	return srv.ListenAndServe()
}
