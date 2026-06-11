package xendit

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// --------------------------------- ENUMS ---------------------------------

type SessionType string

const SessionTypePay SessionType = "PAY"

type SessionMode string

const (
	SessionModePaymentLink SessionMode = "PAYMENT_LINK"
	SessionModeComponents  SessionMode = "COMPONENTS"
)

type CaptureMethod string

const (
	CaptureMethodAutomatic CaptureMethod = "AUTOMATIC"
	CaptureMethodManual    CaptureMethod = "MANUAL"
)

type CustomerType string

const CustomerTypeIndividual CustomerType = "INDIVIDUAL"

type SessionItemType string

const (
	SessionItemTypeDigitalProduct  SessionItemType = "DIGITAL_PRODUCT"
	SessionItemTypePhysicalProduct SessionItemType = "PHYSICAL_PRODUCT"
	SessionItemTypeDigitalService  SessionItemType = "DIGITAL_SERVICE"
	SessionItemTypePhysicalService SessionItemType = "PHYSICAL_SERVICE"
	SessionItemTypeFee             SessionItemType = "FEE"
)

var validCurrencies = map[string]struct{}{
	"IDR": {},
	"PHP": {},
	"VND": {},
	"THB": {},
	"SGD": {},
	"MYR": {},
	"USD": {},
}

var validCountries = map[string]struct{}{
	"ID": {},
	"PH": {},
	"VN": {},
	"TH": {},
	"SG": {},
	"MY": {},
}

// --------------------------------- REQUESTS ---------------------------------

// CreateSessionRequest is the request structure for CreateSession.
// Scoped to session_type=PAY, mode=PAYMENT_LINK or COMPONENTS.
type CreateSessionRequest struct {
	ReferenceID             string                   `json:"reference_id"`
	Currency                string                   `json:"currency"`
	Amount                  float64                  `json:"amount"`
	Country                 string                   `json:"country"`
	SessionType             SessionType              `json:"session_type"`
	Mode                    SessionMode              `json:"mode"`
	CustomerID              *string                  `json:"customer_id,omitempty"`
	Customer                *CustomerDetails         `json:"customer,omitempty"`
	CaptureMethod           *CaptureMethod           `json:"capture_method,omitempty"`
	Description             *string                  `json:"description,omitempty"`
	SuccessReturnURL        *string                  `json:"success_return_url,omitempty"`
	CancelReturnURL         *string                  `json:"cancel_return_url,omitempty"`
	Locale                  *string                  `json:"locale,omitempty"`
	Metadata                map[string]string        `json:"metadata,omitempty"`
	Items                   []SessionItem            `json:"items,omitempty"`
	AllowedPaymentChannels  []string                 `json:"allowed_payment_channels,omitempty"`
	ExpiresAt               *time.Time               `json:"expires_at,omitempty"`
	ComponentsConfiguration *ComponentsConfiguration `json:"components_configuration,omitempty"`
}

