package panorama

import "errors"

// --------------------------------- GET FF AVAILABILITY REQUEST ---------------------------------

// GetFFAvailabilityRequest is the request structure for GetFFAvailability.
// The signature is the only parameter; it is sent in the Action/Signature headers
// and echoed in the request body (see the GetBalanceRequest pattern).
type GetFFAvailabilityRequest struct {
	Signature string `json:"signature"`
}

func (r *GetFFAvailabilityRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

// --------------------------------- GET FF AVAILABILITY RESPONSE ---------------------------------

// GetFFAvailabilityResponse is the response structure for GetFFAvailability
type GetFFAvailabilityResponse struct {
	Result struct {
		BasicPanoramaResult
		Response FFAvailabilityResponse `json:"response"`
	} `json:"result"`
}

// FFAvailabilityResponse is the response payload for GetFFAvailability
type FFAvailabilityResponse struct {
	FFAvailabilityProvider []FFAvailabilityProvider `json:"ff_availability_provider"`
}

// FFAvailabilityProvider is a provider block in the GetFFAvailability response
type FFAvailabilityProvider struct {
	Provider       string           `json:"provider"`
	ProgramName    string           `json:"program_name"`
	FFAvailability []FFAvailability `json:"ff_availability"`
	Status         string           `json:"status"`
	ScheduleID     int              `json:"schedule_id"`
	Sequence       int              `json:"sequence"`
	CarrierCodes   []string         `json:"carrier_codes"`
}

// FFAvailability is a frequent flyer program in the GetFFAvailability response
type FFAvailability struct {
	Name   string `json:"name"`
	FFCode string `json:"ff_code"`
}
