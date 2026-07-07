package xendit

import (
	"testing"

	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	testLogger *logger.Logger
}

func (suite *ClientTestSuite) SetupTest() {
	suite.testLogger = logger.New(logger.Config{
		Level:      "info",
		JSONOutput: false,
	})
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (suite *ClientTestSuite) TestNew_MissingSecretAPIKey() {
	client, err := New(Config{}, suite.testLogger)

	suite.Nil(client)
	suite.EqualError(err, "secret API key is required")
}

func (suite *ClientTestSuite) TestNew_MissingLogger() {
	client, err := New(Config{
		SecretAPIKey: "test-secret-key",
	}, nil)

	suite.Nil(client)
	suite.EqualError(err, "logger is required")
}

func (suite *ClientTestSuite) TestNew_DefaultBaseURL() {
	client, err := New(Config{
		SecretAPIKey: "test-secret-key",
	}, suite.testLogger)

	suite.NoError(err)
	suite.NotNil(client)
	suite.Equal(defaultBaseURL, client.config.BaseURL)
}

func (suite *ClientTestSuite) TestNew_CustomBaseURL() {
	customBaseURL := "https://custom.xendit.test"

	client, err := New(Config{
		BaseURL:      customBaseURL,
		SecretAPIKey: "test-secret-key",
	}, suite.testLogger)

	suite.NoError(err)
	suite.NotNil(client)
	suite.Equal(customBaseURL, client.config.BaseURL)
}

func (suite *ClientTestSuite) TestNew_AuthHeader() {
	secretKey := "test-secret-key"

	client, err := New(Config{
		SecretAPIKey: secretKey,
	}, suite.testLogger)

	suite.NoError(err)
	suite.Equal("Basic dGVzdC1zZWNyZXQta2V5Og==", client.authHeader)
}

func validPaymentLinkRequest() CreateSessionRequest {
	successURL := "https://example.com/success"
	cancelURL := "https://example.com/cancel"

	return CreateSessionRequest{
		ReferenceID:      "order_12345",
		Currency:         "IDR",
		Amount:           150000,
		Country:          "ID",
		SessionType:      SessionTypePay,
		Mode:             SessionModePaymentLink,
		SuccessReturnURL: &successURL,
		CancelReturnURL:  &cancelURL,
	}
}

func validComponentsRequest() CreateSessionRequest {
	return CreateSessionRequest{
		ReferenceID: "order_12345",
		Currency:    "PHP",
		Amount:      10000,
		Country:     "PH",
		SessionType: SessionTypePay,
		Mode:        SessionModeComponents,
		ComponentsConfiguration: &ComponentsConfiguration{
			Origins: []string{"https://yourcompany.com"},
		},
	}
}

func TestCreateSessionRequest_Validate_PaymentLink(t *testing.T) {
	req := validPaymentLinkRequest()
	if err := req.Validate(); err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}
}

func TestCreateSessionRequest_Validate_Components(t *testing.T) {
	req := validComponentsRequest()
	if err := req.Validate(); err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}
}

