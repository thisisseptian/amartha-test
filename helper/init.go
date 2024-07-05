package helper

import "amartha-test/model"

type IHelper interface {
	// helper user
	GenerateIncrementalUserID() int64
	InitUsers()
	GetUsers() []model.User
	GetUserByUserID(userID int64) model.User

	// helper loan
	GenerateIncrementalLoanID() int64
	UpsertLoan(loan model.Loan)
	GetLoans() []model.Loan
	GetLoanByLoanID(loanID int64) model.Loan

	// helper agreement
	GenerateIncrementalAgreementID() int64
	UpsertAgreement(agreement model.Aggrement)
	GetAgreements() []model.Aggrement
	GetAgreementByAgreementID(agreementID int64) model.Aggrement
	GenerateBorrowerAgreementPDF(loan *model.Loan) error
	GenerateLenderAgreementPDF(loan *model.Loan) error
	GenerateSignedAgreementPDF(loan *model.Loan, userID int64) error
	CheckAgreementCompletelySignedByLender(loan model.Loan) (bool, error)
	GetAgreementIDByAgreementURL(url string) (int64, error)
}

type Helper struct {
	IHelper
}

func NewHelper() *Helper {
	return &Helper{}
}
