package panorama

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlexibleBool_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected FlexibleBool
		wantErr  bool
	}{
		{
			name:     "bool true",
			input:    `true`,
			expected: true,
		},
		{
			name:     "bool false",
			input:    `false`,
			expected: false,
		},
		{
			name:     "string true",
			input:    `"true"`,
			expected: true,
		},
		{
			name:     "string false",
			input:    `"false"`,
			expected: false,
		},
		{
			name:     "string 1",
			input:    `"1"`,
			expected: true,
		},
		{
			name:     "string 0",
			input:    `"0"`,
			expected: false,
		},
		{
			name:     "null",
			input:    `null`,
			expected: false,
		},
		{
			name:     "number 1",
			input:    `1`,
			expected: true,
		},
		{
			name:     "number 0",
			input:    `0`,
			expected: false,
		},
		{
			name:    "invalid string",
			input:   `"maybe"`,
			wantErr: true,
		},
		{
			name:    "invalid type",
			input:   `{"value":true}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got FlexibleBool
			err := json.Unmarshal([]byte(tt.input), &got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, bool(tt.expected), got.Bool())
		})
	}
}

func TestFlexibleBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    FlexibleBool
		expected string
	}{
		{
			name:     "true",
			input:    true,
			expected: `true`,
		},
		{
			name:     "false",
			input:    false,
			expected: `false`,
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

func TestPriceItineraryJourney_UnmarshalIsVtlFlightString(t *testing.T) {
	const payload = `{
		"origin": "CGK",
		"destination": "SIN",
		"is_vtl_flight": "true",
		"sequence": 1
	}`

	var journey PriceItineraryJourney
	err := json.Unmarshal([]byte(payload), &journey)
	require.NoError(t, err)

	assert.True(t, journey.IsVtlFlight.Bool())
	assert.Equal(t, 1, journey.Sequence)
}

func TestBookingJourney_UnmarshalIsVtlFlightString(t *testing.T) {
	const payload = `{
		"sequence": 1,
		"origin": "CGK",
		"destination": "SIN",
		"is_vtl_flight": "false"
	}`

	var journey BookingJourney
	err := json.Unmarshal([]byte(payload), &journey)
	require.NoError(t, err)

	assert.False(t, journey.IsVtlFlight.Bool())
}
