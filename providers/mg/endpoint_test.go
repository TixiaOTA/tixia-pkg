package mg

import (
	"context"
	"fmt"
	"testing"

	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/stretchr/testify/suite"
)

// Test configuration variables - Update these with your test credentials
var (
	testBaseURL    = "http://test-url.com" // Base URL for the API
	testAgencyCode = "test-agency-code"    // Your agency code
	testUsername   = "test-username"       // Your username
	testPassword   = "test-password"       // Your password

	// Test data for GetHotels
	testNationality  = "ID"
	testCountry      = "ID"
	testCity         = "ID-DPS"
	testCheckIn      = "2026-01-22"
	testCheckOut     = "2026-01-23"
	testCurrency     = "IDR"
	testLanguage     = "EN"
	testAvailFlag    = true
	testMaxNoOfHotel = 10
	testDetailLevel  = "Full"
)

type EndpointTestSuite struct {
	suite.Suite
	client *Client
}

func (suite *EndpointTestSuite) SetupTest() {
	config := Config{
		BaseURL: testBaseURL,
		Credentials: Credentials{
			AgencyCode: testAgencyCode,
			Username:   testUsername,
			Password:   testPassword,
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

func (suite *EndpointTestSuite) TestSearchHotel_Success() {
	rooms := &struct {
		Room []RoomRequest `json:"Room"`
	}{
		Room: []RoomRequest{
			{
				RoomNo:     1,
				NoOfAdults: 2,
				NoOfChild:  0,
				Child1Age:  0,
				Child2Age:  0,
				ExtraBed:   false,
			},
		},
	}

	req := SearchHotelRequest{
		Login: Credentials{
			AgencyCode: testAgencyCode,
			Username:   testUsername,
			Password:   testPassword,
		},
		Nationality:  testNationality,
		Country:      testCountry,
		City:         testCity,
		CheckIn:      testCheckIn,
		CheckOut:     testCheckOut,
		Rooms:        rooms,
		Currency:     testCurrency,
		Language:     testLanguage,
		AvailFlag:    testAvailFlag,
		MaxNoOfHotel: testMaxNoOfHotel,
		DetailLevel:  testDetailLevel,
	}

	resp, err := suite.client.SearchHotel(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(resp)
	suite.NotNil(resp.Hotels)

	fmt.Printf("Received %d hotels\n", len(resp.Hotels.Hotel))
	for i, hotel := range resp.Hotels.Hotel {
		if i < 5 {
			fmt.Printf("Hotel %d: Code=%s, Name=%s\n",
				i+1, hotel.Code, hotel.Name)
		}
	}
}
