package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mo-taki/mem-price-dashboard/internal/store"
)

const dbPath = "./stock_prices.db"

func main() {
	st, err := store.Open(dbPath)
	if err != nil {
		log.Println("Error opening database:", err)
		return
	}
	defer st.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/prices", pricesHandler(st))
	mux.HandleFunc("POST /api/products", productsPostHandler(st))
	mux.HandleFunc("GET /api/products", productsGetHandler(st))
	mux.HandleFunc("POST /api/market", marketPostHandler(st))
	mux.HandleFunc("GET /api/market", marketGetHandler(st))

	log.Fatal(http.ListenAndServe(":8080", cors(mux)))
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

func productsPostHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p store.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		err := s.UpsertProduct(p)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func productsGetHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := s.ListProducts()
		if err != nil {
			log.Printf("list products: %v", err)
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(products)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
	}
}

func marketPostHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m store.MarketPrice
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		err := s.UpsertMarketPrice(m)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func marketGetHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		product := r.URL.Query().Get("product")
		if product == "" {
			http.Error(w, "product is empty", http.StatusBadRequest)
			return
		}

		prices, err := s.ListMarketPrices(product)
		if err != nil {
			log.Printf("list market prices (product=%s): %v", product, err)
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

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // MEMO: 一旦全許可
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
