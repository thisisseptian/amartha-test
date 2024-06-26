package helper

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/jung-kurt/gofpdf/v2"

	"amartha-test/constant"
	"amartha-test/model"
)

var (
	agreementIDCounter int64
	mutexAgreement     sync.Mutex

	agreements = make(map[int64]*model.Aggrement)
)

func GenerateIncrementalAgreementID() int64 {
	mutexAgreement.Lock()
	defer mutexAgreement.Unlock()
	agreementIDCounter++
	return agreementIDCounter
}

func UpsertAgreement(agreement model.Aggrement) {
	agreements[agreement.AggrementID] = &agreement
}

func GetAgreements() []model.Aggrement {
	var listAgreement []model.Aggrement
	for _, v := range agreements {
		listAgreement = append(listAgreement, *v)
	}

	return listAgreement
}

func GetAgreementByAgreementID(agreementID int64) model.Aggrement {
	agreement, exists := agreements[agreementID]
	if exists {
		return *agreement
	}

	return model.Aggrement{}
}

func GenerateBorrowerAgreementPDF(loan *model.Loan) error {
	borrower := GetUserByUserID(loan.BorrowerID)
	if borrower.UserID == 0 {
		log.Println("[GenerateAgreementPDF] borrower is not found")
		return errors.New("borrower is not found")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("ORGANIZER-BORROWER AGREEMENT [Loan ID: %d]", loan.LoanID))
	pdf.Ln(20)
	pdf.Cell(40, 10, fmt.Sprintf("Borrower ID: %d", borrower.UserID))
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("Borrower Name: %s", borrower.UserName))
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("Principal Amount: Rp %.2f", loan.PrincipalAmount))
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("Interest Rate: %.2f%%", loan.InterestRate*100))
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("Amount of debt: Rp %.2f", loan.CalculateReturnAmount()))
	pdf.Ln(20)
	pdf.Cell(40, 10, "Sign: UNSIGNED")
	pdf.Ln(5)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("[GenerateAgreementPDF] failed generate organizer-borrower agreement with error: %+v", err)
		return err
	}

	organizerBorrowerAgreement := model.Aggrement{
		AggrementID:  GenerateIncrementalAgreementID(),
		DocumentData: base64.StdEncoding.EncodeToString(buf.Bytes()),
		UserID:       borrower.UserID,
	}
	UpsertAgreement(organizerBorrowerAgreement)

	loan.OrganizerBorrowerAggrementURL = fmt.Sprintf(constant.AgreementPrefix, organizerBorrowerAgreement.AggrementID)
	UpsertLoan(*loan)

	return nil
}

func GenerateLenderAgreementPDF(loan *model.Loan) error {
	borrower := GetUserByUserID(loan.BorrowerID)
	if borrower.UserID == 0 {
		log.Println("[GenerateAgreementPDF] borrower is not found")
		return errors.New("borrower is not found")
	}

	for i := 0; i < len(loan.Lending); i++ {
		lender := GetUserByUserID(loan.Lending[i].LenderID)
		if lender.UserID == 0 {
			log.Println("[GenerateAgreementPDF] lender is not found")
			return errors.New("lender is not found")
		}

		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, fmt.Sprintf("ORGANIZER-LENDER AGREEMENT [Loan ID: %d]", loan.LoanID))
		pdf.Ln(20)
		pdf.Cell(40, 10, fmt.Sprintf("Borrower ID: %d", loan.BorrowerID))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Borrower Name: %s", borrower.UserName))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Principal Amount: Rp %.2f", loan.PrincipalAmount))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Interest Rate: %.2f%%", loan.InterestRate*100))
		pdf.Ln(10)
		pdf.Cell(40, 10, fmt.Sprintf("Lender ID: %d", lender.UserID))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Lender Name: %s", lender.UserName))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Invested Amount: Rp %.2f", loan.Lending[i].InvestedAmount))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Return Amountt: Rp %.2f", loan.Lending[i].ReturnAmount))
		pdf.Ln(20)
		pdf.Cell(40, 10, "Sign: UNSIGNED")
		pdf.Ln(5)

		var buf bytes.Buffer
		err := pdf.Output(&buf)
		if err != nil {
			log.Printf("[GenerateAgreementPDF] failed generate organizer-lender agreement with error: %+v", err)
			return err
		}

		organizerLenderAgreement := model.Aggrement{
			AggrementID:  GenerateIncrementalAgreementID(),
			DocumentData: base64.StdEncoding.EncodeToString(buf.Bytes()),
			UserID:       lender.UserID,
		}
		UpsertAgreement(organizerLenderAgreement)

		loan.Lending[i].OrganizerLenderAggrementURL = fmt.Sprintf(constant.AgreementPrefix, organizerLenderAgreement.AggrementID)
	}

	UpsertLoan(*loan)

	return nil
}

