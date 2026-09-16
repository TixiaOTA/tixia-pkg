package xendit

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPStatusHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   int
		isStatus func(error) bool
	}{
		{name: "bad request", status: http.StatusBadRequest, isStatus: IsBadRequest},
		{name: "unauthorized", status: http.StatusUnauthorized, isStatus: IsUnauthorized},
		{name: "forbidden", status: http.StatusForbidden, isStatus: IsForbidden},
		{name: "not found", status: http.StatusNotFound, isStatus: IsNotFound},
		{name: "conflict", status: http.StatusConflict, isStatus: IsConflict},
		{name: "unprocessable entity", status: http.StatusUnprocessableEntity, isStatus: IsUnprocessableEntity},
		{name: "internal server error", status: http.StatusInternalServerError, isStatus: IsInternalServerError},
		{name: "service unavailable", status: http.StatusServiceUnavailable, isStatus: IsServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			xenditErr := &XenditError{HTTPStatus: tt.status}

			assert.True(t, tt.isStatus(xenditErr))
			assert.True(t, tt.isStatus(fmt.Errorf("xendit request: %w", xenditErr)))
			assert.False(t, tt.isStatus(&XenditError{HTTPStatus: http.StatusTeapot}))
			assert.False(t, tt.isStatus(errors.New("other error")))
			assert.False(t, tt.isStatus(nil))
		})
	}
}
