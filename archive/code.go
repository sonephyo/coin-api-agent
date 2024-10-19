// Package <package_name> provides <a brief description of the package's purpose>.
//
// This package is responsible for <a more detailed explanation of the core functionality>.
// It includes features such as <list of primary features or responsibilities>.
//
//
// Author: Phone Pyae Sone Phyo (Soney)
// Created: 10/03/2024

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
	loggly "github.com/jamespearly/loggly"
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

type TimeSeriesEntry struct {
    TimePeriodStart time.Time `json:"time_period_start"`
    TimePeriodEnd   time.Time `json:"time_period_end"`
    TimeOpen        time.Time `json:"time_open"`
    TimeClose       time.Time `json:"time_close"`
    ValueOpen       float64   `json:"value_open"`
    ValueHigh       float64   `json:"value_high"`
    ValueLow        float64   `json:"value_low"`
    ValueClose      float64   `json:"value_close"`
    ValueCount      int       `json:"value_count"`
}

type TradeEntry struct {
	TimePeriodStart time.Time `json:"time_period_start"`
	TimePeriodEnd   time.Time `json:"time_period_end"`
	TimeOpen        time.Time `json:"time_open"`
	TimeClose       time.Time `json:"time_close"`
	PriceOpen       float64   `json:"price_open"`
	PriceHigh       float64   `json:"price_high"`
	PriceLow        float64   `json:"price_low"`
	PriceClose      float64   `json:"price_close"`
	VolumeTraded    float64   `json:"volume_traded"`
	TradesCount     int       `json:"trades_count"`
}


func createURL(baseURL string, params map[string]string) (string, error) {

	// Create URL parameters using url.Values
	urlParams := url.Values{}
	for key, value := range params {
        urlParams.Add(key, value) // Add key-value pairs
    }

	// Encoding the parameters
	encodedParams := urlParams.Encode()

	//Checking to make sure the apiKey exists
	if !urlParams.Has("apiKey") || urlParams.Get("apiKey") == "" {
		return "", fmt.Errorf("apiKey is needed as a query parameter to make requests")
	}

	return fmt.Sprintf("%s?%s", baseURL, encodedParams), nil
}

func makeGetRequest(client *http.Client, finalURL string) ([]byte, error) {

	// Generating a GET Request
	request, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate a new Get request: %w", err)
	}
	request.Header.Add("Accept", `application/json`)

	// Initiating a request
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate request: %w", err)
	}

	defer response.Body.Close()

	// Reading the Response
	fmt.Println("Response status: ", response.Status)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response from %s: %w", finalURL, err)
	}

	return body, nil
}

func getAllExchangeRatesByAssetIDBase(client *http.Client, logglyClient *loggly.ClientType, baseURL string, parameters map[string]string) (ExchangeRates, error)  {

	finalURL, err := createURL(baseURL, parameters)
	if err != nil {
		return ExchangeRates{}, err
	}

	body, err := makeGetRequest(client, finalURL)
	if err != nil {
		return ExchangeRates{}, err
	}

	var exchange_rates ExchangeRates

	err = json.Unmarshal([]byte(body), &exchange_rates)
	if err != nil {
		return ExchangeRates{}, err
	}


	err = logglyClient.EchoSend("info", "ExchangeRates successfully fetched")
	if err != nil {
		return exchange_rates, err
	}
	return exchange_rates, nil
}

func getTimeSeriesByIndexId(client *http.Client, logglyClient *loggly.ClientType, baseURL string, parameters map[string]string) ([]TimeSeriesEntry, error) {
	finalURL, err := createURL(baseURL, parameters)
	if err != nil {
		return nil, err
	}

	body, err := makeGetRequest(client, finalURL)
	if err != nil {
		return nil, err
	}

	var timeSeriesEntryArray []TimeSeriesEntry

	err = json.Unmarshal([]byte(body), &timeSeriesEntryArray)
	if err != nil {
		return nil, err
	}


	err = logglyClient.EchoSend("info", "TimeSeriesEntryArray successfully fetched")
	if err != nil {
		return timeSeriesEntryArray, err
	}
	return timeSeriesEntryArray, nil
}

