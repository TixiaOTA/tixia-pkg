package train

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestGetSeatAvailabilityOWAndRT(t *testing.T) {
	for _, tc := range []struct {
		name  string
		fares []string
		body  string
	}{
		{"OW", []string{"fare-out"}, `{"jsonrpc":"2.0","result":{"error_code":0,"error_msg":"","response":[[{"cabin_number":1,"cabin_code":"LUX","cabin_name":"LUX-1","maximum_column_number":6,"seat_rows":[{"row_number":"2","seats":[{"seat_code":15999808,"column":"A","sequence":1,"availability":1},{"seat_code":"","column":"","sequence":2,"availability":-1},{"seat_code":15999809,"column":"B","sequence":3,"availability":0}]}]}]]}}`},
		{"RT", []string{"fare-out", "fare-back"}, `{"jsonrpc":"2.0","result":{"error_code":0,"error_msg":"","response":[[{"cabin_number":1,"cabin_name":"LUX-1","maximum_column_number":6,"seat_rows":[]}],[{"cabin_number":1,"cabin_name":"EKS-1","seat_rows":[]},{"cabin_number":2,"cabin_name":"EKS-2","seat_rows":[]}]]}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var requestCount atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCount.Add(1)
				if r.URL.Path != "/booking/train" || r.Header.Get("Action") != "get_seat_availability" || r.Header.Get("Signature") != "signature" {
					t.Errorf("incorrect method, path or headers: %s %s", r.Method, r.URL.Path)
				}
				var sent GetSeatAvailabilityRequest
				if err := json.NewDecoder(r.Body).Decode(&sent); err != nil {
					t.Errorf("decode request: %v", err)
				}
				if len(sent.FareCodes) != len(tc.fares) || sent.Provider != ProviderKAI {
					t.Errorf("unexpected request: %+v", sent)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()

			c, err := New(Config{BaseURL: server.URL, Credentials: Credentials{User: "u", Password: "p", APIKey: "key"}})
			if err != nil { t.Fatal(err) }
			got, providerErr := c.GetSeatAvailability(context.Background(), "signature", GetSeatAvailabilityRequest{FareCodes: tc.fares, Provider: ProviderKAI})
			if providerErr != nil { t.Fatal(providerErr) }
			if len(got.Result.Response) != len(tc.fares) || requestCount.Load() != 1 {
				t.Fatalf("unexpected outer response length or request count: %+v", got)
			}
			if tc.name == "OW" {
				seats := got.Result.Response[0][0].SeatRows[0].Seats
				if len(seats) != 3 || string(seats[0].SeatCode) != "15999808" || string(seats[1].SeatCode) != `""` || seats[1].Availability != -1 {
					t.Fatalf("seat positions changed: %+v", seats)
				}
			} else if len(got.Result.Response[1]) != 2 {
				t.Fatal("RT inbound cabins missing")
			}
		})
	}
}

func TestCommitBookingOWAndRTNoRetry(t *testing.T) {
	for _, n := range []int{1, 2} {
		t.Run(string(rune('0'+n))+"-way", func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Header.Get("Action") != "commit_booking" || r.Header.Get("Signature") != "signature" {
					t.Errorf("wrong commit headers")
				}
				var req CommitBookingRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil { t.Errorf("decode commit: %v", err) }
				if req.ForceIssued || len(req.Schedules) != n { t.Errorf("invalid commit request: %+v", req) }
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"jsonrpc":"2.0","result":{"error_code":0,"error_msg":"","response":{"order_number":"TN.20081330304","pnr":"BKC3BGS","state":"booked","state_description":"Booked","hold_date":"2020-08-13 02:42:13","provider_bookings":[{"pnr":"BKC3BGS","provider":"kai","sequence":1,"journeys":[{"journey_code":"J","fare_code":"F","seats":[{"seat":"LUX-1,1B","seat_code":"13404192","passenger_sequence":1}]}]}]}}}`))
			}))
			defer server.Close()
			c, err := New(Config{BaseURL:server.URL,Credentials:Credentials{User:"u",Password:"p",APIKey:"key"}})
			if err != nil { t.Fatal(err) }
			req := exampleCommitRequest(n)
			resp, providerErr := c.CommitBooking(context.Background(), "signature", req)
			if providerErr != nil { t.Fatal(providerErr) }
			if resp.Result.Response.OrderNumber != "TN.20081330304" || len(resp.Result.Response.ProviderBookings[0].Journeys[0].Seats) != 1 { t.Fatalf("invalid response: %+v", resp) }
			if calls.Load() != 1 { t.Fatalf("unexpected commit attempts: %d", calls.Load()) }
		})
	}

	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	c, _ := New(Config{BaseURL:server.URL,Credentials:Credentials{User:"u",Password:"p",APIKey:"key"}})
	if _, err := c.CommitBooking(context.Background(), "signature", exampleCommitRequest(1)); err == nil {
		t.Fatal("expected HTTP provider error")
	}
	if calls.Load() != 1 { t.Fatalf("commit retried %d times", calls.Load()) }
}

func TestTrainBookingValidation(t *testing.T) {
	req := exampleCommitRequest(1)
	req.ForceIssued = true
	if err := req.Validate(); err == nil { t.Fatal("force_issued=true must fail") }
	req = exampleCommitRequest(1)
	req.Passengers[0].Identity.IdentityNumber = ""
	if err := req.Validate(); err == nil { t.Fatal("missing identity must fail") }
	if err := (GetSeatAvailabilityRequest{Provider:"kai",FareCodes:[]string{"","fare"}}).Validate(); err == nil { t.Fatal("empty fare must fail") }
	if err := (GetSeatAvailabilityRequest{Provider:"kai",FareCodes:[]string{"fare"}}).Validate(); err != nil { t.Fatal(err) }

	body := []byte(`{"jsonrpc":"2.0","result":{"error_code":123,"error_msg":"sold out","response":[]}}`)
	var response GetSeatAvailabilityResponse
	if err := parseResponseBody(body, &response); err == nil || !strings.Contains(err.Message, "sold out") { t.Fatal("provider error not propagated") }
}

func exampleCommitRequest(n int) CommitBookingRequest {
	booker := BookingBooker{Title:"MR",FirstName:"Budi",LastName:"Santoso",Email:"budi@example.com",CallingCode:"62",Mobile:"81234567890",NationalityCode:"ID"}
	req := CommitBookingRequest{
		Contacts: []BookingContact{{BookingBooker:booker,IsAlsoBooker:true}},
		Booker: booker,
		Passengers: []BookingPassenger{{PaxType:"ADT",Title:"MR",FirstName:"Budi",LastName:"Santoso",BirthDate:"1996-02-20",NationalityCode:"ID",Identity:BookingIdentity{IdentityNumber:"1234567890123456",IdentityType:"ktp"}}},
		Schedules: []BookingSchedule{{Provider:ProviderKAI,Journeys:[]BookingJourney{{JourneyCode:"outbound-journey",FareCode:"outbound-fare"}}}},
	}
	if n == 2 { req.Schedules = append(req.Schedules, BookingSchedule{Provider:ProviderKAI,Journeys:[]BookingJourney{{JourneyCode:"inbound-journey",FareCode:"inbound-fare"}}}) }
	return req
}
