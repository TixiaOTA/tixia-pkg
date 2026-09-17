package panorama

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpclient "github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/logger"
)

func TestSellJourneysOptions(t *testing.T) {
	requestTimeout := 45 * time.Second
	retry := httpclient.RetryConfig{MaxRetries: 2}
	requestOptions := buildRequestOptions(
		sellJourneysTimeout,
		WithTimeout(requestTimeout),
		WithRetry(retry),
	)
	if requestOptions.Timeout == nil || *requestOptions.Timeout != requestTimeout {
		t.Fatalf("expected timeout %s, got %v", requestTimeout, requestOptions.Timeout)
	}
	if requestOptions.Retry == nil || requestOptions.Retry.MaxRetries != retry.MaxRetries {
		t.Fatalf("expected max retries %d, got %v", retry.MaxRetries, requestOptions.Retry)
	}
}

func TestSellJourneys_WithRequestTimeout(t *testing.T) {
	assertRequestTimeout(t, func(client *Client, options ...RequestOption) *PanoramaError {
		_, panoramaErr := client.SellJourneys(
			context.Background(),
			SellJourneysRequest{
				Signature: "signature",
				Adult:     1,
				Schedules: []PriceItinerarySchedule{{
					Provider: "provider",
					Journeys: []PriceItineraryScheduleJourney{{
						Segments: []PriceItineraryScheduleSegment{{
							SegmentCode: "segment",
							FareCode:    "fare",
						}},
					}},
				}},
			},
			options...,
		)
		return panoramaErr
	})
}

func TestUpdateEndpoints_WithRequestTimeout(t *testing.T) {
	t.Run("UpdateContacts", func(t *testing.T) {
		assertRequestTimeout(t, func(client *Client, options ...RequestOption) *PanoramaError {
			_, panoramaErr := client.UpdateContacts(
				context.Background(),
				UpdateContactsRequest{
					Signature: "signature",
					Booker:    testContactPerson(),
					Contacts:  []ContactPerson{testContactPerson()},
				},
				options...,
			)
			return panoramaErr
		})
	})

	t.Run("UpdatePassengers", func(t *testing.T) {
		assertRequestTimeout(t, func(client *Client, options ...RequestOption) *PanoramaError {
			_, panoramaErr := client.UpdatePassengers(
				context.Background(),
				UpdatePassengersRequest{
					Signature:  "signature",
					Passengers: []UpdatePassenger{testUpdatePassenger()},
				},
				options...,
			)
			return panoramaErr
		})
	})
}

func assertRequestTimeout(t *testing.T, call func(*Client, ...RequestOption) *PanoramaError) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
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
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	panoramaErr := call(
		client,
		WithTimeout(10*time.Millisecond),
		WithRetry(httpclient.RetryConfig{MaxRetries: 0}),
	)
	if panoramaErr == nil {
		t.Fatal("expected request timeout error")
	}
}

func testContactPerson() ContactPerson {
	return ContactPerson{
		Title:           "Mr",
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john@example.com",
		CallingCode:     "+62",
		Mobile:          "81234567890",
		NationalityCode: "ID",
	}
}

func testUpdatePassenger() UpdatePassenger {
	return UpdatePassenger{
		Title:           "Mr",
		FirstName:       "John",
		LastName:        "Doe",
		BirthDate:       "1990-01-01",
		PaxType:         "ADT",
		NationalityCode: "ID",
	}
}
