package main

import (
	"log/slog"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rawbil/ecom2/internal/config"
	db "github.com/rawbil/ecom2/internal/database"
	"github.com/rawbil/ecom2/internal/server"
)

func main() {

	err := godotenv.Load()
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

	m := app.Mount()
	if err := app.Run(m); err != nil {
		// log.Printf("Server failed: %s", err)
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
