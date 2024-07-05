package handler

import (
	"net/http"

	"amartha-test/helper"
)

type IHandler interface {
	// handler user
	ListUser(w http.ResponseWriter, r *http.Request)
	DetailUser(w http.ResponseWriter, r *http.Request)

	// handler loan
	ListLoan(w http.ResponseWriter, r *http.Request)
	DetailLoan(w http.ResponseWriter, r *http.Request)
	SubmitLoan(w http.ResponseWriter, r *http.Request)
	ApproveLoan(w http.ResponseWriter, r *http.Request)
	InvestLoan(w http.ResponseWriter, r *http.Request)
	DisburseLoan(w http.ResponseWriter, r *http.Request)

	// handler agreement
	ListAgreement(w http.ResponseWriter, r *http.Request)
	DetailAgreement(w http.ResponseWriter, r *http.Request)
	SignAgreement(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	IHandler
	Helper helper.IHelper
}

func NewHandler(helper helper.IHelper) *Handler {
	return &Handler{
		Helper: helper,
	}
}
