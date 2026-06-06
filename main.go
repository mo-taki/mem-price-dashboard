package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

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

	fmt.Println("Response:", string(body))
}