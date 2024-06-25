package main

import (
	"amartha-test/handler"
	"amartha-test/helper"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func dynamicLoanRouteHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/loan/") {
		path = strings.TrimPrefix(path, "/loan/")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 {
			loanID, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				http.Error(w, "invalid loan id", http.StatusBadRequest)
				return
			}

			if loanID == 0 {
				http.Error(w, "loan id is empty", http.StatusBadRequest)
				return
			}

			action := parts[1]
			switch action {
			case "approve":
				handler.ApproveLoan(w, r, loanID)
			case "invest":
				handler.InvestLoan(w, r, loanID)
			// case "disburse":
			// 	disburseLoan(w, r, id)
			default:
				http.Error(w, "invalid action", http.StatusBadRequest)
			}
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	// for testing purpose
	helper.InitUsers()
	http.HandleFunc("/user/list", handler.ListUser)

	// route
	http.HandleFunc("/loan/list", handler.ListLoan)
	http.HandleFunc("/loan/create", handler.SubmitLoan)
	http.HandleFunc("/loan/", dynamicLoanRouteHandler)
	// http.HandleFunc("/loans/", func(w http.ResponseWriter, r *http.Request) {
	// 	switch {
	// 	case r.URL.Path[len(r.URL.Path)-len("/approve"):] == "/approve":
	// 		ApproveLoan(w, r)
	// 	case r.URL.Path[len(r.URL.Path)-len("/invest"):] == "/invest":
	// 		InvestInLoan(w, r)
	// 	case r.URL.Path[len(r.URL.Path)-len("/disburse"):] == "/disburse":
	// 		DisburseLoan(w, r)
	// 	default:
	// 		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	// 	}
	// })

	fmt.Println("Starting Server on port 8080")
	http.ListenAndServe(":8080", nil)
}
