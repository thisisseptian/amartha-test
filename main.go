package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	hand "amartha-test/handler"
	help "amartha-test/helper"
)

func main() {
	// init helper
	helper := help.NewHelper()
	helper.InitUsers()

	// init handler
	handler := &hand.Handler{
		Helper: helper,
	}

	// init router
	router := mux.NewRouter()

	// list of user routes
	router.HandleFunc("/user/list", handler.Middleware(handler.ListUser)).Methods("GET")
	router.HandleFunc("/user/{user_id}/detail", handler.Middleware(handler.DetailUser)).Methods("GET")

	// list of loan routes
	router.HandleFunc("/loan/list", handler.Middleware(handler.ListLoan)).Methods("GET")
	router.HandleFunc("/loan/{loan_id}/detail", handler.Middleware(handler.DetailLoan)).Methods("GET")
	router.HandleFunc("/loan/submit", handler.Middleware(handler.SubmitLoan)).Methods("POST")
	router.HandleFunc("/loan/{loan_id}/approve", handler.Middleware(handler.ApproveLoan)).Methods("POST")
	router.HandleFunc("/loan/{loan_id}/invest", handler.Middleware(handler.InvestLoan)).Methods("POST")
	router.HandleFunc("/loan/{loan_id}/disburse", handler.Middleware(handler.DisburseLoan)).Methods("POST")

	// list of agreement routes
	router.HandleFunc("/agreement/list", handler.Middleware(handler.ListAgreement)).Methods("GET")
	router.HandleFunc("/agreement/{agreement_id}/view", handler.Middleware(handler.ViewAgreement)).Methods("GET")
	router.HandleFunc("/agreement/{agreement_id}/sign", handler.Middleware(handler.SignAgreement)).Methods("POST")

	fmt.Println("listening server on port :8080")
	http.ListenAndServe(":8080", router)
}
