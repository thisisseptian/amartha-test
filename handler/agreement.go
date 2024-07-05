package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"amartha-test/constant"
	"amartha-test/model"
)

// ListAgreement is handler to get list of agreements
func (h *Handler) ListAgreement(w http.ResponseWriter, r *http.Request) {
	// 1. get agreement list
	agreements := h.Helper.GetAgreements()
	if len(agreements) == 0 {
		log.Println("[ListAgreement] agreement list is empty")
		h.RenderResponse(w, r, "", http.StatusNotFound, "[ListAgreement] agreement list is empty")
		return
	}

	// 2. render response
	h.RenderResponse(w, r, agreements, http.StatusOK, "")
}

// ViewAgreement is handler to view agreement detail
func (h *Handler) ViewAgreement(w http.ResponseWriter, r *http.Request) {
	// 1. get vars
	vars := mux.Vars(r)
	agreementID, err := strconv.ParseInt(vars["agreement_id"], 10, 64)
	if err != nil {
		log.Printf("[ViewAgreement] failed parse int, with error: %+v", err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[ViewAgreement] failed parse int, with error: %+v", err))
		return
	}

	// 2. sanitize payload
	if agreementID == 0 {
		log.Println("[ViewAgreement] agreement id is zero")
		h.RenderResponse(w, r, "", http.StatusBadRequest, "[ViewAgreement] agreement id is zero")
		return
	}

	// 3. get agreement by agreement id
	agreement := h.Helper.GetAgreementByAgreementID(agreementID)
	if agreement.AggrementID == 0 {
		log.Printf("[ViewAgreement][AgreementID: %d] agreement data is not found", agreementID)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[ViewAgreement][AgreementID: %d] agreement data is not found", agreementID))
		return
	}

	// 4. decode agreement base 64
	pdfData, err := base64.StdEncoding.DecodeString(agreement.DocumentData)
	if err != nil {
		log.Printf("[ViewAgreement][AgreementID: %d] failed to decode base64 pdf data with error: %+v", agreementID, err)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[ViewAgreement][AgreementID: %d] failed to decode base64 pdf data with error: %+v", agreementID, err))
		return
	}

	// 5. render response
	h.RenderPDFResponse(w, pdfData, http.StatusOK)
}

// SignAgreement is handler to sign agreement
func (h *Handler) SignAgreement(w http.ResponseWriter, r *http.Request) {
	// 1. get vars
	vars := mux.Vars(r)
	agreementID, err := strconv.ParseInt(vars["agreement_id"], 10, 64)
	if err != nil {
		log.Printf("[SignAgreement] failed parse int, with error: %+v", err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[SignAgreement] failed parse int, with error: %+v", err))
		return
	}

	// 2. decode body
	var sign model.Sign
	err = json.NewDecoder(r.Body).Decode(&sign)
	if err != nil {
		log.Printf("[SignAgreement][AgreementID: %d] fail decode body with error: %+v", agreementID, err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[SignAgreement][AgreementID: %d] fail decode body with error: %+v", agreementID, err))
		return
	}

	// 3. sanitize payload
	if sign.LoanID == 0 {
		log.Printf("[SignAgreement][AgreementID: %d] loan id is empty", agreementID)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[SignAgreement][AgreementID: %d] loan id is empty", agreementID))
		return
	}
	if sign.UserID == 0 {
		log.Printf("[SignAgreement][AgreementID: %d][LoanID: %d] user id is empty", agreementID, sign.LoanID)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[SignAgreement][AgreementID: %d][LoanID: %d] user id is empty", agreementID, sign.LoanID))
		return
	}

	// 4. get loan by loan id
	loan := h.Helper.GetLoanByLoanID(sign.LoanID)
	if loan.LoanID == 0 {
		log.Printf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] loan data not found", agreementID, sign.LoanID, sign.UserID)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] loan data not found", agreementID, sign.LoanID, sign.UserID))
		return
	}

	// 5. check loan status
	if loan.Status != constant.LoanStatusInvested {
		log.Printf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] loan status invalid, current status is: %s", agreementID, sign.LoanID, sign.UserID, constant.GetLoanStatusDesc(loan.Status))
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] loan status invalid, current status is: %s", agreementID, sign.LoanID, sign.UserID, constant.GetLoanStatusDesc(loan.Status)))
		return
	}

	// 6. get user by user id
	user := h.Helper.GetUserByUserID(sign.UserID)
	if user.UserID == 0 {
		log.Printf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] user data not found", agreementID, sign.LoanID, sign.UserID)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] user data not found", agreementID, sign.LoanID, sign.UserID))
		return
	}

	// 7. get agreement by agreement id
	agreement := h.Helper.GetAgreementByAgreementID(agreementID)
	if agreement.AggrementID == 0 {
		log.Printf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] agreement data not found", agreementID, sign.LoanID, sign.UserID)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] agreement data not found", agreementID, sign.LoanID, sign.UserID))
		return
	}

	// 8. wrong user to sign
	if agreement.UserID != sign.UserID {
		log.Printf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] wrong user to sign this agreement", agreementID, sign.LoanID, sign.UserID)
		h.RenderResponse(w, r, "", http.StatusForbidden, fmt.Sprintf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] wrong user to sign this agreement", agreementID, sign.LoanID, sign.UserID))
		return
	}

	// 9. check agreement sign
	if agreement.IsSigned {
		log.Printf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] agreement already signed", agreementID, sign.LoanID, sign.UserID)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] agreement already signed", agreementID, sign.LoanID, sign.UserID))
		return
	}

	// 10. update agreement sign
	agreement.IsSigned = true
	h.Helper.UpsertAgreement(agreement)

	// 11. generate agreement sign pdf
	err = h.Helper.GenerateSignedAgreementPDF(&loan, sign.UserID)
	if err != nil {
		log.Printf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] fail to generate signed agreement pdf with error: %+v", agreementID, sign.LoanID, sign.UserID, err)
		h.RenderResponse(w, r, "", http.StatusInternalServerError, fmt.Sprintf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] fail to generate signed agreement pdf with error: %+v", agreementID, sign.LoanID, sign.UserID, err))
		return
	}

	// 12. check based on user type
	if user.UserType == constant.UserTypeLender {
		// 12a. check agreement is completely signed by all lender
		isCompletelySignedByLender, err := h.Helper.CheckAgreementCompletelySignedByLender(loan)
		if err != nil {
			log.Printf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] check agreement completely signed by lender got fail with error: %+v", agreementID, sign.LoanID, sign.UserID, err)
			h.RenderResponse(w, r, "", http.StatusInternalServerError, fmt.Sprintf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] check agreement completely signed by lender got fail with error: %+v", agreementID, sign.LoanID, sign.UserID, err))
			return
		}
		if isCompletelySignedByLender {
			// 12b. generate borrower agreement pdf
			err = h.Helper.GenerateBorrowerAgreementPDF(&loan)
			if err != nil {
				log.Printf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] fail to generate borrower agreement pdf with error: %+v", agreementID, sign.LoanID, sign.UserID, err)
				h.RenderResponse(w, r, "", http.StatusInternalServerError, fmt.Sprintf("[SignAgreement][AgreementID: %d][LoanID: %d][UserID: %d] fail to generate borrower agreement pdf with error: %+v", agreementID, sign.LoanID, sign.UserID, err))
				return
			}
		}
	} else if user.UserType == constant.UserTypeBorrower {
		// 12c. if borrower sign, must be final sign, and move status to signed
		loan.Status = constant.LoanStatusSigned
		loan.StatusDesc = constant.GetLoanStatusDesc(loan.Status)
		h.Helper.UpsertLoan(loan)
	}

	// 13. render response
	h.RenderResponse(w, r, loan, http.StatusOK, "")
}
