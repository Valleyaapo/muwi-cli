package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func CreateDB(db *sql.DB) {
	// Create a table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ProductSubcategory TEXT,
		Brand TEXT,
		PaymentMethod TEXT,
		CustomerID TEXT,
		ProductPrice REAL,
		ShippingCost REAL,
		TotalAmount REAL,
		CustomerRating REAL,
		DeliveryStatus TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}
}
func Insert(o Order, db *sql.DB) {
	// Insert a row
	_, err := db.Exec(`INSERT INTO orders(
		ProductSubcategory, 
		Brand, 
		PaymentMethod,
		CustomerID, 
		ProductPrice,
		ShippingCost,
		TotalAmount,
		CustomerRating,
		DeliveryStatus
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.ProductSubcategory,
		o.Brand, o.PaymentMethod, o.CustomerID, o.ProductPrice, o.ShippingCost,
		o.TotalAmount, o.CustomerRating, o.DeliveryStatus)
	if err != nil {
		log.Fatal(err)
	}
}
func Scan(db *sql.DB) {
	// Query rows
	rows, err := db.Query("SELECT id, ProductSubcategory, TotalAmount FROM orders")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Print results
	for rows.Next() {
		var id int
		var productSubcategory string
		var totalAmount float64
		if err := rows.Scan(&id, &productSubcategory, &totalAmount); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Order: %d, %s -> %f\n", id, productSubcategory, totalAmount)
	}
}
