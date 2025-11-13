package mg

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	httppkg "github.com/TixiaOTA/tixia-pkg/http/client"
	"github.com/TixiaOTA/tixia-pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MGClientTestSuite struct {
	suite.Suite
	server       *httptest.Server
	client       *MGClient
	logger       *logger.Logger
	config       MGConfig
	checkInDate  string
	checkOutDate string
}

func (suite *MGClientTestSuite) SetupTest() {
	suite.server = nil
	suite.client = nil

	logConfig := logger.Config{
		Level:       "debug",
		Environment: "test",
		JSONOutput:  false,
	}
	suite.logger = logger.New(logConfig)

	suite.config = MGConfig{
		BaseURL:   "https://dummy-api.example.com/api",
		AgentID:   "DUMMY_AGENT_ID",
		LoginName: "dummy_user",
		Password:  "dummy_password",
	}

	// Setup date helpers
	suite.checkInDate = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	suite.checkOutDate = time.Now().AddDate(0, 0, 2).Format("2006-01-02")
}

func (suite *MGClientTestSuite) TearDownTest() {
	if suite.server != nil {
		suite.server.Close()
	}
}

func (suite *MGClientTestSuite) createClient() error {
	client, err := NewMGClient(suite.config, suite.logger)
	if err != nil {
		return err
	}
	suite.client = client
	return nil
}

func TestMGClientTestSuite(t *testing.T) {
	suite.Run(t, new(MGClientTestSuite))
}

// =============================================================================
// Constructor Tests
// =============================================================================

func (suite *MGClientTestSuite) TestNewMGClient_Success() {
	client, err := NewMGClient(suite.config, suite.logger)

	suite.NoError(err)
	suite.NotNil(client)
	suite.Equal(suite.config, client.config)
	suite.Equal(suite.logger, client.logger)
	suite.NotNil(client.httpClient)
}

func (suite *MGClientTestSuite) TestNewMGClient_NilLogger() {
	client, err := NewMGClient(suite.config, nil)

	suite.Error(err)
	suite.Nil(client)
	suite.Contains(err.Error(), "logger is required")
}

func (suite *MGClientTestSuite) TestNewMGClient_InvalidConfig_MissingBaseURL() {
	config := MGConfig{
		AgentID:   "DUMMY_AGENT_ID",
		LoginName: "dummy_user",
		Password:  "dummy_password",
	}

	client, err := NewMGClient(config, suite.logger)

	suite.Error(err)
	suite.Nil(client)
	suite.Contains(err.Error(), "failed to validate MG config")
}

func (suite *MGClientTestSuite) TestNewMGClient_InvalidConfig_MissingAgentID() {
	config := MGConfig{
		BaseURL:   "https://dummy-api.example.com/api",
		LoginName: "dummy_user",
		Password:  "dummy_password",
	}

	client, err := NewMGClient(config, suite.logger)

	suite.Error(err)
	suite.Nil(client)
	suite.Contains(err.Error(), "failed to validate MG config")
}

func (suite *MGClientTestSuite) TestNewMGClient_InvalidConfig_MissingLoginName() {
	config := MGConfig{
		BaseURL:  "https://dummy-api.example.com/api",
		AgentID:  "DUMMY_AGENT_ID",
		Password: "dummy_password",
	}

	client, err := NewMGClient(config, suite.logger)

	suite.Error(err)
	suite.Nil(client)
	suite.Contains(err.Error(), "failed to validate MG config")
}

func (suite *MGClientTestSuite) TestNewMGClient_InvalidConfig_MissingPassword() {
	config := MGConfig{
		BaseURL:   "https://dummy-api.example.com/api",
		AgentID:   "DUMMY_AGENT_ID",
		LoginName: "dummy_user",
	}

	client, err := NewMGClient(config, suite.logger)

	suite.Error(err)
	suite.Nil(client)
	suite.Contains(err.Error(), "failed to validate MG config")
}

func (suite *MGClientTestSuite) TestNewMGClient_WithCustomTimeout() {
	customTimeout := 5 * time.Second
	config := suite.config
	config.Timeout = &customTimeout

	client, err := NewMGClient(config, suite.logger)

	suite.NoError(err)
	suite.NotNil(client)
}

