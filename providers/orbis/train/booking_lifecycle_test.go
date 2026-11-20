package train

import (
	"encoding/json"
	"testing"
)

func TestBookingLifecycleRequestValidation(t *testing.T) {
	if err := (IssueBookingRequest{}).Validate(); err == nil {
		t.Fatal("empty issue order_number should fail")
	}
	if err := (CancelBookingRequest{OrderNumber: "  "}).Validate(); err == nil {
		t.Fatal("blank cancel order_number should fail")
	}
	if err := (IssueBookingRequest{OrderNumber: "TN.123"}).Validate(); err != nil {
		t.Fatalf("valid issue request: %v", err)
	}
	if err := (CancelBookingRequest{OrderNumber: "TN.123"}).Validate(); err != nil {
		t.Fatalf("valid cancel request: %v", err)
	}
}

func TestIssueBookingResponseIssuedTickets(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","result":{"error_code":0,"error_msg":"","response":{"order_number":"TN.123","state":"issued","provider_bookings":[{"state":"issued","sequence":1,"provider":"kai","tickets":[{"ticket_number":"T123","passenger_sequence":1,"pax_type":"ADT"}]}]}}}`)
	var resp IssueBookingResponse
	if err := parseResponseBody(data, &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Result.Response.State != "issued" || resp.Result.Response.ProviderBookings[0].Tickets[0].TicketNumber != "T123" {
		t.Fatalf("unexpected decoded issue response: %+v", resp.Result.Response)
	}
}

func TestCancelBookingResponseTwoProviders(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":null,"result":{"error_code":0,"error_msg":"","response":[{"error_code":0,"error_msg":"","response":"","status":"SUCCESS"},{"error_code":0,"error_msg":"","response":"","status":"SUCCESS"}]}}`)
	var resp CancelBookingResponse
	if err := parseResponseBody(data, &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Result.Response) != 2 || resp.Result.Response[1].Status != "SUCCESS" {
		t.Fatalf("unexpected cancel response: %+v", resp)
	}
}

func TestCancelBookingProviderError4022(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":null,"result":{"error_code":4022,"error_msg":"Cancel booking failed, 7 sequence 1, provider kai","response":""}}`)
	var resp CancelBookingResponse
	err := parseResponseBody(data, &resp)
	if err == nil || err.Code != ErrorCode(4022) {
		t.Fatalf("expected 4022 error, got %v", err)
	}
}

func TestCancelBookingResponsePartialFailure(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","result":{"error_code":0,"error_msg":"","response":[{"error_code":0,"error_msg":"","response":"","status":"SUCCESS"},{"error_code":4022,"error_msg":"failed","response":"","status":"FAILED"}]}}`)
	var resp CancelBookingResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Result.Response) != 2 || resp.Result.Response[1].ErrorCode != 4022 {
		t.Fatalf("expected provider error details: %+v", resp)
	}
}
