package main

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
)

func readCSV(filepath string, ctx context.Context) <-chan Order {
	orders := make(chan Order, 16)
	var order Order
	go func() {
		defer close(orders)
		f, err := os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		r := csv.NewReader(f)
		if err != nil {
			log.Fatal(err)
		}
		_, err = r.Read()
		if err != nil {
			log.Fatal(err)
		}
		parseToFloat := func(s string) float64 {
			val, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0.0
			}
			return val
		}

		parseToInt := func(s string) int {
			val, err := strconv.Atoi(s)
			if err != nil {
				return 0
			}
			return val
		}
		for {
			record, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatal(err)
			}
			order = Order{
				ProductCategory:          record[0],
				ProductSubcategory:       record[1],
				Brand:                    record[2],
				DeliveryStatus:           record[3],
				AssemblyServiceRequested: record[4],
				PaymentMethod:            record[5],
				OrderID:                  record[6],
				CustomerID:               record[7],
				ProductPrice:             parseToFloat(record[8]),
				ShippingCost:             parseToFloat(record[9]),
				AssemblyCost:             parseToFloat(record[10]),
				TotalAmount:              parseToFloat(record[11]),
				DeliveryWindowDays:       parseToInt(record[12]),
				CustomerRating:           parseToFloat(record[13]),
			}
			select {
			case <-ctx.Done():
				return
			case orders <- order:
			}
		}
	}()
	return orders
}

type Order struct {
	ProductCategory          string
	ProductSubcategory       string
	Brand                    string
	DeliveryStatus           string
	AssemblyServiceRequested string
	PaymentMethod            string
	OrderID                  string
	CustomerID               string
	ProductPrice             float64
	ShippingCost             float64
	AssemblyCost             float64
	TotalAmount              float64
	DeliveryWindowDays       int
	CustomerRating           float64
}