func GenerateSignedAgreementPDF(loan *model.Loan, userID int64) error {
	tmpURL := ""

	borrower := GetUserByUserID(loan.BorrowerID)
	if borrower.UserID == 0 {
		log.Println("[GenerateSignedAgreementPDF] borrower is not found")
		return errors.New("borrower is not found")
	}

	if userID == loan.BorrowerID {
		// 1. organizer borrower agreement signed
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, fmt.Sprintf("ORGANIZER-BORROWER AGREEMENT [Loan ID: %d]", loan.LoanID))
		pdf.Ln(20)
		pdf.Cell(40, 10, fmt.Sprintf("Borrower ID: %d", borrower.UserID))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Borrower Name: %s", borrower.UserName))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Principal Amount: Rp %.2f", loan.PrincipalAmount))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Interest Rate: %.2f%%", loan.InterestRate*100))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Amount of debt: Rp %.2f", loan.CalculateReturnAmount()))
		pdf.Ln(20)
		pdf.Cell(40, 10, "Sign: SIGNED")
		pdf.Ln(5)

		var buf bytes.Buffer
		err := pdf.Output(&buf)
		if err != nil {
			log.Printf("[GenerateSignedAgreementPDF] failed generate organizer-borrower agreement signed with error: %+v", err)
			return err
		}

		organizerBorrowerAgreement := model.Aggrement{
			AggrementID:  GenerateIncrementalAgreementID(),
			DocumentData: base64.StdEncoding.EncodeToString(buf.Bytes()),
			UserID:       borrower.UserID,
			IsSigned:     true,
		}
		UpsertAgreement(organizerBorrowerAgreement)

		tmpURL = fmt.Sprintf(constant.AgreementPrefix, organizerBorrowerAgreement.AggrementID)
	} else {
		// 2. organizer lender agreement
		for i := 0; i < len(loan.Lending); i++ {
			if userID == loan.Lending[i].LenderID {
				lender := GetUserByUserID(loan.Lending[i].LenderID)
				if lender.UserID == 0 {
					log.Println("[GenerateSignedAgreementPDF] lender is not found")
					return errors.New("lender is not found")
				}

				pdf := gofpdf.New("P", "mm", "A4", "")
				pdf.AddPage()
				pdf.SetFont("Arial", "B", 16)
				pdf.Cell(40, 10, fmt.Sprintf("ORGANIZER-LENDER AGREEMENT [Loan ID: %d]", loan.LoanID))
				pdf.Ln(20)
				pdf.Cell(40, 10, fmt.Sprintf("Borrower ID: %d", loan.BorrowerID))
				pdf.Ln(5)
				pdf.Cell(40, 10, fmt.Sprintf("Borrower Name: %s", borrower.UserName))
				pdf.Ln(5)
				pdf.Cell(40, 10, fmt.Sprintf("Principal Amount: Rp %.2f", loan.PrincipalAmount))
				pdf.Ln(5)
				pdf.Cell(40, 10, fmt.Sprintf("Interest Rate: %.2f%%", loan.InterestRate*100))
				pdf.Ln(10)
				pdf.Cell(40, 10, fmt.Sprintf("Lender ID: %d", lender.UserID))
				pdf.Ln(5)
				pdf.Cell(40, 10, fmt.Sprintf("Lender Name: %s", lender.UserName))
				pdf.Ln(5)
				pdf.Cell(40, 10, fmt.Sprintf("Invested Amount: Rp %.2f", loan.Lending[i].InvestedAmount))
				pdf.Ln(5)
				pdf.Cell(40, 10, fmt.Sprintf("Return Amountt: Rp %.2f", loan.Lending[i].ReturnAmount))
				pdf.Ln(20)
				pdf.Cell(40, 10, "Sign: SIGNED")
				pdf.Ln(5)

				var buf bytes.Buffer
				err := pdf.Output(&buf)
				if err != nil {
					log.Printf("[GenerateSignedAgreementPDF] failed generate organizer-lender agreement signed with error: %+v", err)
					return err
				}

				organizerLenderAgreement := model.Aggrement{
					AggrementID:  GenerateIncrementalAgreementID(),
					DocumentData: base64.StdEncoding.EncodeToString(buf.Bytes()),
					UserID:       lender.UserID,
					IsSigned:     true,
				}
				UpsertAgreement(organizerLenderAgreement)

				tmpURL = fmt.Sprintf(constant.AgreementPrefix, organizerLenderAgreement.AggrementID)
			}
		}
	}

	if tmpURL != "" {
		loan.DisbursementInfo.AgreementSignedURLs = append(loan.DisbursementInfo.AgreementSignedURLs, tmpURL)
		UpsertLoan(*loan)
	}

	return nil
}

func CheckAgreementCompletelySignedByLender(loan model.Loan) (bool, error) {
	for _, v := range loan.Lending {
		organizerLenderAgreementID, err := GetAgreementIDByAgreementURL(v.OrganizerLenderAggrementURL)
		if err != nil {
			log.Printf("[CheckLoanCompletelySigned] failed get organizer-lender agreement id with error: %+v", err)
			return false, err
		}

		organizerLenderAgreement := GetAgreementByAgreementID(organizerLenderAgreementID)
		if !organizerLenderAgreement.IsSigned {
			log.Printf("[CheckLoanCompletelySigned] agremeent with agreement id: %d in loan id: %d is still unsigned", organizerLenderAgreement.AggrementID, loan.LoanID)
			return false, nil
		}
	}

	return true, nil
}

func GetAgreementIDByAgreementURL(url string) (int64, error) {
	parts := strings.Split(url, "/")
	agreementIDString := parts[len(parts)-2]

	return strconv.ParseInt(agreementIDString, 10, 64)
}
