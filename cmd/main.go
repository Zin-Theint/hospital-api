package main

import (
	"log"
	"os"

	"github.com/Zin-Theint/hospital-api/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.Setup()

	log.Printf("Server running on :", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
