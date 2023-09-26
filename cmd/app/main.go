package main

import (
	"log"
	"net/http"

	"web-service/internal/api"
	postgresql "web-service/pkg/database"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	err := postgresql.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
		return
	}
	defer postgresql.CloseDB()

	r := mux.NewRouter()

	api.InitRoutes(r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	h := c.Handler(r)

	log.Println("Server is running on port http://localhost:8088")
	log.Fatal(http.ListenAndServe(":8088", h))
}
