package vervotech

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// Test configuration variables - Update these with your test credentials
var (
	testBaseURL    = "https://hotelmapping.vervotech.com/api" // Base URL for the API
	testAPIKey     = "test-api-key"                           // Your API key
	testAccountID  = "test-account-id"                        // Your account ID
	testSyncID     = "test-sync-id"                           // Your sync ID (optional)
	testAPIVersion = "3.0"                                    // API version

	// Test data for GetNewMappings
	testNewMappingsLimit          = 10
	testNewMappingsProviderFamily = ""
	testNewMappingsLastUpdateDate = time.Date(2025, 1, 0, 0, 0, 0, 0, time.UTC)

	// Test data for GetUpdatedMappings
	testUpdatedMappingsLimit          = 10
	testUpdatedMappingsProviderFamily = ""
	testUpdatedMappingsLastUpdateDate = time.Date(2025, 1, 0, 0, 0, 0, 0, time.UTC)
	testUpdatedMappingsResumeKey      = ""

	// Test data for GetDeletedMappings
	testDeletedMappingsLimit          = 10
	testDeletedMappingsProviderFamily = ""
	testDeletedMappingsLastUpdateDate = time.Date(2025, 1, 0, 0, 0, 0, 0, time.UTC)

	// Test data for GetCuratedContentByVervotechIds
	testCuratedContentVervotechIDs = []int64{15069020}

	// Test data for GetProviderContentByVervotechIds
	testProviderContentVervotechIDs        = []int64{15069020}
	testProviderContentNumberOfCopies      = 1
	testProviderContentProviderPreferences = []string{}
)

type EndpointTestSuite struct {
	suite.Suite
	client *Client
}

func (suite *EndpointTestSuite) SetupTest() {
	config := Config{
		BaseURL:    testBaseURL,
		APIKey:     testAPIKey,
		AccountID:  testAccountID,
		APIVersion: testAPIVersion,
		SyncID:     testSyncID,
	}

	client, err := New(config)
	suite.NoError(err)
	suite.client = client
}

func TestEndpointTestSuite(t *testing.T) {
	suite.Run(t, new(EndpointTestSuite))
}

