package xendit

import (
	"errors"
	"strings"
	"time"
)

// PaymentRequestStatus represents the status of a payment request.
type PaymentRequestStatus string

const (
	PaymentRequestStatusRequiresAction PaymentRequestStatus = "REQUIRES_ACTION"
	PaymentRequestStatusPending        PaymentRequestStatus = "PENDING"
	PaymentRequestStatusSucceeded      PaymentRequestStatus = "SUCCEEDED"
	PaymentRequestStatusFailed         PaymentRequestStatus = "FAILED"
	PaymentRequestStatusCanceled       PaymentRequestStatus = "CANCELED"
	PaymentRequestStatusExpired        PaymentRequestStatus = "EXPIRED"
)

// PaymentRequestActionType describes the action required to complete a payment request.
type PaymentRequestActionType string

const PaymentRequestActionTypeRedirectCustomer PaymentRequestActionType = "REDIRECT_CUSTOMER"

// PaymentRequestActionDescriptor describes the format of a payment request action value.
type PaymentRequestActionDescriptor string

const PaymentRequestActionDescriptorWebURL PaymentRequestActionDescriptor = "WEB_URL"

// PaymentRequestAction describes a customer action required to complete a payment request.
type PaymentRequestAction struct {
	Type       PaymentRequestActionType       `json:"type"`
	Value      string                         `json:"value"`
	Descriptor PaymentRequestActionDescriptor `json:"descriptor"`
}

// CreatePaymentRequestRequest is the request structure for CreatePaymentRequest.
// Scoped to type=PAY with currency=IDR and country=ID.
type CreatePaymentRequestRequest struct {
	ReferenceID       string                          `json:"reference_id"`
	Type              PaymentType                     `json:"type"`
	Country           string                          `json:"country"`
	Currency          string                          `json:"currency"`
	RequestAmount     DecimalAmount                   `json:"request_amount"`
	ChannelCode       string                          `json:"channel_code"`
	CaptureMethod     *CaptureMethod                  `json:"capture_method,omitempty"`
	ChannelProperties PaymentRequestChannelProperties `json:"channel_properties"`
	CustomerID        *string                         `json:"customer_id,omitempty"`
	Customer          *CustomerDetails                `json:"customer,omitempty"`
	Description       *string                         `json:"description,omitempty"`
	Metadata          map[string]string               `json:"metadata,omitempty"`
}

func (r *CreatePaymentRequestRequest) Validate() error {
	if strings.TrimSpace(r.ReferenceID) == "" {
		return errors.New("reference_id is required")
	}

	if r.Type == "" {
		return errors.New("type is required")
	}

	if r.Type != PaymentTypePay {
		return errors.New("type must be PAY")
	}

	if strings.TrimSpace(r.Country) == "" {
		return errors.New("country is required")
	}

	if r.Country != "ID" {
		return errors.New("country must be ID")
	}

	if strings.TrimSpace(r.Currency) == "" {
		return errors.New("currency is required")
	}

	if r.Currency != "IDR" {
		return errors.New("currency must be IDR")
	}

	amount, err := r.RequestAmount.Float64()
	if err != nil {
		return errors.New("request_amount must be a valid number")
	}

	if amount <= 0 {
		return errors.New("request_amount must be greater than zero")
	}

	if strings.TrimSpace(r.ChannelCode) == "" {
		return errors.New("channel_code is required")
	}

	if r.CaptureMethod != nil {
		switch *r.CaptureMethod {
		case CaptureMethodAutomatic, CaptureMethodManual:
		default:
			return errors.New("capture_method must be AUTOMATIC or MANUAL")
		}
	}

	if err := validateChannelProperties(r.ChannelCode, r.ChannelProperties); err != nil {
		return err
	}

	if r.CustomerID != nil && strings.TrimSpace(*r.CustomerID) == "" {
		return errors.New("customer_id cannot be empty")
	}

	if r.Customer != nil {
		if err := r.Customer.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// CreatePaymentRequestResponse is the response structure for CreatePaymentRequest.
type CreatePaymentRequestResponse struct {
	PaymentRequestID  string                 `json:"payment_request_id"`
	BusinessID        string                 `json:"business_id"`
	ReferenceID       string                 `json:"reference_id"`
	Type              PaymentType            `json:"type"`
	Country           string                 `json:"country"`
	Currency          string                 `json:"currency"`
	RequestAmount     DecimalAmount          `json:"request_amount"`
	CaptureMethod     *CaptureMethod         `json:"capture_method,omitempty"`
	ChannelCode       string                 `json:"channel_code"`
	ChannelProperties map[string]any         `json:"channel_properties,omitempty"`
	Actions           []PaymentRequestAction `json:"actions,omitempty"`
	Status            string                 `json:"status"`
	Description       *string                `json:"description,omitempty"`
	Metadata          map[string]string      `json:"metadata,omitempty"`
	Created           time.Time              `json:"created"`
	Updated           time.Time              `json:"updated"`
}
