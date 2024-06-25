package model

import (
	"time"
)

type Loan struct {
	LoanID                     int64             `json:"loan_id"`
	BorrowerID                 int64             `json:"borrower_id"`
	PrincipalAmount            float64           `json:"principal_amount"`
	InterestRate               float64           `json:"interest_rate"`
	Status                     int               `json:"status"`
	StatusDesc                 string            `json:"status_desc"`
	OrganizerBorrowerAggrement string            `json:"organizer_borrower_aggrement,omitempty"`
	ApprovalInfo               *ApprovalInfo     `json:"approval_info,omitempty"`
	Lending                    []Lending         `json:"lending,omitempty"`
	DisbursementInfo           *DisbursementInfo `json:"disbursement_info,omitempty"`
}

type ApprovalInfo struct {
	PictureProof             string    `json:"picture_proof"` // base64 encoded string (image)
	FieldValidatorEmployeeID int64     `json:"field_validator_employee_id"`
	ApprovalDate             time.Time `json:"approval_date"`
}

type Lending struct {
	LenderID                 int64   `json:"lender_id"`
	InvestedAmount           float64 `json:"invested_amount"`
	OrganizerLenderAggrement string  `json:"organizer_lender_aggrement"`
	ReturnOfInvestment       float64 `json:"return_of_investment"`
}

type DisbursementInfo struct {
	DisbursementInfoID int64  `json:"disbursement_info_id"`
	FieldOfficerID     string `json:"field_officer_id"`
	DisbursementDate   string `json:"disbursement_date"`
}
