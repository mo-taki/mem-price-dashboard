package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mo-taki/mem-price-dashboard/internal/jquants"
	"github.com/mo-taki/mem-price-dashboard/internal/store"
)

const (
	stockCode = "285A0" // Kioxia
	dbPath    = "./stock_prices.db"
)

func main() {
	godotenv.Load()
	apiKey := os.Getenv("JQUANTS_API_KEY")

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

	for _, data := range prices {
		err := st.UpsertPrice(data)
		if err != nil {
			log.Println("Error upserting stock price:", err)
			continue
		}
		log.Println("Inserted/Updated:", data.Date, data.AdjC)
	}
}
