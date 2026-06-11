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

func strPtr(value string) *string {
	return &value
}
