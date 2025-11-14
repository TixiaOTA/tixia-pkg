package vervotech

import (
	"context"
	"time"
)

type InitializeSyncRequest struct {
	AllVervotechIdsResumeKey *string `json:"allVervotechIdsResumeKey"` // Resume key for pagination of all vervotech ids (optional, default: "")
	NewMappingsResumeKey     *string `json:"newMappingsResumeKey"`     // Resume key for pagination of new mappings (optional, default: "")

	ProviderFamily     *string    `json:"providerFamily"`     // Provider family name (optional, default: "")
	LastUpdateDateTime *time.Time `json:"lastUpdateDateTime"` // Last update date time filter (optional, default: nil)
	Limit              *int       `json:"limit"`              // Maximum number of results to return (optional, default: nil)
}

type InitializeSyncResponse struct {
	VervotechIds   GetAllVervotechIdsResponse              `json:"vervotechIds"`
	NewMappings    GetNewMappingsResponse                  `json:"newMappings"`
	CuratedContent GetCuratedContentByVervotechIdsResponse `json:"curatedContent"`
}

type IncrementalSyncRequest struct {
	DeletedMappingsResumeKey *string `json:"deletedMappingsResumeKey"` // Resume key for pagination of deleted mappings (optional, default: nil)
	NewMappingsResumeKey     *string `json:"newMappingsResumeKey"`     // Resume key for pagination of new mappings (optional, default: nil)
	UpdatedMappingsResumeKey *string `json:"updatedMappingsResumeKey"` // Resume key for pagination of updated mappings (optional, default: nil)
	AllVervotechIdsResumeKey *string `json:"allVervotechIdsResumeKey"` // Resume key for pagination of all vervotech ids (optional, default: nil)

	ProviderFamily     *string    `json:"providerFamily"`     // Provider family name (optional, default: nil)
	LastUpdateDateTime *time.Time `json:"lastUpdateDateTime"` // Last update date time filter (optional, default: nil)
	Limit              *int       `json:"limit"`              // Maximum number of results to return (optional, default: nil)
}

type IncrementalSyncResponse struct {
	DeletedMappings GetDeletedMappingsResponse               `json:"deletedMappings"`
	NewMappings     GetNewMappingsResponse                   `json:"newMappings"`
	UpdatedMappings GetUpdatedMappingsResponse               `json:"updatedMappings"`
	CuratedContent  GetCuratedContentByVervotechIdsResponse  `json:"curatedContent"`
	ProviderContent GetProviderContentByVervotechIdsResponse `json:"providerContent"`
}

// vervotechClient defines the interface for client methods used by the usecase
type vervotechClient interface {
	GetAllVervotechIds(ctx context.Context, req GetAllVervotechIdsRequest) (*GetAllVervotechIdsResponse, error)
	GetNewMappings(ctx context.Context, req GetNewMappingsRequest) (*GetNewMappingsResponse, error)
	GetUpdatedMappings(ctx context.Context, req GetUpdatedMappingsRequest) (*GetUpdatedMappingsResponse, error)
	GetDeletedMappings(ctx context.Context, req GetDeletedMappingsRequest) (*GetDeletedMappingsResponse, error)
	GetCuratedContentByVervotechIds(ctx context.Context, req GetCuratedContentByVervotechIdsRequest) (*GetCuratedContentByVervotechIdsResponse, error)
	GetProviderContentByVervotechIds(ctx context.Context, req GetProviderContentByVervotechIdsRequest) (*GetProviderContentByVervotechIdsResponse, error)
}

type Usecase struct {
	client vervotechClient
}

func NewUsecase(client *Client) *Usecase {
	return &Usecase{
		client: client,
	}
}

