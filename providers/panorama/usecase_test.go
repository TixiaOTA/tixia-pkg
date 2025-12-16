package panorama

import (
	"context"
	"encoding/json"
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

func (suite *UsecaseTestSuite) TestSearchBookingRequest_Success() {
	req := SearchBookingRequest{
		JourneyList: []struct {
			Origin        string "json:\"origin\""
			Destination   string "json:\"destination\""
			DepartureDate string "json:\"departure_date\""
		}{
			{
				Origin:        "CGK",
				Destination:   "SIN",
				DepartureDate: "2025-12-19",
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

	for _, v := range req.JourneyList {
		suite.NotNil(v.Origin)
		suite.NotNil(v.Destination)
		suite.NotNil(v.DepartureDate)
	}
	suite.NotNil(req.Adult)
	suite.NotNil(req.CabinClass)
	suite.NotNil(req.Provider)

	resp, err := suite.usecase.Search(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(resp)

	res, errM := json.Marshal(resp)
	suite.Nil(errM)

	fmt.Printf("SearchBookingRequest successful - Response: %v\n", string(res))
}
