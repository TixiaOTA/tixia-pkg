package panorama

import (
	"errors"
	"regexp"
	"strconv"
)

var panoramaErrorCodePattern = regexp.MustCompile(`(?i)code[ =](\d+)`)

// ErrorCodeFromError extracts a Panorama error code from an error chain.
func ErrorCodeFromError(err error) (PanoramaErrorCode, bool) {
	if err == nil {
		return 0, false
	}

	var panoramaErr *PanoramaError
	if errors.As(err, &panoramaErr) {
		return panoramaErr.Result.ErrorCode, true
	}

	matches := panoramaErrorCodePattern.FindStringSubmatch(err.Error())
	if len(matches) < 2 {
		return 0, false
	}

	code, parseErr := strconv.Atoi(matches[1])
	if parseErr != nil {
		return 0, false
	}

	return PanoramaErrorCode(code), true
}

// IsFlightInventoryUnavailable reports whether the provider error means the
// selected flight inventory can no longer be booked.
func IsFlightInventoryUnavailable(code PanoramaErrorCode) bool {
	switch code {
	case SSRNotAvailable,
		SeatNotAvailable,
		NoSchedules,
		PriceNotAvailable,
		FlightNotAvailable:
		return true
	default:
		return false
	}
}
