package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("PORT")

	cfg := config{
		addr: addr,
	}

	app := application{
		config: cfg,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
