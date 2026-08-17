package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mo-taki/mem-price-dashboard/internal/config"
	"github.com/mo-taki/mem-price-dashboard/internal/jquants"
	"github.com/mo-taki/mem-price-dashboard/internal/store"
)

const (
	stockCode = "285A0" // Kioxia
)

func main() {
	godotenv.Load()
	apiKey := config.Getenv("JQUANTS_API_KEY", "")
	dbPath := config.Getenv("DB_PATH", "./stock_prices.db")

	client := jquants.New(apiKey)
	prices, err := client.FetchStockPrices(stockCode)
	if err != nil {
		log.Println("Error fetching stock prices:", err)
		return
	}

	st, err := store.Open(dbPath)
	if err != nil {
		log.Println("Error opening database:", err)
		return
	}
	defer st.Close()

	var saved, failed int

	for _, data := range prices {
		err := st.UpsertPrice(data)
		if err != nil {
			log.Println("Error upserting stock price:", err)
			failed++
			continue
		}
		saved++
	}
	log.Printf("stock fetch done: %d ok, %d failed (code=%s)", saved, failed, stockCode)
}