func (r *CreateSessionRequest) Validate() error {
	if strings.TrimSpace(r.ReferenceID) == "" {
		return errors.New("reference_id is required")
	}

	if strings.TrimSpace(r.Currency) == "" {
		return errors.New("currency is required")
	}

	if _, ok := validCurrencies[r.Currency]; !ok {
		return errors.New("currency must be one of IDR, PHP, VND, THB, SGD, MYR, USD")
	}

	if r.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if strings.TrimSpace(r.Country) == "" {
		return errors.New("country is required")
	}

	if _, ok := validCountries[r.Country]; !ok {
		return errors.New("country must be one of ID, PH, VN, TH, SG, MY")
	}

	if r.SessionType == "" {
		return errors.New("session_type is required")
	}

	if r.SessionType != SessionTypePay {
		return errors.New("session_type must be PAY")
	}

	if r.Mode == "" {
		return errors.New("mode is required")
	}

	switch r.Mode {
	case SessionModePaymentLink:
		if r.ComponentsConfiguration != nil {
			return errors.New("components_configuration is not allowed when mode is PAYMENT_LINK")
		}
	case SessionModeComponents:
		if r.ComponentsConfiguration == nil {
			return errors.New("components_configuration is required when mode is COMPONENTS")
		}

		if err := r.ComponentsConfiguration.Validate(); err != nil {
			return err
		}
	default:
		return errors.New("mode must be PAYMENT_LINK or COMPONENTS")
	}

	if r.CaptureMethod != nil {
		switch *r.CaptureMethod {
		case CaptureMethodAutomatic, CaptureMethodManual:
		default:
			return errors.New("capture_method must be AUTOMATIC or MANUAL")
		}
	}

	if err := validateHTTPSURL(r.SuccessReturnURL, "success_return_url"); err != nil {
		return err
	}

	if err := validateHTTPSURL(r.CancelReturnURL, "cancel_return_url"); err != nil {
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

	for i := range r.Items {
		if err := r.Items[i].Validate(); err != nil {
			return err
		}
	}

	return nil
}

// ComponentsConfiguration configures CORS origins for Xendit Component SDK sessions.
type ComponentsConfiguration struct {
	Origins []string `json:"origins"`
}

func (c *ComponentsConfiguration) Validate() error {
	if len(c.Origins) == 0 {
		return errors.New("components_configuration.origins is required")
	}

	for i, origin := range c.Origins {
		if err := validateHTTPSOrigin(origin, "components_configuration.origins"); err != nil {
			return fmt.Errorf("components_configuration.origins[%d]: %w", i, err)
		}
	}

	return nil
}

// CustomerDetails contains customer information for a payment session.
type CustomerDetails struct {
	ReferenceID      string           `json:"reference_id"`
	Type             CustomerType     `json:"type"`
	Email            *string          `json:"email,omitempty"`
	MobileNumber     *string          `json:"mobile_number,omitempty"`
	IndividualDetail IndividualDetail `json:"individual_detail"`
}

func (c *CustomerDetails) Validate() error {
	if strings.TrimSpace(c.ReferenceID) == "" {
		return errors.New("customer.reference_id is required")
	}

	if c.Type == "" {
		return errors.New("customer.type is required")
	}

	if c.Type != CustomerTypeIndividual {
		return errors.New("customer.type must be INDIVIDUAL")
	}

	return c.IndividualDetail.Validate()
}

// IndividualDetail contains individual customer name details.
type IndividualDetail struct {
	GivenNames string  `json:"given_names"`
	Surname    *string `json:"surname,omitempty"`
}

func (d *IndividualDetail) Validate() error {
	if strings.TrimSpace(d.GivenNames) == "" {
		return errors.New("customer.individual_detail.given_names is required")
	}

	return nil
}

// SessionItem describes an item attached to a payment session.
type SessionItem struct {
	ReferenceID   string            `json:"reference_id"`
	Type          SessionItemType   `json:"type"`
	Name          string            `json:"name"`
	NetUnitAmount float64           `json:"net_unit_amount"`
	Quantity      int               `json:"quantity"`
	Category      string            `json:"category"`
	Currency      *string           `json:"currency,omitempty"`
	URL           *string           `json:"url,omitempty"`
	ImageURL      *string           `json:"image_url,omitempty"`
	Subcategory   *string           `json:"subcategory,omitempty"`
	Description   *string           `json:"description,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

func (i *SessionItem) Validate() error {
	if strings.TrimSpace(i.ReferenceID) == "" {
		return errors.New("items.reference_id is required")
	}

	if i.Type == "" {
		return errors.New("items.type is required")
	}

	switch i.Type {
	case SessionItemTypeDigitalProduct, SessionItemTypePhysicalProduct,
		SessionItemTypeDigitalService, SessionItemTypePhysicalService, SessionItemTypeFee:
	default:
		return errors.New("items.type is invalid")
	}

	if strings.TrimSpace(i.Name) == "" {
		return errors.New("items.name is required")
	}

	if i.NetUnitAmount <= 0 {
		return errors.New("items.net_unit_amount must be greater than zero")
	}

	if i.Quantity < 1 {
		return errors.New("items.quantity must be at least 1")
	}

	if strings.TrimSpace(i.Category) == "" {
		return errors.New("items.category is required")
	}

	return nil
}

func validateHTTPSURL(value *string, fieldName string) error {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return errors.New(fieldName + " must be a valid URL")
	}

	if parsed.Scheme != "https" {
		return errors.New(fieldName + " must use HTTPS")
	}

	return nil
}

func validateHTTPSOrigin(value string, fieldName string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return errors.New(fieldName + " must not be empty")
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return errors.New("must be a valid URL")
	}

	if parsed.Scheme != "https" {
		return errors.New("must use HTTPS")
	}

	return nil
}

// --------------------------------- RESPONSES ---------------------------------

// CreateSessionResponse is the response structure for CreateSession.
type CreateSessionResponse struct {
	PaymentSessionID        string                   `json:"payment_session_id"`
	Created                 time.Time                `json:"created"`
	Updated                 time.Time                `json:"updated"`
	ReferenceID             string                   `json:"reference_id"`
	CustomerID              *string                  `json:"customer_id"`
	SessionType             SessionType              `json:"session_type"`
	AllowSavePaymentMethod  *string                  `json:"allow_save_payment_method"`
	Currency                string                   `json:"currency"`
	Amount                  float64                  `json:"amount"`
	Country                 string                   `json:"country"`
	Mode                    SessionMode              `json:"mode"`
	CaptureMethod           *CaptureMethod           `json:"capture_method"`
	ChannelProperties       map[string]any           `json:"channel_properties"`
	AllowedPaymentChannels  []string                 `json:"allowed_payment_channels"`
	ExpiresAt               *time.Time               `json:"expires_at"`
	Locale                  *string                  `json:"locale"`
	Metadata                map[string]string        `json:"metadata"`
	Description             *string                  `json:"description"`
	Items                   []SessionItem            `json:"items"`
	SuccessReturnURL        *string                  `json:"success_return_url"`
	CancelReturnURL         *string                  `json:"cancel_return_url"`
	Status                  string                   `json:"status"`
	PaymentLinkURL          *string                  `json:"payment_link_url"`
	PaymentTokenID          *string                  `json:"payment_token_id"`
	PaymentID               *string                  `json:"payment_id"`
	PaymentRequestID        *string                  `json:"payment_request_id"`
	BusinessID              string                   `json:"business_id"`
	ComponentsSDKKey        *string                  `json:"components_sdk_key"`
	ComponentsConfiguration *ComponentsConfiguration `json:"components_configuration"`
}
