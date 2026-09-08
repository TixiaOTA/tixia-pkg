package panorama

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlexibleNumber_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected FlexibleNumber
		wantErr  bool
	}{
		{name: "integer", input: `1521000`, expected: 1521000},
		{name: "decimal", input: `1521000.0`, expected: 1521000},
		{name: "numeric string", input: `"1521000.0"`, expected: 1521000},
		{name: "empty string", input: `""`, expected: 0},
		{name: "null", input: `null`, expected: 0},
		{name: "invalid string", input: `"invalid"`, wantErr: true},
		{name: "invalid type", input: `{"value":1}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got FlexibleNumber
			err := json.Unmarshal([]byte(tt.input), &got)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.expected.Int64(), got.Int64())
		})
	}
}

func TestFlexibleNumber_MarshalJSON(t *testing.T) {
	got, err := json.Marshal(FlexibleNumber(1521000.5))
	require.NoError(t, err)
	assert.Equal(t, `1521000.5`, string(got))
}
