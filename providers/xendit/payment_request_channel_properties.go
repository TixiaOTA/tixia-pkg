package xendit

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// IDR PAY channel codes supported by CreatePaymentRequest.
const (
	ChannelCodeAlfamart               = "ALFAMART"
	ChannelCodeAstraPay               = "ASTRAPAY"
	ChannelCodeBCAVirtualAccount      = "BCA_VIRTUAL_ACCOUNT"
	ChannelCodeBJBVirtualAccount      = "BJB_VIRTUAL_ACCOUNT"
	ChannelCodeBNCVirtualAccount      = "BNC_VIRTUAL_ACCOUNT"
	ChannelCodeBNIVirtualAccount      = "BNI_VIRTUAL_ACCOUNT"
	ChannelCodeBRIVirtualAccount      = "BRI_VIRTUAL_ACCOUNT"
	ChannelCodeBSSVirtualAccount      = "BSS_VIRTUAL_ACCOUNT"
	ChannelCodeBSIVirtualAccount      = "BSI_VIRTUAL_ACCOUNT"
	ChannelCodeCards                  = "CARDS"
	ChannelCodeCIMBVirtualAccount     = "CIMB_VIRTUAL_ACCOUNT"
	ChannelCodeDANA                   = "DANA"
	ChannelCodeGoPay                  = "GOPAY"
	ChannelCodeINAVirtualAccount      = "INA_VIRTUAL_ACCOUNT"
	ChannelCodeIndomaret              = "INDOMARET"
	ChannelCodeJeniusPay              = "JENIUSPAY"
	ChannelCodeLinkAja                = "LINKAJA"
	ChannelCodeMandiriVirtualAccount  = "MANDIRI_VIRTUAL_ACCOUNT"
	ChannelCodeMuamalatVirtualAccount = "MUAMALAT_VIRTUAL_ACCOUNT"
	ChannelCodeNexCash                = "NEXCASH"
	ChannelCodeOVO                    = "OVO"
	ChannelCodePermataVirtualAccount  = "PERMATA_VIRTUAL_ACCOUNT"
	ChannelCodeQRIS                   = "QRIS"
	ChannelCodeShopeePay              = "SHOPEEPAY"
)

var (
	mobileNumberPattern = regexp.MustCompile(`^(\+|)[0-9]\d{1,22}$`)
	cashtagPattern      = regexp.MustCompile(`^[$][a-zA-Z0-9_]{3,15}$`)
	displayNamePattern  = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9 .,']*$`)
	paymentCodePattern  = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	vaNumberPattern     = regexp.MustCompile(`^\d+$`)
)

var virtualAccountChannelCodes = map[string]struct{}{
	ChannelCodeBCAVirtualAccount:      {},
	ChannelCodeBJBVirtualAccount:      {},
	ChannelCodeBNCVirtualAccount:      {},
	ChannelCodeBNIVirtualAccount:      {},
	ChannelCodeBRIVirtualAccount:      {},
	ChannelCodeBSSVirtualAccount:      {},
	ChannelCodeBSIVirtualAccount:      {},
	ChannelCodeCIMBVirtualAccount:     {},
	ChannelCodeINAVirtualAccount:      {},
	ChannelCodeMandiriVirtualAccount:  {},
	ChannelCodeMuamalatVirtualAccount: {},
	ChannelCodePermataVirtualAccount:  {},
}

// PaymentRequestChannelProperties wraps typed channel_properties for IDR PAY requests.
type PaymentRequestChannelProperties struct {
	value any
}

func (p PaymentRequestChannelProperties) MarshalJSON() ([]byte, error) {
	if p.value == nil {
		return []byte("null"), nil
	}

	return json.Marshal(p.value)
}

func (p PaymentRequestChannelProperties) IsZero() bool {
	return p.value == nil
}

func NewOVOChannelProperties(props OVOChannelProperties) PaymentRequestChannelProperties {
	return PaymentRequestChannelProperties{value: props}
}

func NewSingleRedirectChannelProperties(props SingleRedirectChannelProperties) PaymentRequestChannelProperties {
	return PaymentRequestChannelProperties{value: props}
}

func NewDualRedirectChannelProperties(props DualRedirectChannelProperties) PaymentRequestChannelProperties {
	return PaymentRequestChannelProperties{value: props}
}

