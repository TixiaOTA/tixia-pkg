package panorama

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterFeesByType(t *testing.T) {
	fees := []BookingFee{
		{FeeType: BookingFeeTypeSSR, Amount: 50000, Currency: "IDR", Description: "Baggage 20kg"},
		{FeeType: BookingFeeTypeSEAT, Amount: 30000, Currency: "IDR", Description: "Window seat"},
		{FeeType: BookingFeeTypeTax, Amount: 100000, Currency: "IDR"},
		{FeeType: BookingFeeTypeSSR, Amount: 25000, Currency: "IDR", Description: "Meal"},
		{FeeType: BookingFeeTypeFare, Amount: 1500000, Currency: "IDR"},
	}

	t.Run("filter SSR fees", func(t *testing.T) {
		result := FilterFeesByType(fees, BookingFeeTypeSSR)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(50000), result[0].Amount)
		assert.Equal(t, int64(25000), result[1].Amount)
	})

	t.Run("filter SEAT fees", func(t *testing.T) {
		result := FilterFeesByType(fees, BookingFeeTypeSEAT)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(30000), result[0].Amount)
	})

	t.Run("filter TAX fees", func(t *testing.T) {
		result := FilterFeesByType(fees, BookingFeeTypeTax)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(100000), result[0].Amount)
	})

	t.Run("filter FARE fees", func(t *testing.T) {
		result := FilterFeesByType(fees, BookingFeeTypeFare)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1500000), result[0].Amount)
	})

	t.Run("no matches returns nil", func(t *testing.T) {
		result := FilterFeesByType(fees, "UNKNOWN")
		assert.Nil(t, result)
	})

	t.Run("empty fees returns nil", func(t *testing.T) {
		result := FilterFeesByType([]BookingFee{}, BookingFeeTypeSSR)
		assert.Nil(t, result)
	})

	t.Run("all match", func(t *testing.T) {
		sameFees := []BookingFee{
			{FeeType: BookingFeeTypeSSR, Amount: 100},
			{FeeType: BookingFeeTypeSSR, Amount: 200},
		}
		result := FilterFeesByType(sameFees, BookingFeeTypeSSR)
		assert.Len(t, result, 2)
	})
}
