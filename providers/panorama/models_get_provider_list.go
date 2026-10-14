package panorama

import "errors"

// GetProviderListRequest requests the providers for a provider type.
type GetProviderListRequest struct {
	Signature    string `json:"-" log:"redact"`
	ProviderType string `json:"provider_type"`
}

func (r *GetProviderListRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if r.ProviderType == "" {
		return errors.New("provider_type is required")
	}
	return nil
}

// GetProviderListResponse preserves the existing provider-list event model.
type GetProviderListResponse = FlightGetProviderListResultsEventData

type FlightGetProviderListResultsEventData struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
		Response  struct {
			Providers []struct {
				Provider string   `json:"provider"`
				Carriers []string `json:"carriers"`
			} `json:"providers"`
		} `json:"response"`
	} `json:"result"`
}

type FlightResultsErrorEventData struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	} `json:"result"`
}
