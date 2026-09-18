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

func TestSearchRequestValidate(t *testing.T) {
	valid := SearchRequest{
		Signature:  "signature",
		Adult:      1,
		Provider:   "amadeus",
		CabinClass: CabinClassEconomy,
		JourneyList: []SearchJourney{{
			Origin:        "CGK",
			Destination:   "SIN",
			DepartureDate: "2026-10-01",
		}},
	}

	require.NoError(t, valid.Validate())
	require.ErrorContains(t, (&SearchRequest{}).Validate(), "signature is required")
	require.ErrorContains(t, (&SearchRequest{Signature: "signature"}).Validate(), "adult is required")
	require.ErrorContains(t, (&SearchRequest{Signature: "signature", Adult: 1}).Validate(), "provider is required")
	require.ErrorContains(t, (&SearchRequest{Signature: "signature", Adult: 1, Provider: "amadeus"}).Validate(), "cabin_class is required")
	require.ErrorContains(t, (&SearchRequest{Signature: "signature", Adult: 1, Provider: "amadeus", CabinClass: "Z"}).Validate(), "cabin_class must be one of Y, W, C, F")
	require.ErrorContains(t, (&SearchRequest{Signature: "signature", Adult: 1, Provider: "amadeus", CabinClass: CabinClassEconomy}).Validate(), "journey_list is required")
}

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/booking/airline" {
			t.Errorf("expected path /booking/airline, got %s", r.URL.Path)
		}
		if got := r.Header.Get("Action"); got != string(Search) {
			t.Errorf("expected action %q, got %q", Search, got)
		}
		if got := r.Header.Get("Signature"); got != "signature" {
			t.Errorf("expected signature header, got %q", got)
		}

		var body struct {
			Signature    *string         `json:"signature"`
			JourneyList  []SearchJourney `json:"journey_list"`
			CarrierCodes []string        `json:"carrier_codes"`
			PromoCodes   []string        `json:"promo_codes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
			return
		}
		if body.Signature != nil {
			t.Errorf("expected signature to be omitted from body, got %q", *body.Signature)
		}
		if len(body.JourneyList) == 0 {
			t.Error("expected at least one journey")
			return
		}
		assert.Equal(t, "CGK", body.JourneyList[0].Origin)
		assert.Equal(t, []string{}, body.CarrierCodes)
		assert.Equal(t, []string{}, body.PromoCodes)

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"result": {
				"error_code": 0,
				"error_msg": "",
				"response": {"schedules": [], "recommendations": []}
			}
		}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

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

	resp, panoramaErr := client.Search(context.Background(), SearchRequest{
		Signature:  "signature",
		Adult:      1,
		Provider:   "amadeus",
		CabinClass: CabinClassEconomy,
		JourneyList: []SearchJourney{{
			Origin:        "CGK",
			Destination:   "SIN",
			DepartureDate: "2026-10-01",
		}},
	})

	require.Nil(t, panoramaErr)
	require.NotNil(t, resp)
	assert.Equal(t, Success, resp.Result.ErrorCode)
	assert.Empty(t, resp.Result.Response.Schedules)
	assert.Empty(t, resp.Result.Response.Recommendations)
}

func TestSearchResponse_Parse(t *testing.T) {
	const payload = `{
		"result": {
			"error_code": 0,
			"error_msg": "",
			"response": {
				"schedules": [{
					"provider": "amadeus",
					"origin": "CGK",
					"destination": "SIN",
					"departure_date": "2021-08-15",
					"sequence": 1,
					"journeys": [{
						"origin": "CGK",
						"destination": "SIN",
						"journey_code": "J1",
						"provider": "amadeus",
						"sequence": 1,
						"segments": [{
							"segment_code": "S1",
						"provider": "amadeus",
							"sequence": 1,
							"legs": [{
								"leg_code": "L1",
								"leg_code2": "CGK-SIN",
								"sequence": 1
							}],
							"fares": [{
								"available_count": "9",
								"fare_details": [{"amount": "30"}],
								"service_charges": [{
									"charge_type": "FARE",
									"amount": 1400000
								}],
								"service_charge_summary": [{
									"pax_type": "ADT",
									"total_price": 1860776
								}]
							}]
						}]
					}]
				}],
				"recommendations": [{
					"sequence": 1,
						"journey_flight_refs": [{
							"journey_ref_id": "J1",
							"fare_flight_refs": [{
								"fare_ref_id": "F1",
								"class_of_service": "N"
							}]
						}]
				}]
			}
		}
	}`

	var resp SearchResponse
	require.NoError(t, parseResponseBody([]byte(payload), &resp))
	require.Len(t, resp.Result.Response.Schedules, 1)
	require.Len(t, resp.Result.Response.Schedules[0].Journeys, 1)

	journey := resp.Result.Response.Schedules[0].Journeys[0]
	require.Len(t, journey.Segments, 1)
	segment := journey.Segments[0]
	require.Len(t, segment.Legs, 1)
	assert.Equal(t, FlexibleString("L1"), segment.Legs[0].LegCode)
	require.Len(t, segment.Fares, 1)
	fare := segment.Fares[0]
	assert.Equal(t, FlexibleNumber(9), fare.AvailableCount)
	assert.Equal(t, FlexibleString("30"), fare.FareDetails[0].Amount)
	assert.Equal(t, 1400000.0, fare.ServiceCharges[0].Amount)
	assert.Equal(t, 1860776.0, fare.ServiceChargeSummary[0].TotalPrice)

	require.Len(t, resp.Result.Response.Recommendations, 1)
	assert.Equal(t, "F1", resp.Result.Response.Recommendations[0].JourneyFlightRefs[0].FareFlightRefs[0].FareRefID)
}

func TestSearchPanoramaError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"result": {
				"error_code": 4020,
				"error_msg": "no schedules",
				"response": ""
			}
		}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

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

	resp, panoramaErr := client.Search(context.Background(), SearchRequest{
		Signature:  "signature",
		Adult:      1,
		Provider:   "amadeus",
		CabinClass: CabinClassEconomy,
		JourneyList: []SearchJourney{{
			Origin:        "CGK",
			Destination:   "SIN",
			DepartureDate: "2026-10-01",
		}},
	})

	require.Nil(t, resp)
	require.NotNil(t, panoramaErr)
	assert.Equal(t, PanoramaErrorCode(4020), panoramaErr.Result.ErrorCode)
	assert.Equal(t, "no schedules", panoramaErr.Result.ErrorMsg)
}
