package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type ExchangeRate struct {
	TIME string `json:"time"`
	Asset_ID_QUOTE string `json:"asset_id_quote"`
	Rate float32 `json:"rate"`
}

type ExchangeRates struct {
	Asset_ID_BASE string `json:"asset_id_base"`
	Exchange_Rate_Array []ExchangeRate `json:"rates"`
}

func createURL(baseURL string, params ...string) string {

	// Create URL parameters using url.Values
	urlParams := url.Values{}
	for i := 0; i < len(params); i += 2 {
		if i+1 < len(params) {
			urlParams.Add(params[i], params[i+1]) // Add key-value pairs
		}
	}

	// Encoding the parameters
	encodedParams := urlParams.Encode()

	//Checking to make sure the apiKey exists
	if !urlParams.Has("apiKey") {
		panic("The apiKey need to be provided")
	}

	return fmt.Sprintf("%s?%s", baseURL, encodedParams)
}

func makeGetRequest(client *http.Client, finalURL string) ([]byte, error) {

	// Generating a GET Request
	request, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		fmt.Printf("error %s", err)
		panic("error")
	}
	request.Header.Add("Accept", `application/json`)

	// Initiating a request
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		panic("error")
	}

	defer response.Body.Close()

	// Reading the Response
	fmt.Println("Response status: ", response.Status)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

func main() {

	// Load env file
	errEnvFile := godotenv.Load(".env")
	if errEnvFile != nil {
		panic(errEnvFile)
	}

	// Get API key from the env file
	coinApiKey := os.Getenv("CoinApiKey")

	const assetIdBase = "ETH"

	baseURL := "https://rest.coinapi.io/v1/exchangerate/" + assetIdBase

	finalURL := createURL(baseURL, "apiKey", coinApiKey)

	// Starting a http Client
	client := &http.Client{Timeout: time.Duration(10) * time.Second}

	body, err := makeGetRequest(client, finalURL)

	if err != nil {
		panic("Error: " + err.Error())
	}

	var exchange_rates ExchangeRates

	err = json.Unmarshal([]byte(body), &exchange_rates)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", exchange_rates)

	// scanner := bufio.NewScanner(response.Body)
	// for i := 0; scanner.Scan() && i < 5; i++ {
	// 	fmt.Println(scanner.Text())
	// }

	// if err := scanner.Err(); err!= nil {
	// 	panic (err)
	// }
}