func (suite *EndpointTestSuite) TestGetNewMappings_Success() {
	limit := testNewMappingsLimit
	lastUpdateDateTime := testNewMappingsLastUpdateDate
	req := GetNewMappingsRequest{
		Limit:              &limit,
		ProviderFamily:     testNewMappingsProviderFamily,
		LastUpdateDateTime: &lastUpdateDateTime,
	}

	resp, err := suite.client.GetNewMappings(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.True(resp.Status)
	suite.Equal("Success", resp.Message)
	suite.NotEmpty(resp.Mappings)
	suite.NotEmpty(resp.ResumeKey)

	fmt.Printf("Received %d mappings\n", len(resp.Mappings))
	for i, mapping := range resp.Mappings {
		if i < 3 {
			fmt.Printf("Mapping %d: VervotechId=%d, ProviderHotelId=%s, ProviderName=%s, Type=%s\n",
				i+1, mapping.VervotechId, mapping.ProviderHotelId, mapping.ProviderName, mapping.Type)
		}
	}
}

func (suite *EndpointTestSuite) TestGetUpdatedMappings_Success() {
	limit := testUpdatedMappingsLimit
	lastUpdateDateTime := testUpdatedMappingsLastUpdateDate
	req := GetUpdatedMappingsRequest{
		Limit:              &limit,
		ProviderFamily:     testUpdatedMappingsProviderFamily,
		LastUpdateDateTime: &lastUpdateDateTime,
		ResumeKey:          testUpdatedMappingsResumeKey,
	}

	resp, err := suite.client.GetUpdatedMappings(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.True(resp.Status)
	suite.Equal("Success", resp.Message)
	suite.NotEmpty(resp.Mappings)
	suite.NotEmpty(resp.ResumeKey)

	fmt.Printf("Received %d updated mappings\n", len(resp.Mappings))
	for i, mapping := range resp.Mappings {
		if i < 3 {
			fmt.Printf("Mapping %d: VervotechId=%d, ProviderHotelId=%s, ProviderName=%s, Type=%s\n",
				i+1, mapping.VervotechId, mapping.ProviderHotelId, mapping.ProviderName, mapping.Type)
		}
	}
}

func (suite *EndpointTestSuite) TestGetDeletedMappings_Success() {
	limit := testDeletedMappingsLimit
	lastUpdateDateTime := testDeletedMappingsLastUpdateDate
	req := GetDeletedMappingsRequest{
		Limit:              &limit,
		ProviderFamily:     testDeletedMappingsProviderFamily,
		LastUpdateDateTime: &lastUpdateDateTime,
	}

	resp, err := suite.client.GetDeletedMappings(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.True(resp.Status)
	suite.Equal("Success", resp.Message)
	suite.NotEmpty(resp.Mappings)
	suite.NotEmpty(resp.ResumeKey)

	fmt.Printf("Received %d deleted mappings\n", len(resp.Mappings))
	for i, mapping := range resp.Mappings {
		if i < 3 {
			fmt.Printf("Mapping %d: VervotechId=%d, ProviderHotelId=%s, ProviderName=%s, Type=%s\n",
				i+1, mapping.VervotechId, mapping.ProviderHotelId, mapping.ProviderName, mapping.Type)
		}
	}
}

func (suite *EndpointTestSuite) TestGetAllVervotechIds_Success() {
	req := GetAllVervotechIdsRequest{}

	resp, err := suite.client.GetAllVervotechIds(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.True(resp.Status)
	suite.Equal("Success", resp.Message)
	suite.NotEmpty(resp.VervotechIds)
	suite.NotEmpty(resp.ResumeKey)

	fmt.Printf("Received %d Vervotech IDs\n", len(resp.VervotechIds))
	if len(resp.VervotechIds) > 0 {
		maxPrint := 10
		if len(resp.VervotechIds) < maxPrint {
			maxPrint = len(resp.VervotechIds)
		}
		fmt.Printf("First %d IDs: %v\n", maxPrint, resp.VervotechIds[:maxPrint])
	}
}

func (suite *EndpointTestSuite) TestGetCuratedContentByVervotechIds_Success() {
	req := GetCuratedContentByVervotechIdsRequest{
		VervotechIds: testCuratedContentVervotechIDs,
	}

	resp, err := suite.client.GetCuratedContentByVervotechIds(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.True(resp.Status)
	suite.Equal("Success", resp.Message)
	suite.NotEmpty(resp.CuratedHotels)

	fmt.Printf("Received %d curated hotels\n", len(resp.CuratedHotels))
	for i, hotel := range resp.CuratedHotels {
		if i < 3 {
			fmt.Printf("Hotel %d: VervotechId=%d, Name=%s, Rating=%s, City=%s, Country=%s\n",
				i+1, hotel.VervotechID, hotel.Name, hotel.Rating, hotel.CityName, hotel.CountryName)
		}
	}
}

func (suite *EndpointTestSuite) TestGetProviderContentByVervotechIds_Success() {
	req := GetProviderContentByVervotechIdsRequest{
		NumberOfProviderContentCopies: testProviderContentNumberOfCopies,
		ProviderPreferences:           testProviderContentProviderPreferences,
		VervotechIds:                  testProviderContentVervotechIDs,
	}

	resp, err := suite.client.GetProviderContentByVervotechIds(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.True(resp.Status)
	suite.Equal("Success", resp.Message)
	suite.NotEmpty(resp.Hotels)

	fmt.Printf("Received %d hotel groups\n", len(resp.Hotels))
	for i, hotelGroup := range resp.Hotels {
		if i < 3 {
			fmt.Printf("Hotel Group %d: VervotechId=%d, ProviderHotels count=%d\n",
				i+1, hotelGroup.VervotechID, len(hotelGroup.ProviderHotels))

			for j, providerHotel := range hotelGroup.ProviderHotels {
				if j < 2 {
					fmt.Printf("  Provider Hotel %d: ProviderHotelId=%s, Name=%s, Rating=%f\n",
						j+1, providerHotel.ProviderHotelID, providerHotel.Name, providerHotel.Rating)
				}
			}
		}
	}
}
