package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"amartha-test/constant"
	"amartha-test/helper"
	"amartha-test/model"
)

// SubmitLoan ...
func SubmitLoan(w http.ResponseWriter, r *http.Request) {
	// 1. check http method
	if r.Method != http.MethodPost {
		log.Println("[SubmitLoan] invalid request method")
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 2. decode body
	var loan model.Loan
	err := json.NewDecoder(r.Body).Decode(&loan)
	if err != nil {
		log.Printf("[SubmitLoan] fail decode body with error: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. sanitize payload
	if loan.BorrowerID == 0 {
		log.Print("[SubmitLoan] borrower id is empty")
		http.Error(w, "borrower id is empty", http.StatusBadRequest)
		return
	}
	if loan.PrincipalAmount == 0 {
		log.Print("[SubmitLoan] principal amount is empty")
		http.Error(w, "principal amount is empty", http.StatusBadRequest)
		return
	}
	if loan.InterestRate < 0 || loan.InterestRate > 1 {
		log.Print("[SubmitLoan] interest rate is invalid")
		http.Error(w, "interest rate is invalid", http.StatusBadRequest)
		return
	}

	// 4. get borrower by user id
	borrower := helper.GetUserByUserID(loan.BorrowerID)
	if borrower.UserID == 0 {
		log.Println("[SubmitLoan] borrower data is not found")
		http.Error(w, "borrower data is not found", http.StatusBadRequest)
		return
	}

	// 5. check user status
	if borrower.UserType != constant.UserTypeBorrower {
		log.Println("[SubmitLoan] user type is not borrower")
		http.Error(w, "user type is not borrower", http.StatusBadRequest)
		return
	}

	// 6. create loan
	loan.LoanID = helper.GenerateIncrementalLoanID()
	loan.Status = constant.LoanStatusProposed
	loan.StatusDesc = constant.GetLoanStatusDesc(loan.Status)
	helper.UpsertLoan(loan)

	// 7. return response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loan)
}

// ListLoan ...
func ListLoan(w http.ResponseWriter, r *http.Request) {
	// 1. check http method
	if r.Method != http.MethodGet {
		log.Println("[ListLoan] invalid request method")
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 2. get loan list
	loans := helper.GetLoans()
	if len(loans) == 0 {
		log.Println("[ListLoan] loan list is empty")
		http.Error(w, "loan list is empty", http.StatusNotFound)
		return
	}

	// 3. return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loans)
}

// DetailLoan ...
func DetailLoan(w http.ResponseWriter, r *http.Request, loanID int64) {
	// 1. check http method
	if r.Method != http.MethodGet {
		log.Println("[DetailLoan] invalid request method")
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 2. get loan by loan id
	loan := helper.GetLoanByLoanID(loanID)
	if loan.LoanID == 0 {
		log.Println("[DetailLoan] loan data is not found")
		http.Error(w, "loan data is not found", http.StatusNotFound)
		return
	}

	// 3. return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loan)
}

// ApproveLoan ...
func ApproveLoan(w http.ResponseWriter, r *http.Request, loanID int64) {
	// 1. check http method
	if r.Method != http.MethodPost {
		log.Println("[ApproveLoan] invalid request method")
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 2. decode body
	var approvalInfo model.ApprovalInfo
	err := json.NewDecoder(r.Body).Decode(&approvalInfo)
	if err != nil {
		log.Printf("[ApproveLoan] fail decode body with error: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. sanitize payload
	if approvalInfo.PictureProof == "" {
		log.Print("[ApproveLoan] picture proof is empty")
		http.Error(w, "picture proof is empty", http.StatusBadRequest)
		return
	}
	if approvalInfo.FieldValidatorEmployeeID == 0 {
		log.Print("[ApproveLoan] field validator employee id is empty")
		http.Error(w, "field validator employee id is empty", http.StatusBadRequest)
		return
	}

	// 4. get loan by loan id
	loan := helper.GetLoanByLoanID(loanID)
	if loan.LoanID == 0 {
		log.Println("[ApproveLoan] loan data is not found")
		http.Error(w, "loan data is not found", http.StatusBadRequest)
		return
	}

	// 5. check loan status
	if loan.Status != constant.LoanStatusProposed {
		log.Printf("[ApproveLoan] loan status is invalid, current status is: %s", constant.GetLoanStatusDesc(loan.Status))
		http.Error(w, fmt.Sprintf("[ApproveLoan] loan status is invalid, current status is: %s", constant.GetLoanStatusDesc(loan.Status)), http.StatusBadRequest)
		return
	}

	// 6. get field validator employee by user id
	fieldValidatorEmployee := helper.GetUserByUserID(approvalInfo.FieldValidatorEmployeeID)
	if fieldValidatorEmployee.UserID == 0 {
		log.Println("[ApproveLoan] field validator employee data is not found")
		http.Error(w, "field validator employee data is not found", http.StatusBadRequest)
		return
	}

	// 7. check user type
	if fieldValidatorEmployee.UserType != constant.UserTypeFieldValidatorEmployee {
		log.Println("[ApproveLoan] user type is not field validator employee")
		http.Error(w, "user type is not field validator employee", http.StatusBadRequest)
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
	helper.UpsertLoan(loan)

	// 9. return response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loan)
}

// InvestLoan ...
func InvestLoan(w http.ResponseWriter, r *http.Request, loanID int64) {
	// 1. check http method
	if r.Method != http.MethodPost {
		log.Println("[InvestLoan] invalid request method")
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 2. decode body
	var lending model.Lending
	err := json.NewDecoder(r.Body).Decode(&lending)
	if err != nil {
		log.Printf("[InvestLoan] fail decode body with error: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. sanitize payload
	if lending.LenderID == 0 {
		log.Print("[InvestLoan] lender id is empty")
		http.Error(w, "lender id is empty", http.StatusBadRequest)
		return
	}
	if lending.InvestedAmount == 0 {
		log.Print("[InvestLoan] invested amount is empty")
		http.Error(w, "invested amount is empty", http.StatusBadRequest)
		return
	}

	// 4. get loan by loan id
	loan := helper.GetLoanByLoanID(loanID)
	if loan.LoanID == 0 {
		log.Println("[InvestLoan] loan data not found")
		http.Error(w, "loan data not found", http.StatusBadRequest)
		return
	}

	// 5. check loan status
	if loan.Status != constant.LoanStatusApproved {
		log.Printf("[InvestLoan] loan status invalid, current status is: %s", constant.GetLoanStatusDesc(loan.Status))
		http.Error(w, fmt.Sprintf("loan status invalid, current status is: %s", constant.GetLoanStatusDesc(loan.Status)), http.StatusBadRequest)
		return
	}

	// 6. get lender by user id
	lender := helper.GetUserByUserID(lending.LenderID)
	if lender.UserID == 0 {
		log.Println("[InvestLoan] lender data is not found")
		http.Error(w, "lender data is not found", http.StatusBadRequest)
		return
	}

	// 7. check user type
	if lender.UserType != constant.UserTypeLender {
		log.Println("[InvestLoan] user type is not lender")
		http.Error(w, "user type is not lender", http.StatusBadRequest)
		return
	}

	// 8. check invested amount
	if lending.InvestedAmount > loan.GetRemainingRequiredAmount() {
		log.Printf("[InvestLoan] invested amount is bigger than remaining required amount: %.2f", loan.GetRemainingRequiredAmount())
		http.Error(w, fmt.Sprintf("invested amount is bigger than remaining required amount: %.2f", loan.GetRemainingRequiredAmount()), http.StatusBadRequest)
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
		err = helper.GenerateLenderAgreementPDF(&loan)
		if err != nil {
			log.Printf("[InvestLoan] fail to generate lender agreement pdf with error: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 11. update loan lending
	helper.UpsertLoan(loan)

	// 12. return response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loan)
}
