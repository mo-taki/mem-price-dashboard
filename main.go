package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	_ "embed"

	_ "modernc.org/sqlite"

	"github.com/joho/godotenv"
)

//go:embed schema.sql
var schemaSQL string

const upsertStockPriceSQL = `INSERT INTO stock_prices (code, date, open, high, low, close, volume, adj_close, adj_factor)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(code, date)
		DO UPDATE SET open=excluded.open, high=excluded.high, low=excluded.low, close=excluded.close, volume=excluded.volume, adj_close=excluded.adj_close, adj_factor=excluded.adj_factor`

type StockPriceResponse struct {
	Data []StockPrice `json:"data"`
}

type StockPrice struct {
	Date      string  `json:"Date"`
	Code      string  `json:"Code"`
	O         float64 `json:"O"`
	H         float64 `json:"H"`
	L         float64 `json:"L"`
	C         float64 `json:"C"`
	Ul        string  `json:"UL"`
	Ll        string  `json:"LL"`
	Vo        float64 `json:"Vo"`
	Va        float64 `json:"Va"`
	AdjFactor float64 `json:"AdjFactor"`
	AdjO      float64 `json:"AdjO"`
	AdjH      float64 `json:"AdjH"`
	AdjL      float64 `json:"AdjL"`
	AdjC      float64 `json:"AdjC"`
	AdjVo     float64 `json:"AdjVo"`
}

func main(){
	godotenv.Load()
	apiKey := os.Getenv("JQUANTS_API_KEY")
	stockCode := "285A0"	// Kioxia
	dbPath := "./stock_prices.db"

	stockPriceResponse, err := fetchStockPrices(apiKey, stockCode)
	if err != nil {
		fmt.Println("Error fetching stock prices:", err)
		return
	}

	db, err := setupDB(dbPath)
	if err != nil {
		fmt.Println("Error setting up database:", err)
		return
	}
	defer db.Close()

	for _, data := range stockPriceResponse.Data {
		err := upsertStockPrice(db, data)
		if err != nil {
			fmt.Println("Error upserting stock price:", err)
			continue
		}
	}
}

func fetchStockPrices(apiKey, code string) (*StockPriceResponse, error) {
	api_url := "https://api.jquants.com/v2/equities/bars/daily"
	params := url.Values{}
	params.Set("code", code)

	reqURL := api_url + "?" + params.Encode()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check if the response status is OK
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var stockPriceResponse StockPriceResponse
	err = json.Unmarshal(body, &stockPriceResponse)
	if err != nil {
		return nil, err
	}

	return &stockPriceResponse, nil
}

func setupDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schemaSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func upsertStockPrice(db *sql.DB, p StockPrice) error {
	_, err := db.Exec(upsertStockPriceSQL,
		p.Code, p.Date, p.O, p.H, p.L, p.C, p.Vo, p.AdjC, p.AdjFactor)
	if err != nil {
		return err
	}
	fmt.Println("Inserted/Updated:", p.Date, p.AdjC)
	return nil
}