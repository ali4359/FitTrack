package main

import (
	"log"
	"os"

	"github.com/ali4359/iron-and-spice/backend/internal/api"
	"github.com/ali4359/iron-and-spice/backend/internal/store"
)

func main() {
	db := store.Open()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := api.New(db)
	log.Printf("iron-and-spice backend listening on :%s", port)
	if err := srv.Router().Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
