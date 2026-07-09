package panorama

import "errors"

// --------------------------------- ISSUED REQUEST ---------------------------------

// IssuedRequest is the request structure for Issued
type IssuedRequest struct {
	Signature   string `json:"signature"`
	OrderNumber string `json:"order_number"`
}

func (r *IssuedRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if r.OrderNumber == "" {
		return errors.New("order_number is required")
	}
	return nil
}

// --------------------------------- ISSUED RESPONSE ---------------------------------

// IssuedResponse is the response structure for Issued
type IssuedResponse struct {
	Result struct {
		BasicPanoramaResult
		Response BookingResponse `json:"response"`
	} `json:"result"`
}
