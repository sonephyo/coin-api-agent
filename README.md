## CSC482 - Software Deployment

In today’s fast-paced business world, developing an **effective** and **efficient information infrastructure** is crucial. This infrastructure must ensure **high availability**, support **multiple client platforms**, provide **secure authenticated access**, and allow for **efficient updates**, all while respecting **cost constraints**.

This course focuses on **modern software deployment strategies** that address these challenges, equipping students with practical skills to design and implement robust infrastructures that meet operational needs and align with budgetary requirements.

# Go Polling Worker and Loggly

## Purpose
The purpose of this project is to design and build a Go program that periodically collects formatted data via API requests, displays this data on the console, and reports its results (or errors) to Loggly.

Go, a programming language developed by Google, is ideal for writing back-end server code due to its simplicity and performance. In this project, we utilize the `net/http` library to make API requests and Loggly, an online logging service with extensive tagging and search capabilities, to handle logs.

The data collected from the API requests will be stored in a Go `struct` appropriate for the data source. It will also be printed on the console with clear key/value pairs for easy readability.

The worker will log messages to Loggly indicating the success or failure of each request, as well as the amount of data collected. It will also make use of Loggly's built-in tagging feature for easy tracking and management.

## Features
- Periodically fetch data using Go's `net/http` library.
- Display formatted data on the console.
- Report success, failure, and data collection details to Loggly.
- Use Loggly's tagging feature for log categorization.

## Prerequisites
- Go installed: [Install Go](https://go.dev/doc/install)
- A Loggly account: [Loggly Sign Up](https://www.loggly.com/)

## Setup
1. Clone the repository.
2. Install necessary Go dependencies using:
    ```bash
    go mod tidy
    ```
3. Create <b>.env</b> file in the same directory and import API Key for <i>CoinApiKey</i> and <i>LOGGLY_TOKEN</i>


## Usage
1. Run the Go program:
    ```bash
    go run code.go
    ```
2. The program will periodically fetch data from the specified API and display it on the console. (Adjust accordingly on how you want your data to be displayed)
3. Log messages will be sent to Loggly, indicating whether the requests were successful or failed, along with the data size collected.

## Example
Here’s a simple example of the console output from 

IDX_REFRATE_PRIMKT_ETH_USDT TimeSeries

```
[{2024-09-28 14:00:00 +0000 UTC 2024-09-28 15:00:00 +0000 UTC 2024-09-28 14:00:00 +0000 UTC 2024-09-28 14:59:00 +0000 UTC 2659.61 2675.29 2659.61 2672.85 0} {2024-09-28 15:00:00 +0000 UTC 2024-09-28 16:00:00 +0000 UTC 2024-09-28 15:00:00 +0000 UTC 2024-09-28 15:59:00 +0000 UTC 2673.56 2673.56 2667.75 2669.62 0} {2024-09-28 16:00:00 +0000 UTC 2024-09-28 17:00:00 +0000 UTC 2024-09-28 16:00:00 +0000 UTC 2024-09-28 16:59:00 +0000 UTC 2669.96 2671.73 2662.6 2665.14 0} {2024-09-28 17:00:00 +0000 UTC 2024-09-28 18:00:00 +0000 UTC 2024-09-28 17:00:00 +0000 UTC 2024-09-28 17:59:00 +0000 UTC 2665.07 2682.15 2663.82 2674.77 0} {2024-09-28 18:00:00 +0000 UTC 2024-09-28 19:00:00 +0000 UTC 2024-09-28 18:00:00 +0000 UTC 2024-09-28 18:59:00 +0000 UTC 2675.47 2683.3 2675 2675 0} {2024-09-28 19:00:00 +0000 UTC 2024-09-28 20:00:00 +0000 UTC 2024-09-28 19:00:00 +0000 UTC 2024-09-28 19:59:00 +0000 UTC 2675.3 2679.88 2674.71 2678.16 0} {2024-09-28 20:00:00 +0000 UTC 2024-09-28 21:00:00 +0000 UTC 2024-09-28 20:00:00 +0000 UTC 2024-09-28 20:59:00 +0000 UTC 2677.89 2679.88 2673.68 2675.82 0} {2024-09-28 21:00:00 +0000 UTC 2024-09-28 22:00:00 +0000 UTC 2024-09-28 21:00:00 +0000 UTC 2024-09-28 21:59:00 +0000 UTC 2675.51 2676.13 2657.76 2670.45 0} {2024-09-28 22:00:00 +0000 UTC 2024-09-28 23:00:00 +0000 UTC 2024-09-28 22:00:00 +0000 UTC 2024-09-28 22:59:00 +0000 UTC 2670.65 2678.68 2665.99 2676.66 0} {2024-09-28 23:00:00 +0000 UTC 2024-09-29 00:00:00 +0000 UTC 2024-09-28 23:00:00 +0000 UTC 2024-09-28 23:59:00 +0000 UTC 2675 2679.68 2669.23 2675.6 0}]
```