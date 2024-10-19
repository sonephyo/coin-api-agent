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
4. Using AWS DynamoDB
   1. If you plan to directly run code in the terminal, you are expected to do "aws configure" and a table set up in DynamoDB. (NOTE: make sure to uncomment godotenv package and its respective code in the main method)
   2. If you are using Docker, you will need to add AWS Credentials into env file (NOTE: also build the docker image)

## Usage
1. Running the program
   1.  In terminal
    ```bash
    go run code.go -count=<int> -time-interval=<int>
    ```
    2. In docker
    ```
    docker run -d --env-file .env <image_name>
    ```

## Result
- The program will update AWS DyanmoDB with Trades Data
- Every action will be logged in Loggly for monitoring