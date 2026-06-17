package jquants

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

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

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) FetchStockPrices(code string) ([]StockPrice, error) {
	api_url := "https://api.jquants.com/v2/equities/bars/daily"
	params := url.Values{}
	params.Set("code", code)

	reqURL := api_url + "?" + params.Encode()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
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

	return stockPriceResponse.Data, nil
}
