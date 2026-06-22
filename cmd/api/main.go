package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mo-taki/mem-price-dashboard/internal/store"
)

const dbPath = "./stock_prices.db"

func main() {
	st, err := store.Open(dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer st.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/prices", pricesHandler(st))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func pricesHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "code is empty", http.StatusBadRequest)
			return
		}

		prices, err := s.ListPrices(code)
		if err != nil {
			log.Printf("list prices (code=%s): %v", code, err)
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(prices)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
	}
}
