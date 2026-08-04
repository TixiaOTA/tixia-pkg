package panorama

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexibleBool unmarshals JSON values that Panorama may return as either a
// bool or a string (e.g. "true", "false", "0", "1").
type FlexibleBool bool

func (b FlexibleBool) Bool() bool {
	return bool(b)
}

func (b FlexibleBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(b))
}

func (b *FlexibleBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*b = false
		return nil
	}

	var asBool bool
	if err := json.Unmarshal(data, &asBool); err == nil {
		*b = FlexibleBool(asBool)
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		if asString == "" {
			*b = false
			return nil
		}
		parsed, err := strconv.ParseBool(asString)
		if err != nil {
			return fmt.Errorf("flexible bool string %q is not a valid bool", asString)
		}
		*b = FlexibleBool(parsed)
		return nil
	}

	var asNumber json.Number
	if err := json.Unmarshal(data, &asNumber); err == nil {
		i, err := asNumber.Int64()
		if err != nil {
			return fmt.Errorf("flexible bool number %q is not a valid integer", asNumber)
		}
		*b = FlexibleBool(i != 0)
		return nil
	}

	return fmt.Errorf("flexible bool must be a bool, string, number, or null")
}
