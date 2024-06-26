package constant

const AgreementPrefix = "http://localhost:8080/agreement/%d/view"

const (
	UserTypeBorrower               = 1
	UserTypeLender                 = 2
	UserTypeFieldValidatorEmployee = 3
	UserTypeFieldOfficerEmployee   = 4
)

const (
	LoanStatusProposed  = 1
	LoanStatusApproved  = 2
	LoanStatusInvested  = 3
	LoanStatusDisbursed = 4
)

var LoanStatusDesc = map[int]string{
	LoanStatusProposed:  "proposed",
	LoanStatusApproved:  "approved",
	LoanStatusInvested:  "invested",
	LoanStatusDisbursed: "disbursed",
}

func GetLoanStatusDesc(status int) string {
	desc, ok := LoanStatusDesc[status]
	if ok {
		return desc
	}

	return ""
}