func (suite *MGClientTestSuite) TestNewMGClient_WithCustomRetry() {
	customRetry := httppkg.RetryConfig{
		MaxRetries:        5,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        10 * time.Second,
		BackoffMultiplier: 3.0,
	}
	config := suite.config
	config.Retry = &customRetry

	client, err := NewMGClient(config, suite.logger)

	suite.NoError(err)
	suite.NotNil(client)
}

func (suite *MGClientTestSuite) TestNewMGClient_WithDefaultTimeout() {
	client, err := NewMGClient(suite.config, suite.logger)

	suite.NoError(err)
	suite.NotNil(client)
}

func (suite *MGClientTestSuite) TestNewMGClient_WithDefaultRetry() {
	client, err := NewMGClient(suite.config, suite.logger)

	suite.NoError(err)
	suite.NotNil(client)
}

// =============================================================================
// SearchHotelAvailability Tests
// =============================================================================

func (suite *MGClientTestSuite) TestSearchHotelAvailability_AdultsOnly() {
	err := suite.createClient()
	require.NoError(suite.T(), err)

	ctx := context.Background()
	req := &HotelSearchRequest{
		Nationality: "ID",
		Country:     "ID",
		City:        "ID-DPS",
		Hotels: &Hotels{
			Code: []string{""},
		},
		CheckIn:  suite.checkInDate,
		CheckOut: suite.checkOutDate,
		Rooms: &Rooms{
			Room: []RoomDetailSearchRequest{
				{
					RoomNo:     1,
					NoOfAdults: 1,
					ExtraBed:   false,
				},
			},
		},
		Currency:     "IDR",
		Language:     "EN",
		AvailFlag:    true,
		MaxNoOfHotel: 2,
		DetailLevel:  "FULL",
	}

	response, err := suite.client.SearchHotelAvailability(ctx, req)
	require.NoError(suite.T(), err)
	require.True(suite.T(), response.Status)

	assert.Len(suite.T(), response.Hotels.Hotel, 2)
	assert.Equal(suite.T(), 2, response.NoOfHotels)
	assert.Equal(suite.T(), "IDR", response.Currency)
	assert.Equal(suite.T(), suite.checkInDate, response.CheckIn)
	assert.Equal(suite.T(), suite.checkOutDate, response.CheckOut)
}

func (suite *MGClientTestSuite) TestSearchHotelAvailability_AdultsAndChildsWithTwoHotelCode() {
	err := suite.createClient()
	require.NoError(suite.T(), err)

	ctx := context.Background()
	req := &HotelSearchRequest{
		Nationality: "ID",
		Country:     "ID",
		City:        "ID-DPS",
		Hotels: &Hotels{
			Code: []string{"ID10008902", "ID10000736"},
		},
		CheckIn:  suite.checkInDate,
		CheckOut: suite.checkOutDate,
		Rooms: &Rooms{
			Room: []RoomDetailSearchRequest{
				{
					RoomNo:     1,
					NoOfAdults: 1,
					NoOfChild:  2,
					Child1Age:  11,
					Child2Age:  9,
					ExtraBed:   false,
				},
			},
		},
		Currency:     "IDR",
		Language:     "EN",
		AvailFlag:    true,
		MaxNoOfHotel: 2,
		DetailLevel:  "FULL",
	}

	response, err := suite.client.SearchHotelAvailability(ctx, req)
	require.NoError(suite.T(), err)
	require.True(suite.T(), response.Status)

	assert.Len(suite.T(), response.Hotels.Hotel, 2)
	assert.Equal(suite.T(), 2, response.NoOfHotels)
	assert.Equal(suite.T(), "IDR", response.Currency)
	assert.Equal(suite.T(), suite.checkInDate, response.CheckIn)
	assert.Equal(suite.T(), suite.checkOutDate, response.CheckOut)
	assert.Equal(suite.T(), "ID10008902", response.Hotels.Hotel[0].Code)
	assert.Equal(suite.T(), "ID10000736", response.Hotels.Hotel[1].Code)
	assert.Equal(suite.T(), false, response.Hotels.Hotel[0].RoomDetails[0].Rooms.Room[0].ExtraBed)
	assert.Equal(suite.T(), false, response.Hotels.Hotel[1].RoomDetails[0].Rooms.Room[0].ExtraBed)
	assert.Equal(suite.T(), float32(9), response.Hotels.Hotel[0].RoomDetails[0].Rooms.Room[0].Child1Age)
	assert.Equal(suite.T(), float32(0), response.Hotels.Hotel[1].RoomDetails[0].Rooms.Room[0].Child2Age)
}

