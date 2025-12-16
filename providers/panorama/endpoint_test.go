package panorama

import (
	"context"
	"fmt"
	"testing"

	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/stretchr/testify/suite"
)

// Test configuration variables - Update these with your test credentials
var (
	testBaseURL  = "test-base-url" // Base URL for the API
	testApiKey   = "test-api-key"  // Your API key
	testUser     = "test-user"     // Your user
	testPassword = "test-password" // Your password

	// Test data for GetBooking
	testOrderNumber = "test-order-number" // Test order number
)

type EndpointTestSuite struct {
	suite.Suite
	client *Client
}

func (suite *EndpointTestSuite) SetupTest() {
	config := Config{
		BaseURL: testBaseURL,
		Credentials: Credentials{
			APIKey:   testApiKey,
			User:     testUser,
			Password: testPassword,
		},
	}

	logger := logger.New(logger.Config{
		Level:      "info",
		JSONOutput: false,
		FilePath:   "",
	})

	client, err := New(config, logger)
	suite.NoError(err)
	suite.client = client
}

func TestEndpointTestSuite(t *testing.T) {
	suite.Run(t, new(EndpointTestSuite))
}

func (suite *EndpointTestSuite) TestSignIn_Success() {
	resp, err := suite.client.SignIn(context.Background())

	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(resp.Result.ErrorCode, Success)
	suite.NotEmpty(resp.Result.Response.Signature)

	fmt.Printf("SignIn successful - Signature: %s\n", resp.Result.Response.Signature)
}

func (suite *EndpointTestSuite) TestSignIn_Error() {
	logger := logger.New(logger.Config{
		Level:      "info",
		JSONOutput: false,
		FilePath:   "",
	})

	client, clientErr := New(
		Config{
			BaseURL: suite.client.config.BaseURL,
			Credentials: Credentials{
				APIKey:   "test-error",
				User:     "test-error",
				Password: "test-error",
			},
		},
		logger,
	)

	suite.NoError(clientErr)

	resp, err := client.SignIn(context.Background())

	suite.Nil(resp)
	suite.NotNil(err)
	suite.NotEmpty(err.Result.ErrorMsg)
	suite.Equal(err.Result.ErrorCode, AccessDenied)

	fmt.Printf("SignIn error - Error: %v\n", err.Result)
}

func (suite *EndpointTestSuite) TestSearch_Success() {
	// First, sign in to get the signature
	signInResp, signInErr := suite.client.SignIn(context.Background())

	suite.Nil(signInErr)
	suite.NotNil(signInResp)
	suite.Equal(signInResp.Result.ErrorCode, Success)
	suite.NotEmpty(signInResp.Result.Response.Signature)

	// Now test Search with the obtained signature
	req := SearchRequest{
		Signature: signInResp.Result.Response.Signature,
		JourneyList: []JourneyListRequest{
			{
				Origin:        "CGK",
				Destination:   "SIN",
				DepartureDate: "2025-12-19",
			},
			{
				Origin:        "SIN",
				Destination:   "CGK",
				DepartureDate: "2025-12-30",
			},
		},
		Adult:        1,
		Child:        0,
		Infant:       0,
		CabinClass:   "Y",
		CarrierCodes: []string{"SQ"},
		IsComboPrice: true,
		Provider:     "sia",
		PromoCodes:   []string{},
	}

	resp, err := suite.client.Search(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(resp.Result.ErrorCode, Success)
	for _, s := range resp.Result.Response.Schedules {
		for _, j := range s.Journeys {
			suite.NotEmpty(j.JourneyCode)
			suite.Equal(j.JourneyCode, j.JourneyCode)

			fmt.Printf("Search successful - Sequence: %d, Journey Code: %s, Journey Ref ID: %s\n",
				j.Sequence, j.JourneyCode, j.JourneyRefID)
		}
	}
}

func (suite *EndpointTestSuite) TestGetBooking_Success() {
	// First, sign in to get the signature
	signInResp, signInErr := suite.client.SignIn(context.Background())

	suite.Nil(signInErr)
	suite.NotNil(signInResp)
	suite.Equal(signInResp.Result.ErrorCode, Success)
	suite.NotEmpty(signInResp.Result.Response.Signature)

	// Now test GetBooking with the obtained signature
	req := GetBookingRequest{
		Signature:   signInResp.Result.Response.Signature,
		OrderNumber: testOrderNumber,
	}

	resp, err := suite.client.GetBooking(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(resp.Result.ErrorCode, Success)
	suite.NotEmpty(resp.Result.Response.OrderNumber)
	suite.Equal(resp.Result.Response.OrderNumber, testOrderNumber)

	fmt.Printf("GetBooking successful - Order Number: %s, State: %s\n",
		resp.Result.Response.OrderNumber, resp.Result.Response.State)
}

func (suite *EndpointTestSuite) TestCancel_Success() {
	// First, sign in to get the signature
	signInResp, signInErr := suite.client.SignIn(context.Background())

	suite.Nil(signInErr)
	suite.NotNil(signInResp)
	suite.Equal(signInResp.Result.ErrorCode, Success)
	suite.NotEmpty(signInResp.Result.Response.Signature)

	// Now test GetBooking with the obtained signature
	req := CancelVoidRequest{
		Signature:   signInResp.Result.Response.Signature,
		OrderNumber: testOrderNumber,
	}

	resp, err := suite.client.CancelVoid(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(resp.Result.ErrorCode, Success)
	suite.NotEmpty(resp.Result.Response.OrderNumber)
	suite.Equal(resp.Result.Response.OrderNumber, testOrderNumber)
	suite.Equal(resp.Result.Response.State, "void")

	fmt.Printf("GetBooking successful - Order Number: %s, State: %s\n",
		resp.Result.Response.OrderNumber, resp.Result.Response.State)
}
