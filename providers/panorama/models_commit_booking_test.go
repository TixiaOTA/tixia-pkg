package panorama

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBookingFeeDescriptionAndMetadata(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		description string
		items       []string
	}{
		{name: "string", payload: `{"fee_type":"SSR","description":"5 Kg Baggage"}`, description: "5 Kg Baggage"},
		{name: "array", payload: `{"fee_type":"SSR","description":["5 Kg Baggage","Outbound"]}`, description: "5 Kg Baggage, Outbound", items: []string{"5 Kg Baggage", "Outbound"}},
		{name: "empty array", payload: `{"fee_type":"SSR","description":[]}`, items: []string{}},
		{name: "null", payload: `{"fee_type":"SSR","description":null}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fee BookingFee
			require.NoError(t, json.Unmarshal([]byte(tt.payload), &fee))
			require.Equal(t, tt.description, fee.Description)
			require.Equal(t, tt.items, fee.DescriptionItems)
		})
	}
}

func TestBookingFeeMetadataRoundTrip(t *testing.T) {
	var fee BookingFee
	require.NoError(t, json.Unmarshal([]byte(`{
		"fee_type":"SSR","fee_code":"XBAG_1","fee_name":"5 Kg Baggage",
		"fee_value":"5","fee_category":"baggage","amount":125000,"currency":"IDR",
		"service_charges":[{"amount":125000,"currency":"IDR","charge_code":"BASE"}],
		"journey_code":"J1","legs":[{"leg_code":"L1","origin":"CGK","destination":"DPS"}],"sequence":1
	}`), &fee))
	require.Equal(t, "XBAG_1", fee.FeeCode)
	require.Equal(t, "5 Kg Baggage", fee.FeeName)
	require.Equal(t, "5", fee.FeeValue)
	require.Equal(t, "baggage", fee.FeeCategory)
	require.Len(t, fee.ServiceCharges, 1)
	require.Len(t, fee.Legs, 1)
	require.Equal(t, "J1", fee.JourneyCode)

	encoded, err := json.Marshal(fee)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), "description_items")
}

func TestBookingFeeAmountAcceptsDecimalNumber(t *testing.T) {
	var fee BookingFee
	require.NoError(t, json.Unmarshal([]byte(`{"fee_type":"SEAT","amount":1521000.0,"currency":"IDR"}`), &fee))
	require.Equal(t, int64(1521000), fee.Amount.Int64())
}

func TestBookingFeeMinimalPayload(t *testing.T) {
	var fee BookingFee
	require.NoError(t, json.Unmarshal([]byte(`{"fee_type":"SEAT","amount":30000,"currency":"IDR"}`), &fee))
	require.Equal(t, BookingFeeTypeSEAT, fee.FeeType)
}
