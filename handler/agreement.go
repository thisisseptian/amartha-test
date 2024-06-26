package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"amartha-test/constant"
	"amartha-test/helper"
	"amartha-test/model"
)

// ListAgreement ...
func ListAgreement(w http.ResponseWriter, r *http.Request) {
	// 1. check http method
	if r.Method != http.MethodGet {
		log.Println("[ListAgreement] invalid request method")
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 2. get agreement list
	agreements := helper.GetAgreements()
	if len(agreements) == 0 {
		log.Println("[ListAgreement] agreement list is empty")
		http.Error(w, "agreement list is empty", http.StatusNotFound)
		return
	}

	// 3. return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(agreements)
}

// DetailAgreement ...
func DetailAgreement(w http.ResponseWriter, r *http.Request, agreementID int64) {
	// 1. check http method
	if r.Method != http.MethodGet {
		log.Println("[DetailAgreement] invalid request method")
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 2. get agreement by agreement id
	agreement := helper.GetAgreementByAgreementID(agreementID)
	if agreement.AggrementID == 0 {
		log.Println("[DetailAgreement] agreement data is not found")
		http.Error(w, "agreement data is not found", http.StatusNotFound)
		return
	}

	// 3. decode agreement base 64
	pdfData, err := base64.StdEncoding.DecodeString(agreement.DocumentData)
	if err != nil {
		log.Printf("[DetailAgreement] failed to decode base64 pdf data with error: %+v", err)
		http.Error(w, fmt.Sprintf("[DetailAgreement] failed to decode base64 pdf data with error: %+v", err), http.StatusInternalServerError)
		return
	}

	// 4. return response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=agreement.pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfData)
}

// SignAgreement ...
func SignAgreement(w http.ResponseWriter, r *http.Request, agreementID int64) {
	// 1. check http method
	if r.Method != http.MethodPost {
		log.Println("[SignAgreement] invalid request method")
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 2. decode body
	var disbursement model.Disbursement
	err := json.NewDecoder(r.Body).Decode(&disbursement)
	if err != nil {
		log.Printf("[SignAgreement] fail decode body with error: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. sanitize payload
	if disbursement.LoanID == 0 {
		log.Print("[SignAgreement] loan id is empty")
		http.Error(w, "loan id is empty", http.StatusBadRequest)
		return
	}
	if disbursement.UserID == 0 {
		log.Print("[SignAgreement] user id is empty")
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return
	}

	// 4. get loan by loan id
	loan := helper.GetLoanByLoanID(disbursement.LoanID)
	if loan.LoanID == 0 {
		log.Println("[SignAgreement] loan data not found")
		http.Error(w, "loan data not found", http.StatusBadRequest)
		return
	}

	// 5. check loan status
	if loan.Status != constant.LoanStatusInvested {
		log.Printf("[SignAgreement] loan status invalid, current status is: %s", constant.GetLoanStatusDesc(loan.Status))
		http.Error(w, fmt.Sprintf("loan status invalid, current status is: %s", constant.GetLoanStatusDesc(loan.Status)), http.StatusBadRequest)
		return
	}

	// 6. get agreement by agreement id
	agreement := helper.GetAgreementByAgreementID(agreementID)
	if agreement.AggrementID == 0 {
		log.Println("[SignAgreement] agreement data not found")
		http.Error(w, "agreement data not found", http.StatusBadRequest)
		return
	}

	// 7. wrong user to sign
	if agreement.UserID != disbursement.UserID {
		log.Println("[SignAgreement] wrong user to sign this agreement")
		http.Error(w, "wrong user to sign this agreement", http.StatusBadRequest)
		return
	}

	// 8. check agreement sign
	if agreement.IsSigned {
		log.Println("[SignAgreement] agreement already signed")
		http.Error(w, "agreement already signed", http.StatusBadRequest)
		return
	}

	// 9. update agreement sign
	agreement.IsSigned = true
	helper.UpsertAgreement(agreement)

	// 10. generate agreement sign pdf
	err = helper.GenerateSignedAgreementPDF(&loan, disbursement.UserID)
	if err != nil {
		log.Printf("[SignAgreement] fail to generate signed agreement pdf with error: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 11. get user by user id
	user := helper.GetUserByUserID(disbursement.UserID)
	if user.UserID == 0 {
		log.Println("[SignAgreement] user data not found")
		http.Error(w, "user data not found", http.StatusBadRequest)
		return
	}

	// 12. check based on user type
	if user.UserType == constant.UserTypeLender {
		// 12a. check agreement is completely signed by all lender
		isCompletelySignedByLender, err := helper.CheckAgreementCompletelySignedByLender(loan)
		if err != nil {
			log.Printf("[SignAgreement] check agreement completely signed by lender got fail with error: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if isCompletelySignedByLender {
			// 12b. generate borrower agreement pdf
			err = helper.GenerateBorrowerAgreementPDF(&loan)
			if err != nil {
				log.Printf("[SignAgreement] fail to generate borrower agreement pdf with error: %+v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else if user.UserType == constant.UserTypeBorrower { // final sign must be from borrower
		// 12c. update loan disbursement and status
		loan.Status = constant.LoanStatusDisbursed
		loan.StatusDesc = constant.GetLoanStatusDesc(loan.Status)
		fieldOfficerID := disbursement.FieldOfficerID
		if disbursement.FieldOfficerID != 0 {
			// 12d. get field officer employee by user id
			fieldOfficerEmployee := helper.GetUserByUserID(disbursement.FieldOfficerID)
			if fieldOfficerEmployee.UserID == 0 {
				if agreement.AggrementID == 0 {
					// just log the error and continue with field officer employee id = 0
					log.Println("[SignAgreement] field officer employee data not found")
					fieldOfficerID = 0
				}
			}

			// 12e. check user type
			if fieldOfficerEmployee.UserType != constant.UserTypeFieldOfficerEmployee {
				// just log the error and continue with field officer employee id = 0
				log.Println("[SignAgreement] user type is not field officer employee")
				fieldOfficerID = 0
			}
		}
		loan.DisbursementInfo.FieldOfficerID = fieldOfficerID
		disbursementDate := disbursement.DisbursementDate
		if disbursement.DisbursementDate.IsZero() {
			disbursementDate = time.Now()
		}
		loan.DisbursementInfo.DisbursementDate = disbursementDate
		helper.UpsertLoan(loan)
	}

	// 12. return response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loan)
}
