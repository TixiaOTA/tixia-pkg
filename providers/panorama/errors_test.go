package panorama

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorCodeFromError(t *testing.T) {
	t.Parallel()

	t.Run("from panorama error", func(t *testing.T) {
		t.Parallel()

		err := &PanoramaError{
			Result: PanoramaErrorResult{
				BasicPanoramaResult: BasicPanoramaResult{
					ErrorCode: SeatNotAvailable,
					ErrorMsg:  "Seat not available",
				},
			},
		}

		code, ok := ErrorCodeFromError(err)
		require.True(t, ok)
		assert.Equal(t, SeatNotAvailable, code)
	})

	t.Run("from wrapped message", func(t *testing.T) {
		t.Parallel()

		err := fmt.Errorf("panorama error (code=4025): Flight is not available")

		code, ok := ErrorCodeFromError(err)
		require.True(t, ok)
		assert.Equal(t, FlightNotAvailable, code)
	})

	t.Run("missing code", func(t *testing.T) {
		t.Parallel()

		_, ok := ErrorCodeFromError(errors.New("provider timeout"))
		assert.False(t, ok)
	})
}

func TestIsFlightInventoryUnavailable(t *testing.T) {
	t.Parallel()

	assert.True(t, IsFlightInventoryUnavailable(SeatNotAvailable))
	assert.True(t, IsFlightInventoryUnavailable(FlightNotAvailable))
	assert.False(t, IsFlightInventoryUnavailable(SessionExpired))
	assert.False(t, IsFlightInventoryUnavailable(PriceChanged))
}
