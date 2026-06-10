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

type StockPriceResponse struct {
	Data []struct {
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
	} `json:"data"`
}

func main(){
	godotenv.Load()
	apiKey := os.Getenv("JQUANTS_API_KEY")
	fmt.Println("API Key:", apiKey)

	api_url := "https://api.jquants.com/v2/equities/bars/daily"
	params := url.Values{}
	params.Set("code", "285A0")

	reqURL := api_url + "?" + params.Encode()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	db, err := sql.Open("sqlite", "./stock_prices.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(schemaSQL)
	if err != nil {
		fmt.Println("Error executing schema:", err)
		return
	}

	var stockPriceResponse StockPriceResponse
	err = json.Unmarshal(body, &stockPriceResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	insertSQL := `INSERT INTO stock_prices (code, date, open, high, low, close, volume, adj_close, adj_factor)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(code, date)
		DO UPDATE SET open=excluded.open, high=excluded.high, low=excluded.low, close=excluded.close, volume=excluded.volume, adj_close=excluded.adj_close, adj_factor=excluded.adj_factor`
	for _, data := range stockPriceResponse.Data {
		_, err = db.Exec(insertSQL,
			data.Code, data.Date, data.O, data.H, data.L, data.C, data.Vo, data.AdjC, data.AdjFactor)
		if err != nil {
			fmt.Println("Error inserting data:", err)
			return
		}
		fmt.Println("Inserted/Updated:", data.Date, data.AdjC)
	}
}