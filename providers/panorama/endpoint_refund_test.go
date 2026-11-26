package panorama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRefundTestClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()

	log := logger.New(logger.Config{Level: "error", JSONOutput: false})
	client, err := New(Config{
		BaseURL: server.URL,
		Credentials: Credentials{
			APIKey:   "api-key",
			User:     "user",
			Password: "password",
		},
	}, log)
	require.NoError(t, err)

	return client
}

func refundTestProviderBookings() []RefundProviderBooking {
	return []RefundProviderBooking{{
		Pnr:      "ABC123",
		Journeys: []RefundJourney{{Origin: "CGK", Destination: "DPS", DepartureDate: "2026-10-01", Pax: []int{1}}},
	}}
}

func TestPreRefundLoginList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/booking/airline" {
			t.Errorf("expected path /booking/airline, got %s", r.URL.Path)
		}
		if got := r.Header.Get("Action"); got != string(PreRefundLoginList) {
			t.Errorf("expected action %q, got %q", PreRefundLoginList, got)
		}
		if got := r.Header.Get("Signature"); got != "sig-123" {
			t.Errorf("expected signature header, got %q", got)
		}

		var body PreRefundLoginListRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
			return
		}
		if body.OrderNumber != "order-1" {
			t.Errorf("expected order_number order-1, got %q", body.OrderNumber)
		}
		if len(body.Provider) != 1 || body.Provider[0] != "sia" {
			t.Errorf("expected provider [sia], got %v", body.Provider)
		}
		if len(body.Pnr) != 1 || body.Pnr[0] != "ABC123" {
			t.Errorf("expected pnr [ABC123], got %v", body.Pnr)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"result": {
				"error_code": 0,
				"error_msg": "",
				"response": {"pre_refund_login": [{"pnr": "ABC123", "img": "base64-captcha-image"}]}
			}
		}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := newRefundTestClient(t, server)

	resp, panoramaErr := client.PreRefundLoginList(context.Background(), PreRefundLoginListRequest{
		Signature:   "sig-123",
		OrderNumber: "order-1",
		Provider:    []string{"sia"},
		Pnr:         []string{"ABC123"},
	})

	require.Nil(t, panoramaErr)
	require.NotNil(t, resp)
	assert.Equal(t, Success, resp.Result.ErrorCode)
	require.Len(t, resp.Result.Response.PreRefundLogin, 1)
	assert.Equal(t, "ABC123", resp.Result.Response.PreRefundLogin[0].Pnr)
	assert.Equal(t, "base64-captcha-image", resp.Result.Response.PreRefundLogin[0].Img)
}

func TestGetRefundBooking(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/booking/airline" {
			t.Errorf("expected path /booking/airline, got %s", r.URL.Path)
		}
		if got := r.Header.Get("Action"); got != string(GetRefundBooking) {
			t.Errorf("expected action %q, got %q", GetRefundBooking, got)
		}
		if got := r.Header.Get("Signature"); got != "sig-123" {
			t.Errorf("expected signature header, got %q", got)
		}

		var body GetRefundBookingRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
			return
		}
		if body.OrderNumber != "order-1" {
			t.Errorf("expected order_number order-1, got %q", body.OrderNumber)
		}
		if body.Captcha["ABC123"] != "Refund" {
			t.Errorf("expected captcha {ABC123: Refund}, got %v", body.Captcha)
		}
		if len(body.Passengers) != 2 || body.Passengers[1].PassengerNumber != 2 {
			t.Errorf("expected passengers 1 and 2, got %v", body.Passengers)
		}
		if len(body.ProviderBookings) != 1 || body.ProviderBookings[0].Pnr != "ABC123" {
			t.Errorf("expected provider booking ABC123, got %v", body.ProviderBookings)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
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
					"provider_bookings": [{
						"pnr": "ABC123",
						"provider": "sia",
						"status": "confirm",
						"total_price_cost": 2500000.5,
						"penalty_amount": "250000",
						"refundable_amount": 2250000.5,
						"admin_fee": 25000,
						"resv_total_price": 2500000.5
					}]
				}
			}
		}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := newRefundTestClient(t, server)

	resp, panoramaErr := client.GetRefundBooking(context.Background(), GetRefundBookingRequest{
		Signature:        "sig-123",
		OrderNumber:      "order-1",
		Passengers:       []RefundPassengerNumber{{PassengerNumber: 1}, {PassengerNumber: 2}},
		ProviderBookings: refundTestProviderBookings(),
		Captcha:          map[string]string{"ABC123": "Refund"},
	})

	require.Nil(t, panoramaErr)
	require.NotNil(t, resp)
	assert.Equal(t, Success, resp.Result.ErrorCode)
	assert.Equal(t, "order-1", resp.Result.Response.OrderNumber)
	assert.Equal(t, "issued", resp.Result.Response.State)
	require.Len(t, resp.Result.Response.ProviderBookings, 1)

	quote := resp.Result.Response.ProviderBookings[0]
	assert.Equal(t, "ABC123", quote.Pnr)
	assert.Equal(t, "confirm", quote.Status)
	assert.Equal(t, 2500000.5, float64(quote.TotalPriceCost))
	assert.Equal(t, int64(250000), quote.PenaltyAmount.Int64())
	assert.Equal(t, 2250000.5, float64(quote.RefundableAmount))
	assert.Equal(t, int64(25000), quote.AdminFee.Int64())
	assert.Equal(t, 2500000.5, float64(quote.ResvTotalPrice))
}

