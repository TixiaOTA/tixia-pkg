package panorama

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefundRequestValidate(t *testing.T) {
	providerBookings := []RefundProviderBooking{{
		Pnr:      "ABC123",
		Journeys: []RefundJourney{{Origin: "CGK", Destination: "DPS", DepartureDate: "2026-10-01", Pax: []int{1}}},
	}}
	passengers := []RefundPassengerNumber{{PassengerNumber: 1}}

	full := GetRefundBookingRequest{
		Signature:        "signature",
		OrderNumber:      "order-1",
		Passengers:       passengers,
		ProviderBookings: providerBookings,
		Captcha:          map[string]string{"ABC123": "Refund"},
	}
	require.NoError(t, full.Validate())
	require.ErrorContains(t, (&GetRefundBookingRequest{}).Validate(), "signature is required")
	require.ErrorContains(t, (&GetRefundBookingRequest{Signature: "signature"}).Validate(), "order_number is required")
	require.ErrorContains(t, (&GetRefundBookingRequest{Signature: "signature", OrderNumber: "order-1"}).Validate(), "passengers is required")
	require.ErrorContains(t, (&GetRefundBookingRequest{Signature: "signature", OrderNumber: "order-1", Passengers: passengers}).Validate(), "provider_bookings is required")
	require.ErrorContains(t, (&GetRefundBookingRequest{Signature: "signature", OrderNumber: "order-1", Passengers: passengers, ProviderBookings: providerBookings}).Validate(), "captcha is required")

	cancelFull := CancelRefundBookingRequest{Signature: "signature", OrderNumber: "order-1", ProviderBookings: providerBookings}
	require.NoError(t, cancelFull.Validate())
	require.ErrorContains(t, (&CancelRefundBookingRequest{}).Validate(), "signature is required")
	require.ErrorContains(t, (&CancelRefundBookingRequest{Signature: "signature"}).Validate(), "order_number is required")
	require.ErrorContains(t, (&CancelRefundBookingRequest{Signature: "signature", OrderNumber: "order-1"}).Validate(), "provider_bookings is required")

	preRefundFull := PreRefundLoginListRequest{Signature: "signature", OrderNumber: "order-1", Provider: []string{"sia"}, Pnr: []string{"ABC123"}}
	require.NoError(t, preRefundFull.Validate())
	require.ErrorContains(t, (&PreRefundLoginListRequest{}).Validate(), "signature is required")
	require.ErrorContains(t, (&PreRefundLoginListRequest{Signature: "signature"}).Validate(), "order_number is required")
}

func TestGetRefundBookingRequest_Marshal(t *testing.T) {
	req := GetRefundBookingRequest{
		Signature:        "sig-123",
		OrderNumber:      "order-1",
		Passengers:       []RefundPassengerNumber{{PassengerNumber: 1}, {PassengerNumber: 2}},
		ProviderBookings: []RefundProviderBooking{{Pnr: "ABC123", Journeys: []RefundJourney{{Origin: "CGK", Destination: "DPS", DepartureDate: "2026-10-01", Pax: []int{1, 2}}}}},
		Captcha:          map[string]string{"ABC123": "Refund"},
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(body, &decoded))
	assert.Equal(t, "sig-123", decoded["signature"])
	assert.Equal(t, "order-1", decoded["order_number"])
	assert.Equal(t, map[string]any{"ABC123": "Refund"}, decoded["captcha"])

	passengers, ok := decoded["passengers"].([]any)
	require.True(t, ok)
	assert.Equal(t, map[string]any{"passenger_number": float64(1)}, passengers[0])

	providerBookings, ok := decoded["provider_bookings"].([]any)
	require.True(t, ok)
	assert.Equal(t, map[string]any{
		"pnr": "ABC123",
		"journeys": []any{map[string]any{
			"origin":         "CGK",
			"destination":    "DPS",
			"departure_date": "2026-10-01",
			"pax":            []any{1.0, 2.0},
		}},
	}, providerBookings[0])
}

func TestPreRefundLoginListRequest_Marshal(t *testing.T) {
	req := PreRefundLoginListRequest{
		Signature:   "sig-123",
		OrderNumber: "order-1",
		Provider:    []string{"sia"},
		Pnr:         []string{"ABC123"},
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)
	assert.JSONEq(t, `{"signature":"sig-123","order_number":"order-1","provider":["sia"],"pnr":["ABC123"]}`, string(body))

	// Provider and Pnr are optional and omitted when empty
	body, err = json.Marshal(PreRefundLoginListRequest{Signature: "sig-123", OrderNumber: "order-1"})
	require.NoError(t, err)
	assert.JSONEq(t, `{"signature":"sig-123","order_number":"order-1"}`, string(body))
}

func TestCancelRefundBookingRequest_Marshal(t *testing.T) {
	req := CancelRefundBookingRequest{
		Signature:        "sig-123",
		OrderNumber:      "order-1",
		ProviderBookings: []RefundProviderBooking{{Pnr: "ABC123", Journeys: []RefundJourney{{Origin: "CGK", Destination: "DPS", DepartureDate: "2026-10-01", Pax: []int{1}}}}},
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(body, &decoded))
	assert.Equal(t, "sig-123", decoded["signature"])
	assert.Equal(t, "order-1", decoded["order_number"])
	providerBookings, ok := decoded["provider_bookings"].([]any)
	require.True(t, ok)
	assert.Equal(t, "ABC123", providerBookings[0].(map[string]any)["pnr"])
}

