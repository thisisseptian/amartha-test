package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"amartha-test/constant"
	"amartha-test/helper/mocks"
	"amartha-test/model"
)

func TestListLoan(t *testing.T) {
	mockHelper := new(mocks.IHelper)
	mockHandler := &Handler{
		Helper: mockHelper,
	}

	tests := []struct {
		name         string
		isError      bool
		expectedCode int
		mocks        func()
	}{
		{
			name:         "error - list empty",
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoans").Return([]model.Loan{}).Once()
			},
		},
		{
			name:         "success",
			isError:      false,
			expectedCode: http.StatusOK,
			mocks: func() {
				mockHelper.On("GetLoans").Return([]model.Loan{
					{
						LoanID: 1,
					},
				}).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks()

			r, err := http.NewRequest("GET", "/loans/list", nil)
			if err != nil {
				t.Fatal(err)
			}
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.ListLoan(w, r)

			isErr := false
			if w.Code != http.StatusOK && w.Code != http.StatusCreated {
				isErr = true
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.isError, isErr)
			mockHelper.AssertExpectations(t)
		})
	}
}

func TestDetailLoan(t *testing.T) {
	mockHelper := new(mocks.IHelper)
	mockHandler := &Handler{
		Helper: mockHelper,
	}

	tests := []struct {
		name         string
		vars         string
		isError      bool
		expectedCode int
		mocks        func()
	}{
		{
			name:         "error - convert string to int64",
			vars:         "?",
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name:         "error - sanitize payload",
			vars:         "0",
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name:         "error - loan data not found",
			vars:         "1",
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{}).Once()
			},
		},
		{
			name:         "success",
			vars:         "1",
			isError:      false,
			expectedCode: http.StatusOK,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{
					LoanID: 1,
				}).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks()

			r, err := http.NewRequest("GET", "loan/1/detail", nil)
			if err != nil {
				t.Fatal(err)
			}
			vars := map[string]string{"loan_id": tt.vars}
			r = mux.SetURLVars(r, vars)
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.DetailLoan(w, r)

			isErr := false
			if w.Code != http.StatusOK && w.Code != http.StatusCreated {
				isErr = true
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.isError, isErr)
			mockHelper.AssertExpectations(t)
		})
	}
}

func TestSubmitLoan(t *testing.T) {
	mockHelper := new(mocks.IHelper)
	mockHandler := &Handler{
		Helper: mockHelper,
	}

	tests := []struct {
		name         string
		requestBody  interface{}
		isError      bool
		expectedCode int
		mocks        func()
	}{
		{
			name:         "error - decode body",
			requestBody:  "invalid body",
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name:         "error - borrower id is empty",
			requestBody:  model.Loan{},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - principal amount is empty",
			requestBody: model.Loan{
				BorrowerID: 1,
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - interest rate < 0",
			requestBody: model.Loan{
				BorrowerID:      1,
				PrincipalAmount: 1000000,
				InterestRate:    -1,
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - user data not found",
			requestBody: model.Loan{
				BorrowerID:      1,
				PrincipalAmount: 1000000,
				InterestRate:    0.5,
			},
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{}).Once()
			},
		},
		{
			name: "error - user type not borrower",
			requestBody: model.Loan{
				BorrowerID:      1,
				PrincipalAmount: 1000000,
				InterestRate:    0.5,
			},
			isError:      true,
			expectedCode: http.StatusForbidden,
			mocks: func() {
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 1, UserType: constant.UserTypeLender}).Once()
			},
		},
		{
			name: "success",
			requestBody: model.Loan{
				BorrowerID:      1,
				PrincipalAmount: 1000000,
				InterestRate:    0.5,
			},
			isError:      false,
			expectedCode: http.StatusCreated,
			mocks: func() {
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 1, UserType: constant.UserTypeBorrower}).Once()
				mockHelper.On("GenerateIncrementalLoanID").Return(int64(1)).Once()
				mockHelper.On("UpsertLoan", mock.Anything).Return().Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks()

			var requestBody bytes.Buffer
			if err := json.NewEncoder(&requestBody).Encode(tt.requestBody); err != nil {
				t.Fatalf("could not encode request body: %v", err)
			}

			r, err := http.NewRequest("POST", "loan/submit", &requestBody)
			if err != nil {
				t.Fatal(err)
			}
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.SubmitLoan(w, r)

			isErr := false
			if w.Code != http.StatusOK && w.Code != http.StatusCreated {
				isErr = true
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.isError, isErr)
			mockHelper.AssertExpectations(t)
		})
	}
}

