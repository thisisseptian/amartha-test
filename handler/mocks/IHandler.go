// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// IHandler is an autogenerated mock type for the IHandler type
type IHandler struct {
	mock.Mock
}

// ApproveLoan provides a mock function with given fields: w, r
func (_m *IHandler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// DetailAgreement provides a mock function with given fields: w, r
func (_m *IHandler) DetailAgreement(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// DetailLoan provides a mock function with given fields: w, r
func (_m *IHandler) DetailLoan(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// DetailUser provides a mock function with given fields: w, r
func (_m *IHandler) DetailUser(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// DisburseLoan provides a mock function with given fields: w, r
func (_m *IHandler) DisburseLoan(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// InvestLoan provides a mock function with given fields: w, r
func (_m *IHandler) InvestLoan(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// ListAgreement provides a mock function with given fields: w, r
func (_m *IHandler) ListAgreement(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// ListLoan provides a mock function with given fields: w, r
func (_m *IHandler) ListLoan(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// ListUser provides a mock function with given fields: w, r
func (_m *IHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// SignAgreement provides a mock function with given fields: w, r
func (_m *IHandler) SignAgreement(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// SubmitLoan provides a mock function with given fields: w, r
func (_m *IHandler) SubmitLoan(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// NewIHandler creates a new instance of IHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *IHandler {
	mock := &IHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
