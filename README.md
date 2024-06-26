# Amartha Test

This repository was created as a test that can solve the simple problem of dummy loan service in Amartha (system design and abstraction)

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
    ```

## Usage

To run the application, use the following command:
```sh
go run main.go
```

## Project Structure

```sh
amartha-test/
│
├── main.go        # The main entry point of the application
├── collection     # Contains Postman collection for testing purposes
├── constant       # Contains constants used in the repository, such as loan statuses or user types
├── handler        # Contains handler functions for REST API endpoints
├── helper         # Contains helper functions; since no database is used, these functions are used to access data in memory
├── model          # Contains object structs and their associated methods
└── README.md      # Project documentation
```

## Flow

To test the flow, you can use Postman collection (in folder collection) to hit the endpoints in the following order:
```sh
1. Hit Loan Create
    - Requires borrower_id, principal_amount, and interest_rate
    - Creates a new loan with the status proposed
2. Hit Loan Approve
    - Requires field_validator_employee, picture_proof, and approval_date
    - Changes the loan status to approved
3. Hit Loan Invest 
    - Requires loan_id, lender_id, and invested_amount
    - Can handle multiple lenders, with the amounts accumulating but not exceeding the loan limit
    - After each lender invests, they receive their own organizer-lender agreement URL, and the loan status changes to invested
4. Hit Agreement Sign
    - Requires agreement_id, loan_id, and user_id
    - The agreement URL can be clicked to display the PDF
    - Each lender must sign their organizer-lender agreement URL
    - Once all lenders sign their agreements, the borrower receives the organizer-borrower agreement URL
    - The borrower must sign the organizer-borrower agreement URL for the loan status to change to loan disbursed
```

Note: I have also created several APIs to assist in debugging, mostly for getting lists and details:
```sh
1. Loan List
2. Loan Detail
3. User List
4. Agreement List
```

## Dependencies

To facilitate PDF creation, this repository uses the `gofpdf` package. Alternatively, we could use `wkhtmltopdf`, but it would take longer since i need to design the HTML first.

