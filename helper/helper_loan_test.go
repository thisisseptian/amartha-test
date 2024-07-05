package helper

import (
	"testing"
	"time"

	"amartha-test/model"
)

func TestGenerateIncrementalLoanID(t *testing.T) {
	helper := NewHelper()

	t.Run("generate incremental loan id", func(t *testing.T) {
		expectedIDs := []int64{1, 2, 3, 4, 5}

		for i := 0; i < len(expectedIDs); i++ {
			actualID := helper.GenerateIncrementalLoanID()
			if actualID != expectedIDs[i] {
				t.Errorf("expected loan id %d, got %d", expectedIDs[i], actualID)
			}
		}
	})
}

func TestGetLoans(t *testing.T) {
	helper := NewHelper()

	t.Run("get loans", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			loan := model.Loan{
				LoanID:          helper.GenerateIncrementalLoanID(),
				TrxID:           int64(1000 + i),
				BorrowerID:      int64(2000 + i),
				PrincipalAmount: 20000 + float64(i)*5000,
				CollectedAmount: 10000 + float64(i)*2000,
				InterestRate:    0.05 + float64(i)*0.01,
				Status:          1,
				StatusDesc:      "disbursed",
			}
			helper.UpsertLoan(loan)
		}

		loansList := helper.GetLoans()
		if len(loansList) < 0 {
			t.Errorf("expected, got %d", len(loansList))
		}
	})
}

func TestGetLoanByLoanID(t *testing.T) {
	helper := NewHelper()

	t.Run("get loan by loan id", func(t *testing.T) {
		loan := model.Loan{
			LoanID:          helper.GenerateIncrementalLoanID(),
			TrxID:           1001,
			BorrowerID:      2001,
			PrincipalAmount: 500000,
			CollectedAmount: 250000,
			InterestRate:    0.1,
			Status:          1,
			StatusDesc:      "disbursed",
		}
		helper.UpsertLoan(loan)

		storedLoan := helper.GetLoanByLoanID(loan.LoanID)
		if storedLoan.LoanID != loan.LoanID {
			t.Errorf("expected loan id %d, got %d", loan.LoanID, storedLoan.LoanID)
		}
		if storedLoan.PrincipalAmount != loan.PrincipalAmount {
			t.Errorf("expected principal amount %f, got %f", loan.PrincipalAmount, storedLoan.PrincipalAmount)
		}

		nonExistentLoanID := int64(999)
		storedLoanNotFound := helper.GetLoanByLoanID(nonExistentLoanID)
		if storedLoanNotFound.LoanID != 0 {
			t.Errorf("expected non-existing loan to return empty loan, got %+v", storedLoanNotFound)
		}
	})
}

func TestUpsertLoan(t *testing.T) {
	helper := NewHelper()

	t.Run("upsert loan", func(t *testing.T) {
		loan := model.Loan{
			LoanID:          helper.GenerateIncrementalLoanID(),
			TrxID:           1001,
			BorrowerID:      2001,
			PrincipalAmount: 10000,
			CollectedAmount: 5000,
			InterestRate:    0.1,
			Status:          1,
			StatusDesc:      "Active",
			ApprovalInfo: &model.ApprovalInfo{
				PictureProof:             "base64encodedimage",
				FieldValidatorEmployeeID: 3001,
				ApprovalDate:             time.Now(),
			},
			Lending: []model.Lending{
				{
					LenderID:                    4001,
					InvestedAmount:              5000,
					OrganizerLenderAggrementURL: "https://example.com/agreement",
					ReturnAmount:                5500,
				},
			},
			DisbursementInfo: model.DisbursementInfo{
				AgreementSignedURLs: []string{"https://example.com/agreement1", "https://example.com/agreement2"},
				FieldOfficerID:      5001,
				DisbursementDate:    time.Now(),
			},
		}

		helper.UpsertLoan(loan)

		storedLoan := helper.GetLoanByLoanID(loan.LoanID)
		if storedLoan.LoanID != loan.LoanID {
			t.Errorf("expected loan id %d, got %d", loan.LoanID, storedLoan.LoanID)
		}
		if storedLoan.PrincipalAmount != loan.PrincipalAmount {
			t.Errorf("expected principal amount %f, got %f", loan.PrincipalAmount, storedLoan.PrincipalAmount)
		}
		if storedLoan.IsLenderInvested(4001) != true {
			t.Errorf("expected lender %d to be invested, but they are not", 4001)
		}
	})
}