func NewTripleRedirectChannelProperties(props TripleRedirectChannelProperties) PaymentRequestChannelProperties {
	return PaymentRequestChannelProperties{value: props}
}

func NewQRISChannelProperties(props QRISChannelProperties) PaymentRequestChannelProperties {
	return PaymentRequestChannelProperties{value: props}
}

func NewCardsChannelProperties(props CardsChannelProperties) PaymentRequestChannelProperties {
	return PaymentRequestChannelProperties{value: props}
}

func NewVirtualAccountChannelProperties(props VirtualAccountChannelProperties) PaymentRequestChannelProperties {
	return PaymentRequestChannelProperties{value: props}
}

func NewRetailOutletChannelProperties(props RetailOutletChannelProperties) PaymentRequestChannelProperties {
	return PaymentRequestChannelProperties{value: props}
}

func NewJeniusPayChannelProperties(props JeniusPayChannelProperties) PaymentRequestChannelProperties {
	return PaymentRequestChannelProperties{value: props}
}

// OVOChannelProperties is the channel_properties payload for OVO payments.
type OVOChannelProperties struct {
	AccountMobileNumber string `json:"account_mobile_number"`
}

func (p OVOChannelProperties) Validate() error {
	mobile := strings.TrimSpace(p.AccountMobileNumber)
	if mobile == "" {
		return errors.New("channel_properties.account_mobile_number is required for OVO")
	}

	if !mobileNumberPattern.MatchString(mobile) {
		return errors.New("channel_properties.account_mobile_number must be a valid mobile number")
	}

	return nil
}

// SingleRedirectChannelProperties is used by DANA, LINKAJA, and SHOPEEPAY.
type SingleRedirectChannelProperties struct {
	SuccessReturnURL string `json:"success_return_url"`
}

func (p SingleRedirectChannelProperties) Validate() error {
	return validateChannelHTTPSURL(p.SuccessReturnURL, "channel_properties.success_return_url")
}

// DualRedirectChannelProperties is used by ASTRAPAY and NEXCASH.
type DualRedirectChannelProperties struct {
	SuccessReturnURL string `json:"success_return_url"`
	FailureReturnURL string `json:"failure_return_url"`
}

func (p DualRedirectChannelProperties) Validate() error {
	if err := validateChannelHTTPSURL(p.SuccessReturnURL, "channel_properties.success_return_url"); err != nil {
		return err
	}

	return validateChannelHTTPSURL(p.FailureReturnURL, "channel_properties.failure_return_url")
}

// TripleRedirectChannelProperties is used by GOPAY.
type TripleRedirectChannelProperties struct {
	SuccessReturnURL string `json:"success_return_url"`
	FailureReturnURL string `json:"failure_return_url"`
	CancelReturnURL  string `json:"cancel_return_url"`
}

func (p TripleRedirectChannelProperties) Validate() error {
	if err := validateChannelHTTPSURL(p.SuccessReturnURL, "channel_properties.success_return_url"); err != nil {
		return err
	}

	if err := validateChannelHTTPSURL(p.FailureReturnURL, "channel_properties.failure_return_url"); err != nil {
		return err
	}

	return validateChannelHTTPSURL(p.CancelReturnURL, "channel_properties.cancel_return_url")
}

