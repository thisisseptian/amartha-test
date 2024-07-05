package model

import "time"

type Disbursement struct {
	FieldOfficerID   int64     `json:"field_officer_id"`
	DisbursementDate time.Time `json:"disbursement_date"`
}