func TestApproveLoan(t *testing.T) {
	mockHelper := new(mocks.IHelper)
	mockHandler := &Handler{
		Helper: mockHelper,
	}

	tests := []struct {
		name         string
		vars         string
		requestBody  interface{}
		isError      bool
		expectedCode int
		mocks        func()
	}{
		{
			name:         "error - get loan id from vars",
			vars:         "?",
			requestBody:  "invalid body",
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name:         "error - decode body",
			vars:         "1",
			requestBody:  "invalid body",
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name:         "error - picture of proof is empty",
			vars:         "1",
			requestBody:  model.ApprovalInfo{},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - employee id is empty",
			vars: "1",
			requestBody: model.ApprovalInfo{
				PictureProof: "test",
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - loan data not found",
			vars: "1",
			requestBody: model.ApprovalInfo{
				PictureProof:             "test",
				FieldValidatorEmployeeID: 4,
			},
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{}).Once()
			},
		},
		{
			name: "error - loan status not proposed",
			vars: "1",
			requestBody: model.ApprovalInfo{
				PictureProof:             "test",
				FieldValidatorEmployeeID: 4,
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusApproved}).Once()
			},
		},
		{
			name: "error - user id not found",
			vars: "1",
			requestBody: model.ApprovalInfo{
				PictureProof:             "test",
				FieldValidatorEmployeeID: 4,
			},
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusProposed}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{}).Once()
			},
		},
		{
			name: "error - user type is not field validator employee",
			vars: "1",
			requestBody: model.ApprovalInfo{
				PictureProof:             "test",
				FieldValidatorEmployeeID: 4,
			},
			isError:      true,
			expectedCode: http.StatusForbidden,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusProposed}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 4, UserType: constant.UserTypeBorrower}).Once()
			},
		},
		{
			name: "success",
			vars: "1",
			requestBody: model.ApprovalInfo{
				PictureProof:             "test",
				FieldValidatorEmployeeID: 4,
			},
			isError:      false,
			expectedCode: http.StatusOK,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusProposed}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 4, UserType: constant.UserTypeFieldValidatorEmployee}).Once()
				mockHelper.On("UpsertLoan", mock.Anything).Return().Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks()

			var requestBody bytes.Buffer
			if err := json.NewEncoder(&requestBody).Encode(tt.requestBody); err != nil {
				t.Fatalf("could not encode request body: %v", err)
			}

			r, err := http.NewRequest("POST", "loan/1/approve", &requestBody)
			if err != nil {
				t.Fatal(err)
			}
			vars := map[string]string{"loan_id": tt.vars}
			r = mux.SetURLVars(r, vars)
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.ApproveLoan(w, r)

			isErr := false
			if w.Code != http.StatusOK && w.Code != http.StatusCreated {
				isErr = true
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.isError, isErr)
			mockHelper.AssertExpectations(t)
		})
	}
}

func TestInvestLoan(t *testing.T) {
	mockHelper := new(mocks.IHelper)
	mockHandler := &Handler{
		Helper: mockHelper,
	}

	tests := []struct {
		name         string
		vars         string
		requestBody  interface{}
		isError      bool
		expectedCode int
		mocks        func()
	}{
		{
			name:         "error - get loan id from vars",
			vars:         "?",
			requestBody:  "invalid body",
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name:         "error - decode body",
			vars:         "1",
			requestBody:  "invalid body",
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name:         "error - lender id is empty",
			vars:         "1",
			requestBody:  model.Lending{},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - invested amount is empty",
			vars: "1",
			requestBody: model.Lending{
				LenderID: 2,
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - loan data not found",
			vars: "1",
			requestBody: model.Lending{
				LenderID:       2,
				InvestedAmount: 1000000,
			},
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{}).Once()
			},
		},
		{
			name: "error - loan status is not approved",
			vars: "1",
			requestBody: model.Lending{
				LenderID:       2,
				InvestedAmount: 1000000,
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusProposed}).Once()
			},
		},
		{
			name: "error - user data not found",
			vars: "1",
			requestBody: model.Lending{
				LenderID:       2,
				InvestedAmount: 1000000,
			},
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusApproved}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{}).Once()
			},
		},
		{
			name: "error - user type is not lender",
			vars: "1",
			requestBody: model.Lending{
				LenderID:       2,
				InvestedAmount: 1000000,
			},
			isError:      true,
			expectedCode: http.StatusForbidden,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusApproved}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 2, UserType: constant.UserTypeBorrower}).Once()
			},
		},
		{
			name: "error - invested amount is bigger than remaining required amount",
			vars: "1",
			requestBody: model.Lending{
				LenderID:       2,
				InvestedAmount: 1000000,
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusApproved, PrincipalAmount: 1000000, CollectedAmount: 500000}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 2, UserType: constant.UserTypeLender}).Once()
			},
		},
		{
			name: "error - fail generate lender agreements",
			vars: "1",
			requestBody: model.Lending{
				LenderID:       2,
				InvestedAmount: 1000000,
			},
			isError:      true,
			expectedCode: http.StatusInternalServerError,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusApproved, PrincipalAmount: 1000000, CollectedAmount: 0}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 2, UserType: constant.UserTypeLender}).Once()
				mockHelper.On("GenerateLenderAgreementPDF", mock.Anything).Return(errors.New("fail")).Once()
			},
		},
		{
			name: "success",
			vars: "1",
			requestBody: model.Lending{
				LenderID:       2,
				InvestedAmount: 1000000,
			},
			isError:      false,
			expectedCode: http.StatusOK,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusApproved, PrincipalAmount: 1500000, CollectedAmount: 500000, Lending: []model.Lending{{LenderID: 2, InvestedAmount: 500000}}}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 2, UserType: constant.UserTypeLender}).Once()
				mockHelper.On("GenerateLenderAgreementPDF", mock.Anything).Return(nil).Once()
				mockHelper.On("UpsertLoan", mock.Anything).Return().Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks()

			var requestBody bytes.Buffer
			if err := json.NewEncoder(&requestBody).Encode(tt.requestBody); err != nil {
				t.Fatalf("could not encode request body: %v", err)
			}

			r, err := http.NewRequest("POST", "loan/1/invest", &requestBody)
			if err != nil {
				t.Fatal(err)
			}
			vars := map[string]string{"loan_id": tt.vars}
			r = mux.SetURLVars(r, vars)
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.InvestLoan(w, r)

			isErr := false
			if w.Code != http.StatusOK && w.Code != http.StatusCreated {
				isErr = true
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.isError, isErr)
			mockHelper.AssertExpectations(t)
		})
	}
}

