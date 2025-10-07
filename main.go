package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/martialanouman/personal-library/internal/app"
	"github.com/martialanouman/personal-library/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("error loading .env file: %v", err))
	}

	
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.Db.Close()

	r := routes.SetupRoutes(app)

	server := http.Server{
		Addr:         ":3000",
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("server running on :3000")

	if err := server.ListenAndServe(); err != nil {
		app.Logger.Fatal(err)
	}
}
