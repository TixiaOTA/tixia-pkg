package panorama

import (
	"context"
	"fmt"
	"testing"

	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/stretchr/testify/suite"
)

// Test data
var (
	testOrderNumbers = []string{"test-order-number"}
)

type UsecaseTestSuite struct {
	suite.Suite
	usecase *Usecase
	client  *Client
}

func (suite *UsecaseTestSuite) SetupTest() {
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

	suite.usecase = NewUsecase(client)
}

func TestUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseTestSuite))
}

func (suite *UsecaseTestSuite) TestCancelBookingRequest_Success() {
	req := CancelBookingRequest{
		OrderNumbers: testOrderNumbers,
	}

	suite.NotNil(req.OrderNumbers)

	resp, err := suite.usecase.CancelBooking(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(resp)

	fmt.Printf("CancelBookingRequest successful - Results: %v\n", resp.Results)
}
