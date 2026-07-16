package xendit

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const paymentIDLength = 39

// PaymentType is the payment collection intent type for a payment.
type PaymentType string

const (
	PaymentTypePay                 PaymentType = "PAY"
	PaymentTypePayAndSave          PaymentType = "PAY_AND_SAVE"
	PaymentTypeReusablePaymentCode PaymentType = "REUSABLE_PAYMENT_CODE"
)

// PaymentStatus represents the status of a payment.
type PaymentStatus string

const (
	PaymentStatusAuthorized PaymentStatus = "AUTHORIZED"
	PaymentStatusCanceled   PaymentStatus = "CANCELED"
	PaymentStatusSucceeded  PaymentStatus = "SUCCEEDED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusExpired    PaymentStatus = "EXPIRED"
	PaymentStatusPending    PaymentStatus = "PENDING"
)

// DecimalAmount unmarshals Xendit amount fields that may be JSON strings or numbers.
type DecimalAmount struct {
	value float64
}

func roundDecimalAmount(value float64) float64 {
	return math.Round(value*100) / 100
}

func (a DecimalAmount) Float64() (float64, error) {
	return a.value, nil
}

func (a DecimalAmount) String() string {
	return strconv.FormatFloat(a.value, 'f', 2, 64)
}

func (a *DecimalAmount) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		a.value = 0
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		parsed, err := strconv.ParseFloat(asString, 64)
		if err != nil {
			return fmt.Errorf("decimal amount must be a string or number")
		}

		a.value = roundDecimalAmount(parsed)
		return nil
	}

	var asFloat float64
	if err := json.Unmarshal(data, &asFloat); err == nil {
		a.value = roundDecimalAmount(asFloat)
		return nil
	}

	return fmt.Errorf("decimal amount must be a string or number")
}

func (a DecimalAmount) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.value)
}

// NewDecimalAmountFromFloat creates a DecimalAmount from a float64 value rounded to 2 decimal places.
func NewDecimalAmountFromFloat(value float64) DecimalAmount {
	return DecimalAmount{value: roundDecimalAmount(value)}
}

func validatePaymentID(paymentID string) error {
	trimmed := strings.TrimSpace(paymentID)
	if trimmed == "" {
		return errors.New("payment_id is required")
	}

	if len(trimmed) != paymentIDLength {
		return errors.New("payment_id must be 39 characters")
	}

	return nil
}

// GetPaymentRequest is the request structure for GetPayment.
type GetPaymentRequest struct {
	PaymentID string `json:"-"`
}

func (r *GetPaymentRequest) Validate() error {
	return validatePaymentID(r.PaymentID)
}

// CancelPaymentRequest is the request structure for CancelPayment.
type CancelPaymentRequest struct {
	PaymentID string `json:"-"`
}

func (r *CancelPaymentRequest) Validate() error {
	return validatePaymentID(r.PaymentID)
}

// CapturePaymentRequest is the request structure for CapturePayment.
type CapturePaymentRequest struct {
	PaymentID     string        `json:"-"`
	CaptureAmount DecimalAmount `json:"capture_amount"`
}

func (r *CapturePaymentRequest) Validate() error {
	if err := validatePaymentID(r.PaymentID); err != nil {
		return err
	}

	amount, err := r.CaptureAmount.Float64()
	if err != nil {
		return errors.New("capture_amount must be a valid number")
	}

	if amount <= 0 {
		return errors.New("capture_amount must be greater than zero")
	}

	return nil
}

// PaymentCapture describes a capture performed on a payment.
type PaymentCapture struct {
	CaptureTimestamp time.Time     `json:"capture_timestamp"`
	CaptureID        string        `json:"capture_id"`
	CaptureAmount    DecimalAmount `json:"capture_amount"`
}

// PaymentAuthorizationData contains card authorization processing details.
type PaymentAuthorizationData struct {
	AuthorizationCode             *string `json:"authorization_code,omitempty"`
	CVNVerificationResult         *string `json:"cvn_verification_result,omitempty"`
	AddressVerificationResult     *string `json:"address_verification_result,omitempty"`
	RetrievalReferenceNumber      *string `json:"retrieval_reference_number,omitempty"`
	NetworkResponseCode           *string `json:"network_response_code,omitempty"`
	NetworkResponseCodeDescriptor *string `json:"network_response_code_descriptor,omitempty"`
	NetworkTransactionID          *string `json:"network_transaction_id,omitempty"`
	AcquirerMerchantID            *string `json:"acquirer_merchant_id,omitempty"`
	ReconciliationID              *string `json:"reconciliation_id,omitempty"`
}

// PaymentAuthenticationARes contains 3DS authentication response details.
type PaymentAuthenticationARes struct {
	ECI                 *string `json:"eci,omitempty"`
	MessageVersion      *string `json:"message_version,omitempty"`
	AuthenticationValue *string `json:"authentication_value,omitempty"`
	DSTransID           *string `json:"ds_trans_id,omitempty"`
}

// PaymentAuthenticationData contains card authentication (3DS) details.
type PaymentAuthenticationData struct {
	Flow string                     `json:"flow,omitempty"`
	ARes *PaymentAuthenticationARes `json:"a_res,omitempty"`
}

// PaymentDetails contains payment information from the payment method provider.
type PaymentDetails struct {
	AuthorizationData  *PaymentAuthorizationData  `json:"authorization_data,omitempty"`
	AuthenticationData *PaymentAuthenticationData `json:"authentication_data,omitempty"`
	IssuerName         *string                    `json:"issuer_name,omitempty"`
	PayerAccountNumber *string                    `json:"payer_account_number,omitempty"`
	PayerName          *string                    `json:"payer_name,omitempty"`
	ReceiptID          *string                    `json:"receipt_id,omitempty"`
	Remark             *string                    `json:"remark,omitempty"`
	Network            *string                    `json:"network,omitempty"`
	FundSource         *string                    `json:"fund_source,omitempty"`
}

// GetPaymentResponse is the response structure for GetPayment.
type GetPaymentResponse struct {
	PaymentID         string                          `json:"payment_id"`
	BusinessID        string                          `json:"business_id"`
	ReferenceID       string                          `json:"reference_id"`
	PaymentRequestID  *string                         `json:"payment_request_id,omitempty"`
	PaymentTokenID    *string                         `json:"payment_token_id,omitempty"`
	CustomerID        *string                         `json:"customer_id,omitempty"`
	Type              PaymentType                     `json:"type"`
	Country           string                          `json:"country"`
	Currency          string                          `json:"currency"`
	RequestAmount     DecimalAmount                   `json:"request_amount"`
	CaptureMethod     *CaptureMethod                  `json:"capture_method,omitempty"`
	ChannelCode       string                          `json:"channel_code"`
	ChannelProperties PaymentRequestChannelProperties `json:"channel_properties,omitempty"`
	Captures          []PaymentCapture                `json:"captures,omitempty"`
	Status            string                          `json:"status"`
	PaymentDetails    *PaymentDetails                 `json:"payment_details,omitempty"`
	FailureCode       *string                         `json:"failure_code,omitempty"`
	Metadata          map[string]string               `json:"metadata,omitempty"`
	Created           time.Time                       `json:"created"`
	Updated           time.Time                       `json:"updated"`
}

// CancelPaymentResponse is the response structure for CancelPayment.
type CancelPaymentResponse = GetPaymentResponse

// CapturePaymentResponse is the response structure for CapturePayment.
type CapturePaymentResponse = GetPaymentResponse