func (suite *MGClientTestSuite) TestSearchHotelAvailability_TechnicalError() {
	err := suite.createClient()
	require.NoError(suite.T(), err)

	ctx := context.Background()
	req := &HotelSearchRequest{
		Nationality: "ID",
		Country:     "ID",
		City:        "ID-DPS",
		Hotels: &Hotels{
			Code: []string{""},
		},
		CheckIn:  suite.checkInDate,
		CheckOut: suite.checkOutDate,
		Rooms: &Rooms{
			Room: []RoomDetailSearchRequest{
				{
					RoomNo:     0,
					NoOfAdults: 0,
					ExtraBed:   false,
				},
			},
		},
		Currency:     "IDR",
		Language:     "EN",
		AvailFlag:    true,
		MaxNoOfHotel: 2,
		DetailLevel:  "FULL",
	}

	response, err := suite.client.SearchHotelAvailability(ctx, req)
	require.NoError(suite.T(), err)
	require.False(suite.T(), response.Status)

	assert.NotNil(suite.T(), response.ErrorCode)
	assert.Equal(suite.T(), "JRVXML168", *response.ErrorCode)
	assert.NotNil(suite.T(), response.ErrorMessage)
	assert.Contains(suite.T(), *response.ErrorMessage, "RoomNo")
}

func (suite *MGClientTestSuite) TestSearchHotelAvailability_InternalError() {
	err := suite.createClient()
	require.NoError(suite.T(), err)

	ctx := context.Background()
	req := &HotelSearchRequest{
		Nationality: "ID",
		Country:     "ID",
		City:        "ID-DPS",
		Hotels: &Hotels{
			Code: []string{""},
		},
		CheckIn:  suite.checkInDate,
		CheckOut: suite.checkOutDate,
		Rooms: &Rooms{
			Room: []RoomDetailSearchRequest{
				{
					RoomNo:     1,
					NoOfAdults: 1,
					ExtraBed:   false,
				},
			},
		},
		Currency:     "IDR",
		Language:     "EN",
		AvailFlag:    true,
		MaxNoOfHotel: 0,
		DetailLevel:  "FULL",
	}

	response, err := suite.client.SearchHotelAvailability(ctx, req)
	require.Error(suite.T(), err)
	assert.False(suite.T(), response.Status)
}

func (suite *MGClientTestSuite) TestSearchHotelAvailability_InvalidDetailLevelError() {
	err := suite.createClient()
	require.NoError(suite.T(), err)

	ctx := context.Background()
	req := &HotelSearchRequest{
		Nationality: "ID",
		Country:     "ID",
		City:        "ID-DPS",
		Hotels: &Hotels{
			Code: []string{""},
		},
		CheckIn:  suite.checkInDate,
		CheckOut: suite.checkOutDate,
		Rooms: &Rooms{
			Room: []RoomDetailSearchRequest{
				{
					RoomNo:     1,
					NoOfAdults: 1,
					ExtraBed:   false,
				},
			},
		},
		Currency:     "IDR",
		Language:     "EN",
		AvailFlag:    true,
		MaxNoOfHotel: 2,
		DetailLevel:  "ABC",
	}

	response, err := suite.client.SearchHotelAvailability(ctx, req)
	require.NoError(suite.T(), err)
	require.False(suite.T(), response.Status)

	assert.NotNil(suite.T(), response.ErrorCode)
	assert.Equal(suite.T(), "JRVXML157", *response.ErrorCode)
}
