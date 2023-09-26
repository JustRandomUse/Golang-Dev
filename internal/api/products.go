package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	postgresql "web-service/pkg/database"

	"github.com/gorilla/mux"
)

func GET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Запрос к базе данных
	products, err := postgresql.GetProducts()
	if err != nil {
		log.Printf("Error getting products: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

func GET_BY_ID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Извлечь ID товара из параметров запроса
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error parsing product ID: %v", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := postgresql.GetProductByID(id)
	if err != nil {
		log.Printf("Error getting product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func POST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Декодировать JSON-запрос и создать новый товар
	var newProduct postgresql.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Запрос к базе данных для создания нового товара
	createdProduct, err := postgresql.CreateProduct(newProduct)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(createdProduct)
}

func PUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Декодировать JSON-запрос и получить обновленный товар
	var updatedProduct postgresql.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Вызов функции UpdateProduct из пакета database
	if err := postgresql.UpdateProduct(updatedProduct); err != nil {
		log.Printf("Error updating product: %v", err)
		http.Error(w, "Error updating product", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedProduct)
}

func DEL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// ID товара
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error parsing product ID: %v", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Запрос к базе данных для удаления товара по ID
	err = postgresql.DeleteProduct(id)
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
