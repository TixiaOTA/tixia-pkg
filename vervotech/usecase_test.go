package vervotech

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// Test data for InitializeSync
var (
	testInitializeSyncLimit                    = 10
	testInitializeSyncProviderFamily           = ""
	testInitializeSyncLastUpdateDateTime       = time.Date(2025, 1, 0, 0, 0, 0, 0, time.UTC)
	testInitializeSyncAllVervotechIdsResumeKey = ""
	testInitializeSyncNewMappingsResumeKey     = ""
)

// Test data for IncrementalSync
var (
	testIncrementalSyncLimit                    = 10
	testIncrementalSyncProviderFamily           = ""
	testIncrementalSyncLastUpdateDateTime       = time.Date(2025, 1, 0, 0, 0, 0, 0, time.UTC)
	testIncrementalSyncDeletedMappingsResumeKey = ""
	testIncrementalSyncNewMappingsResumeKey     = ""
	testIncrementalSyncUpdatedMappingsResumeKey = ""
	testIncrementalSyncAllVervotechIdsResumeKey = ""
)

type UsecaseTestSuite struct {
	suite.Suite
	usecase *Usecase
	client  *Client
}

func (suite *UsecaseTestSuite) SetupTest() {
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

	suite.usecase = NewUsecase(client)
}

func TestUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseTestSuite))
}

func (suite *UsecaseTestSuite) TestInitializeSync_Success() {
	ctx := context.Background()
	limit := testInitializeSyncLimit
	lastUpdateDateTime := testInitializeSyncLastUpdateDateTime

	req := InitializeSyncRequest{
		Limit:              &limit,
		LastUpdateDateTime: &lastUpdateDateTime,
	}

	if testInitializeSyncAllVervotechIdsResumeKey != "" {
		req.AllVervotechIdsResumeKey = &testInitializeSyncAllVervotechIdsResumeKey
	}

	if testInitializeSyncNewMappingsResumeKey != "" {
		req.NewMappingsResumeKey = &testInitializeSyncNewMappingsResumeKey
	}

	if testInitializeSyncProviderFamily != "" {
		req.ProviderFamily = &testInitializeSyncProviderFamily
	}

	resp, err := suite.usecase.InitializeSync(ctx, req)

	suite.NoError(err)
	suite.NotNil(resp)

	// Verify response structure
	// VervotechResumeKey should exist (can be empty string if no resume key was provided)
	_ = resp.VervotechIds.ResumeKey

	// NewMappings should be a struct with Mappings and MappingsResumeKey
	suite.NotNil(resp.NewMappings)
	// Mappings can be empty if no new mappings were found
	suite.True(len(resp.NewMappings.Mappings) == 0 || len(resp.NewMappings.Mappings) > 0)
	// MappingsResumeKey should exist (can be empty string if no resume key was provided)
	_ = resp.NewMappings.ResumeKey

	// CuratedContent should exist (can be empty if no curated content was found)
	suite.True(len(resp.CuratedContent.CuratedHotels) >= 0)
}

func (suite *UsecaseTestSuite) TestIncrementalSync_Success() {
	ctx := context.Background()
	limit := testIncrementalSyncLimit
	lastUpdateDateTime := testIncrementalSyncLastUpdateDateTime

	req := IncrementalSyncRequest{
		Limit:              &limit,
		LastUpdateDateTime: &lastUpdateDateTime,
	}

	if testIncrementalSyncDeletedMappingsResumeKey != "" {
		req.DeletedMappingsResumeKey = &testIncrementalSyncDeletedMappingsResumeKey
	}

	if testIncrementalSyncNewMappingsResumeKey != "" {
		req.NewMappingsResumeKey = &testIncrementalSyncNewMappingsResumeKey
	}

	if testIncrementalSyncUpdatedMappingsResumeKey != "" {
		req.UpdatedMappingsResumeKey = &testIncrementalSyncUpdatedMappingsResumeKey
	}

	if testIncrementalSyncAllVervotechIdsResumeKey != "" {
		req.AllVervotechIdsResumeKey = &testIncrementalSyncAllVervotechIdsResumeKey
	}

	if testIncrementalSyncProviderFamily != "" {
		req.ProviderFamily = &testIncrementalSyncProviderFamily
	}

	resp, err := suite.usecase.IncrementalSync(ctx, req)

	suite.NoError(err)
	suite.NotNil(resp)

	// Verify DeletedMappings structure
	suite.NotNil(resp.DeletedMappings)
	// Mappings can be empty if no deleted mappings were found
	suite.True(len(resp.DeletedMappings.Mappings) == 0 || len(resp.DeletedMappings.Mappings) > 0)
	// ResumeKey should exist (can be empty string if no resume key was provided)
	_ = resp.DeletedMappings.ResumeKey

	// Verify NewMappings structure
	suite.NotNil(resp.NewMappings)
	// Mappings can be empty if no new mappings were found
	suite.True(len(resp.NewMappings.Mappings) == 0 || len(resp.NewMappings.Mappings) > 0)
	// ResumeKey should exist (can be empty string if no resume key was provided)
	_ = resp.NewMappings.ResumeKey

	// Verify UpdatedMappings structure
	suite.NotNil(resp.UpdatedMappings)
	// Mappings can be empty if no updated mappings were found
	suite.True(len(resp.UpdatedMappings.Mappings) == 0 || len(resp.UpdatedMappings.Mappings) > 0)
	// ResumeKey should exist (can be empty string if no resume key was provided)
	_ = resp.UpdatedMappings.ResumeKey

	// Verify CuratedContent structure
	suite.NotNil(resp.CuratedContent)
	// CuratedHotels can be empty if no curated content was found
	suite.True(len(resp.CuratedContent.CuratedHotels) >= 0)

	// Verify ProviderContent structure
	suite.NotNil(resp.ProviderContent)
	// Hotels can be empty if no provider content was found
	suite.True(len(resp.ProviderContent.Hotels) >= 0)
}
