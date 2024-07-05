package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"amartha-test/helper/mocks"
)

func TestNewHandler(t *testing.T) {
	mockHelper := new(mocks.IHelper)
	handler := NewHandler(mockHelper)

	assert.NotNil(t, handler)
	assert.Equal(t, mockHelper, handler.Helper)
}