func TestCreateSessionRequest_Validate_MissingReferenceID(t *testing.T) {
	req := validPaymentLinkRequest()
	req.ReferenceID = ""

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "reference_id is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateSessionRequest_Validate_InvalidSessionType(t *testing.T) {
	req := validPaymentLinkRequest()
	req.SessionType = "SAVE"

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "session_type must be PAY" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateSessionRequest_Validate_InvalidMode(t *testing.T) {
	req := validPaymentLinkRequest()
	req.Mode = "INVALID"

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "mode must be PAYMENT_LINK or COMPONENTS" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateSessionRequest_Validate_InvalidReturnURL(t *testing.T) {
	httpURL := "http://example.com/success"

	req := validPaymentLinkRequest()
	req.SuccessReturnURL = &httpURL

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "success_return_url must use HTTPS" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateSessionRequest_Validate_ComponentsMissingConfiguration(t *testing.T) {
	req := validComponentsRequest()
	req.ComponentsConfiguration = nil

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "components_configuration is required when mode is COMPONENTS" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateSessionRequest_Validate_ComponentsEmptyOrigins(t *testing.T) {
	req := validComponentsRequest()
	req.ComponentsConfiguration = &ComponentsConfiguration{
		Origins: []string{},
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "components_configuration.origins is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateSessionRequest_Validate_ComponentsNonHTTPSOrigin(t *testing.T) {
	req := validComponentsRequest()
	req.ComponentsConfiguration = &ComponentsConfiguration{
		Origins: []string{"http://yourcompany.com"},
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "components_configuration.origins[0]: must use HTTPS" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateSessionRequest_Validate_PaymentLinkWithComponentsConfiguration(t *testing.T) {
	req := validPaymentLinkRequest()
	req.ComponentsConfiguration = &ComponentsConfiguration{
		Origins: []string{"https://yourcompany.com"},
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "components_configuration is not allowed when mode is PAYMENT_LINK" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func validGetSessionRequest() GetSessionRequest {
	return GetSessionRequest{
		PaymentSessionID: "ps-661f87c614802d6c402cd82d",
	}
}

func TestGetSessionRequest_Validate(t *testing.T) {
	req := validGetSessionRequest()

	err := req.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetSessionRequest_Validate_MissingPaymentSessionID(t *testing.T) {
	req := validGetSessionRequest()
	req.PaymentSessionID = ""

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "payment_session_id is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetSessionRequest_Validate_InvalidLength(t *testing.T) {
	req := validGetSessionRequest()
	req.PaymentSessionID = "ps-too-short"

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "payment_session_id must be 27 characters" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func validCancelSessionRequest() CancelSessionRequest {
	return CancelSessionRequest{
		PaymentSessionID: "ps-661f87c614802d6c402cd82d",
	}
}

func TestCancelSessionRequest_Validate(t *testing.T) {
	req := validCancelSessionRequest()

	err := req.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCancelSessionRequest_Validate_MissingPaymentSessionID(t *testing.T) {
	req := validCancelSessionRequest()
	req.PaymentSessionID = ""

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "payment_session_id is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCancelSessionRequest_Validate_InvalidLength(t *testing.T) {
	req := validCancelSessionRequest()
	req.PaymentSessionID = "ps-too-short"

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "payment_session_id must be 27 characters" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func validGetPaymentRequest() GetPaymentRequest {
	return GetPaymentRequest{
		PaymentID: "py-1402feb0-bb79-47ae-9d1e-e69394d3949c",
	}
}

func TestGetPaymentRequest_Validate(t *testing.T) {
	req := validGetPaymentRequest()

	err := req.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetPaymentRequest_Validate_MissingPaymentID(t *testing.T) {
	req := validGetPaymentRequest()
	req.PaymentID = ""

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "payment_id is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetPaymentRequest_Validate_InvalidLength(t *testing.T) {
	req := validGetPaymentRequest()
	req.PaymentID = "py-too-short"

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "payment_id must be 39 characters" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func validCancelPaymentRequest() CancelPaymentRequest {
	return CancelPaymentRequest{
		PaymentID: "py-1402feb0-bb79-47ae-9d1e-e69394d3949c",
	}
}

func TestCancelPaymentRequest_Validate(t *testing.T) {
	req := validCancelPaymentRequest()

	err := req.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCancelPaymentRequest_Validate_MissingPaymentID(t *testing.T) {
	req := validCancelPaymentRequest()
	req.PaymentID = ""

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "payment_id is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func validCapturePaymentRequest() CapturePaymentRequest {
	return CapturePaymentRequest{
		PaymentID:     "py-1402feb0-bb79-47ae-9d1e-e69394d3949c",
		CaptureAmount: NewDecimalAmountFromFloat(1999.01),
	}
}

func TestCapturePaymentRequest_Validate(t *testing.T) {
	req := validCapturePaymentRequest()

	err := req.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCapturePaymentRequest_Validate_MissingPaymentID(t *testing.T) {
	req := validCapturePaymentRequest()
	req.PaymentID = ""

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "payment_id is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCapturePaymentRequest_Validate_InvalidCaptureAmount(t *testing.T) {
	req := validCapturePaymentRequest()
	req.CaptureAmount = NewDecimalAmountFromFloat(0)

	err := req.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	if err.Error() != "capture_amount must be greater than zero" {
		t.Fatalf("unexpected error: %v", err)
	}
}
