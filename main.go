package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	UnitCost float64 `json:"unit_cost"`
}

var products []Product
var ID int

func main() {
	r := mux.NewRouter()

	// Определение маршрутов
	r.HandleFunc("/product", GET).Methods("GET")
	r.HandleFunc("/product/{id}", GET_ID).Methods("GET")
	r.HandleFunc("/product", POST).Methods("POST")
	r.HandleFunc("/product/{id}", PUT).Methods("PUT")
	r.HandleFunc("/product/{id}", DELETE).Methods("DELETE")

	fmt.Println("Запуск сервера на порту 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

// Функция GET возвращает список всех товаров
func GET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Функция GET_ID возвращает товар по ID
func GET_ID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := mux.Vars(r)
	id, err := strconv.Atoi(p["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	// _ - игнорирования значения
	for _, product := range products {
		if product.ID == id {
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// Функция POST создает новый товар
func POST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var postProduct Product
	_ = json.NewDecoder(r.Body).Decode(&postProduct)
	postProduct.ID = NextID()
	products = append(products, postProduct)
	json.NewEncoder(w).Encode(postProduct)
}

// Функция PUT обновляет товар по ID
func PUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := mux.Vars(r)
	id, err := strconv.Atoi(p["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var putProduct Product
	_ = json.NewDecoder(r.Body).Decode(&putProduct)
	for index, product := range products {
		if product.ID == id {
			putProduct.ID = id
			products[index] = putProduct
			json.NewEncoder(w).Encode(putProduct)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// Функция DELETE удаляет товар по ID
func DELETE(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	id, err := strconv.Atoi(p["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for index, product := range products {
		if product.ID == id {
			products = append(products[:index], products[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// Вспомогательная функция для получения следующего доступного ID товара
func NextID() int {
	ID++
	return ID
}
