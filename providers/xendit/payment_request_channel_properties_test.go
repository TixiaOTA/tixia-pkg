package xendit

import (
	"encoding/json"
	"testing"
)

func TestPaymentRequestChannelProperties_UnmarshalJSON_Object(t *testing.T) {
	var props PaymentRequestChannelProperties

	err := json.Unmarshal([]byte(`{"account_mobile_number":"+628123456789","display_name":"John"}`), &props)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, ok := props.value.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", props.value)
	}

	if m["account_mobile_number"] != "+628123456789" {
		t.Fatalf("unexpected account_mobile_number: %v", m["account_mobile_number"])
	}

	if m["display_name"] != "John" {
		t.Fatalf("unexpected display_name: %v", m["display_name"])
	}
}

func TestPaymentRequestChannelProperties_UnmarshalJSON_Null(t *testing.T) {
	var props PaymentRequestChannelProperties

	err := json.Unmarshal([]byte(`null`), &props)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !props.IsZero() {
		t.Fatal("expected zero value after null unmarshal")
	}
}

func TestPaymentRequestChannelProperties_UnmarshalJSON_Invalid(t *testing.T) {
	var props PaymentRequestChannelProperties

	err := json.Unmarshal([]byte(`["not","an","object"]`), &props)
	if err == nil {
		t.Fatal("expected error for non-object JSON")
	}
}

func TestPaymentRequestChannelProperties_UnmarshalMarshalRoundTrip(t *testing.T) {
	input := `{"account_mobile_number":"+628123456789","verification_data":{"customer_name":"Jane"}}`

	var props PaymentRequestChannelProperties
	if err := json.Unmarshal([]byte(input), &props); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	out, err := json.Marshal(props)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("unmarshal round-trip output: %v", err)
	}

	var want map[string]any
	if err := json.Unmarshal([]byte(input), &want); err != nil {
		t.Fatalf("unmarshal want: %v", err)
	}

	if got["account_mobile_number"] != want["account_mobile_number"] {
		t.Fatalf("account_mobile_number mismatch: got %v want %v", got["account_mobile_number"], want["account_mobile_number"])
	}

	verification, ok := got["verification_data"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map, got %T", got["verification_data"])
	}

	if verification["customer_name"] != "Jane" {
		t.Fatalf("unexpected customer_name: %v", verification["customer_name"])
	}
}

func TestPaymentRequestChannelProperties_MarshalTypedOutbound(t *testing.T) {
	props := NewOVOChannelProperties(OVOChannelProperties{
		AccountMobileNumber: "+628123456789",
	})

	out, err := json.Marshal(props)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	const want = `{"account_mobile_number":"+628123456789"}`
	if string(out) != want {
		t.Fatalf("got %s want %s", out, want)
	}
}

func TestCreatePaymentRequestResponse_UnmarshalChannelProperties(t *testing.T) {
	payload := `{
		"payment_request_id": "pr-90392f42-d98a-49ef-a7f3-abcezas123",
		"business_id": "5f27a14a9bf05c73dd040bc8",
		"reference_id": "order_12345",
		"type": "PAY",
		"country": "ID",
		"currency": "IDR",
		"request_amount": "150000",
		"channel_code": "OVO",
		"channel_properties": {"account_mobile_number": "+628123456789"},
		"status": "SUCCEEDED",
		"created": "2021-12-31T23:59:59Z",
		"updated": "2021-12-31T23:59:59Z"
	}`

	var resp CreatePaymentRequestResponse
	if err := json.Unmarshal([]byte(payload), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	m, ok := resp.ChannelProperties.value.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", resp.ChannelProperties.value)
	}

	if m["account_mobile_number"] != "+628123456789" {
		t.Fatalf("unexpected account_mobile_number: %v", m["account_mobile_number"])
	}
}
