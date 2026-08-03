package xendit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannelCodeToPaymentMethodType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		channelCode string
		want        PaymentMethodType
		wantErr     bool
	}{
		{channelCode: ChannelCodeQRIS, want: PaymentMethodTypeQRCode},
		{channelCode: ChannelCodeBCAVirtualAccount, want: PaymentMethodTypeBankTransfer},
		{channelCode: ChannelCodeOVO, want: PaymentMethodTypeEWallet},
		{channelCode: ChannelCodeAlfamart, want: PaymentMethodTypeOverTheCounter},
		{channelCode: ChannelCodeCards, want: PaymentMethodTypeCards},
		{channelCode: ChannelCodeJeniusPay, want: PaymentMethodTypeOnlineBanking},
		{channelCode: "", wantErr: true},
		{channelCode: "UNKNOWN_CHANNEL", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.channelCode, func(t *testing.T) {
			t.Parallel()

			got, err := ChannelCodeToPaymentMethodType(tt.channelCode)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateTransactionAmount(t *testing.T) {
	t.Parallel()

	t.Run("QR_CODE within limits", func(t *testing.T) {
		t.Parallel()
		err := ValidateTransactionAmount(ChannelCodeQRIS, "IDR", 1_000_000)
		assert.NoError(t, err)
	})

	t.Run("QR_CODE below minimum", func(t *testing.T) {
		t.Parallel()
		err := ValidateTransactionAmount(ChannelCodeQRIS, "IDR", 0)
		assert.ErrorIs(t, err, ErrMinimumTransactionNotMet)
	})

	t.Run("QR_CODE above maximum", func(t *testing.T) {
		t.Parallel()
		err := ValidateTransactionAmount(ChannelCodeQRIS, "IDR", 20_000_000)
		assert.ErrorIs(t, err, ErrMaximumTransactionExceeded)
	})

	t.Run("BANK_TRANSFER below minimum", func(t *testing.T) {
		t.Parallel()
		err := ValidateTransactionAmount(ChannelCodeBCAVirtualAccount, "IDR", 5_000)
		assert.ErrorIs(t, err, ErrMinimumTransactionNotMet)
	})

	t.Run("EWALLET below minimum", func(t *testing.T) {
		t.Parallel()
		err := ValidateTransactionAmount(ChannelCodeOVO, "IDR", 50)
		assert.ErrorIs(t, err, ErrMinimumTransactionNotMet)
	})

	t.Run("non-IDR currency skipped", func(t *testing.T) {
		t.Parallel()
		err := ValidateTransactionAmount(ChannelCodeQRIS, "PHP", 1)
		assert.NoError(t, err)
	})

	t.Run("unknown channel skipped", func(t *testing.T) {
		t.Parallel()
		err := ValidateTransactionAmount("UNKNOWN", "IDR", 1_000_000)
		assert.NoError(t, err)
	})
}

func TestIsMinimumTransactionNotMet(t *testing.T) {
	t.Parallel()

	assert.True(t, IsMinimumTransactionNotMet(ErrMinimumTransactionNotMet))
	assert.True(t, IsMinimumTransactionNotMet(&XenditError{
		ErrorCode: InvalidRequest,
		Message:   ErrMinimumTransactionNotMet.Error(),
	}))
	assert.False(t, IsMinimumTransactionNotMet(&XenditError{
		ErrorCode: InvalidRequest,
		Message:   "reference_id is required",
	}))
	assert.False(t, IsMinimumTransactionNotMet(errors.New("other error")))
}

func TestIsMaximumTransactionExceeded(t *testing.T) {
	t.Parallel()

	assert.True(t, IsMaximumTransactionExceeded(ErrMaximumTransactionExceeded))
	assert.True(t, IsMaximumTransactionExceeded(&XenditError{
		ErrorCode: InvalidRequest,
		Message:   ErrMaximumTransactionExceeded.Error(),
	}))
	assert.False(t, IsMaximumTransactionExceeded(&XenditError{
		ErrorCode: InvalidRequest,
		Message:   "reference_id is required",
	}))
}

func TestCreatePaymentRequestRequest_Validate_TransactionLimits(t *testing.T) {
	t.Parallel()

	t.Run("below minimum returns InvalidRequest from client", func(t *testing.T) {
		t.Parallel()

		req := validCreatePaymentRequestRequest()
		req.RequestAmount = NewDecimalAmountFromFloat(50)

		err := req.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMinimumTransactionNotMet)
	})

	t.Run("above maximum returns error", func(t *testing.T) {
		t.Parallel()

		req := validCreatePaymentRequestRequest()
		req.RequestAmount = NewDecimalAmountFromFloat(25_000_000)

		err := req.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMaximumTransactionExceeded)
	})
}

func TestCreateSessionRequest_Validate_TransactionLimits(t *testing.T) {
	t.Parallel()

	t.Run("IDR session requires allowed payment channels", func(t *testing.T) {
		t.Parallel()

		req := validPaymentLinkRequest()
		req.AllowedPaymentChannels = nil

		err := req.Validate()
		require.Error(t, err)
		assert.Equal(t, "allowed_payment_channels is required for IDR sessions", err.Error())
	})

	t.Run("IDR session above maximum", func(t *testing.T) {
		t.Parallel()

		req := validPaymentLinkRequest()
		req.Amount = 20_000_000

		err := req.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMaximumTransactionExceeded)
	})
}