// QRISChannelProperties is the channel_properties payload for QRIS payments.
type QRISChannelProperties struct {
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func (p QRISChannelProperties) Validate() error {
	return nil
}

// PaymentRequestCardDetails contains card input for CARDS channel_properties.
type PaymentRequestCardDetails struct {
	CardNumber            string  `json:"card_number,omitempty"`
	ExpiryYear            string  `json:"expiry_year,omitempty"`
	ExpiryMonth           string  `json:"expiry_month,omitempty"`
	CardholderFirstName   string  `json:"cardholder_first_name,omitempty"`
	CardholderLastName    string  `json:"cardholder_last_name,omitempty"`
	CardholderEmail       *string `json:"cardholder_email,omitempty"`
	CardholderPhoneNumber *string `json:"cardholder_phone_number,omitempty"`
	CVN                   *string `json:"cvn,omitempty"`
}

// PaymentRequestBillingInformation contains billing details for CARDS channel_properties.
type PaymentRequestBillingInformation struct {
	FirstName     string  `json:"first_name,omitempty"`
	LastName      string  `json:"last_name,omitempty"`
	Email         string  `json:"email,omitempty"`
	PhoneNumber   *string `json:"phone_number,omitempty"`
	City          *string `json:"city,omitempty"`
	Country       *string `json:"country,omitempty"`
	PostalCode    *string `json:"postal_code,omitempty"`
	StreetLine1   *string `json:"street_line1,omitempty"`
	StreetLine2   *string `json:"street_line2,omitempty"`
	ProvinceState *string `json:"province_state,omitempty"`
}

// PaymentRequestInstallmentConfiguration contains installment options for CARDS channel_properties.
type PaymentRequestInstallmentConfiguration struct {
	Terms    *int    `json:"terms,omitempty"`
	Interval *string `json:"interval,omitempty"`
	Code     *string `json:"code,omitempty"`
}

// PaymentRequestRecurringConfiguration contains recurring options for CARDS channel_properties.
type PaymentRequestRecurringConfiguration struct {
	RecurringExpiry    *string `json:"recurring_expiry,omitempty"`
	RecurringFrequency *string `json:"recurring_frequency,omitempty"`
}

// CardsChannelProperties is the channel_properties payload for CARDS payments.
type CardsChannelProperties struct {
	CardDetails              *PaymentRequestCardDetails              `json:"card_details,omitempty"`
	SkipThreeDS              *bool                                   `json:"skip_three_ds,omitempty"`
	MidLabel                 *string                                 `json:"mid_label,omitempty"`
	CardOnFileType           *string                                 `json:"card_on_file_type,omitempty"`
	BillingInformation       *PaymentRequestBillingInformation       `json:"billing_information,omitempty"`
	StatementDescriptor      *string                                 `json:"statement_descriptor,omitempty"`
	TransactionSequence      *string                                 `json:"transaction_sequence,omitempty"`
	NetworkTransactionID     *string                                 `json:"network_transaction_id,omitempty"`
	InstallmentConfiguration *PaymentRequestInstallmentConfiguration `json:"installment_configuration,omitempty"`
	RecurringConfiguration   *PaymentRequestRecurringConfiguration   `json:"recurring_configuration,omitempty"`
	SuccessReturnURL         *string                                 `json:"success_return_url,omitempty"`
	FailureReturnURL         *string                                 `json:"failure_return_url,omitempty"`
}

func (p CardsChannelProperties) Validate() error {
	if err := validateOptionalChannelHTTPSURL(p.SuccessReturnURL, "channel_properties.success_return_url"); err != nil {
		return err
	}

	return validateOptionalChannelHTTPSURL(p.FailureReturnURL, "channel_properties.failure_return_url")
}

// PaymentRequestVerificationData contains payer verification hints for virtual account channels.
type PaymentRequestVerificationData struct {
	CustomerName           *string  `json:"customer_name,omitempty"`
	AcceptedNameVariations []string `json:"accepted_name_variations,omitempty"`
	AllowedBankAccounts    []string `json:"allowed_bank_accounts,omitempty"`
}

// VirtualAccountChannelProperties is used by Indonesia virtual account channels.
type VirtualAccountChannelProperties struct {
	DisplayName          string                          `json:"display_name"`
	ExpiresAt            *time.Time                      `json:"expires_at,omitempty"`
	VirtualAccountNumber *string                         `json:"virtual_account_number,omitempty"`
	Description          *string                         `json:"description,omitempty"`
	VerificationData     *PaymentRequestVerificationData `json:"verification_data,omitempty"`
	SuccessReturnURL     *string                         `json:"success_return_url,omitempty"`
}

func (p VirtualAccountChannelProperties) Validate() error {
	displayName := strings.TrimSpace(p.DisplayName)
	if displayName == "" {
		return errors.New("channel_properties.display_name is required for virtual account channels")
	}

	if len(displayName) > 255 || !displayNamePattern.MatchString(displayName) {
		return errors.New("channel_properties.display_name must be 1-255 alphanumeric characters")
	}

	if p.VirtualAccountNumber != nil {
		vaNumber := strings.TrimSpace(*p.VirtualAccountNumber)
		if vaNumber == "" {
			return errors.New("channel_properties.virtual_account_number cannot be empty")
		}

		if len(vaNumber) > 32 || !vaNumberPattern.MatchString(vaNumber) {
			return errors.New("channel_properties.virtual_account_number must be numeric and at most 32 digits")
		}
	}

	return validateOptionalChannelHTTPSURL(p.SuccessReturnURL, "channel_properties.success_return_url")
}

// RetailOutletChannelProperties is used by ALFAMART and INDOMARET.
type RetailOutletChannelProperties struct {
	PayerName   string     `json:"payer_name"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	PaymentCode *string    `json:"payment_code,omitempty"`
}

func (p RetailOutletChannelProperties) Validate() error {
	payerName := strings.TrimSpace(p.PayerName)
	if payerName == "" {
		return errors.New("channel_properties.payer_name is required for retail outlet channels")
	}

	if len(payerName) > 250 || !displayNamePattern.MatchString(payerName) {
		return errors.New("channel_properties.payer_name must be 1-250 alphanumeric characters")
	}

	if p.PaymentCode != nil {
		code := strings.TrimSpace(*p.PaymentCode)
		if code == "" {
			return errors.New("channel_properties.payment_code cannot be empty")
		}

		if len(code) > 14 || !paymentCodePattern.MatchString(code) {
			return errors.New("channel_properties.payment_code must be alphanumeric and at most 14 characters")
		}
	}

	return nil
}

// JeniusPayChannelProperties is the channel_properties payload for JENIUSPAY payments.
type JeniusPayChannelProperties struct {
	Cashtag string `json:"cashtag"`
}

func (p JeniusPayChannelProperties) Validate() error {
	cashtag := strings.TrimSpace(p.Cashtag)
	if cashtag == "" {
		return errors.New("channel_properties.cashtag is required for JENIUSPAY")
	}

	if !cashtagPattern.MatchString(cashtag) {
		return errors.New("channel_properties.cashtag must match $[a-zA-Z0-9_]{3,15}")
	}

	return nil
}

func validateChannelProperties(channelCode string, properties PaymentRequestChannelProperties) error {
	if properties.IsZero() {
		return errors.New("channel_properties is required")
	}

	switch channelCode {
	case ChannelCodeOVO:
		return assertChannelPropertiesType[OVOChannelProperties](channelCode, properties)
	case ChannelCodeDANA, ChannelCodeLinkAja, ChannelCodeShopeePay:
		return assertChannelPropertiesType[SingleRedirectChannelProperties](channelCode, properties)
	case ChannelCodeAstraPay, ChannelCodeNexCash:
		return assertChannelPropertiesType[DualRedirectChannelProperties](channelCode, properties)
	case ChannelCodeGoPay:
		return assertChannelPropertiesType[TripleRedirectChannelProperties](channelCode, properties)
	case ChannelCodeQRIS:
		return assertChannelPropertiesType[QRISChannelProperties](channelCode, properties)
	case ChannelCodeCards:
		return assertChannelPropertiesType[CardsChannelProperties](channelCode, properties)
	case ChannelCodeAlfamart, ChannelCodeIndomaret:
		return assertChannelPropertiesType[RetailOutletChannelProperties](channelCode, properties)
	case ChannelCodeJeniusPay:
		return assertChannelPropertiesType[JeniusPayChannelProperties](channelCode, properties)
	default:
		if _, ok := virtualAccountChannelCodes[channelCode]; ok {
			return assertChannelPropertiesType[VirtualAccountChannelProperties](channelCode, properties)
		}

		return fmt.Errorf("unsupported channel_code %q for IDR PAY payment requests", channelCode)
	}
}

type channelPropertiesValidator interface {
	Validate() error
}

func assertChannelPropertiesType[T channelPropertiesValidator](channelCode string, properties PaymentRequestChannelProperties) error {
	props, ok := properties.value.(T)
	if !ok {
		return fmt.Errorf("channel_properties type mismatch for channel_code %q", channelCode)
	}

	return props.Validate()
}

func validateChannelHTTPSURL(value string, fieldName string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fmt.Errorf("%s is required", fieldName)
	}

	return validateHTTPSURL(&trimmed, fieldName)
}

func validateOptionalChannelHTTPSURL(value *string, fieldName string) error {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return validateHTTPSURL(&trimmed, fieldName)
}
