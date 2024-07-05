package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"amartha-test/constant"
)

// Middleware is middleware handler to initialize start time
func (h *Handler) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		time.Sleep(123 * time.Millisecond) // add sleep for simulate latency
		ctx := context.WithValue(r.Context(), constant.CtxStartTimeKey, startTime)
		next(w, r.WithContext(ctx))
	}
}

type Response struct {
	Code         int         `json:"code"`
	Latency      string      `json:"latency"`
	Data         interface{} `json:"data,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
}

func (h *Handler) RenderResponse(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int, errMsg string) {
	startTime, ok := r.Context().Value(constant.CtxStartTimeKey).(time.Time)
	if !ok {
		log.Println("[RenderResponse] error retrieving start time")
		http.Error(w, "error retrieving start time", http.StatusInternalServerError)
		return
	}

	latency := time.Since(startTime).Milliseconds()

	response := Response{
		Code:         statusCode,
		Latency:      fmt.Sprintf("%dms", latency),
		Data:         data,
		ErrorMessage: errMsg,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("[RenderResponse] error marshalling JSON:", err)
		http.Error(w, "error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

func (h *Handler) RenderPDFResponse(w http.ResponseWriter, pdfData []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=agreement.pdf")
	w.WriteHeader(statusCode)
	w.Write(pdfData)
}
