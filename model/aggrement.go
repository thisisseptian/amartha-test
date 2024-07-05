package model

type Aggrement struct {
	AggrementID  int64  `json:"aggrement_id"`
	DocumentData string `json:"document_data"`
	UserID       int64  `json:"user_id"`
	IsSigned     bool   `json:"is_signed"`
}

type Sign struct {
	LoanID int64 `json:"loan_id"`
	UserID int64 `json:"user_id"`
}
