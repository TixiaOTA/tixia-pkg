package xendit

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/stretchr/testify/suite"
)

type EndpointTestSuite struct {
	suite.Suite
	testLogger *logger.Logger
}

func (suite *EndpointTestSuite) SetupTest() {
	suite.testLogger = logger.New(logger.Config{
		Level:      "info",
		JSONOutput: false,
	})
}

func TestEndpointTestSuite(t *testing.T) {
	suite.Run(t, new(EndpointTestSuite))
}

func (suite *EndpointTestSuite) newTestClient(serverURL string) *Client {
	client, err := New(Config{
		BaseURL:      serverURL,
		SecretAPIKey: "test-secret-key",
	}, suite.testLogger)
	suite.NoError(err)

	return client
}

func (suite *EndpointTestSuite) TestCreateSession_Success201() {
	var receivedAuth string
	var receivedBody CreateSessionRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodPost, r.Method)
		suite.Equal("/sessions", r.URL.Path)

		receivedAuth = r.Header.Get("Authorization")

		bodyBytes, err := io.ReadAll(r.Body)
		suite.NoError(err)

		err = json.Unmarshal(bodyBytes, &receivedBody)
		suite.NoError(err)

		response := CreateSessionResponse{
			PaymentSessionID: "ps-661f87c614802d6c402cd82d",
			Created:          time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated:          time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			ReferenceID:      receivedBody.ReferenceID,
			SessionType:      SessionTypePay,
			Currency:         "IDR",
			Amount:           150000,
			Country:          "ID",
			Mode:             SessionModePaymentLink,
			Status:           "ACTIVE",
			BusinessID:       "661f87c614802d6c402cd82d",
			PaymentLinkURL:   strPtr("https://dev.xen.to/qZx5RD_7"),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	req := validPaymentLinkRequest()

	resp, xerr := client.CreateSession(context.Background(), req)

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal("ps-661f87c614802d6c402cd82d", resp.PaymentSessionID)
	suite.Equal("ACTIVE", resp.Status)
	suite.NotNil(resp.PaymentLinkURL)
	suite.Equal("https://dev.xen.to/qZx5RD_7", *resp.PaymentLinkURL)

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("test-secret-key:"))
	suite.Equal(expectedAuth, receivedAuth)
	suite.Equal("order_12345", receivedBody.ReferenceID)
}

func (suite *EndpointTestSuite) TestCreateSession_ComponentsSuccess201() {
	var receivedBody CreateSessionRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		suite.NoError(err)

		err = json.Unmarshal(bodyBytes, &receivedBody)
		suite.NoError(err)

		response := CreateSessionResponse{
			PaymentSessionID: "ps-661f87c614802d6c402cd82e",
			Created:          time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated:          time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			ReferenceID:      receivedBody.ReferenceID,
			SessionType:      SessionTypePay,
			Currency:         "PHP",
			Amount:           10000,
			Country:          "PH",
			Mode:             SessionModeComponents,
			Status:           "ACTIVE",
			BusinessID:       "661f87c614802d6c402cd82d",
			ComponentsSDKKey: strPtr("f5RBHw4S3SeZV8KGlBxjy8tyN2Vg4klA5Bww7"),
			ComponentsConfiguration: &ComponentsConfiguration{
				Origins: []string{"https://yourcompany.com"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	req := validComponentsRequest()

	resp, xerr := client.CreateSession(context.Background(), req)

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal(SessionModeComponents, resp.Mode)
	suite.NotNil(resp.ComponentsSDKKey)
	suite.Equal("f5RBHw4S3SeZV8KGlBxjy8tyN2Vg4klA5Bww7", *resp.ComponentsSDKKey)
	suite.NotNil(receivedBody.ComponentsConfiguration)
	suite.Equal([]string{"https://yourcompany.com"}, receivedBody.ComponentsConfiguration.Origins)
}

func (suite *EndpointTestSuite) TestCreateSession_ValidationFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	req := validPaymentLinkRequest()
	req.ReferenceID = ""

	resp, xerr := client.CreateSession(context.Background(), req)

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
	suite.Equal("reference_id is required", xerr.Message)
}

func (suite *EndpointTestSuite) TestCreateSession_InvalidMode() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	req := validPaymentLinkRequest()
	req.Mode = "INVALID"

	resp, xerr := client.CreateSession(context.Background(), req)

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
	suite.Equal("mode must be PAYMENT_LINK or COMPONENTS", xerr.Message)
}

func (suite *EndpointTestSuite) TestCreateSession_APIValidationError() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(w, `{"error_code":"API_VALIDATION_ERROR","message":"request/body must have required property 'session_type'"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CreateSession(context.Background(), validPaymentLinkRequest())

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(APIValidationError, xerr.ErrorCode)
	suite.Equal(http.StatusBadRequest, xerr.HTTPStatus)
	suite.True(strings.Contains(xerr.Message, "session_type"))
}

func (suite *EndpointTestSuite) TestCreateSession_InvalidAPIKey() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := io.WriteString(w, `{"error_code":"INVALID_API_KEY","message":"API key is not authorized for this API"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CreateSession(context.Background(), validPaymentLinkRequest())

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidAPIKey, xerr.ErrorCode)
	suite.Equal(http.StatusUnauthorized, xerr.HTTPStatus)
}

func (suite *EndpointTestSuite) TestCreateSession_NetworkError() {
	client := suite.newTestClient("http://127.0.0.1:1")

	resp, xerr := client.CreateSession(context.Background(), validPaymentLinkRequest())

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InternalServerError, xerr.ErrorCode)
	suite.NotEmpty(xerr.Message)
}

