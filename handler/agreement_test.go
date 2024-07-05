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

func TestListAgreement(t *testing.T) {
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
				mockHelper.On("GetAgreements").Return([]model.Aggrement{}).Once()
			},
		},
		{
			name:         "success",
			isError:      false,
			expectedCode: http.StatusOK,
			mocks: func() {
				mockHelper.On("GetAgreements").Return([]model.Aggrement{
					{
						AggrementID: 1,
					},
				}).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks()

			r, err := http.NewRequest("GET", "/agreement/list", nil)
			if err != nil {
				t.Fatal(err)
			}
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.ListAgreement(w, r)

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

func TestViewAgreement(t *testing.T) {
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
			name:         "error - agreement data not found",
			vars:         "1",
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetAgreementByAgreementID", mock.Anything).Return(model.Aggrement{}).Once()
			},
		},
		{
			name:         "success",
			vars:         "1",
			isError:      false,
			expectedCode: http.StatusOK,
			mocks: func() {
				mockHelper.On("GetAgreementByAgreementID", mock.Anything).Return(model.Aggrement{
					AggrementID: 1,
				}).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks()

			r, err := http.NewRequest("GET", "agreement/1/view", nil)
			if err != nil {
				t.Fatal(err)
			}
			vars := map[string]string{"agreement_id": tt.vars}
			r = mux.SetURLVars(r, vars)
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.ViewAgreement(w, r)

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

func TestSignAgreement(t *testing.T) {
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
			name:         "error - get agreement id from vars",
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
			name:         "error - loan id is empty",
			vars:         "1",
			requestBody:  model.Sign{},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - user id is empty",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks:        func() {},
		},
		{
			name: "error - loan data not found",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
				UserID: 2,
			},
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{}).Once()
			},
		},
		{
			name: "error - loan status is not invested",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
				UserID: 2,
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusApproved}).Once()
			},
		},
		{
			name: "error - user data is not found",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
				UserID: 2,
			},
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusInvested}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{}).Once()
			},
		},
		{
			name: "error - agreement data is not found",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
				UserID: 2,
			},
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusInvested}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 1}).Once()
				mockHelper.On("GetAgreementByAgreementID", mock.Anything).Return(model.Aggrement{}).Once()
			},
		},
		{
			name: "error - forbidden user to sign",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
				UserID: 2,
			},
			isError:      true,
			expectedCode: http.StatusForbidden,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusInvested}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 2}).Once()
				mockHelper.On("GetAgreementByAgreementID", mock.Anything).Return(model.Aggrement{AggrementID: 1, UserID: 1}).Once()
			},
		},
		{
			name: "error - agreement is already signed",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
				UserID: 2,
			},
			isError:      true,
			expectedCode: http.StatusBadRequest,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusInvested}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 2}).Once()
				mockHelper.On("GetAgreementByAgreementID", mock.Anything).Return(model.Aggrement{AggrementID: 1, UserID: 2, IsSigned: true}).Once()
			},
		},
		{
			name: "error - fail generate signed agreement",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
				UserID: 2,
			},
			isError:      true,
			expectedCode: http.StatusInternalServerError,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusInvested}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 2}).Once()
				mockHelper.On("GetAgreementByAgreementID", mock.Anything).Return(model.Aggrement{AggrementID: 1, UserID: 2, IsSigned: false}).Once()
				mockHelper.On("UpsertAgreement", mock.Anything).Return().Once()
				mockHelper.On("GenerateSignedAgreementPDF", mock.Anything, mock.Anything).Return(errors.New("fail")).Once()
			},
		},
		{
			name: "error - lender - fail check all lender signed",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
				UserID: 2,
			},
			isError:      true,
			expectedCode: http.StatusInternalServerError,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusInvested}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 2, UserType: constant.UserTypeLender}).Once()
				mockHelper.On("GetAgreementByAgreementID", mock.Anything).Return(model.Aggrement{AggrementID: 1, UserID: 2, IsSigned: false}).Once()
				mockHelper.On("UpsertAgreement", mock.Anything).Return().Once()
				mockHelper.On("GenerateSignedAgreementPDF", mock.Anything, mock.Anything).Return(nil).Once()
				mockHelper.On("CheckAgreementCompletelySignedByLender", mock.Anything).Return(false, errors.New("fail")).Once()
			},
		},
		{
			name: "error - lender - generate borrower agreement",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
				UserID: 2,
			},
			isError:      true,
			expectedCode: http.StatusInternalServerError,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusInvested}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 2, UserType: constant.UserTypeLender}).Once()
				mockHelper.On("GetAgreementByAgreementID", mock.Anything).Return(model.Aggrement{AggrementID: 1, UserID: 2, IsSigned: false}).Once()
				mockHelper.On("UpsertAgreement", mock.Anything).Return().Once()
				mockHelper.On("GenerateSignedAgreementPDF", mock.Anything, mock.Anything).Return(nil).Once()
				mockHelper.On("CheckAgreementCompletelySignedByLender", mock.Anything).Return(true, nil).Once()
				mockHelper.On("GenerateBorrowerAgreementPDF", mock.Anything).Return(errors.New("fail")).Once()
			},
		},
		{
			name: "success - borrower",
			vars: "1",
			requestBody: model.Sign{
				LoanID: 1,
				UserID: 2,
			},
			isError:      false,
			expectedCode: http.StatusOK,
			mocks: func() {
				mockHelper.On("GetLoanByLoanID", mock.Anything).Return(model.Loan{LoanID: 1, Status: constant.LoanStatusInvested}).Once()
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{UserID: 2, UserType: constant.UserTypeBorrower}).Once()
				mockHelper.On("GetAgreementByAgreementID", mock.Anything).Return(model.Aggrement{AggrementID: 1, UserID: 2, IsSigned: false}).Once()
				mockHelper.On("UpsertAgreement", mock.Anything).Return().Once()
				mockHelper.On("GenerateSignedAgreementPDF", mock.Anything, mock.Anything).Return(nil).Once()
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

			r, err := http.NewRequest("POST", "agreement/1/sign", &requestBody)
			if err != nil {
				t.Fatal(err)
			}
			vars := map[string]string{"agreement_id": tt.vars}
			r = mux.SetURLVars(r, vars)
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.SignAgreement(w, r)

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
