package vervotech

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/TixiaOTA/tixia-pkg/http/client"
)

// GetNewMappings retrieves new mappings from the Vervotech API
func (c *Client) GetNewMappings(ctx context.Context, req GetNewMappingsRequest) (*GetNewMappingsResponse, error) {
	path := "/" + c.config.APIVersion + "/mappings/GetNewMappings"

	params := make(map[string]string)
	if req.Limit != nil {
		params["limit"] = strconv.Itoa(*req.Limit)
	}
	if req.ProviderFamily != "" {
		params["providerFamily"] = req.ProviderFamily
	}
	if req.LastUpdateDateTime != nil {
		params["lastUpdateDateTime"] = formatDateTime(req.LastUpdateDateTime)
	}
	if req.ResumeKey != "" {
		params["resumeKey"] = req.ResumeKey
	}

	path += buildQueryString(params)

	headers := c.buildHeaders(true)

	httpResp, err := c.httpClient.Get(ctx, client.GetRequest{
		Path:    path,
		Headers: &headers,
	})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	var response GetNewMappingsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetUpdatedMappings retrieves updated mappings from the Vervotech API
func (c *Client) GetUpdatedMappings(ctx context.Context, req GetUpdatedMappingsRequest) (*GetUpdatedMappingsResponse, error) {
	path := "/" + c.config.APIVersion + "/mappings/GetUpdatedMappings"

	params := make(map[string]string)
	if req.Limit != nil {
		params["limit"] = strconv.Itoa(*req.Limit)
	}
	if req.ProviderFamily != "" {
		params["providerFamily"] = req.ProviderFamily
	}
	if req.LastUpdateDateTime != nil {
		params["lastUpdateDateTime"] = formatDateTime(req.LastUpdateDateTime)
	}
	if req.ResumeKey != "" {
		params["resumeKey"] = req.ResumeKey
	}

	path += buildQueryString(params)

	headers := c.buildHeaders(true)

	httpResp, err := c.httpClient.Get(ctx, client.GetRequest{
		Path:    path,
		Headers: &headers,
	})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	var response GetUpdatedMappingsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetDeletedMappings retrieves deleted mappings from the Vervotech API
func (c *Client) GetDeletedMappings(ctx context.Context, req GetDeletedMappingsRequest) (*GetDeletedMappingsResponse, error) {
	path := "/" + c.config.APIVersion + "/mappings/GetDeletedMappings"

	params := make(map[string]string)
	if req.Limit != nil {
		params["limit"] = strconv.Itoa(*req.Limit)
	}
	if req.ProviderFamily != "" {
		params["providerFamily"] = req.ProviderFamily
	}
	if req.LastUpdateDateTime != nil {
		params["lastUpdateDateTime"] = formatDateTime(req.LastUpdateDateTime)
	}
	if req.ResumeKey != "" {
		params["resumeKey"] = req.ResumeKey
	}

	path += buildQueryString(params)

	headers := c.buildHeaders(true)

	httpResp, err := c.httpClient.Get(ctx, client.GetRequest{
		Path:    path,
		Headers: &headers,
	})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	var response GetDeletedMappingsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetAllVervotechIds retrieves all Vervotech IDs from the API
func (c *Client) GetAllVervotechIds(ctx context.Context, req GetAllVervotechIdsRequest) (*GetAllVervotechIdsResponse, error) {
	path := "/" + c.config.APIVersion + "/hotels/GetAllVervotechIds"

	params := make(map[string]string)
	if req.LastUpdateDateTime != nil {
		params["lastUpdateDateTime"] = formatDateTime(req.LastUpdateDateTime)
	}
	if req.ResumeKey != "" {
		params["resumeKey"] = req.ResumeKey
	}

	path += buildQueryString(params)

	headers := c.buildHeaders(true)

	httpResp, err := c.httpClient.Get(ctx, client.GetRequest{
		Path:    path,
		Headers: &headers,
	})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	var response GetAllVervotechIdsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCuratedContentByVervotechIds retrieves curated content for the given Vervotech IDs
func (c *Client) GetCuratedContentByVervotechIds(ctx context.Context, req GetCuratedContentByVervotechIdsRequest) (*GetCuratedContentByVervotechIdsResponse, error) {
	path := "/" + c.config.APIVersion + "/content/GetCuratedContentByVervotechIds"

	headers := c.buildHeaders(false)

	httpResp, err := c.httpClient.PostJSON(ctx, client.PostJSONRequest{
		Path:    path,
		Body:    req,
		Headers: &headers,
	})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	var response GetCuratedContentByVervotechIdsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProviderContentByVervotechIds retrieves provider content for the given Vervotech IDs
func (c *Client) GetProviderContentByVervotechIds(ctx context.Context, req GetProviderContentByVervotechIdsRequest) (*GetProviderContentByVervotechIdsResponse, error) {
	path := "/" + c.config.APIVersion + "/content/GetProviderContentByVervotechIds"

	headers := c.buildHeaders(false)

	httpResp, err := c.httpClient.PostJSON(ctx, client.PostJSONRequest{
		Path:    path,
		Body:    req,
		Headers: &headers,
	})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	var response GetProviderContentByVervotechIdsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMasterHotelFacilityGroups retrieves master hotel facility groups from the Vervotech API
func (c *Client) GetMasterHotelFacilityGroups(ctx context.Context, req GetMasterHotelFacilityGroupsRequest) (*GetMasterHotelFacilityGroupsResponse, error) {
	path := "/" + c.config.APIVersion + "/content/GetMasterHotelFacilityGroups"

	headers := c.buildHeaders(true)

	httpResp, err := c.httpClient.Get(ctx, client.GetRequest{
		Path:    path,
		Headers: &headers,
	})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	var response GetMasterHotelFacilityGroupsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMasterHotelChains retrieves master hotel chains from the Vervotech API
func (c *Client) GetMasterHotelChains(ctx context.Context, req GetMasterHotelChainsRequest) (*GetMasterHotelChainsResponse, error) {
	path := "/" + c.config.APIVersion + "/content/GetMasterHotelChains"

	headers := c.buildHeaders(true)

	httpResp, err := c.httpClient.Get(ctx, client.GetRequest{
		Path:    path,
		Headers: &headers,
	})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	var response GetMasterHotelChainsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMasterPropertyTypes retrieves master property types from the Vervotech API
func (c *Client) GetMasterPropertyTypes(ctx context.Context, req GetMasterPropertyTypesRequest) (*GetMasterPropertyTypesResponse, error) {
	path := "/" + c.config.APIVersion + "/content/GetMasterPropertyTypes"

	headers := c.buildHeaders(true)

	httpResp, err := c.httpClient.Get(ctx, client.GetRequest{
		Path:    path,
		Headers: &headers,
	})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	var response GetMasterPropertyTypesResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMasterCountryCodes retrieves master country codes from the Vervotech API
func (c *Client) GetMasterCountryCodes(ctx context.Context, req GetMasterCountryCodesRequest) (*GetMasterCountryCodesResponse, error) {
	path := "/" + c.config.APIVersion + "/content/GetMasterCountryCodes"

	headers := c.buildHeaders(true)

	httpResp, err := c.httpClient.Get(ctx, client.GetRequest{
		Path:    path,
		Headers: &headers,
	})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, &APIError{
			Code:    httpResp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	var response GetMasterCountryCodesResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
