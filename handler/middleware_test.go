package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"amartha-test/constant"
)

func TestMiddleware(t *testing.T) {
	mockHandler := &Handler{} // Assuming no need for mocks in the Middleware itself

	tests := []struct {
		name         string
		expectedCode int
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}
			startTime := time.Now()
			time.Sleep(123 * time.Millisecond) // add sleep to simulate latency
			ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// main func
			mockHandler.Middleware(handlerFunc)(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestRenderResponse(t *testing.T) {
	mockHandler := &Handler{} // Assuming no need for mocks in RenderResponse

	tests := []struct {
		name         string
		data         interface{}
		statusCode   int
		expectedCode int
	}{
		{
			name:         "error - got start time from context",
			data:         "test data",
			statusCode:   http.StatusOK,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "success",
			data:         "test data",
			statusCode:   http.StatusOK,
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.name != "error - got start time from context" {
				startTime := time.Now()
				time.Sleep(123 * time.Millisecond) // add sleep to simulate latency
				ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
				r = r.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			// main func
			mockHandler.RenderResponse(w, r, tt.data, tt.statusCode, "")

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestRenderPDFResponse(t *testing.T) {
	mockHandler := &Handler{} // Assuming no need for mocks in RenderResponse

	tests := []struct {
		name         string
		pdfData      []byte
		statusCode   int
		expectedCode int
	}{
		{
			name:         "success",
			pdfData:      []byte{},
			statusCode:   http.StatusOK,
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.name != "error - got start time from context" {
				startTime := time.Now()
				time.Sleep(123 * time.Millisecond) // add sleep to simulate latency
				ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
				r = r.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			// main func
			mockHandler.RenderPDFResponse(w, tt.pdfData, tt.statusCode)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
