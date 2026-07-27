package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-sql-driver/mysql"
	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	"github.com/rawbil/ecom2/internal/config"
	db "github.com/rawbil/ecom2/internal/database"
	"github.com/rawbil/ecom2/internal/seed"
	"github.com/rawbil/ecom2/internal/server"
	"github.com/rawbil/ecom2/internal/utils"
)

func main() {

	err := config.LoadEnv()
	if err != nil {
		slog.Warn("No .env file found")
	}

	cfg := mysql.Config{
		Addr:                 config.InitConfig().DBAddress,
		User:                 config.InitConfig().DBUser,
		Passwd:               config.InitConfig().DBPassword,
		DBName:               config.InitConfig().DBName,
		ParseTime:            config.InitConfig().ParseTime,
		AllowNativePasswords: true,
	}

	dbConfig := server.DBConfig{
		DSN: cfg.FormatDSN(),
	}

	config := server.Config{
		Addr: config.GetServerAddr(),
		DB:   dbConfig,
	}

	//slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := db.DbConnection(dbConfig)
	if err != nil {
		slog.Error("Database Connection failed", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	app := server.Application{
		Config: config,
		DB:     db,
	}

	// Seed functions
	if err := seed.SeedPermissions(db); err != nil {
		utils.Log.Error(err.Error())
		return
	}
	if err := seed.SeedRoles(db); err != nil {
		utils.Log.Error(err.Error())
		return
	}

	if err := seed.SeedRolePermissions(context.Background(), *repository.New(app.DB), db); err != nil {
		utils.Log.Error(err.Error())
		return
	}

	m := app.Mount()
	if err := app.Run(m); err != nil {
		// log.Printf("Server failed: %s", err)
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
