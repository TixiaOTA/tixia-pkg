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
	testBaseURL    = "http://test-base-url.com" // Base URL for the API
	testAgencyCode = "test-agency-code"         // Your agency code
	testUsername   = "test-username"            // Your username
	testPassword   = "test-password"            // Your password

	// Test data for SearchHotel
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
	testRoomRequest  = RoomRequest{
		RoomNo:     1,
		NoOfAdults: 2,
		NoOfChild:  0,
		Child1Age:  0,
		Child2Age:  0,
		ExtraBed:   false,
	}
	testAgencyBookingID = "TIXIA-BOOKING-12345"
	testReservationDate = "2025-11-12"
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
			testRoomRequest,
		},
	}

	req := SearchHotelRequestDTO{
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

func (suite *EndpointTestSuite) TestRecheckHotel_Success() {
	// First, run SearchHotel to get a session ID
	searchRooms := &struct {
		Room []RoomRequest `json:"Room"`
	}{
		Room: []RoomRequest{
			testRoomRequest,
		},
	}

	searchReq := SearchHotelRequestDTO{
		Nationality:  testNationality,
		Country:      testCountry,
		City:         testCity,
		CheckIn:      testCheckIn,
		CheckOut:     testCheckOut,
		Rooms:        searchRooms,
		Currency:     testCurrency,
		Language:     testLanguage,
		AvailFlag:    testAvailFlag,
		MaxNoOfHotel: testMaxNoOfHotel,
		DetailLevel:  testDetailLevel,
	}

	searchResp, err := suite.client.SearchHotel(context.Background(), searchReq)
	suite.Nil(err)
	suite.NotNil(searchResp)
	suite.NotEmpty(searchResp.SessionID)

	// Extract session ID and hotel/room details from SearchHotel response
	sessionID := searchResp.SessionID
	var hotelCode string
	var roomDetails *RoomDetailsRequest

	if len(searchResp.Hotels.Hotel) > 0 {
		hotelCode = searchResp.Hotels.Hotel[0].Code
		if len(searchResp.Hotels.Hotel[0].RoomDetails) > 0 {
			roomDetail := searchResp.Hotels.Hotel[0].RoomDetails[0]
			roomDetails = &RoomDetailsRequest{
				Code:                   roomDetail.Code,
				MealPlan:               roomDetail.MealPlan,
				CancellationPolicyType: roomDetail.CancellationPolicyType,
				PackageRate:            roomDetail.PackageRate,
			}
		}
	}

	// Now run RecheckHotel with the session ID from SearchHotel
	testRoomRequest.RateKey = searchResp.Hotels.Hotel[0].RoomDetails[0].Rooms.Room[0].RateKey
	rooms := struct {
		Room []RoomRequest `json:"Room"`
	}{
		Room: []RoomRequest{
			testRoomRequest,
		},
	}

	req := RecheckHotelRequestDTO{
		SessionID:   sessionID,
		Nationality: testNationality,
		Country:     testCountry,
		City:        testCity,
		HotelCode:   hotelCode,
		CheckIn:     testCheckIn,
		CheckOut:    testCheckOut,
		RoomDetails: roomDetails,
		Rooms:       rooms,
		Currency:    testCurrency,
		Language:    testLanguage,
		AvailFlag:   testAvailFlag,
		DetailLevel: testDetailLevel,
	}

	resp, err := suite.client.RecheckHotel(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(resp)
	suite.True(resp.Status)

	if resp.Hotels.Hotel.Code != "" {
		fmt.Printf("Hotel Code: %s, Name: %s\n",
			resp.Hotels.Hotel.Code, resp.Hotels.Hotel.Name)
	}
}

func (suite *EndpointTestSuite) TestBookHotel_Success() {
	// First, run SearchHotel to get a session ID and hotel/room details
	searchRooms := &struct {
		Room []RoomRequest `json:"Room"`
	}{
		Room: []RoomRequest{
			testRoomRequest,
		},
	}

	searchReq := SearchHotelRequestDTO{
		Nationality:  testNationality,
		Country:      testCountry,
		City:         testCity,
		CheckIn:      testCheckIn,
		CheckOut:     testCheckOut,
		Rooms:        searchRooms,
		Currency:     testCurrency,
		Language:     testLanguage,
		AvailFlag:    testAvailFlag,
		MaxNoOfHotel: testMaxNoOfHotel,
		DetailLevel:  testDetailLevel,
	}

	searchResp, err := suite.client.SearchHotel(context.Background(), searchReq)
	suite.Nil(err)
	suite.NotNil(searchResp)
	suite.NotEmpty(searchResp.SessionID)
	suite.NotEmpty(searchResp.Hotels.Hotel)

	// Extract hotel/room details from SearchHotel response
	hotelCode := searchResp.Hotels.Hotel[0].Code
	var roomDetails RoomDetailsRequest

	if len(searchResp.Hotels.Hotel[0].RoomDetails) > 0 {
		roomDetail := searchResp.Hotels.Hotel[0].RoomDetails[0]
		roomDetails = RoomDetailsRequest{
			Code:                   roomDetail.Code,
			MealPlan:               roomDetail.MealPlan,
			CancellationPolicyType: roomDetail.CancellationPolicyType,
			PackageRate:            roomDetail.PackageRate,
		}
	}

	// Set RateKey for room request
	roomRequestWithRateKey := testRoomRequest
	if len(searchResp.Hotels.Hotel[0].RoomDetails) > 0 &&
		len(searchResp.Hotels.Hotel[0].RoomDetails[0].Rooms.Room) > 0 {
		roomRequestWithRateKey.RateKey = searchResp.Hotels.Hotel[0].RoomDetails[0].Rooms.Room[0].RateKey
	}

	// Second, run RecheckHotel to verify availability
	recheckRooms := struct {
		Room []RoomRequest `json:"Room"`
	}{
		Room: []RoomRequest{
			roomRequestWithRateKey,
		},
	}

	recheckReq := RecheckHotelRequestDTO{
		SessionID:   searchResp.SessionID,
		Nationality: testNationality,
		Country:     testCountry,
		City:        testCity,
		HotelCode:   hotelCode,
		CheckIn:     testCheckIn,
		CheckOut:    testCheckOut,
		RoomDetails: &roomDetails,
		Rooms:       recheckRooms,
		Currency:    testCurrency,
		Language:    testLanguage,
		AvailFlag:   testAvailFlag,
		DetailLevel: testDetailLevel,
	}

	recheckResp, err := suite.client.RecheckHotel(context.Background(), recheckReq)
	suite.Nil(err)
	suite.NotNil(recheckResp)
	suite.True(recheckResp.Status)
	suite.NotEmpty(recheckResp.SessionID)

	// Third, run BookHotel with the session ID from RecheckHotel
	// Create RoomBookingRequest with PaxDetails
	rateKey := ""
	if len(searchResp.Hotels.Hotel[0].RoomDetails) > 0 &&
		len(searchResp.Hotels.Hotel[0].RoomDetails[0].Rooms.Room) > 0 {
		rateKey = searchResp.Hotels.Hotel[0].RoomDetails[0].Rooms.Room[0].RateKey
	}

	bookRooms := struct {
		Room []RoomBookingRequest `json:"Room"`
	}{
		Room: []RoomBookingRequest{
			{
				PaxDetails: struct {
					Pax []Pax `json:"Pax"`
				}{
					Pax: []Pax{
						{
							Salutation: "Mr",
							FirstName:  "John",
							LastName:   "Doe",
							Age:        30,
						},
						{
							Salutation: "Mrs",
							FirstName:  "Jane",
							LastName:   "Doe",
							Age:        28,
						},
					},
				},
				RoomNo:     testRoomRequest.RoomNo,
				NoOfAdults: testRoomRequest.NoOfAdults,
				NoOfChild:  testRoomRequest.NoOfChild,
				ExtraBed:   testRoomRequest.ExtraBed,
				RateKey:    rateKey,
			},
		},
	}

	bookReq := BookHotelRequestDTO{
		SessionID:       recheckResp.SessionID,
		AgencyBookingID: testAgencyBookingID,
		Nationality:     testNationality,
		CheckIn:         testCheckIn,
		CheckOut:        testCheckOut,
		HotelCode:       hotelCode,
		RoomDetails:     roomDetails,
		Rooms:           bookRooms,
		Currency:        testCurrency,
		Language:        testLanguage,
		AvailFlag:       testAvailFlag,
		DetailLevel:     testDetailLevel,
	}

	resp, err := suite.client.BookHotel(context.Background(), bookReq)

	suite.Nil(err)
	suite.NotNil(resp)
	suite.True(resp.Status)
	suite.NotEmpty(resp.SessionID)
	suite.NotEmpty(resp.BookingDetail.MgBookingID)

	fmt.Printf("Booking successful - MG Booking ID: %s, Agency Booking ID: %s\n",
		resp.BookingDetail.MgBookingID, resp.BookingDetail.AgencyBookingID)
}

func (suite *EndpointTestSuite) TestGetRSVNDetails_Success() {
	req := GetRSVNDetailsRequestDTO{
		AgencyBookingID: testAgencyBookingID,
	}

	resp, err := suite.client.GetRSVNDetails(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(resp)
	suite.True(resp.Status)
	suite.NotEmpty(resp.SessionID)

	suite.Equal(testAgencyBookingID, resp.BookingDetail.AgencyBookingID)
	fmt.Printf("GetRSVNDetails successful - MG Booking ID: %s, Status: %s\n",
		resp.BookingDetail.MgBookingID, resp.BookingDetail.Status)
}
func (suite *EndpointTestSuite) TestCancelReservation_Success() {
	req := CancelReservationRequestDTO{
		AgencyBookingID: testAgencyBookingID,
	}

	resp, err := suite.client.CancelReservation(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(resp)
	suite.True(resp.Status)
	suite.NotEmpty(resp.SessionID)
	suite.NotEmpty(resp.BookingDetail.MgBookingID)

	suite.Equal(testAgencyBookingID, resp.BookingDetail.AgencyBookingID)
	fmt.Printf("CancelReservation successful - MG Booking ID: %s, Agency Booking ID: %s, Status: %s\n",
		resp.BookingDetail.MgBookingID, resp.BookingDetail.AgencyBookingID, resp.BookingDetail.Status)
}

func (suite *EndpointTestSuite) TestGetRSVNList_Success() {
	req := GetRSVNListRequestDTO{
		Date:     testReservationDate,
		DateType: BookingDateType,
		Status:   ConfirmedStatus,
	}

	resp, err := suite.client.GetRSVNList(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(resp)
	suite.True(resp.Status)
	suite.NotEmpty(resp.SessionID)
	suite.NotEmpty(resp.ReservationDetails)

	fmt.Printf("GetRSVNList successful - Agency Booking ID: %s, MG Booking ID: %s, Status: %s\n",
		resp.ReservationDetails[0].AgencyBookingID, resp.ReservationDetails[0].MGBookingID, resp.ReservationDetails[0].Status)
}

func (suite *EndpointTestSuite) TestCancelReservation_Error() {
	req := CancelReservationRequestDTO{
		AgencyBookingID: "123", // Not a valid agency booking ID
	}

	resp, err := suite.client.CancelReservation(context.Background(), req)

	suite.NotNil(err)
	suite.Nil(resp)
	suite.NotNil(err.RawBody)

	fmt.Printf("CancelReservation error - Error: %v, RawBody: %s\n", err.Error(), err.RawBody)
}
