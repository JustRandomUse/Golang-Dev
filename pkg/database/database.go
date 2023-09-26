package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	UnitCost float64 `json:"unit_cost"`
}

var db *pgxpool.Pool

func InitDB() error {
	connStr := "postgresql://postgres:21031975@localhost:5432/TestTask"
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Error parsing connection string: %v", err)
		return err
	}
	// Макс кол-во соединений в пуле
	poolConfig.MaxConns = 3

	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return err
	}

	log.Println("Connected to database")
	return nil
}

func GetDB() *pgxpool.Pool {
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// Получение списка всех товаров
func GetProducts() ([]Product, error) {
	// Пул соединений с базой данных
	pool, err := pgxpool.Connect(context.Background(), "postgresql://postgres:21031975@localhost:5432/TestTask")
	if err != nil {
		return nil, err
	}
	defer pool.Close()

	// Запрос к базе данных для получения списка товаров
	rows, err := pool.Query(context.Background(), "SELECT id, name, quantity, unit_cost FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Срез для хранения товаров
	var products []Product

	// Результаты запроса и добавления товаров в срез
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.UnitCost); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// Проверка на наличие ошибок после прохода по результатам запроса
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// Получение товара по ID
func GetProductByID(id int) (Product, error) {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://postgres:21031975@localhost:5432/TestTask")
	if err != nil {
		return Product{}, err
	}
	defer pool.Close()

	// Запрос к базе данных для получения товара по ID
	var product Product
	err = pool.QueryRow(context.Background(), "SELECT id, name, quantity, unit_cost FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Quantity, &product.UnitCost) // переделать, много элементов !!!
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

// Создание нового товара
func CreateProduct(newProduct Product) (Product, error) {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://postgres:21031975@localhost:5432/TestTask")
	if err != nil {
		return Product{}, err
	}
	defer pool.Close()

	// Запрос к базе данных для создания нового товара
	var createdProduct Product
	err = pool.QueryRow(context.Background(), "INSERT INTO products (name, quantity, unit_cost) VALUES ($1, $2, $3) RETURNING id", newProduct.Name, newProduct.Quantity, newProduct.UnitCost).Scan(&createdProduct.ID) // переделать, много элементов !!!
	if err != nil {
		return Product{}, err
	}

	// Установка значений остальных полей созданного товара
	createdProduct.Name = newProduct.Name
	createdProduct.Quantity = newProduct.Quantity
	createdProduct.UnitCost = newProduct.UnitCost

	return createdProduct, nil
}

// Обновление товара
func UpdateProduct(updatedProduct Product) error {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://postgres:21031975@localhost:5432/TestTask")
	if err != nil {
		return err
	}
	defer pool.Close()

	// Запрос к базе данных для обновления товара
	_, err = pool.Exec(context.Background(), "UPDATE products SET name = $1, quantity = $2, unit_cost = $3 WHERE id = $4", updatedProduct.Name, updatedProduct.Quantity, updatedProduct.UnitCost, updatedProduct.ID) // переделать, много элементов !!!
	if err != nil {
		return err
	}

	return nil
}

// Удаление товара по ID
func DeleteProduct(id int) error {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://postgres:21031975@localhost:5432/TestTask")
	if err != nil {
		return err
	}
	defer pool.Close()

	// Запрос к базе данных для удаления товара
	_, err = pool.Exec(context.Background(), "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