func TestGetRefundBookingResponse_Parse(t *testing.T) {
	const payload = `{
		"result": {
			"error_code": 0,
			"error_msg": "",
			"response": {
				"order_number": "order-1",
				"pnr": "XYZ789",
				"state": "issued",
				"departure_date": "2026-10-01 08:00:00",
				"arrival_date": "2026-10-01 10:00:00",
				"provider_type": "airline",
				"provider_bookings": [
					{
						"pnr": "ABC123",
						"provider": "sia",
						"total_price_cost": "2500000.50",
						"penalty_amount": "250000",
						"status": "confirm",
						"refundable_amount": 2250000.5,
						"admin_fee": 25000,
						"resv_total_price": "2500000.50"
					}
				]
			}
		}
	}`

	var resp GetRefundBookingResponse
	require.NoError(t, parseResponseBody([]byte(payload), &resp))
	require.Equal(t, Success, resp.Result.ErrorCode)

	quote := resp.Result.Response
	assert.Equal(t, "order-1", quote.OrderNumber)
	assert.Equal(t, "XYZ789", quote.Pnr)
	assert.Equal(t, "issued", quote.State)
	assert.Equal(t, "2026-10-01 08:00:00", quote.DepartureDate)
	assert.Equal(t, "2026-10-01 10:00:00", quote.ArrivalDate)
	assert.Equal(t, "airline", quote.ProviderType)
	require.Len(t, quote.ProviderBookings, 1)

	providerBooking := quote.ProviderBookings[0]
	assert.Equal(t, "ABC123", providerBooking.Pnr)
	assert.Equal(t, "sia", providerBooking.Provider)
	assert.Equal(t, "confirm", providerBooking.Status)
	// Amounts accept both number and decimal-string tokens (FlexibleNumber)
	assert.Equal(t, 2500000.5, float64(providerBooking.TotalPriceCost))
	assert.Equal(t, int64(250000), providerBooking.PenaltyAmount.Int64())
	assert.Equal(t, 2250000.5, float64(providerBooking.RefundableAmount))
	assert.Equal(t, int64(25000), providerBooking.AdminFee.Int64())
	assert.Equal(t, 2500000.5, float64(providerBooking.ResvTotalPrice))
}

func TestRefundProviderBookingQuote_FlexibleAmounts(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		expected float64
	}{
		{name: "number", payload: `{"pnr":"ABC123","refundable_amount":2250000.5}`, expected: 2250000.5},
		{name: "number string", payload: `{"pnr":"ABC123","refundable_amount":"2250000.5"}`, expected: 2250000.5},
		{name: "integer string", payload: `{"pnr":"ABC123","refundable_amount":"2250000"}`, expected: 2250000},
		{name: "null", payload: `{"pnr":"ABC123","refundable_amount":null}`, expected: 0},
		{name: "missing", payload: `{"pnr":"ABC123"}`, expected: 0},
		{name: "empty string", payload: `{"pnr":"ABC123","refundable_amount":""}`, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var quote RefundProviderBookingQuote
			require.NoError(t, json.Unmarshal([]byte(tt.payload), &quote))
			assert.Equal(t, tt.expected, float64(quote.RefundableAmount))
		})
	}

	var quote RefundProviderBookingQuote
	require.Error(t, json.Unmarshal([]byte(`{"pnr":"ABC123","refundable_amount":"abc"}`), &quote))
}

func TestPreRefundLoginListResponse_Parse(t *testing.T) {
	const payload = `{
		"result": {
			"error_code": 0,
			"error_msg": "",
			"response": {
				"pre_refund_login": [
					{"pnr": "ABC123", "img": "base64-captcha-image"},
					{"pnr": "DEF456", "img": "base64-captcha-image-2"}
				]
			}
		}
	}`

	var resp PreRefundLoginListResponse
	require.NoError(t, parseResponseBody([]byte(payload), &resp))
	require.Equal(t, Success, resp.Result.ErrorCode)
	require.Len(t, resp.Result.Response.PreRefundLogin, 2)

	assert.Equal(t, "ABC123", resp.Result.Response.PreRefundLogin[0].Pnr)
	assert.Equal(t, "base64-captcha-image", resp.Result.Response.PreRefundLogin[0].Img)
	assert.Equal(t, "DEF456", resp.Result.Response.PreRefundLogin[1].Pnr)
	assert.Equal(t, "base64-captcha-image-2", resp.Result.Response.PreRefundLogin[1].Img)
}

func TestCancelRefundBookingResponse_Parse(t *testing.T) {
	const payload = `{
		"result": {
			"error_code": 0,
			"error_msg": "",
			"response": {"order_number": "order-1"}
		}
	}`

	var resp CancelRefundBookingResponse
	require.NoError(t, parseResponseBody([]byte(payload), &resp))
	require.Equal(t, Success, resp.Result.ErrorCode)
	assert.Equal(t, "order-1", resp.Result.Response.OrderNumber)
}
