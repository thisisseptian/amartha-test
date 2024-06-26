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
			case "detail":
				handler.DetailLoan(w, r, loanID)
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

func dynamicAgreementRouterHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/agreement/") {
		path = strings.TrimPrefix(path, "/agreement/")
		parts := strings.Split(path, "/")
		if len(parts) >= 1 {
			agreementID, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				http.Error(w, "invalid agreement id", http.StatusBadRequest)
				return
			}

			if agreementID == 0 {
				http.Error(w, "agreement id is empty", http.StatusBadRequest)
				return
			}

			action := parts[1]
			switch action {
			case "view":
				handler.DetailAgreement(w, r, agreementID)
			case "sign":
				handler.SignAgreement(w, r, agreementID)
			default:
				http.Error(w, "invalid action", http.StatusBadRequest)
			}
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	// route user
	helper.InitUsers()
	http.HandleFunc("/user/list", handler.ListUser)

	// route loan
	http.HandleFunc("/loan/list", handler.ListLoan)
	http.HandleFunc("/loan/create", handler.SubmitLoan)
	http.HandleFunc("/loan/", dynamicLoanRouteHandler)

	// route agreement
	http.HandleFunc("/agreement/list", handler.ListAgreement)
	http.HandleFunc("/agreement/", dynamicAgreementRouterHandler)

	fmt.Println("listening server on port :8080")
	http.ListenAndServe(":8080", nil)
}
