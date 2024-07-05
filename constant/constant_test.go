package constant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLoanStatusDesc(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		expectedDesc string
	}{
		{
			name:         "status proposed",
			status:       LoanStatusProposed,
			expectedDesc: "proposed",
		},
		{
			name:         "status approved",
			status:       LoanStatusApproved,
			expectedDesc: "approved",
		},
		{
			name:         "status invested",
			status:       LoanStatusInvested,
			expectedDesc: "invested",
		},
		{
			name:         "status disbursed",
			status:       LoanStatusDisbursed,
			expectedDesc: "disbursed",
		},
		{
			name:         "status unknown",
			status:       999, // assuming 999 is not a valid status
			expectedDesc: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetLoanStatusDesc(tt.status)
			assert.Equal(t, tt.expectedDesc, desc)
		})
	}
}
