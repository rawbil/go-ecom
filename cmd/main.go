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
	app := Application{
		Config: config,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	m := app.mount()
	if err := app.run(m); err != nil {
		// log.Printf("Server failed: %s", err)
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