func (suite *EndpointTestSuite) TestGetSession_Success200() {
	sessionID := "ps-661f87c614802d6c402cd82d"
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodGet, r.Method)
		suite.Equal("/sessions/"+sessionID, r.URL.Path)

		receivedAuth = r.Header.Get("Authorization")

		response := CreateSessionResponse{
			PaymentSessionID: sessionID,
			Created:          time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated:          time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			ReferenceID:      "order_12345",
			SessionType:      SessionTypePay,
			Currency:         "IDR",
			Amount:           150000,
			Country:          "ID",
			Mode:             SessionModePaymentLink,
			Status:           string(SessionStatusCompleted),
			BusinessID:       "661f87c614802d6c402cd82d",
			PaymentLinkURL:   strPtr("https://dev.xen.to/qZx5RD_7"),
			PaymentID:        strPtr("py-661f87c614802d6c402cd82e"),
			PaymentRequestID: strPtr("pr-661f87c614802d6c402cd82f"),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.GetSession(context.Background(), GetSessionRequest{
		PaymentSessionID: sessionID,
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal(sessionID, resp.PaymentSessionID)
	suite.Equal(string(SessionStatusCompleted), resp.Status)
	suite.NotNil(resp.PaymentID)
	suite.Equal("py-661f87c614802d6c402cd82e", *resp.PaymentID)

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("test-secret-key:"))
	suite.Equal(expectedAuth, receivedAuth)
}

func (suite *EndpointTestSuite) TestGetSession_ValidationFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.GetSession(context.Background(), GetSessionRequest{
		PaymentSessionID: "",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
	suite.Equal("payment_session_id is required", xerr.Message)
}

func (suite *EndpointTestSuite) TestGetSession_NotFound404() {
	sessionID := "ps-661f87c614802d6c402cd82d"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := io.WriteString(w, `{"error_code":"NOT_FOUND","message":"Session not found"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.GetSession(context.Background(), GetSessionRequest{
		PaymentSessionID: sessionID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(NotFound, xerr.ErrorCode)
	suite.Equal(http.StatusNotFound, xerr.HTTPStatus)
	suite.Equal("Session not found", xerr.Message)
}

func (suite *EndpointTestSuite) TestGetSession_InvalidAPIKey401() {
	sessionID := "ps-661f87c614802d6c402cd82d"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := io.WriteString(w, `{"error_code":"INVALID_API_KEY","message":"API key is not authorized for this API"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.GetSession(context.Background(), GetSessionRequest{
		PaymentSessionID: sessionID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidAPIKey, xerr.ErrorCode)
	suite.Equal(http.StatusUnauthorized, xerr.HTTPStatus)
}

func (suite *EndpointTestSuite) TestGetSession_NetworkError() {
	client := suite.newTestClient("http://127.0.0.1:1")

	resp, xerr := client.GetSession(context.Background(), GetSessionRequest{
		PaymentSessionID: "ps-661f87c614802d6c402cd82d",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InternalServerError, xerr.ErrorCode)
	suite.NotEmpty(xerr.Message)
}

func (suite *EndpointTestSuite) TestCancelSession_Success200() {
	sessionID := "ps-661f87c614802d6c402cd82d"
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodPost, r.Method)
		suite.Equal("/sessions/"+sessionID+"/cancel", r.URL.Path)

		receivedAuth = r.Header.Get("Authorization")

		response := CreateSessionResponse{
			PaymentSessionID: sessionID,
			Created:          time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated:          time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			ReferenceID:      "order_12345",
			SessionType:      SessionTypePay,
			Currency:         "IDR",
			Amount:           150000,
			Country:          "ID",
			Mode:             SessionModePaymentLink,
			Status:           string(SessionStatusCanceled),
			BusinessID:       "661f87c614802d6c402cd82d",
			PaymentLinkURL:   strPtr("https://dev.xen.to/qZx5RD_7"),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.CancelSession(context.Background(), CancelSessionRequest{
		PaymentSessionID: sessionID,
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal(sessionID, resp.PaymentSessionID)
	suite.Equal(string(SessionStatusCanceled), resp.Status)

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("test-secret-key:"))
	suite.Equal(expectedAuth, receivedAuth)
}

func (suite *EndpointTestSuite) TestCancelSession_ValidationFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.CancelSession(context.Background(), CancelSessionRequest{
		PaymentSessionID: "",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
	suite.Equal("payment_session_id is required", xerr.Message)
}

func (suite *EndpointTestSuite) TestCancelSession_NotFound404() {
	sessionID := "ps-661f87c614802d6c402cd82d"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := io.WriteString(w, `{"error_code":"NOT_FOUND","message":"Session not found"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CancelSession(context.Background(), CancelSessionRequest{
		PaymentSessionID: sessionID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(NotFound, xerr.ErrorCode)
	suite.Equal(http.StatusNotFound, xerr.HTTPStatus)
	suite.Equal("Session not found", xerr.Message)
}

func (suite *EndpointTestSuite) TestCancelSession_InvalidSessionStatus422() {
	sessionID := "ps-661f87c614802d6c402cd82d"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err := io.WriteString(w, `{"error_code":"INVALID_SESSION_STATUS","message":"Session is not in ACTIVE status"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CancelSession(context.Background(), CancelSessionRequest{
		PaymentSessionID: sessionID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidSessionStatus, xerr.ErrorCode)
	suite.Equal(http.StatusUnprocessableEntity, xerr.HTTPStatus)
}

func (suite *EndpointTestSuite) TestCancelSession_InvalidAPIKey401() {
	sessionID := "ps-661f87c614802d6c402cd82d"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := io.WriteString(w, `{"error_code":"INVALID_API_KEY","message":"API key is not authorized for this API"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CancelSession(context.Background(), CancelSessionRequest{
		PaymentSessionID: sessionID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidAPIKey, xerr.ErrorCode)
	suite.Equal(http.StatusUnauthorized, xerr.HTTPStatus)
}

func (suite *EndpointTestSuite) TestCancelSession_NetworkError() {
	client := suite.newTestClient("http://127.0.0.1:1")

	resp, xerr := client.CancelSession(context.Background(), CancelSessionRequest{
		PaymentSessionID: "ps-661f87c614802d6c402cd82d",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InternalServerError, xerr.ErrorCode)
	suite.NotEmpty(xerr.Message)
}

func (suite *EndpointTestSuite) TestGetPayment_Success200() {
	paymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodGet, r.Method)
		suite.Equal("/v3/payments/"+paymentID, r.URL.Path)

		receivedAuth = r.Header.Get("Authorization")

		response := GetPaymentResponse{
			PaymentID:        paymentID,
			BusinessID:       "5f27a14a9bf05c73dd040bc8",
			ReferenceID:      "order_12345",
			PaymentRequestID: strPtr("pr-1102feb0-bb79-47ae-9d1e-e69394d3949c"),
			Type:             PaymentTypePay,
			Country:          "ID",
			Currency:         "IDR",
			RequestAmount:    DecimalAmount{value: "1999.01"},
			ChannelCode:      ChannelCodeOVO,
			Status:           string(PaymentStatusSucceeded),
			Created:          time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated:          time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.GetPayment(context.Background(), GetPaymentRequest{
		PaymentID: paymentID,
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal(paymentID, resp.PaymentID)
	suite.Equal(string(PaymentStatusSucceeded), resp.Status)
	suite.Equal("1999.01", resp.RequestAmount.String())

	amount, err := resp.RequestAmount.Float64()
	suite.NoError(err)
	suite.Equal(1999.01, amount)

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("test-secret-key:"))
	suite.Equal(expectedAuth, receivedAuth)
}

func (suite *EndpointTestSuite) TestGetPayment_Success200_StringRequestAmount() {
	paymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{
			"payment_id": "py-1402feb0-bb79-47ae-9d1e-e69394d3949c",
			"business_id": "5f27a14a9bf05c73dd040bc8",
			"reference_id": "order_12345",
			"type": "PAY",
			"country": "ID",
			"currency": "IDR",
			"request_amount": "1999.01",
			"channel_code": "OVO",
			"status": "SUCCEEDED",
			"created": "2029-12-31T23:59:59Z",
			"updated": "2029-12-31T23:59:59Z"
		}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.GetPayment(context.Background(), GetPaymentRequest{
		PaymentID: paymentID,
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal("1999.01", resp.RequestAmount.String())
}

func (suite *EndpointTestSuite) TestGetPayment_ValidationFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.GetPayment(context.Background(), GetPaymentRequest{
		PaymentID: "",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
	suite.Equal("payment_id is required", xerr.Message)
}

func (suite *EndpointTestSuite) TestGetPayment_NotFound404() {
	paymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := io.WriteString(w, `{"error_code":"DATA_NOT_FOUND","message":"ID specified in request cannot be found."}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.GetPayment(context.Background(), GetPaymentRequest{
		PaymentID: paymentID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(DataNotFound, xerr.ErrorCode)
	suite.Equal(http.StatusNotFound, xerr.HTTPStatus)
}

func (suite *EndpointTestSuite) TestGetPayment_InvalidAPIKey401() {
	paymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := io.WriteString(w, `{"error_code":"INVALID_API_KEY","message":"API key is not authorized for this API"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.GetPayment(context.Background(), GetPaymentRequest{
		PaymentID: paymentID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidAPIKey, xerr.ErrorCode)
	suite.Equal(http.StatusUnauthorized, xerr.HTTPStatus)
}

func (suite *EndpointTestSuite) TestGetPayment_NetworkError() {
	client := suite.newTestClient("http://127.0.0.1:1")

	resp, xerr := client.GetPayment(context.Background(), GetPaymentRequest{
		PaymentID: "py-1402feb0-bb79-47ae-9d1e-e69394d3949c",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InternalServerError, xerr.ErrorCode)
	suite.NotEmpty(xerr.Message)
}

func (suite *EndpointTestSuite) TestCancelPayment_Success200() {
	paymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodPost, r.Method)
		suite.Equal("/v3/payments/"+paymentID+"/cancel", r.URL.Path)

		receivedAuth = r.Header.Get("Authorization")

		response := CancelPaymentResponse{
			PaymentID:     paymentID,
			BusinessID:    "5f27a14a9bf05c73dd040bc8",
			ReferenceID:   "order_12345",
			Type:          PaymentTypePay,
			Country:       "ID",
			Currency:      "IDR",
			RequestAmount: DecimalAmount{value: "1999.01"},
			ChannelCode:   "CARDS",
			Status:        string(PaymentStatusCanceled),
			Created:       time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated:       time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.CancelPayment(context.Background(), CancelPaymentRequest{
		PaymentID: paymentID,
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal(paymentID, resp.PaymentID)
	suite.Equal(string(PaymentStatusCanceled), resp.Status)

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("test-secret-key:"))
	suite.Equal(expectedAuth, receivedAuth)
}

func (suite *EndpointTestSuite) TestCancelPayment_ValidationFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.CancelPayment(context.Background(), CancelPaymentRequest{
		PaymentID: "",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCancelPayment_IneligibleStatus400() {
	paymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(w, `{"error_code":"INELIGIBLE_TRANSACTION_STATUS","message":"Payment is not in AUTHORIZED status"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CancelPayment(context.Background(), CancelPaymentRequest{
		PaymentID: paymentID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(IneligibleTransactionStatus, xerr.ErrorCode)
	suite.Equal(http.StatusBadRequest, xerr.HTTPStatus)
}

func (suite *EndpointTestSuite) TestCancelPayment_NotFound404() {
	paymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := io.WriteString(w, `{"error_code":"DATA_NOT_FOUND","message":"ID specified in request cannot be found."}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CancelPayment(context.Background(), CancelPaymentRequest{
		PaymentID: paymentID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(DataNotFound, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCancelPayment_NetworkError() {
	client := suite.newTestClient("http://127.0.0.1:1")

	resp, xerr := client.CancelPayment(context.Background(), CancelPaymentRequest{
		PaymentID: "py-1402feb0-bb79-47ae-9d1e-e69394d3949c",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InternalServerError, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCapturePayment_Success200() {
	paymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"
	var receivedBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodPost, r.Method)
		suite.Equal("/v3/payments/"+paymentID+"/capture", r.URL.Path)

		bodyBytes, err := io.ReadAll(r.Body)
		suite.NoError(err)

		err = json.Unmarshal(bodyBytes, &receivedBody)
		suite.NoError(err)

		response := CapturePaymentResponse{
			PaymentID:     paymentID,
			BusinessID:    "5f27a14a9bf05c73dd040bc8",
			ReferenceID:   "order_12345",
			Type:          PaymentTypePay,
			Country:       "ID",
			Currency:      "IDR",
			RequestAmount: DecimalAmount{value: "1999.01"},
			ChannelCode:   "CARDS",
			Status:        string(PaymentStatusSucceeded),
			Captures: []PaymentCapture{
				{
					CaptureID:        "cap-1502feb0-bb79-47ae-9d1e-e69394d3949c",
					CaptureTimestamp: time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
					CaptureAmount:    DecimalAmount{value: "1999.01"},
				},
			},
			Created: time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated: time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.CapturePayment(context.Background(), CapturePaymentRequest{
		PaymentID:     paymentID,
		CaptureAmount: NewDecimalAmountFromFloat(1999.01),
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal(paymentID, resp.PaymentID)
	suite.Equal(string(PaymentStatusSucceeded), resp.Status)
	suite.Len(resp.Captures, 1)
	suite.Equal("1999.01", receivedBody["capture_amount"])
}

func (suite *EndpointTestSuite) TestCapturePayment_ValidationFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.CapturePayment(context.Background(), CapturePaymentRequest{
		PaymentID:     "py-1402feb0-bb79-47ae-9d1e-e69394d3949c",
		CaptureAmount: NewDecimalAmountFromFloat(0),
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
	suite.Equal("capture_amount must be greater than zero", xerr.Message)
}

func (suite *EndpointTestSuite) TestCapturePayment_CaptureAmountExceeded400() {
	paymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(w, `{"error_code":"CAPTURE_AMOUNT_EXCEEDED","message":"Capture amount exceeds authorized amount"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CapturePayment(context.Background(), CapturePaymentRequest{
		PaymentID:     paymentID,
		CaptureAmount: NewDecimalAmountFromFloat(9999.99),
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(CaptureAmountExceeded, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCapturePayment_NotFound404() {
	paymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := io.WriteString(w, `{"error_code":"DATA_NOT_FOUND","message":"ID specified in request cannot be found."}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CapturePayment(context.Background(), CapturePaymentRequest{
		PaymentID:     paymentID,
		CaptureAmount: NewDecimalAmountFromFloat(1999.01),
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(DataNotFound, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCapturePayment_NetworkError() {
	client := suite.newTestClient("http://127.0.0.1:1")

	resp, xerr := client.CapturePayment(context.Background(), CapturePaymentRequest{
		PaymentID:     "py-1402feb0-bb79-47ae-9d1e-e69394d3949c",
		CaptureAmount: NewDecimalAmountFromFloat(1999.01),
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InternalServerError, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCreatePaymentRequest_Success201() {
	var receivedAuth string

	type createPaymentRequestBody struct {
		ReferenceID       string          `json:"reference_id"`
		ChannelCode       string          `json:"channel_code"`
		ChannelProperties json.RawMessage `json:"channel_properties"`
	}

	var receivedBody createPaymentRequestBody

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodPost, r.Method)
		suite.Equal("/v3/payment_requests", r.URL.Path)

		receivedAuth = r.Header.Get("Authorization")

		bodyBytes, err := io.ReadAll(r.Body)
		suite.NoError(err)

		err = json.Unmarshal(bodyBytes, &receivedBody)
		suite.NoError(err)

		response := CreatePaymentRequestResponse{
			PaymentRequestID: "pr-90392f42-d98a-49ef-a7f3-abcezas123",
			BusinessID:       "5f27a14a9bf05c73dd040bc8",
			ReferenceID:      receivedBody.ReferenceID,
			Type:             PaymentTypePay,
			Country:          "ID",
			Currency:         "IDR",
			RequestAmount:    DecimalAmount{value: "150000"},
			ChannelCode:      ChannelCodeOVO,
			ChannelProperties: PaymentRequestChannelProperties{
				value: map[string]any{
					"account_mobile_number": "+628123456789",
				},
			},
			Actions: []PaymentRequestAction{
				{
					Type:       PaymentRequestActionTypeRedirectCustomer,
					Value:      "https://xendit.co/example",
					Descriptor: PaymentRequestActionDescriptorWebURL,
				},
			},
			Status:  string(PaymentRequestStatusRequiresAction),
			Created: time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated: time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	req := validCreatePaymentRequestRequest()

	resp, xerr := client.CreatePaymentRequest(context.Background(), req)

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal("pr-90392f42-d98a-49ef-a7f3-abcezas123", resp.PaymentRequestID)
	suite.Equal(string(PaymentRequestStatusRequiresAction), resp.Status)
	suite.Len(resp.Actions, 1)
	suite.Equal(PaymentRequestActionTypeRedirectCustomer, resp.Actions[0].Type)
	suite.Equal("https://xendit.co/example", resp.Actions[0].Value)

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("test-secret-key:"))
	suite.Equal(expectedAuth, receivedAuth)
	suite.Equal("order_12345", receivedBody.ReferenceID)
	suite.Equal(ChannelCodeOVO, receivedBody.ChannelCode)
	suite.JSONEq(`{"account_mobile_number":"+628123456789"}`, string(receivedBody.ChannelProperties))
}

func (suite *EndpointTestSuite) TestCreatePaymentRequest_Success200() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := CreatePaymentRequestResponse{
			PaymentRequestID: "pr-90392f42-d98a-49ef-a7f3-abcezas123",
			BusinessID:       "5f27a14a9bf05c73dd040bc8",
			ReferenceID:      "order_12345",
			Type:             PaymentTypePay,
			Country:          "ID",
			Currency:         "IDR",
			RequestAmount:    DecimalAmount{value: "150000"},
			ChannelCode:      ChannelCodeOVO,
			Status:           string(PaymentRequestStatusPending),
			Created:          time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated:          time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CreatePaymentRequest(context.Background(), validCreatePaymentRequestRequest())

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal(string(PaymentRequestStatusPending), resp.Status)
}

func (suite *EndpointTestSuite) TestCreatePaymentRequest_ValidationFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	req := validCreatePaymentRequestRequest()
	req.ChannelCode = ""

	resp, xerr := client.CreatePaymentRequest(context.Background(), req)

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCreatePaymentRequest_InvalidChannel() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(w, `{"error_code":"INVALID_PAYMENT_CHANNEL","message":"The payment channel is invalid."}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CreatePaymentRequest(context.Background(), validCreatePaymentRequestRequest())

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidPaymentChannel, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCreatePaymentRequest_NetworkError() {
	client := suite.newTestClient("http://127.0.0.1:1")

	resp, xerr := client.CreatePaymentRequest(context.Background(), validCreatePaymentRequestRequest())

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InternalServerError, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestGetPaymentRequest_Success200() {
	paymentRequestID := "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c"
	latestPaymentID := "py-1402feb0-bb79-47ae-9d1e-e69394d3949c"
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodGet, r.Method)
		suite.Equal("/v3/payment_requests/"+paymentRequestID, r.URL.Path)

		receivedAuth = r.Header.Get("Authorization")

		response := CreatePaymentRequestResponse{
			PaymentRequestID: paymentRequestID,
			BusinessID:       "5f27a14a9bf05c73dd040bc8",
			ReferenceID:      "order_12345",
			LatestPaymentID:  strPtr(latestPaymentID),
			Type:             PaymentTypePay,
			Country:          "ID",
			Currency:         "IDR",
			RequestAmount:    DecimalAmount{value: "150000"},
			ChannelCode:      ChannelCodeOVO,
			ChannelProperties: PaymentRequestChannelProperties{
				value: map[string]any{
					"account_mobile_number": "+628123456789",
				},
			},
			Status:  string(PaymentRequestStatusSucceeded),
			Created: time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated: time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.GetPaymentRequest(context.Background(), GetPaymentRequestStatusRequest{
		PaymentRequestID: paymentRequestID,
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal(paymentRequestID, resp.PaymentRequestID)
	suite.Equal(string(PaymentRequestStatusSucceeded), resp.Status)
	suite.NotNil(resp.LatestPaymentID)
	suite.Equal(latestPaymentID, *resp.LatestPaymentID)
	suite.Equal("150000", resp.RequestAmount.String())

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("test-secret-key:"))
	suite.Equal(expectedAuth, receivedAuth)
}

func (suite *EndpointTestSuite) TestGetPaymentRequest_Success200_StringRequestAmount() {
	paymentRequestID := "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{
			"payment_request_id": "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c",
			"business_id": "5f27a14a9bf05c73dd040bc8",
			"reference_id": "order_12345",
			"type": "PAY",
			"country": "ID",
			"currency": "IDR",
			"request_amount": "150000.50",
			"channel_code": "OVO",
			"status": "SUCCEEDED",
			"created": "2029-12-31T23:59:59Z",
			"updated": "2029-12-31T23:59:59Z"
		}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.GetPaymentRequest(context.Background(), GetPaymentRequestStatusRequest{
		PaymentRequestID: paymentRequestID,
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal("150000.50", resp.RequestAmount.String())
}

func (suite *EndpointTestSuite) TestGetPaymentRequest_ValidationFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.GetPaymentRequest(context.Background(), GetPaymentRequestStatusRequest{
		PaymentRequestID: "",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
	suite.Equal("payment_request_id is required", xerr.Message)
}

func (suite *EndpointTestSuite) TestGetPaymentRequest_NotFound404() {
	paymentRequestID := "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := io.WriteString(w, `{"error_code":"DATA_NOT_FOUND","message":"ID specified in request cannot be found."}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.GetPaymentRequest(context.Background(), GetPaymentRequestStatusRequest{
		PaymentRequestID: paymentRequestID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(DataNotFound, xerr.ErrorCode)
	suite.Equal(http.StatusNotFound, xerr.HTTPStatus)
}

func (suite *EndpointTestSuite) TestGetPaymentRequest_InvalidAPIKey401() {
	paymentRequestID := "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := io.WriteString(w, `{"error_code":"INVALID_API_KEY","message":"API key is not authorized for this API"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.GetPaymentRequest(context.Background(), GetPaymentRequestStatusRequest{
		PaymentRequestID: paymentRequestID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidAPIKey, xerr.ErrorCode)
	suite.Equal(http.StatusUnauthorized, xerr.HTTPStatus)
}

func (suite *EndpointTestSuite) TestGetPaymentRequest_NetworkError() {
	client := suite.newTestClient("http://127.0.0.1:1")

	resp, xerr := client.GetPaymentRequest(context.Background(), GetPaymentRequestStatusRequest{
		PaymentRequestID: "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InternalServerError, xerr.ErrorCode)
	suite.NotEmpty(xerr.Message)
}

func (suite *EndpointTestSuite) TestCancelPaymentRequest_Success200() {
	paymentRequestID := "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c"
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodPost, r.Method)
		suite.Equal("/v3/payment_requests/"+paymentRequestID+"/cancel", r.URL.Path)

		receivedAuth = r.Header.Get("Authorization")

		response := CreatePaymentRequestResponse{
			PaymentRequestID: paymentRequestID,
			BusinessID:       "5f27a14a9bf05c73dd040bc8",
			ReferenceID:      "order_12345",
			Type:             PaymentTypePay,
			Country:          "ID",
			Currency:         "IDR",
			RequestAmount:    DecimalAmount{value: "150000"},
			ChannelCode:      ChannelCodeDANA,
			Status:           string(PaymentRequestStatusCanceled),
			Created:          time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
			Updated:          time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.CancelPaymentRequest(context.Background(), CancelPaymentRequestStatusRequest{
		PaymentRequestID: paymentRequestID,
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal(paymentRequestID, resp.PaymentRequestID)
	suite.Equal(string(PaymentRequestStatusCanceled), resp.Status)

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("test-secret-key:"))
	suite.Equal(expectedAuth, receivedAuth)
}

func (suite *EndpointTestSuite) TestCancelPaymentRequest_ValidationFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.CancelPaymentRequest(context.Background(), CancelPaymentRequestStatusRequest{
		PaymentRequestID: "",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCancelPaymentRequest_IneligibleStatus400() {
	paymentRequestID := "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(w, `{"error_code":"INELIGIBLE_TRANSACTION_STATUS","message":"Payment request is not in a cancelable status"}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CancelPaymentRequest(context.Background(), CancelPaymentRequestStatusRequest{
		PaymentRequestID: paymentRequestID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(IneligibleTransactionStatus, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCancelPaymentRequest_NotFound404() {
	paymentRequestID := "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := io.WriteString(w, `{"error_code":"DATA_NOT_FOUND","message":"ID specified in request cannot be found."}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	resp, xerr := client.CancelPaymentRequest(context.Background(), CancelPaymentRequestStatusRequest{
		PaymentRequestID: paymentRequestID,
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(DataNotFound, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestCancelPaymentRequest_NetworkError() {
	client := suite.newTestClient("http://127.0.0.1:1")

	resp, xerr := client.CancelPaymentRequest(context.Background(), CancelPaymentRequestStatusRequest{
		PaymentRequestID: "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InternalServerError, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestSimulatePaymentRequest_Success200() {
	paymentRequestID := "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c"
	var receivedAuth string
	var receivedBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal(http.MethodPost, r.Method)
		suite.Equal("/v3/payment_requests/"+paymentRequestID+"/simulate", r.URL.Path)

		receivedAuth = r.Header.Get("Authorization")

		bodyBytes, err := io.ReadAll(r.Body)
		suite.NoError(err)
		err = json.Unmarshal(bodyBytes, &receivedBody)
		suite.NoError(err)

		response := SimulatePaymentRequestResponse{
			Status:  "PENDING",
			Message: "A simulated payment for the specified payment request id is being processed. You will be informed of the result via a webhook.",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)
	amount := NewDecimalAmountFromFloat(150000)

	resp, xerr := client.SimulatePaymentRequest(context.Background(), SimulatePaymentRequestRequest{
		PaymentRequestID: paymentRequestID,
		Amount:           &amount,
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal("PENDING", resp.Status)
	suite.NotEmpty(resp.Message)
	suite.Equal("150000", receivedBody["amount"])

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("test-secret-key:"))
	suite.Equal(expectedAuth, receivedAuth)
}

func (suite *EndpointTestSuite) TestSimulatePaymentRequest_Success200_EmptyBody() {
	paymentRequestID := "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c"
	var receivedBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		suite.NoError(err)
		err = json.Unmarshal(bodyBytes, &receivedBody)
		suite.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = io.WriteString(w, `{"status":"PENDING","message":"A simulated payment for the specified payment request id is being processed. You will be informed of the result via a webhook."}`)
		suite.NoError(err)
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.SimulatePaymentRequest(context.Background(), SimulatePaymentRequestRequest{
		PaymentRequestID: paymentRequestID,
	})

	suite.Nil(xerr)
	suite.NotNil(resp)
	suite.Equal("PENDING", resp.Status)
	suite.Empty(receivedBody)
}

func (suite *EndpointTestSuite) TestSimulatePaymentRequest_ValidationFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("request should not reach server when validation fails")
	}))
	defer server.Close()

	client := suite.newTestClient(server.URL)

	resp, xerr := client.SimulatePaymentRequest(context.Background(), SimulatePaymentRequestRequest{
		PaymentRequestID: "",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InvalidRequest, xerr.ErrorCode)
}

func (suite *EndpointTestSuite) TestSimulatePaymentRequest_NetworkError() {
	client := suite.newTestClient("http://127.0.0.1:1")

	resp, xerr := client.SimulatePaymentRequest(context.Background(), SimulatePaymentRequestRequest{
		PaymentRequestID: "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c",
	})

	suite.Nil(resp)
	suite.NotNil(xerr)
	suite.Equal(InternalServerError, xerr.ErrorCode)
}

func strPtr(value string) *string {
	return &value
}