func getTradeDataBySymbolID(client *http.Client, logglyClient *loggly.ClientType, baseURL string, parameters map[string]string) ([]TradeEntry, error) {
	finalURL, err := createURL(baseURL, parameters)
	if err != nil {
		return nil, err
	}

	body, err := makeGetRequest(client, finalURL)
	if err != nil {
		return nil, err
	}

	var tradeEntryArray []TradeEntry

	err = json.Unmarshal([]byte(body), &tradeEntryArray)
	if err != nil {
		return nil, err
	}

	err = logglyClient.EchoSend("info", "TradeEntryArray successfully fetched")
	if err != nil {
		return tradeEntryArray, err
	}
	return tradeEntryArray, nil
}


func main() {

	// Load env file
	errEnvFile := godotenv.Load(".env")
	if errEnvFile != nil {
		panic(errEnvFile)
	}

	// Set up Loggly
	tag := "CoinApiLoggly"
	logglyClient := loggly.New(tag)

	// Get API key from the env file
	coinApiKey := os.Getenv("CoinApiKey")

	// Starting a http Client
	client := &http.Client{Timeout: time.Duration(10) * time.Second}


	var baseURL string
	var parameters map[string]string
	
	/* ----------------------------------------------------------------
	* Get Request for exchange rates based on assetIdBase
	*/

	assetIdBase := "ETH"
	parameters = map[string]string{
		"apiKey": coinApiKey,
	}

	baseURL = "https://rest.coinapi.io/v1/exchangerate/" + assetIdBase
	_, err := getAllExchangeRatesByAssetIDBase(client, logglyClient, baseURL, parameters)
	if err != nil {
		logglyClient.EchoSend("error", err.Error())
		fmt.Println(err)
		return
	}
	fmt.Printf("--- ExchangeRates based on %s ---\n", assetIdBase)
	// fmt.Println(exchangeRates.Exchange_Rate_Array[0:10])


	/* -----------------------------------------------------------------
	* Get Request for ETH_USDT_TimeSeries
	*/
	indexId := "IDX_REFRATE_PRIMKT_ETH_USDT"
	baseURL = "https://rest.coinapi.io/v1/indexes/" + indexId + "/timeseries"

	parameters = map[string]string{
		"apiKey": coinApiKey,
		"period_id": "1HRS",
		"time_start": "2024-09-28T14:00:00Z",
		"time_end": "2024-09-30T14:00:00Z",
	}

	timeSeriesEntryArray, err := getTimeSeriesByIndexId(client, logglyClient, baseURL, parameters)
	if err != nil {
		logglyClient.EchoSend("error", err.Error())
		fmt.Println(err)
		return
	}
	fmt.Printf("--- %s TimeSeries ---\n", indexId)
	fmt.Println(timeSeriesEntryArray[0:10])
	// for index, element := range timeSeriesEntryArray {
	// 	fmt.Println("At index " , index, ", value: ", element)
	// }

	/* -----------------------------------------------------------------
	* Get Request for ETH_USDT_Binance
	* Note: Get OHLCV latest timeseries data returned in time descending order. Data can be requested by the period and for the specific symbol eg BITSTAMP_SPOT_BTC_USD, if you need to query timeseries by asset pairs eg. BTC/USD, then please reffer to the Exchange Rates Timeseries data
	*/

	symbolId := "BINANCE_SPOT_ETH_USDT"


	baseURL = "https://rest.coinapi.io/v1/ohlcv/" + symbolId + "/latest"

	parameters = map[string]string{
		"apiKey": coinApiKey,
		"period_id": "1HRS",
		"time_start": "2024-09-28T14:00:00Z",
		"time_end": "2024-09-30T14:00:00Z",
	}

	tradeData, err := getTradeDataBySymbolID(client, logglyClient, baseURL, parameters)
	if err != nil {
		logglyClient.EchoSend("error", err.Error())
		fmt.Println(err)
		return
	}
	fmt.Printf("--- %s OHLCV ---\n", symbolId)
	fmt.Println(tradeData[0:10])
	// for index, element := range tradeData {
	// 	fmt.Println("At index " , index, ", value: ", element)
	// }
	
}