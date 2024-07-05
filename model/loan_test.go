package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRemainingRequiredAmount(t *testing.T) {
	tests := []struct {
		name           string
		loan           Loan
		expectedAmount float64
	}{
		{
			name: "remaining amount positive",
			loan: Loan{
				PrincipalAmount: 1000.0,
				CollectedAmount: 600.0,
			},
			expectedAmount: 400.0,
		},
		{
			name: "remaining amount zero",
			loan: Loan{
				PrincipalAmount: 1000.0,
				CollectedAmount: 1000.0,
			},
			expectedAmount: 0.0,
		},
		{
			name: "collected amount exceeds principal",
			loan: Loan{
				PrincipalAmount: 1000.0,
				CollectedAmount: 1200.0,
			},
			expectedAmount: -200.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedAmount, tt.loan.GetRemainingRequiredAmount())
		})
	}
}

func TestIsAmountFulfilled(t *testing.T) {
	tests := []struct {
		name          string
		loan          Loan
		expectedValue bool
	}{
		{
			name: "amount fulfilled",
			loan: Loan{
				PrincipalAmount: 1000.0,
				CollectedAmount: 1000.0,
			},
			expectedValue: true,
		},
		{
			name: "amount not fulfilled",
			loan: Loan{
				PrincipalAmount: 1000.0,
				CollectedAmount: 900.0,
			},
			expectedValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedValue, tt.loan.IsAmountFulfilled())
		})
	}
}

func TestCalculateLenderReturnAmount(t *testing.T) {
	tests := []struct {
		name          string
		lending       Lending
		interestRate  float64
		expectedValue float64
	}{
		{
			name: "positive return amount",
			lending: Lending{
				InvestedAmount: 1000.0,
			},
			interestRate:  0.1,
			expectedValue: 1100.0,
		},
		{
			name: "zero interest rate",
			lending: Lending{
				InvestedAmount: 1000.0,
			},
			interestRate:  0.0,
			expectedValue: 1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedValue, tt.lending.CalculateLenderReturnAmount(tt.interestRate))
		})
	}
}

func TestCalculateReturnAmount(t *testing.T) {
	tests := []struct {
		name          string
		loan          Loan
		expectedValue float64
	}{
		{
			name: "positive return amount",
			loan: Loan{
				PrincipalAmount: 1000.0,
				InterestRate:    0.1,
			},
			expectedValue: 1100.0,
		},
		{
			name: "zero interest rate",
			loan: Loan{
				PrincipalAmount: 1000.0,
				InterestRate:    0.0,
			},
			expectedValue: 1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedValue, tt.loan.CalculateReturnAmount())
		})
	}
}

func TestIsLenderInvested(t *testing.T) {
	tests := []struct {
		name          string
		loan          Loan
		lenderID      int64
		expectedValue bool
	}{
		{
			name: "lender invested",
			loan: Loan{
				Lending: []Lending{
					{LenderID: 1},
					{LenderID: 2},
				},
			},
			lenderID:      1,
			expectedValue: true,
		},
		{
			name: "lender not invested",
			loan: Loan{
				Lending: []Lending{
					{LenderID: 1},
					{LenderID: 2},
				},
			},
			lenderID:      3,
			expectedValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedValue, tt.loan.IsLenderInvested(tt.lenderID))
		})
	}
}