func TestDisburseLoan(t *testing.T) {
	mockHelper := new(mocks.IHelper)
	mockHandler := &Handler{
		Helper: mockHelper,
	}

	tests := []struct {
		name         string
		vars         string
		requestBody  interface{}
		isError      bool
		expectedCode int
		mocks        func()
	}{
		{
			name:         "error - get loan id from vars",
			vars:         "?",
			requestBody:  "invalid body",
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name:         "error - decode body",
			vars:         "1",
			requestBody:  "invalid body",
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name:         "error - field officer id is empty",
			vars:         "1",
			requestBody:  model.Disbursement{},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - disbursement date is invalid",
			vars: "1",
			requestBody: model.Disbursement{
				FieldOfficerID: 5,
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - loan data not found",
			vars: "1",
			requestBody: model.Disbursement{
				FieldOfficerID:   5,
				DisbursementDate: time.Now(),
			},
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{}).Once()
			},
		},
		{
			name: "error - loan status is not signed",
			vars: "1",
			requestBody: model.Disbursement{
				FieldOfficerID:   5,
				DisbursementDate: time.Now(),
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusInvested}).Once()
			},
		},
		{
			name: "error - user data not found",
			vars: "1",
			requestBody: model.Disbursement{
				FieldOfficerID:   5,
				DisbursementDate: time.Now(),
			},
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusSigned}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{}).Once()
			},
		},
		{
			name: "error - user type is not field officer employee",
			vars: "1",
			requestBody: model.Disbursement{
				FieldOfficerID:   5,
				DisbursementDate: time.Now(),
			},
			isError:      true,
			expectedCode: http.StatusForbidden,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusSigned}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 4, UserType: constant.UserTypeFieldValidatorEmployee}).Once()
			},
		},
		{
			name: "success",
			vars: "1",
			requestBody: model.Disbursement{
				FieldOfficerID:   5,
				DisbursementDate: time.Now(),
			},
			isError:      false,
			expectedCode: http.StatusOK,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusSigned}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 5, UserType: constant.UserTypeFieldOfficerEmployee}).Once()
				mockHelper.On("UpsertLoan", mock.Anything).Return().Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks()

			var requestBody bytes.Buffer
			if err := json.NewEncoder(&requestBody).Encode(tt.requestBody); err != nil {
				t.Fatalf("could not encode request body: %v", err)
			}

			r, err := http.NewRequest("POST", "loan/1/disburse", &requestBody)
			if err != nil {
				t.Fatal(err)
			}
			vars := map[string]string{"loan_id": tt.vars}
			r = mux.SetURLVars(r, vars)
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.DisburseLoan(w, r)

			isErr := false
			if w.Code != http.StatusOK && w.Code != http.StatusCreated {
				isErr = true
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.isError, isErr)
			mockHelper.AssertExpectations(t)
		})
	}
}
