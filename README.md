# Amartha Test

This repository contains a Go application that demonstrates the usage of various packages, including `goget` and `gofpdf`. The application includes functionalities to fetch and generate PDF documents.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Flow](#flow)
- [Dependencies](#dependencies)

## Installation

To get started with this project, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/thisisseptian/amartha-test.git
    cd amartha-test
    ```

2. Install the required Go packages:
    ```sh
    go get github.com/jung-kurt/gofpdf
    go get github.com/thisisseptian/goget
    ```

## Usage

To run the application, use the following command:
```sh
go run main.go
```

## Project Structure

amartha-test/
│
├── main.go        # The main entry point of the application
├── fetch.go       # Contains logic for fetching data using goget
├── pdf.go         # Contains logic for generating PDFs using gofpdf
├── utils.go       # Utility functions
└── README.md      # Project documentation