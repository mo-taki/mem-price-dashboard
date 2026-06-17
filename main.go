package main

import (
	"fmt"
	"os"

	_ "embed"

	_ "modernc.org/sqlite"

	"github.com/joho/godotenv"
	"github.com/mo-taki/mem-price-dashboard/internal/jquants"
	"github.com/mo-taki/mem-price-dashboard/internal/store"
)

func main(){
	godotenv.Load()
	apiKey := os.Getenv("JQUANTS_API_KEY")
	stockCode := "285A0"	// Kioxia
	dbPath := "./stock_prices.db"

	client := jquants.New(apiKey)
	prices, err := client.FetchStockPrices(stockCode)
    if err != nil {
        fmt.Println("Error fetching stock prices:", err)
        return
    }

	st, err := store.Open(dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer st.Close()

	for _, data := range prices {
		err := st.UpsertPrice(data)
		if err != nil {
			fmt.Println("Error upserting stock price:", err)
			continue
		}
	}
}