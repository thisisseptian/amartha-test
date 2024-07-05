package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"amartha-test/constant"
	"amartha-test/model"
)

// ListLoan is handler to get list of loans
func (h *Handler) ListLoan(w http.ResponseWriter, r *http.Request) {
	// 1. get loan list
	loans := h.Helper.GetLoans()
	if len(loans) == 0 {
		log.Println("[ListLoan] loan list is empty")
		h.RenderResponse(w, r, "", http.StatusNotFound, "[ListLoan] loan list is empty")
		return
	}

	// 2. render response
	h.RenderResponse(w, r, loans, http.StatusOK, "")
}

// DetailLoan is handler to get loan detail
func (h *Handler) DetailLoan(w http.ResponseWriter, r *http.Request) {
	// 1. get vars
	vars := mux.Vars(r)
	loanID, err := strconv.ParseInt(vars["loan_id"], 10, 64)
	if err != nil {
		log.Printf("[DetailLoan] failed parse int, with error: %+v", err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[DetailLoan] failed parse int, with error: %+v", err))
		return
	}

	// 2. sanitize payload
	if loanID == 0 {
		log.Println("[DetailLoan] loan id is zero")
		h.RenderResponse(w, r, "", http.StatusBadRequest, "[DetailLoan] loan id is zero")
		return
	}

	// 3. get loan by loan id
	loan := h.Helper.GetLoanByLoanID(loanID)
	if loan.LoanID == 0 {
		log.Printf("[DetailLoan][LoanID: %d] loan data is not found", loanID)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[DetailLoan][LoanID: %d] loan data is not found", loanID))
		return
	}

	// 4. render response
	h.RenderResponse(w, r, loan, http.StatusOK, "")
}

// SubmitLoan is handler to create new loan
func (h *Handler) SubmitLoan(w http.ResponseWriter, r *http.Request) {
	// 1. decode body
	var loan model.Loan
	err := json.NewDecoder(r.Body).Decode(&loan)
	if err != nil {
		log.Printf("[SubmitLoan] fail decode body with error: %+v", err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[SubmitLoan] fail decode body with error: %+v", err))
		return
	}

	// 2. sanitize payload
	if loan.BorrowerID == 0 {
		log.Println("[SubmitLoan] borrower id is empty")
		h.RenderResponse(w, r, "", http.StatusBadRequest, "[SubmitLoan] borrower id is empty")
		return
	}
	if loan.PrincipalAmount == 0 {
		log.Printf("[SubmitLoan][BorrowerID: %d] principal amount is empty", loan.BorrowerID)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[SubmitLoan][BorrowerID: %d] principal amount is empty", loan.BorrowerID))
		return
	}
	if loan.InterestRate < 0 || loan.InterestRate > 1 {
		log.Printf("[SubmitLoan][BorrowerID: %d][Amount: %.2f] interest rate is invalid", loan.BorrowerID, loan.PrincipalAmount)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[SubmitLoan][BorrowerID: %d][Amount: %.2f] interest rate is invalid", loan.BorrowerID, loan.PrincipalAmount))
		return
	}

	// 3. get borrower by user id
	borrower := h.Helper.GetUserByUserID(loan.BorrowerID)
	if borrower.UserID == 0 {
		log.Printf("[SubmitLoan][BorrowerID: %d][Amount: %.2f][Rate: %.2f] borrower data is not found", loan.BorrowerID, loan.PrincipalAmount, loan.InterestRate)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[SubmitLoan][BorrowerID: %d][Amount: %.2f][Rate: %.2f] borrower data is not found", loan.BorrowerID, loan.PrincipalAmount, loan.InterestRate))
		return
	}

	// 4. check user status
	if borrower.UserType != constant.UserTypeBorrower {
		log.Printf("[SubmitLoan][BorrowerID: %d][Amount: %.2f][Rate: %.2f] user type is not borrower", loan.BorrowerID, loan.PrincipalAmount, loan.InterestRate)
		h.RenderResponse(w, r, "", http.StatusForbidden, fmt.Sprintf("[SubmitLoan][BorrowerID: %d][Amount: %.2f][Rate: %.2f] user type is not borrower", loan.BorrowerID, loan.PrincipalAmount, loan.InterestRate))
		return
	}

	// 5. create loan
	loan.LoanID = h.Helper.GenerateIncrementalLoanID()
	loan.Status = constant.LoanStatusProposed
	loan.StatusDesc = constant.GetLoanStatusDesc(loan.Status)
	h.Helper.UpsertLoan(loan)

	// 6. render response
	h.RenderResponse(w, r, loan, http.StatusCreated, "")
}

// ApproveLoan is handler to approve loan
func (h *Handler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	// 1. get vars
	vars := mux.Vars(r)
	loanID, err := strconv.ParseInt(vars["loan_id"], 10, 64)
	if err != nil {
		log.Printf("[ApproveLoan] failed parse int, with error: %+v", err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[ApproveLoan] failed parse int, with error: %+v", err))
		return
	}

	// 2. decode body
	var approvalInfo model.ApprovalInfo
	err = json.NewDecoder(r.Body).Decode(&approvalInfo)
	if err != nil {
		log.Printf("[ApproveLoan][LoanID: %d] fail decode body with error: %+v", loanID, err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[ApproveLoan][LoanID: %d] fail decode body with error: %+v", loanID, err))
		return
	}

	// 3. sanitize payload
	if approvalInfo.PictureProof == "" {
		log.Printf("[ApproveLoan][LoanID: %d] picture proof is empty", loanID)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[ApproveLoan][LoanID: %d] picture proof is empty", loanID))
		return
	}
	if approvalInfo.FieldValidatorEmployeeID == 0 {
		log.Printf("[ApproveLoan][LoanID: %d] field validator employee id is empty", loanID)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[ApproveLoan][LoanID: %d] field validator employee id is empty", loanID))
		return
	}

	// 4. get loan by loan id
	loan := h.Helper.GetLoanByLoanID(loanID)
	if loan.LoanID == 0 {
		log.Printf("[ApproveLoan][LoanID: %d] loan data is not found", loanID)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[ApproveLoan][LoanID: %d] loan data is not found", loanID))
		return
	}

	// 5. check loan status
	if loan.Status != constant.LoanStatusProposed {
		log.Printf("[ApproveLoan][LoanID: %d] loan status is invalid, current status is: %s", loanID, constant.GetLoanStatusDesc(loan.Status))
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[ApproveLoan][LoanID: %d] loan status is invalid, current status is: %s", loanID, constant.GetLoanStatusDesc(loan.Status)))
		return
	}

	// 6. get field validator employee by user id
	fieldValidatorEmployee := h.Helper.GetUserByUserID(approvalInfo.FieldValidatorEmployeeID)
	if fieldValidatorEmployee.UserID == 0 {
		log.Printf("[ApproveLoan][LoanID: %d][EmployeeID: %d] field validator employee data is not found", loanID, approvalInfo.FieldValidatorEmployeeID)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[ApproveLoan][LoanID: %d][EmployeeID: %d] field validator employee data is not found", loanID, approvalInfo.FieldValidatorEmployeeID))
		return
	}

	// 7. check user type
	if fieldValidatorEmployee.UserType != constant.UserTypeFieldValidatorEmployee {
		log.Printf("[ApproveLoan][LoanID: %d][EmployeeID: %d] user type is not field validator employee", loanID, approvalInfo.FieldValidatorEmployeeID)
		h.RenderResponse(w, r, "", http.StatusForbidden, fmt.Sprintf("[ApproveLoan][LoanID: %d][EmployeeID: %d] user type is not field validator employee", loanID, approvalInfo.FieldValidatorEmployeeID))
		return
	}

	// 8. update loan status to approve
	loan.Status = constant.LoanStatusApproved
	loan.StatusDesc = constant.GetLoanStatusDesc(loan.Status)
	loan.ApprovalInfo = &model.ApprovalInfo{
		PictureProof:             approvalInfo.PictureProof,
		FieldValidatorEmployeeID: approvalInfo.FieldValidatorEmployeeID,
		ApprovalDate:             approvalInfo.ApprovalDate,
	}
	approvalDate := approvalInfo.ApprovalDate
	if approvalInfo.ApprovalDate.IsZero() {
		approvalDate = time.Now()
	}
	loan.ApprovalInfo.ApprovalDate = approvalDate
	h.Helper.UpsertLoan(loan)

	// 9. render response
	h.RenderResponse(w, r, loan, http.StatusOK, "")
}

// InvestLoan is handler to invest loan
func (h *Handler) InvestLoan(w http.ResponseWriter, r *http.Request) {
	// 1. get vars
	vars := mux.Vars(r)
	loanID, err := strconv.ParseInt(vars["loan_id"], 10, 64)
	if err != nil {
		log.Printf("[InvestLoan] failed parse int, with error: %+v", err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[InvestLoan] failed parse int, with error: %+v", err))
		return
	}

	// 2. decode body
	var lending model.Lending
	err = json.NewDecoder(r.Body).Decode(&lending)
	if err != nil {
		log.Printf("[InvestLoan][LoanID: %d] fail decode body with error: %+v", loanID, err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[InvestLoan][LoanID: %d] fail decode body with error: %+v", loanID, err))
		return
	}

	// 3. sanitize payload
	if lending.LenderID == 0 {
		log.Printf("[InvestLoan][LoanID: %d] lender id is empty", loanID)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[InvestLoan][LoanID: %d] lender id is empty", loanID))
		return
	}
	if lending.InvestedAmount == 0 {
		log.Printf("[InvestLoan][LoanID: %d][LenderID: %d] invested amount is empty", loanID, lending.LenderID)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[InvestLoan][LoanID: %d][LenderID: %d] invested amount is empty", loanID, lending.LenderID))
		return
	}

	// 4. get loan by loan id
	loan := h.Helper.GetLoanByLoanID(loanID)
	if loan.LoanID == 0 {
		log.Printf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] loan data not found", loanID, lending.LenderID, lending.InvestedAmount)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] loan data not found", loanID, lending.LenderID, lending.InvestedAmount))
		return
	}

	// 5. check loan status
	if loan.Status != constant.LoanStatusApproved {
		log.Printf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] loan status invalid, current status is: %s", loanID, lending.LenderID, lending.InvestedAmount, constant.GetLoanStatusDesc(loan.Status))
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] loan status invalid, current status is: %s", loanID, lending.LenderID, lending.InvestedAmount, constant.GetLoanStatusDesc(loan.Status)))
		return
	}

	// 6. get lender by user id
	lender := h.Helper.GetUserByUserID(lending.LenderID)
	if lender.UserID == 0 {
		log.Printf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] lender data is not found", loanID, lending.LenderID, lending.InvestedAmount)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] lender data is not found", loanID, lending.LenderID, lending.InvestedAmount))
		return
	}

	// 7. check user type
	if lender.UserType != constant.UserTypeLender {
		log.Printf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] user type is not lender", loanID, lending.LenderID, lending.InvestedAmount)
		h.RenderResponse(w, r, "", http.StatusForbidden, fmt.Sprintf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] user type is not lender", loanID, lending.LenderID, lending.InvestedAmount))
		return
	}

	// 8. check invested amount
	if lending.InvestedAmount > loan.GetRemainingRequiredAmount() {
		log.Printf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] invested amount is bigger than remaining required amount: %.2f", loanID, lending.LenderID, lending.InvestedAmount, loan.GetRemainingRequiredAmount())
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] invested amount is bigger than remaining required amount: %.2f", loanID, lending.LenderID, lending.InvestedAmount, loan.GetRemainingRequiredAmount()))
		return
	}

	// 9. update loan lending data
	if loan.IsLenderInvested(lender.UserID) {
		for i := 0; i < len(loan.Lending); i++ {
			if loan.Lending[i].LenderID == lender.UserID {
				loan.Lending[i].InvestedAmount += lending.InvestedAmount
				loan.Lending[i].ReturnAmount = loan.Lending[i].CalculateLenderReturnAmount(loan.InterestRate)
			}
		}
	} else {
		lendingdata := model.Lending{
			LenderID:       lending.LenderID,
			InvestedAmount: lending.InvestedAmount,
		}
		lendingdata.ReturnAmount = lendingdata.CalculateLenderReturnAmount(loan.InterestRate)
		loan.Lending = append(loan.Lending, lendingdata)
	}
	loan.CollectedAmount += lending.InvestedAmount

	// 10. check principal amount is fulfilled?
	if loan.IsAmountFulfilled() {
		// 10a. update loan status to invested
		loan.Status = constant.LoanStatusInvested
		loan.StatusDesc = constant.GetLoanStatusDesc(loan.Status)

		// 10b. generate lender agreement pdf
		err = h.Helper.GenerateLenderAgreementPDF(&loan)
		if err != nil {
			log.Printf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] fail to generate lender agreement pdf with error: %+v", loanID, lending.LenderID, lending.InvestedAmount, err)
			h.RenderResponse(w, r, "", http.StatusInternalServerError, fmt.Sprintf("[InvestLoan][LoanID: %d][LenderID: %d][Amount: %.2f] fail to generate lender agreement pdf with error: %+v", loanID, lending.LenderID, lending.InvestedAmount, err))
			return
		}
	}

	// 11. update loan lending
	h.Helper.UpsertLoan(loan)

	// 12. render response
	h.RenderResponse(w, r, loan, http.StatusOK, "")
}

