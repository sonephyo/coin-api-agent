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
3. Set up your Loggly API key and other configuration settings.

## Usage
1. Run the Go program:
    ```bash
    go run main.go
    ```
2. The program will periodically fetch data from the specified API and display it on the console.
3. Log messages will be sent to Loggly, indicating whether the requests were successful or failed, along with the data size collected.

## Example
Here’s a simple example of the console output:

```plaintext
Request successful.
Data collected:
{
  "key1": "value1",
  "key2": "value2"
}
```