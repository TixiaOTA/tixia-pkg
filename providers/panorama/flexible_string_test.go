package panorama

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlexibleString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected FlexibleString
		wantErr  bool
	}{
		{
			name:     "string value",
			input:    `"CGK-DPS"`,
			expected: "CGK-DPS",
		},
		{
			name:     "bool false",
			input:    `false`,
			expected: "",
		},
		{
			name:     "bool true",
			input:    `true`,
			expected: "",
		},
		{
			name:     "null",
			input:    `null`,
			expected: "",
		},
		{
			name:     "number",
			input:    `123`,
			expected: "123",
		},
		{
			name:    "invalid type",
			input:   `{"code":"CGK-DPS"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got FlexibleString
			err := json.Unmarshal([]byte(tt.input), &got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestFlexibleString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    FlexibleString
		expected string
	}{
		{
			name:     "non-empty string",
			input:    "CGK-DPS",
			expected: `"CGK-DPS"`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: `""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(got))
		})
	}
}

func TestBookingResponse_UnmarshalLegCodeBool(t *testing.T) {
	const payload = `{
		"result": {
			"error_code": 0,
			"error_msg": "",
			"response": {
				"order_number": "ORD-123",
				"provider_bookings": [
					{
						"pnr": "ABC123",
						"journeys": [
							{
								"sequence": 1,
								"segments": [
									{
										"segment_code": "SEG1",
										"legs": [
											{
												"leg_code": false,
												"origin": "CGK",
												"destination": "DPS"
											}
										]
									}
								]
							}
						]
					}
				]
			}
		}
	}`

	var resp GetBookingResponse
	err := json.Unmarshal([]byte(payload), &resp)
	require.NoError(t, err)

	require.Len(t, resp.Result.Response.ProviderBookings, 1)
	booking := resp.Result.Response.ProviderBookings[0]
	require.Len(t, booking.Journeys, 1)
	require.Len(t, booking.Journeys[0].Segments, 1)
	require.Len(t, booking.Journeys[0].Segments[0].Legs, 1)

	leg := booking.Journeys[0].Segments[0].Legs[0]
	assert.Equal(t, FlexibleString(""), leg.LegCode)
	assert.Equal(t, "CGK", leg.Origin)
	assert.Equal(t, "DPS", leg.Destination)
}

func TestPriceItineraryLeg_UnmarshalLegCodeBool(t *testing.T) {
	const payload = `{
		"leg_code": false,
		"leg_code2": true,
		"origin": "CGK",
		"destination": "DPS"
	}`

	var leg PriceItineraryLeg
	err := json.Unmarshal([]byte(payload), &leg)
	require.NoError(t, err)

	assert.Equal(t, FlexibleString(""), leg.LegCode)
	assert.Equal(t, FlexibleString(""), leg.LegCode2)
	assert.Equal(t, "CGK", leg.Origin)
	assert.Equal(t, "DPS", leg.Destination)
}
