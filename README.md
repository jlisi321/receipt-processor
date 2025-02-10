# Receipt Processor Project

This project is a simple receipt processing API that calculates loyalty points for a given receipt. The points are determined based on several factors such as retailer name, total amount, time of purchase, and item descriptions.

## Overview

The project consists of two main functionalities:

- **Process Receipts**: A POST endpoint that accepts an incoming receipt and calculates loyalty points based on various criteria.
- **Get Points**: A GET endpoint that retrieves the calculated points for a given receipt using the receipt's unique ID.

## Features

- **Process a Receipt**: When a receipt is posted, the application processes it and generates a unique receipt ID, along with calculated loyalty points.
- **Get Points**: Using the receipt's ID, users can fetch the calculated loyalty points.
- **Loyalty Points Calculation**: Points are calculated based on several conditions such as:
  - Alphanumeric characters in the retailer name.
  - Round dollar amounts or multiples of 0.25 for the total.
  - Number of items in the receipt.
  - Description length of items and price calculations.
  - Odd days in the purchase date and specific time frames.

## Project Structure

- **main.go**: Contains the API logic and routes for processing and fetching receipt data.

## Setup Instructions

### Prerequisites

Go 1.18+ installed on your system.

### Running the Application

1. **Clone Repository**
   ```bash
   git clone https://github.com/jlisi321/receipt-processor.git
   cd receipt-processor

2. **Install the dependencies**
    ```bash
   go mod tidy

3. **Run the application**
    ```bash 
    go run main.go

4. The API will start running at http://localhost:8080

## API Endpoints

- **POST /receipts/process**: Process a new receipt and generate points.
  - Request body:
    ```json
    {
      "retailer": "Retailer Name",
      "purchaseDate": "2025-02-10",
      "purchaseTime": "15:30",
      "items": [
        {
          "shortDescription": "Item 1",
          "price": "12.50"
        },
        {
          "shortDescription": "Item 2",
          "price": "5.00"
        }
      ],
      "total": "17.50"
    }
  - Response:
    ```json
    {
      "id": "generated-receipt-id"
    }
    
- **GET /receipts/{id}/points**: Retrieve points for a specific receipt.
    - Example request: GET /receipts/generated-receipt-id/points
    - Response:
      ```json
      {
        "points": 100
      }

## Example curl Commands

Here are some examples of how you can interact with the API using curl:

- **Posting a Receipt to Process**
  This example submits a receipt for processing
  ```bash
  curl -X POST http://localhost:8080/receipts/process -d '{
    "retailer": "M&M Corner Market",
    "purchaseDate": "2022-03-20",
    "purchaseTime": "14:33",
    "items": [
      {
        "shortDescription": "Gatorade",
        "price": "2.25"
      },
      {
        "shortDescription": "Gatorade",
        "price": "2.25"
      },
      {
        "shortDescription": "Gatorade",
        "price": "2.25"
      },
      {
        "shortDescription": "Gatorade",
        "price": "2.25"
      }
    ],
    "total": "9.00"
  }' -H "Content-Type: application/json"
  ```
  **Expected Response**
  ```json
  {
    "id": "generated-receipt-id
  }

- **Retrieving Points for a Receipt**
  Once a receipt has been processed (using the POST method above), you can fetch the points for that receipt by using the following GET request.
  ```bash
  curl http://localhost:8080/receipts/27891132-8bf3-4b49-b334-bbc54a97a045/points
  ```
  **Expected Response**
  ```json
  {
    "points": 100
  }

## Learn More

For more details on the Go programming language and how it works with HTTP APIs, you can refer to the official documentation:

- Gorilla Mux Documentation: https://www.gorillatoolkit.org/pkg/mux
- Go Documentation: https://golang.org/doc/