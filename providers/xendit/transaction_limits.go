package xendit

import (
	"errors"
	"fmt"
	"strings"
)

// PaymentMethodType groups Xendit payment channels by checkout category.
type PaymentMethodType string

const (
	PaymentMethodTypePayLater         PaymentMethodType = "PAYLATER"
	PaymentMethodTypeOverTheCounter   PaymentMethodType = "OVER_THE_COUNTER"
	PaymentMethodTypeEWallet          PaymentMethodType = "EWALLET"
	PaymentMethodTypeBankTransfer     PaymentMethodType = "BANK_TRANSFER"
	PaymentMethodTypeOnlineBanking    PaymentMethodType = "ONLINE_BANKING"
	PaymentMethodTypeCards            PaymentMethodType = "CARDS"
	PaymentMethodTypeQRCode           PaymentMethodType = "QR_CODE"
)

var (
	ErrMinimumTransactionNotMet   = errors.New("minimum transaction amount not met")
	ErrMaximumTransactionExceeded = errors.New("maximum transaction amount exceeded")
)

// TransactionLimit defines optional min/max bounds for a payment method type.
type TransactionLimit struct {
	Min *int64
	Max *int64
}

// idrLimitsByPaymentMethodType holds strictest Xendit IDR limits per payment type.
// Sources: https://docs.xendit.co/docs/available-payment-channels
var idrLimitsByPaymentMethodType = map[PaymentMethodType]TransactionLimit{
	PaymentMethodTypeQRCode: {
		Min: int64Ptr(1),
		Max: int64Ptr(10_000_000),
	},
	PaymentMethodTypeBankTransfer: {
		Min: int64Ptr(10_000),
		Max: int64Ptr(50_000_000),
	},
	PaymentMethodTypeEWallet: {
		Min: int64Ptr(100),
		Max: int64Ptr(20_000_000),
	},
	PaymentMethodTypeOverTheCounter: {
		Min: int64Ptr(50_000),
		Max: int64Ptr(5_000_000),
	},
	PaymentMethodTypeCards: {
		Min: int64Ptr(5_000),
		Max: int64Ptr(200_000_000),
	},
	PaymentMethodTypeOnlineBanking: {
		Min: int64Ptr(1_000),
		Max: int64Ptr(20_000_000),
	},
}

func int64Ptr(v int64) *int64 {
	return &v
}

// ChannelCodeToPaymentMethodType maps a Xendit channel code to a payment method type.
func ChannelCodeToPaymentMethodType(channelCode string) (PaymentMethodType, error) {
	channelCode = strings.ToUpper(strings.TrimSpace(channelCode))
	if channelCode == "" {
		return "", fmt.Errorf("channel_code is required")
	}

	switch channelCode {
	case ChannelCodeQRIS:
		return PaymentMethodTypeQRCode, nil
	case ChannelCodeCards:
		return PaymentMethodTypeCards, nil
	case ChannelCodeAlfamart, ChannelCodeIndomaret:
		return PaymentMethodTypeOverTheCounter, nil
	case ChannelCodeJeniusPay:
		return PaymentMethodTypeOnlineBanking, nil
	case ChannelCodeOVO, ChannelCodeDANA, ChannelCodeLinkAja, ChannelCodeShopeePay,
		ChannelCodeGoPay, ChannelCodeAstraPay, ChannelCodeNexCash:
		return PaymentMethodTypeEWallet, nil
	default:
		if _, ok := virtualAccountChannelCodes[channelCode]; ok {
			return PaymentMethodTypeBankTransfer, nil
		}

		return "", fmt.Errorf("unsupported channel_code %q for transaction limit validation", channelCode)
	}
}

// ValidateTransactionAmount checks amount against Xendit provider limits for the channel's payment type.
// Non-IDR currencies are skipped. Unknown or unconfigured types are skipped.
func ValidateTransactionAmount(channelCode, currency string, amount float64) error {
	if strings.ToUpper(strings.TrimSpace(currency)) != "IDR" {
		return nil
	}

	paymentMethodType, err := ChannelCodeToPaymentMethodType(channelCode)
	if err != nil {
		return nil
	}

	limits, ok := idrLimitsByPaymentMethodType[paymentMethodType]
	if !ok {
		return nil
	}

	amountInt := int64(amount)
	if limits.Min != nil && amountInt < *limits.Min {
		return ErrMinimumTransactionNotMet
	}

	if limits.Max != nil && amountInt > *limits.Max {
		return ErrMaximumTransactionExceeded
	}

	return nil
}

// IsMinimumTransactionNotMet reports whether err is a client-side minimum amount validation failure.
func IsMinimumTransactionNotMet(err error) bool {
	if errors.Is(err, ErrMinimumTransactionNotMet) {
		return true
	}

	xerr, ok := err.(*XenditError)
	if !ok {
		return false
	}

	return xerr.ErrorCode == InvalidRequest && strings.Contains(xerr.Message, ErrMinimumTransactionNotMet.Error())
}

// IsMaximumTransactionExceeded reports whether err is a client-side maximum amount validation failure.
func IsMaximumTransactionExceeded(err error) bool {
	if errors.Is(err, ErrMaximumTransactionExceeded) {
		return true
	}

	xerr, ok := err.(*XenditError)
	if !ok {
		return false
	}

	return xerr.ErrorCode == InvalidRequest && strings.Contains(xerr.Message, ErrMaximumTransactionExceeded.Error())
}