func TestCancelRefundBooking(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/booking/airline" {
			t.Errorf("expected path /booking/airline, got %s", r.URL.Path)
		}
		// The refund execution reuses the cancel action with a refund payload
		if got := r.Header.Get("Action"); got != string(Cancel) {
			t.Errorf("expected action %q, got %q", Cancel, got)
		}
		if got := r.Header.Get("Signature"); got != "sig-123" {
			t.Errorf("expected signature header, got %q", got)
		}

		var body CancelRefundBookingRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
			return
		}
		if body.OrderNumber != "order-1" {
			t.Errorf("expected order_number order-1, got %q", body.OrderNumber)
		}
		if len(body.ProviderBookings) != 1 || body.ProviderBookings[0].Pnr != "ABC123" {
			t.Errorf("expected provider booking ABC123, got %v", body.ProviderBookings)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"result": {
				"error_code": 0,
				"error_msg": "",
				"response": {"order_number": "order-1"}
			}
		}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := newRefundTestClient(t, server)

	resp, panoramaErr := client.CancelRefundBooking(context.Background(), CancelRefundBookingRequest{
		Signature:        "sig-123",
		OrderNumber:      "order-1",
		ProviderBookings: refundTestProviderBookings(),
	})

	require.Nil(t, panoramaErr)
	require.NotNil(t, resp)
	assert.Equal(t, Success, resp.Result.ErrorCode)
	assert.Equal(t, "order-1", resp.Result.Response.OrderNumber)
}

func TestRefundEndpoints_ProviderError(t *testing.T) {
	tests := []struct {
		name   string
		action PanoramaHeaderAction
		call   func(c *Client) (*PanoramaError, bool)
	}{
		{
			name:   "pre refund login list",
			action: PreRefundLoginList,
			call: func(c *Client) (*PanoramaError, bool) {
				resp, err := c.PreRefundLoginList(context.Background(), PreRefundLoginListRequest{
					Signature:   "sig-123",
					OrderNumber: "order-1",
				})
				return err, resp == nil
			},
		},
		{
			name:   "get refund booking",
			action: GetRefundBooking,
			call: func(c *Client) (*PanoramaError, bool) {
				resp, err := c.GetRefundBooking(context.Background(), GetRefundBookingRequest{
					Signature:        "sig-123",
					OrderNumber:      "order-1",
					Passengers:       []RefundPassengerNumber{{PassengerNumber: 1}},
					ProviderBookings: refundTestProviderBookings(),
					Captcha:          map[string]string{"ABC123": "Refund"},
				})
				return err, resp == nil
			},
		},
		{
			name:   "cancel refund booking",
			action: Cancel,
			call: func(c *Client) (*PanoramaError, bool) {
				resp, err := c.CancelRefundBooking(context.Background(), CancelRefundBookingRequest{
					Signature:        "sig-123",
					OrderNumber:      "order-1",
					ProviderBookings: refundTestProviderBookings(),
				})
				return err, resp == nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("Action"); got != string(tt.action) {
					t.Errorf("expected action %q, got %q", tt.action, got)
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{
					"result": {
						"error_code": 4043,
						"error_msg": "cancel not available",
						"response": ""
					}
				}`)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer server.Close()

			client := newRefundTestClient(t, server)

			panoramaErr, nilResp := tt.call(client)
			require.NotNil(t, panoramaErr)
			assert.True(t, nilResp)
			assert.Equal(t, CancelNotAvailable, panoramaErr.Result.ErrorCode)
			assert.Equal(t, "cancel not available", panoramaErr.Result.ErrorMsg)
		})
	}
}

func TestRefundEndpoints_RequestValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("expected no HTTP call for invalid requests")
	}))
	defer server.Close()

	client := newRefundTestClient(t, server)

	_, err := client.PreRefundLoginList(context.Background(), PreRefundLoginListRequest{})
	require.NotNil(t, err)
	assert.Equal(t, InvalidRequest, err.Result.ErrorCode)

	_, err = client.GetRefundBooking(context.Background(), GetRefundBookingRequest{Signature: "sig-123", OrderNumber: "order-1"})
	require.NotNil(t, err)
	assert.Equal(t, InvalidRequest, err.Result.ErrorCode)

	_, err = client.CancelRefundBooking(context.Background(), CancelRefundBookingRequest{Signature: "sig-123", OrderNumber: "order-1"})
	require.NotNil(t, err)
	assert.Equal(t, InvalidRequest, err.Result.ErrorCode)
}
