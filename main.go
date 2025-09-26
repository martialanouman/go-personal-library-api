package main

import (
	"net/http"
	"time"

	"github.com/martialanouman/personal-library/internal/app"
)

func main() {
	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr:         ":3000",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("server running on :3000")

	if err := server.ListenAndServe(); err != nil {
		app.Logger.Fatal(err)
	}
}
