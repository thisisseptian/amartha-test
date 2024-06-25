package constant

const (
	UserTypeBorrower = 1
	UserTypeLender   = 2
	UserTypeEmployee = 3
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
