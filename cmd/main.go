package main

import (
	"log/slog"
	"os"
)

func main() {

	dbConfig := DBConfig{
		DSN: "root:rawbil@/go1?parseTime=true",
	}
	config := Config{
		Addr: ":8080",
		DB:   dbConfig,
	}

	//slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := dbConnection(dbConfig)
	if err != nil {
		slog.Error("Database Connection failed", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	app := Application{
		Config: config,
		db:     db,
	}

	m := app.mount()
	if err := app.run(m); err != nil {
		// log.Printf("Server failed: %s", err)
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
