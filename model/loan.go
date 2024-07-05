package model

import (
	"time"
)

type Loan struct {
	LoanID                        int64            `json:"loan_id"`
	TrxID                         int64            `json:"trx_id"`
	BorrowerID                    int64            `json:"borrower_id"`
	PrincipalAmount               float64          `json:"principal_amount"`
	CollectedAmount               float64          `json:"collected_amount"`
	InterestRate                  float64          `json:"interest_rate"`
	Status                        int              `json:"status"`
	StatusDesc                    string           `json:"status_desc"`
	OrganizerBorrowerAggrementURL string           `json:"organizer_borrower_aggrement_url,omitempty"`
	ApprovalInfo                  *ApprovalInfo    `json:"approval_info,omitempty"`
	Lending                       []Lending        `json:"lending,omitempty"`
	DisbursementInfo              DisbursementInfo `json:"disbursement_info,omitempty"`
}

type ApprovalInfo struct {
	PictureProof             string    `json:"picture_proof"` // base64 encoded string (image)
	FieldValidatorEmployeeID int64     `json:"field_validator_employee_id"`
	ApprovalDate             time.Time `json:"approval_date"`
}

type Lending struct {
	LenderID                    int64   `json:"lender_id"`
	InvestedAmount              float64 `json:"invested_amount"`
	OrganizerLenderAggrementURL string  `json:"organizer_lender_aggrement_url"`
	ReturnAmount                float64 `json:"return_amount"`
}

type DisbursementInfo struct {
	AgreementSignedURLs []string  `json:"agreement_signed_urls"`
	FieldOfficerID      int64     `json:"field_officer_id"`
	DisbursementDate    time.Time `json:"disbursement_date"`
}

func (l *Loan) GetRemainingRequiredAmount() float64 {
	return l.PrincipalAmount - l.CollectedAmount
}

func (l *Loan) IsAmountFulfilled() bool {
	return l.PrincipalAmount == l.CollectedAmount
}

func (l *Lending) CalculateLenderReturnAmount(interestRate float64) float64 {
	return l.InvestedAmount + (l.InvestedAmount * interestRate)
}

func (l *Loan) CalculateReturnAmount() float64 {
	return l.PrincipalAmount + (l.PrincipalAmount * l.InterestRate)
}

func (l *Loan) IsLenderInvested(lenderID int64) bool {
	for _, v := range l.Lending {
		if v.LenderID == lenderID {
			return true
		}
	}

	return false
}