func (u *Usecase) InitializeSync(ctx context.Context, req InitializeSyncRequest) (*InitializeSyncResponse, error) {
	// Default values
	vervotechResumeKey := ""
	if req.AllVervotechIdsResumeKey != nil {
		vervotechResumeKey = *req.AllVervotechIdsResumeKey
	}

	newMappingsResumeKey := ""
	if req.NewMappingsResumeKey != nil {
		newMappingsResumeKey = *req.NewMappingsResumeKey
	}

	providerFamily := ""
	if req.ProviderFamily != nil {
		providerFamily = *req.ProviderFamily
	}

	lastUpdateDateTime := time.Time{}
	if req.LastUpdateDateTime != nil {
		lastUpdateDateTime = *req.LastUpdateDateTime
	}

	limit := 100
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	// Get all vervotech ids
	vervotechIds, err := u.client.GetAllVervotechIds(ctx, GetAllVervotechIdsRequest{
		LastUpdateDateTime: &lastUpdateDateTime,
		ResumeKey:          vervotechResumeKey,
	})
	if err != nil {
		return nil, err
	}

	// Get new mappings
	newMappings, err := u.client.GetNewMappings(ctx, GetNewMappingsRequest{
		Limit:              &limit,
		ProviderFamily:     providerFamily,
		LastUpdateDateTime: &lastUpdateDateTime,
		ResumeKey:          newMappingsResumeKey,
	})
	if err != nil {
		return nil, err
	}

	// Get curated content by vervotech ids
	curatedContent, err := u.client.GetCuratedContentByVervotechIds(ctx, GetCuratedContentByVervotechIdsRequest{
		VervotechIds: vervotechIds.VervotechIds,
	})
	if err != nil {
		return nil, err
	}

	// Response
	response := InitializeSyncResponse{
		VervotechIds:   *vervotechIds,
		NewMappings:    *newMappings,
		CuratedContent: *curatedContent,
	}

	return &response, nil
}

func (u *Usecase) IncrementalSync(ctx context.Context, req IncrementalSyncRequest) (*IncrementalSyncResponse, error) {
	// Default values
	deletedMappingsResumeKey := ""
	if req.DeletedMappingsResumeKey != nil {
		deletedMappingsResumeKey = *req.DeletedMappingsResumeKey
	}

	newMappingsResumeKey := ""
	if req.NewMappingsResumeKey != nil {
		newMappingsResumeKey = *req.NewMappingsResumeKey
	}

	updatedMappingsResumeKey := ""
	if req.UpdatedMappingsResumeKey != nil {
		updatedMappingsResumeKey = *req.UpdatedMappingsResumeKey
	}

	allVervotechIdsResumeKey := ""
	if req.AllVervotechIdsResumeKey != nil {
		allVervotechIdsResumeKey = *req.AllVervotechIdsResumeKey
	}

	providerFamily := ""
	if req.ProviderFamily != nil {
		providerFamily = *req.ProviderFamily
	}

	lastUpdateDateTime := time.Time{}
	if req.LastUpdateDateTime != nil {
		lastUpdateDateTime = *req.LastUpdateDateTime
	}

	limit := 100
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	// Get deleted mappings
	deletedMappings, err := u.client.GetDeletedMappings(ctx, GetDeletedMappingsRequest{
		Limit:              &limit,
		ProviderFamily:     providerFamily,
		LastUpdateDateTime: &lastUpdateDateTime,
		ResumeKey:          deletedMappingsResumeKey,
	})
	if err != nil {
		return nil, err
	}

	// Get new mappings
	newMappings, err := u.client.GetNewMappings(ctx, GetNewMappingsRequest{
		Limit:              &limit,
		ProviderFamily:     providerFamily,
		LastUpdateDateTime: &lastUpdateDateTime,
		ResumeKey:          newMappingsResumeKey,
	})
	if err != nil {
		return nil, err
	}

	// Get updated mappings
	updatedMappings, err := u.client.GetUpdatedMappings(ctx, GetUpdatedMappingsRequest{
		Limit:              &limit,
		ProviderFamily:     providerFamily,
		LastUpdateDateTime: &lastUpdateDateTime,
		ResumeKey:          updatedMappingsResumeKey,
	})
	if err != nil {
		return nil, err
	}

	// Get all vervotech ids
	vervotechIds, err := u.client.GetAllVervotechIds(ctx, GetAllVervotechIdsRequest{
		LastUpdateDateTime: &lastUpdateDateTime,
		ResumeKey:          allVervotechIdsResumeKey,
	})
	if err != nil {
		return nil, err
	}

	// Get curated content by vervotech ids
	curatedContent, err := u.client.GetCuratedContentByVervotechIds(ctx, GetCuratedContentByVervotechIdsRequest{
		VervotechIds: vervotechIds.VervotechIds,
	})
	if err != nil {
		return nil, err
	}

	// Get provider content by vervotech ids
	providerContent, err := u.client.GetProviderContentByVervotechIds(ctx, GetProviderContentByVervotechIdsRequest{
		NumberOfProviderContentCopies: 1,
		ProviderPreferences:           []string{providerFamily},
		VervotechIds:                  vervotechIds.VervotechIds,
	})
	if err != nil {
		return nil, err
	}

	// Response
	response := IncrementalSyncResponse{
		DeletedMappings: *deletedMappings,
		NewMappings:     *newMappings,
		UpdatedMappings: *updatedMappings,
		CuratedContent:  *curatedContent,
		ProviderContent: *providerContent,
	}

	return &response, nil
}
