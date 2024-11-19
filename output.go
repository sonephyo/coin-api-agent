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
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	loggly "github.com/jamespearly/loggly"
	// "github.com/joho/godotenv"
)

type TradeEntry struct {
    SymbolID     string    `json:"symbol_id"`
    TimeExchange time.Time `json:"time_exchange"`
    TimeCoinAPI  time.Time `json:"time_coinapi"`
    UUID         string    `json:"uuid"`
    Price        float64   `json:"price"`
    Size         float64   `json:"size"`
    TakerSide    string    `json:"taker_side"`
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

	err = logglyClient.EchoSend("info", "TradeEntryArray successfully fetched with size " + strconv.Itoa(len(tradeEntryArray)))
	if err != nil {
		return tradeEntryArray, err
	}
	return tradeEntryArray, nil
}

func startOperation(coinApiKey string, client *http.Client, logglyClient *loggly.ClientType, dataCount *int, svc *dynamodb.Client) {
	var baseURL string
	var parameters map[string]string

	/* -----------------------------------------------------------------
	* Get Request for ETH_USDT_Binance
	* Note: Trades snapshots at a specific time interval	*/

	symbolId := "BINANCE_SPOT_ETH_USDT"

	baseURL = "https://rest.coinapi.io/v1/trades/" + symbolId + "/latest"

	parameters = map[string]string{
		"apiKey": coinApiKey,
		"limit": strconv.Itoa(*dataCount),
	}

	tradeData, err := getTradeDataBySymbolID(client, logglyClient, baseURL, parameters)
	if err != nil {
		logglyClient.EchoSend("error", err.Error())
		fmt.Println(err)
		return
	}

	if len(tradeData) == 0 {
		logglyClient.EchoSend("notice", "The API Request done returned an empty data")
	}

	for _, element := range tradeData {
		//Creating a hashmap for AWS
		av := map[string]types.AttributeValue{
			"time_exchange": &types.AttributeValueMemberS{Value: element.TimeExchange.Format(time.RFC3339)},
			"time_coinapi":  &types.AttributeValueMemberS{Value: element.TimeCoinAPI.Format(time.RFC3339)},
			"uuid":          &types.AttributeValueMemberS{Value: element.UUID},
			"price":         &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", element.Price)},
			"size":          &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", element.Size)},
			"taker_side":    &types.AttributeValueMemberS{Value: element.TakerSide},
		}
		
		// AWS DyanmoDB put operations
		_, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String("pphyo_ETH_tradeEntries"),
			Item: av,
		})

		if err != nil {
			logglyClient.EchoSend("error", err.Error())
			fmt.Println(err)
			return
		}



	}


	fmt.Printf("--- First entry of %s OHLCV ---\n", symbolId)
	if len(tradeData) > 0 {
		fmt.Println(tradeData[0].UUID + " " + tradeData[0].TimeCoinAPI.GoString() + " " + tradeData[0].TakerSide)
	} else {
		fmt.Println("The API request returned an empty data")
	}
	// for index, element := range tradeData {
	// 	fmt.Println("At index " , index, ", value: ", element)
	// }
}


func main() {

	// Command Line Flags
	dataCountPtr := flag.Int("count", 50, "dataCount for each request")
	timeIntervalPtr := flag.Float64("time-interval", 15, "time interval in minutes for each request")

	flag.Parse()

	// AWS DyanmoDB Setup
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
        o.Region = "us-east-1"
        return nil
    })
    if err != nil {
        panic(err)
    }
	svc := dynamodb.NewFromConfig(cfg)

	// Load env file
	// errEnvFile := godotenv.Load(".env")
	// if errEnvFile != nil {
	// 	panic(errEnvFile)
	// }

	// Set up Loggly
	tag := "CoinApiLoggly"
	logglyClient := loggly.New(tag)

	// Get API key from the env file
	coinApiKey := os.Getenv("CoinApiKey")


	// Starting a http Client
	client := &http.Client{Timeout: time.Duration(10) * time.Second}

	// Waiting to start until a new minute
	now := time.Now()
	fmt.Println(now)

    // nextMinute := now.Truncate(time.Minute).Add(time.Minute)
    // duration := nextMinute.Sub(now)
    // fmt.Printf("Waiting for %v to reach the next minute to start the operation ...\n", duration)
    // time.Sleep(duration)


	// Tickers - Second
	ticker := time.NewTicker(time.Duration(*timeIntervalPtr) * time.Minute)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <- ctx.Done():
				return
			case t := <-ticker.C:
				startOperation(coinApiKey, client, logglyClient, dataCountPtr, svc)
				fmt.Println("Ticker time : " + t.GoString())
			}
		}
	}()

	select{}
}