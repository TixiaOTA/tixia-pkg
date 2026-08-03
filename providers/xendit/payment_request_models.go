package xendit

import (
	"errors"
	"strings"
	"time"
)

// PaymentRequestStatus represents the status of a payment request.
type PaymentRequestStatus string

const (
	PaymentRequestStatusAcceptingPayments PaymentRequestStatus = "ACCEPTING_PAYMENTS"
	PaymentRequestStatusRequiresAction    PaymentRequestStatus = "REQUIRES_ACTION"
	PaymentRequestStatusAuthorized        PaymentRequestStatus = "AUTHORIZED"
	PaymentRequestStatusPending           PaymentRequestStatus = "PENDING"
	PaymentRequestStatusSucceeded         PaymentRequestStatus = "SUCCEEDED"
	PaymentRequestStatusFailed            PaymentRequestStatus = "FAILED"
	PaymentRequestStatusCanceled          PaymentRequestStatus = "CANCELED"
	PaymentRequestStatusExpired           PaymentRequestStatus = "EXPIRED"
)

// PaymentRequestActionType describes the action required to complete a payment request.
type PaymentRequestActionType string

const (
	PaymentRequestActionTypePresentToCustomer PaymentRequestActionType = "PRESENT_TO_CUSTOMER"
	PaymentRequestActionTypeRedirectCustomer  PaymentRequestActionType = "REDIRECT_CUSTOMER"
	PaymentRequestActionTypeAPIPostRequest    PaymentRequestActionType = "API_POST_REQUEST"
)

// PaymentRequestActionDescriptor describes the format of a payment request action value.
type PaymentRequestActionDescriptor string

const (
	PaymentRequestActionDescriptorCapturePayment       PaymentRequestActionDescriptor = "CAPTURE_PAYMENT"
	PaymentRequestActionDescriptorPaymentCode          PaymentRequestActionDescriptor = "PAYMENT_CODE"
	PaymentRequestActionDescriptorQRString             PaymentRequestActionDescriptor = "QR_STRING"
	PaymentRequestActionDescriptorVirtualAccountNumber PaymentRequestActionDescriptor = "VIRTUAL_ACCOUNT_NUMBER"
	PaymentRequestActionDescriptorWebURL               PaymentRequestActionDescriptor = "WEB_URL"
	PaymentRequestActionDescriptorDeeplinkURL          PaymentRequestActionDescriptor = "DEEPLINK_URL"
	PaymentRequestActionDescriptorValidateOTP          PaymentRequestActionDescriptor = "VALIDATE_OTP"
	PaymentRequestActionDescriptorResendOTP            PaymentRequestActionDescriptor = "RESEND_OTP"
)

// PaymentRequestItem describes a line item attached to a payment request.
type PaymentRequestItem struct {
	ReferenceID   string            `json:"reference_id"`
	Type          string            `json:"type"`
	Name          string            `json:"name"`
	NetUnitAmount DecimalAmount     `json:"net_unit_amount"`
	Quantity      int               `json:"quantity"`
	URL           *string           `json:"url,omitempty"`
	ImageURL      *string           `json:"image_url,omitempty"`
	Category      string            `json:"category"`
	Subcategory   *string           `json:"subcategory,omitempty"`
	Description   *string           `json:"description,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// PaymentRequestShippingInformation describes shipping details for a payment request.
type PaymentRequestShippingInformation struct {
	Country       string  `json:"country"`
	StreetLine1   *string `json:"street_line1,omitempty"`
	StreetLine2   *string `json:"street_line2,omitempty"`
	City          *string `json:"city,omitempty"`
	ProvinceState *string `json:"province_state,omitempty"`
	PostalCode    *string `json:"postal_code,omitempty"`
}

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

	return ValidateTransactionAmount(r.ChannelCode, r.Currency, amount)
}

const paymentRequestIDLength = 39

func validatePaymentRequestID(paymentRequestID string) error {
	trimmed := strings.TrimSpace(paymentRequestID)
	if trimmed == "" {
		return errors.New("payment_request_id is required")
	}

	if len(trimmed) != paymentRequestIDLength {
		return errors.New("payment_request_id must be 39 characters")
	}

	return nil
}

// GetPaymentRequestStatusRequest is the request structure for GetPaymentRequest.
type GetPaymentRequestStatusRequest struct {
	PaymentRequestID string `json:"-"`
}

func (r *GetPaymentRequestStatusRequest) Validate() error {
	return validatePaymentRequestID(r.PaymentRequestID)
}

// CancelPaymentRequestStatusRequest is the request structure for CancelPaymentRequest.
type CancelPaymentRequestStatusRequest struct {
	PaymentRequestID string `json:"-"`
}

func (r *CancelPaymentRequestStatusRequest) Validate() error {
	return validatePaymentRequestID(r.PaymentRequestID)
}

// SimulatePaymentRequestRequest is the request structure for SimulatePaymentRequest.
// Test mode only — simulates payment completion for a payment request.
type SimulatePaymentRequestRequest struct {
	PaymentRequestID string         `json:"-"`
	Amount           *DecimalAmount `json:"amount,omitempty"`
}

func (r *SimulatePaymentRequestRequest) Validate() error {
	if err := validatePaymentRequestID(r.PaymentRequestID); err != nil {
		return err
	}

	if r.Amount != nil {
		amount, err := r.Amount.Float64()
		if err != nil {
			return errors.New("amount must be a valid number")
		}

		if amount <= 0 {
			return errors.New("amount must be greater than zero")
		}
	}

	return nil
}

// SimulatePaymentRequestResponse is the response structure for SimulatePaymentRequest.
type SimulatePaymentRequestResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// CreatePaymentRequestResponse is the response structure for CreatePaymentRequest, GetPaymentRequest, and CancelPaymentRequest.
type CreatePaymentRequestResponse struct {
	PaymentRequestID    string                             `json:"payment_request_id"`
	BusinessID          string                             `json:"business_id"`
	ReferenceID         string                             `json:"reference_id"`
	PaymentTokenID      *string                            `json:"payment_token_id,omitempty"`
	CustomerID          *string                            `json:"customer_id,omitempty"`
	LatestPaymentID     *string                            `json:"latest_payment_id,omitempty"`
	Type                PaymentType                        `json:"type"`
	Country             string                             `json:"country"`
	Currency            string                             `json:"currency"`
	RequestAmount       DecimalAmount                      `json:"request_amount"`
	CaptureMethod       *CaptureMethod                     `json:"capture_method,omitempty"`
	ChannelCode         string                             `json:"channel_code"`
	ChannelProperties   PaymentRequestChannelProperties    `json:"channel_properties,omitempty"`
	Actions             []PaymentRequestAction             `json:"actions,omitempty"`
	Status              string                             `json:"status"`
	FailureCode         *string                            `json:"failure_code,omitempty"`
	Description         *string                            `json:"description,omitempty"`
	Metadata            map[string]string                  `json:"metadata,omitempty"`
	Items               []PaymentRequestItem               `json:"items,omitempty"`
	ShippingInformation *PaymentRequestShippingInformation `json:"shipping_information,omitempty"`
	Created             time.Time                          `json:"created"`
	Updated             time.Time                          `json:"updated"`
}