// DisburseLoan is handler to disburse loan
func (h *Handler) DisburseLoan(w http.ResponseWriter, r *http.Request) {
	// 1. get vars
	vars := mux.Vars(r)
	loanID, err := strconv.ParseInt(vars["loan_id"], 10, 64)
	if err != nil {
		log.Printf("[DisburseLoan] failed parse int, with error: %+v", err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[DisburseLoan] failed parse int, with error: %+v", err))
		return
	}

	// 2. decode body
	var disbursement model.Disbursement
	err = json.NewDecoder(r.Body).Decode(&disbursement)
	if err != nil {
		log.Printf("[DisburseLoan][LoanID: %d] fail decode body with error: %+v", loanID, err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[DisburseLoan][LoanID: %d] fail decode body with error: %+v", loanID, err))
		return
	}

	// 3. sanitize payload
	if disbursement.FieldOfficerID == 0 {
		log.Printf("[DisburseLoan][LoanID: %d] field officer id is empty", loanID)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[DisburseLoan][LoanID: %d] field officer id is empty", loanID))
		return
	}
	if disbursement.DisbursementDate.IsZero() {
		log.Printf("[DisburseLoan][LoanID: %d][OfficerID: %d] invalid disbursement date", loanID, disbursement.FieldOfficerID)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[DisburseLoan][LoanID: %d][OfficerID: %d] invalid disbursement date", loanID, disbursement.FieldOfficerID))
		return
	}

	// 4. get loan by loan id
	loan := h.Helper.GetLoanByLoanID(loanID)
	if loan.LoanID == 0 {
		log.Printf("[DisburseLoan][LoanID: %d][OfficerID: %d] loan data not found", loanID, disbursement.FieldOfficerID)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[DisburseLoan][LoanID: %d][OfficerID: %d] loan data not found", loanID, disbursement.FieldOfficerID))
		return
	}

	// 5. check loan status
	if loan.Status != constant.LoanStatusSigned {
		log.Printf("[DisburseLoan][LoanID: %d][OfficerID: %d] loan status invalid, current status is: %s", loanID, disbursement.FieldOfficerID, constant.GetLoanStatusDesc(loan.Status))
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[DisburseLoan][LoanID: %d][OfficerID: %d] loan status invalid, current status is: %s", loanID, disbursement.FieldOfficerID, constant.GetLoanStatusDesc(loan.Status)))
		return
	}

	// 6. get field officer employee by user id
	fieldOfficerEmployee := h.Helper.GetUserByUserID(disbursement.FieldOfficerID)
	if fieldOfficerEmployee.UserID == 0 {
		log.Printf("[DisburseLoan][LoanID: %d][OfficerID: %d] field officer employee data is not found", loanID, disbursement.FieldOfficerID)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[DisburseLoan][LoanID: %d][OfficerID: %d] field officer employee data is not found", loanID, disbursement.FieldOfficerID))
		return
	}

	// 7. check user type
	if fieldOfficerEmployee.UserType != constant.UserTypeFieldOfficerEmployee {
		log.Printf("[DisburseLoan][LoanID: %d][OfficerID: %d] user type is not field officer employee", loanID, disbursement.FieldOfficerID)
		h.RenderResponse(w, r, "", http.StatusForbidden, fmt.Sprintf("[DisburseLoan][LoanID: %d][OfficerID: %d] user type is not field officer employee", loanID, disbursement.FieldOfficerID))
		return
	}

	// 8. update loan disbursement and status
	loan.Status = constant.LoanStatusDisbursed
	loan.StatusDesc = constant.GetLoanStatusDesc(loan.Status)
	loan.DisbursementInfo.FieldOfficerID = disbursement.FieldOfficerID
	disbursementDate := disbursement.DisbursementDate
	if disbursement.DisbursementDate.IsZero() {
		disbursementDate = time.Now()
	}
	loan.DisbursementInfo.DisbursementDate = disbursementDate

	// 9. upsert loan
	h.Helper.UpsertLoan(loan)

	// 10. render response
	h.RenderResponse(w, r, loan, http.StatusOK, "")
}
