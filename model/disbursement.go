package model

import "time"

type Disbursement struct {
	LoanID           int64     `json:"loan_id"`
	UserID           int64     `json:"user_id"`
	FieldOfficerID   int64     `json:"field_officer_id"`
	DisbursementDate time.Time `json:"disbursement_date"`
}
