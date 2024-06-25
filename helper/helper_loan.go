package helper

import (
	"amartha-test/model"
	"sync"
)

var (
	loanIDCounter int64
	mutexLoan     sync.Mutex

	loans = make(map[int64]*model.Loan)
)

func GenerateIncrementalLoanID() int64 {
	mutexLoan.Lock()
	defer mutexLoan.Unlock()
	loanIDCounter++
	return loanIDCounter
}

func UpsertLoan(loan model.Loan) {
	loans[loan.LoanID] = &loan
}

func GetLoans() []model.Loan {
	var listLoan []model.Loan
	for _, v := range loans {
		listLoan = append(listLoan, *v)
	}

	return listLoan
}

func GetLoanByLoanID(loanID int64) model.Loan {
	loan, exists := loans[loanID]
	if exists {
		return *loan
	}

	return model.Loan{}
}
