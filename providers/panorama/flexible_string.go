package panorama

import (
	"encoding/json"
	"fmt"
)

// FlexibleString unmarshals JSON values that Panorama may return as either a
// string or a bool (typically false meaning "no value").
type FlexibleString string

func (s FlexibleString) String() string {
	return string(s)
}

func (s FlexibleString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *FlexibleString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*s = ""
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		*s = FlexibleString(asString)
		return nil
	}

	var asBool bool
	if err := json.Unmarshal(data, &asBool); err == nil {
		*s = ""
		return nil
	}

	var asNumber json.Number
	if err := json.Unmarshal(data, &asNumber); err == nil {
		*s = FlexibleString(asNumber.String())
		return nil
	}

	return fmt.Errorf("flexible string must be a string, bool, number, or null")
}
