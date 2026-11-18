package train

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetSeatAvailabilityRequest is the Orbis train seat availability payload.
// FareCodes must follow the order of the requested journeys (outbound, inbound).
type GetSeatAvailabilityRequest struct {
	FareCodes []string `json:"fare_codes"`
	Provider  string   `json:"provider"`
}

func (r GetSeatAvailabilityRequest) Validate() error {
	if r.Provider != ProviderKAI {
		return fmt.Errorf("provider must be %q", ProviderKAI)
	}
	if len(r.FareCodes) < 1 || len(r.FareCodes) > 2 {
		return fmt.Errorf("fare_codes must contain one or two fares")
	}
	for i, code := range r.FareCodes {
		if strings.TrimSpace(code) == "" {
			return fmt.Errorf("fare_codes[%d] is required", i)
		}
	}
	return nil
}

// GetSeatAvailabilityResponse contains one set of cabins per input fare code.
// Result.Response[0] belongs to FareCodes[0], etc.; a set may contain multiple cabins.
type GetSeatAvailabilityResponse struct {
	JSONRPC string                    `json:"jsonrpc"`
	Result  GetSeatAvailabilityResult `json:"result"`
}

type GetSeatAvailabilityResult struct {
	BasicResult
	Response [][]SeatCabin `json:"response"`
}

type SeatCabin struct {
	CabinNumber         int       `json:"cabin_number"`
	CabinCode           string    `json:"cabin_code"`
	CabinName           string    `json:"cabin_name"`
	MaximumColumnNumber int       `json:"maximum_column_number"`
	SeatRows            []SeatRow `json:"seat_rows"`
}

type SeatRow struct {
	RowNumber string `json:"row_number"`
	Seats     []Seat `json:"seats"`
}

// SeatCode may be a JSON number or an empty string in Orbis responses.
// Keep the raw value instead of coercing it and risking precision loss.
type Seat struct {
	SeatCode     json.RawMessage `json:"seat_code"`
	Column       string          `json:"column"`
	Sequence     int             `json:"sequence"`
	Availability int             `json:"availability"`
}
