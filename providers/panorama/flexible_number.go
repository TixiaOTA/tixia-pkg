package panorama

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexibleNumber unmarshals JSON numbers returned as either numbers or strings.
type FlexibleNumber float64

func (n FlexibleNumber) Int64() int64 {
	return int64(n)
}

func (n FlexibleNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(n))
}

func (n *FlexibleNumber) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*n = 0
		return nil
	}

	var asNumber float64
	if err := json.Unmarshal(data, &asNumber); err == nil {
		*n = FlexibleNumber(asNumber)
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		if asString == "" {
			*n = 0
			return nil
		}
		parsed, err := strconv.ParseFloat(asString, 64)
		if err != nil {
			return fmt.Errorf("flexible number string %q is not a valid number", asString)
		}
		*n = FlexibleNumber(parsed)
		return nil
	}

	return fmt.Errorf("flexible number must be a number, string, or null")
}
