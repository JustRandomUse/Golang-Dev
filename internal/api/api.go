package api

import (
	"github.com/gorilla/mux"
)

// API
func InitRoutes(r *mux.Router) {
	r.HandleFunc("/products", GET).Methods("GET")
	r.HandleFunc("/products/{id}", GET_BY_ID).Methods("GET")
	r.HandleFunc("/products", POST).Methods("POST")
	r.HandleFunc("/products/{id}", PUT).Methods("PUT")
	r.HandleFunc("/products/{id}", DEL).Methods("DELETE")
}
