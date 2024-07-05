package handler

import (
	"context"
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

func TestListUser(t *testing.T) {
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
				mockHelper.On("GetUsers").Return([]model.User{}).Once()
			},
		},
		{
			name:         "success",
			isError:      false,
			expectedCode: http.StatusOK,
			mocks: func() {
				mockHelper.On("GetUsers").Return([]model.User{
					{
						UserID:   1,
						UserName: "Septian",
						UserType: constant.UserTypeBorrower,
					},
				}).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks()

			r, err := http.NewRequest("GET", "/users/list", nil)
			if err != nil {
				t.Fatal(err)
			}
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.ListUser(w, r)

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

func TestDetailUser(t *testing.T) {
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
			name:         "error - user data not found",
			vars:         "1",
			isError:      true,
			expectedCode: http.StatusNotFound,
			mocks: func() {
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{}).Once()
			},
		},
		{
			name:         "success",
			vars:         "1",
			isError:      false,
			expectedCode: http.StatusOK,
			mocks: func() {
				mockHelper.On("GetUserByUserID", mock.Anything).Return(model.User{
					UserID:   1,
					UserName: "Septian",
					UserType: constant.UserTypeBorrower,
				}).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks()

			r, err := http.NewRequest("GET", "user/1/detail", nil)
			if err != nil {
				t.Fatal(err)
			}
			vars := map[string]string{"user_id": tt.vars}
			r = mux.SetURLVars(r, vars)
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// main func
			mockHandler.DetailUser(w, r)

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
