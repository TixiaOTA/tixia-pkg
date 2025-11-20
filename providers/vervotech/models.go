package vervotech

import (
	"fmt"
	"time"
)

// --------------------------------- REQUESTS ---------------------------------

// GetNewMappingsRequest represents the request parameters for GetNewMappings
type GetNewMappingsRequest struct {
	Limit              *int       // Maximum number of results to return
	ProviderFamily     string     // Provider family name
	LastUpdateDateTime *time.Time // Last update date time filter (optional)
	ResumeKey          string     // Resume key for pagination (optional)
}

// GetUpdatedMappingsRequest represents the request parameters for GetUpdatedMappings
type GetUpdatedMappingsRequest struct {
	Limit              *int       // Maximum number of results to return
	ProviderFamily     string     // Provider family name
	LastUpdateDateTime *time.Time // Last update date time filter (optional)
	ResumeKey          string     // Resume key for pagination (optional)
}

// GetDeletedMappingsRequest represents the request parameters for GetDeletedMappings
type GetDeletedMappingsRequest struct {
	Limit              *int       // Maximum number of results to return
	ProviderFamily     string     // Provider family name
	LastUpdateDateTime *time.Time // Last update date time filter (optional)
	ResumeKey          string     // Resume key for pagination (optional)
}

// GetAllVervotechIdsRequest represents the request parameters for GetAllVervotechIds
type GetAllVervotechIdsRequest struct {
	LastUpdateDateTime *time.Time // Last update date time filter (optional)
	ResumeKey          string     // Resume key for pagination (optional)
}

// GetCuratedContentByVervotechIdsRequest represents the request body for GetCuratedContentByVervotechIds
type GetCuratedContentByVervotechIdsRequest struct {
	VervotechIds []int64 `json:"VervotechIds"` // List of Vervotech IDs
}

// GetProviderContentByVervotechIdsRequest represents the request body for GetProviderContentByVervotechIds
type GetProviderContentByVervotechIdsRequest struct {
	NumberOfProviderContentCopies int      `json:"NumberOfProviderContentCopies"` // Number of provider content copies
	ProviderPreferences           []string `json:"ProviderPreferences"`           // Provider preferences list
	VervotechIds                  []int64  `json:"VervotechIds"`                  // List of Vervotech IDs
}

// GetMasterHotelFacilityGroupsRequest represents the request parameters for GetMasterHotelFacilityGroups
type GetMasterHotelFacilityGroupsRequest struct {
	// No parameters for this endpoint
}

// GetMasterHotelChainsRequest represents the request parameters for GetMasterHotelChains
type GetMasterHotelChainsRequest struct {
	// No parameters for this endpoint
}

// GetMasterPropertyTypesRequest represents the request parameters for GetMasterPropertyTypes
type GetMasterPropertyTypesRequest struct {
	// No parameters for this endpoint
}

// GetMasterCountryCodesRequest represents the request parameters for GetMasterCountryCodes
type GetMasterCountryCodesRequest struct {
	// No parameters for this endpoint
}

// --------------------------------- RESPONSES ---------------------------------

// APIError represents an error response from the API
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.Code, e.Message)
}

// Mapping represents a single hotel mapping
type Mapping struct {
	VervotechId          int64    `json:"VervotechId"`
	ProviderHotelId      string   `json:"ProviderHotelId"`
	ProviderName         string   `json:"ProviderName"`
	ChannelIds           []string `json:"ChannelIds"`
	ProviderLocationCode *string  `json:"ProviderLocationCode"`
	Type                 string   `json:"Type"`
}

// GetNewMappingsResponse represents the response from GetNewMappings endpoint
type GetNewMappingsResponse struct {
	ResumeKey  string    `json:"ResumeKey"`
	Mappings   []Mapping `json:"Mappings"`
	Status     bool      `json:"Status"`
	Message    string    `json:"Message"`
	StatusCode int       `json:"StatusCode"`
}

// GetUpdatedMappingsResponse represents the response from GetUpdatedMappings endpoint
type GetUpdatedMappingsResponse struct {
	ResumeKey  string    `json:"ResumeKey"`
	Mappings   []Mapping `json:"Mappings"`
	Status     bool      `json:"Status"`
	Message    string    `json:"Message"`
	StatusCode int       `json:"StatusCode"`
}

// GetDeletedMappingsResponse represents the response from GetDeletedMappings endpoint
type GetDeletedMappingsResponse struct {
	ResumeKey  string    `json:"ResumeKey"`
	Mappings   []Mapping `json:"Mappings"`
	Status     bool      `json:"Status"`
	Message    string    `json:"Message"`
	StatusCode int       `json:"StatusCode"`
}

// GetAllVervotechIdsResponse represents the response from GetAllVervotechIds endpoint
type GetAllVervotechIdsResponse struct {
	ResumeKey    string  `json:"ResumeKey"`
	VervotechIds []int64 `json:"VervotechIds"`
	Status       bool    `json:"Status"`
	Message      string  `json:"Message"`
	StatusCode   int     `json:"StatusCode"`
}

// GetMasterHotelFacilityGroupsResponse represents the response from GetMasterHotelFacilityGroups endpoint
type GetMasterHotelFacilityGroupsResponse struct {
	Status         bool            `json:"Status"`
	Message        string          `json:"Message"`
	StatusCode     int             `json:"StatusCode"`
	FacilityGroups []FacilityGroup `json:"FacilityGroups"`
}

// GetMasterHotelChainsResponse represents the response from GetMasterHotelChains endpoint
type GetMasterHotelChainsResponse struct {
	Status      bool         `json:"Status"`
	Message     string       `json:"Message"`
	StatusCode  int          `json:"StatusCode"`
	HotelChains []HotelChain `json:"HotelChains"`
}

// GetMasterPropertyTypesResponse represents the response from GetMasterPropertyTypes endpoint
type GetMasterPropertyTypesResponse struct {
	Status        bool     `json:"Status"`
	Message       string   `json:"Message"`
	StatusCode    int      `json:"StatusCode"`
	PropertyTypes []string `json:"PropertyTypes"` // Example: ["Hotel", "Apartment", "Resort", "Motel", "Hostel"]
}

// GetMasterCountryCodesResponse represents the response from GetMasterCountryCodes endpoint
type GetMasterCountryCodesResponse struct {
	Status     bool    `json:"Status"`
	Message    string  `json:"Message"`
	StatusCode int     `json:"StatusCode"`
	Countries  Country `json:"Countries"`
}

// --------------------------------- SHARED TYPES ---------------------------------

// Facility represents a hotel facility
type Facility struct {
	ID      int    `json:"Id"`
	Name    string `json:"Name"`
	GroupID int    `json:"GroupId"`
}

// FacilityGroup represents a facility group
type FacilityGroup struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// HotelChain represents a hotel chain
type HotelChain struct {
	ID        int    `json:"Id"`
	ChainCode string `json:"ChainCode"`
	ChainName string `json:"ChainName"`
}

// CountryCode represents a country code
//
// Example: "ID": "Indonesia"
type Country map[string]string

// Location represents a location entity (airport, train station, etc.)
type Location struct {
	LocationID int    `json:"LocationId"`
	Code       string `json:"Code"`
	Name       string `json:"Name"`
}

// LocationMapping represents location mapping information
type LocationMapping struct {
	CityLocationID     string     `json:"CityLocationId"`
	CityLocationName   string     `json:"CityLocationName"`
	Source             string     `json:"Source"`
	TrainStation       Location   `json:"TrainStation"`
	ConnectedLocations []Location `json:"ConnectedLocations"`
	Airports           []Location `json:"Airports"`
	PointOfInterests   []Location `json:"PointOfInterests"`
	Stadiums           []Location `json:"Stadiums"`
	Neighbourhoods     []Location `json:"Neighbourhoods"`
}

// CuratedHotel represents a curated hotel content
type CuratedHotel struct {
	VervotechID            int               `json:"VervotechId"`
	Name                   string            `json:"Name"`
	Rating                 string            `json:"Rating"`
	AddressLine1           string            `json:"AddressLine1"`
	AddressLine2           string            `json:"AddressLine2"`
	CityName               string            `json:"CityName"`
	CityCode               string            `json:"CityCode"`
	StateName              string            `json:"StateName"`
	StateCode              string            `json:"StateCode"`
	CountryName            string            `json:"CountryName"`
	CountryCode            string            `json:"CountryCode"`
	PostalCode             string            `json:"PostalCode"`
	Latitude               float64           `json:"Latitude"`
	Longitude              float64           `json:"Longitude"`
	Website                string            `json:"Website"`
	Fax                    string            `json:"Fax"`
	PropertyType           string            `json:"PropertyType"`
	ChainName              string            `json:"ChainName"`
	BrandName              string            `json:"BrandName"`
	CityLocationID         string            `json:"CityLocationId"`
	MasterCityLocationName string            `json:"MasterCityLocationName"`
	GooglePlusCode         string            `json:"GooglePlusCode"`
	LocationIds            string            `json:"LocationIds"`
	Phones                 []string          `json:"Phones"`
	Emails                 []string          `json:"Emails"`
	Facilities             []Facility        `json:"Facilities"`
	Airports               []Location        `json:"Airports"`
	TrainStations          []Location        `json:"TrainStations"`
	PointOfInterests       []Location        `json:"PointOfInterests"`
	ConnectedLocations     []Location        `json:"ConnectedLocations"`
	Neighbourhoods         []Location        `json:"Neighbourhoods"`
	Stadiums               []Location        `json:"Stadiums"`
	LocationMappings       []LocationMapping `json:"LocationMappings"`
}

// GetCuratedContentByVervotechIdsResponse represents the response from GetCuratedContentByVervotechIds endpoint
type GetCuratedContentByVervotechIdsResponse struct {
	CuratedHotels []CuratedHotel `json:"CuratedHotels"`
	Status        bool           `json:"Status"`
	Message       string         `json:"Message"`
	StatusCode    int            `json:"StatusCode"`
}

// ProviderFacility represents a facility in provider hotel content
type ProviderFacility struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	MasterFacilityID string `json:"masterFacilityId"`
}

// RoomBed represents a bed in a room
type RoomBed struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// RoomArea represents the area of a room
type RoomArea struct {
	SquareMeters string `json:"SquareMeters"`
	SquareFeet   string `json:"SquareFeet"`
	Text         string `json:"Text"`
}

// RoomFacility represents a facility in a room
type RoomFacility struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	MasterFacilityID string `json:"masterFacilityId"`
}

// ProviderRoom represents a room in provider hotel content
type ProviderRoom struct {
	RoomID             string         `json:"RoomId"`
	Type               string         `json:"Type"`
	Category           string         `json:"Category"`
	Occupancy          string         `json:"Occupancy"`
	BedroomCount       int            `json:"BedroomCount"`
	Description        string         `json:"Description"`
	Facilities         []RoomFacility `json:"Facilities"`
	Beds               []RoomBed      `json:"Beds"`
	MaxGuestAllowed    int            `json:"MaxGuestAllowed"`
	MaxAdultAllowed    int            `json:"MaxAdultAllowed"`
	MaxChildrenAllowed int            `json:"MaxChildrenAllowed"`
	Views              []string       `json:"Views"`
	Area               RoomArea       `json:"Area"`
}

// ProviderGeoCode represents geographic coordinates
type ProviderGeoCode struct {
	Lat  float64 `json:"Lat"`
	Long float64 `json:"Long"`
}

// ProviderCheckout represents checkout information
type ProviderCheckout struct {
	Time string `json:"Time"`
}

// ProviderDescription represents a description
type ProviderDescription struct {
	Type string `json:"Type"`
	Text string `json:"Text"`
}

// ProviderPolicy represents a policy
type ProviderPolicy struct {
	Text  string `json:"Text"`
	Title string `json:"Title"`
}

// ProviderCheckin represents checkin information
type ProviderCheckin struct {
	BeginTime           string   `json:"BeginTime"`
	EndTime             string   `json:"EndTime"`
	MinAge              int      `json:"MinAge"`
	Instructions        []string `json:"Instructions"`
	SpecialInstructions []any    `json:"SpecialInstructions"`
}

// ProviderAttribute represents an attribute
type ProviderAttribute struct {
	Key   string `json:"Key"`
	Desc  string `json:"Desc"`
	Value string `json:"Value"`
}

// ProviderStatistic represents a statistic
type ProviderStatistic struct {
	ID    int    `json:"Id"`
	Name  string `json:"Name"`
	Value int    `json:"Value"`
}

// CategoryRating represents a category rating
type CategoryRating struct {
	Category string  `json:"Category"`
	Rating   float64 `json:"Rating"`
}

// GuestReview represents a guest review
type GuestReview struct {
	ReviewerName    string `json:"ReviewerName"`
	Source          string `json:"Source"`
	Title           string `json:"Title"`
	Text            string `json:"Text"`
	Rating          string `json:"Rating"`
	TravelCompanion string `json:"TravelCompanion"`
}

// ProviderReview represents a review
type ProviderReview struct {
	Rating          float64          `json:"Rating"`
	Count           int              `json:"Count"`
	RatingImage     any              `json:"RatingImage"`
	Provider        string           `json:"Provider"`
	CategoryRatings []CategoryRating `json:"CategoryRatings"`
	GuestReviews    []GuestReview    `json:"GuestReviews"`
}

// PointOfInterest represents a point of interest
type PointOfInterest struct {
	Name     string `json:"Name"`
	Type     string `json:"Type"`
	Distance string `json:"Distance"`
	Geocode  any    `json:"Geocode"`
	Unit     string `json:"Unit"`
}

// ProviderFee represents a fee
type ProviderFee struct {
	Text       string `json:"Text"`
	Fees       string `json:"Fees"`
	Frequency  string `json:"Frequency"`
	Assignment string `json:"Assignment"`
	Type       string `json:"Type"`
}

// ProviderAddress represents an address
type ProviderAddress struct {
	Line1           string `json:"Line1"`
	Line2           string `json:"Line2"`
	City            string `json:"City"`
	Code            string `json:"Code"`
	DestinationCode string `json:"DestinationCode"`
	StateCode       string `json:"StateCode"`
	State           string `json:"State"`
	PostalCode      string `json:"PostalCode"`
	Country         string `json:"Country"`
	CountryCode     string `json:"CountryCode"`
}

// ProviderContact represents contact information
type ProviderContact struct {
	Fax     string          `json:"Fax"`
	Web     string          `json:"Web"`
	Phones  []string        `json:"Phones"`
	Emails  []string        `json:"Emails"`
	Address ProviderAddress `json:"Address"`
}

// ImageLink represents a link to an image
type ImageLink struct {
	Href         string `json:"Href"`
	ProviderHref string `json:"ProviderHref"`
	Size         string `json:"Size"`
}

// ProviderImage represents an image
type ProviderImage struct {
	Caption  string      `json:"Caption"`
	Category string      `json:"Category"`
	RoomID   string      `json:"RoomId"`
	Links    []ImageLink `json:"Links"`
}

// ProviderHotel represents a provider hotel
type ProviderHotel struct {
	ProviderHotelID       string                `json:"ProviderHotelId"`
	VervotechID           int                   `json:"VervotechId"`
	LanguageCode          string                `json:"LanguageCode"`
	Name                  string                `json:"Name"`
	ProviderFamily        string                `json:"ProviderFamily"`
	Rating                float64               `json:"Rating"`
	BrandName             string                `json:"BrandName"`
	BrandCode             string                `json:"BrandCode"`
	ChainCode             string                `json:"ChainCode"`
	ChainName             string                `json:"ChainName"`
	MasterPropertyType    string                `json:"MasterPropertyType"`
	PropertyType          string                `json:"PropertyType"`
	Facilities            []ProviderFacility    `json:"Facilities"`
	HeroImage             string                `json:"HeroImage"`
	ProviderHeroImageHref string                `json:"ProviderHeroImageHref"`
	Rooms                 []ProviderRoom        `json:"Rooms"`
	MasterChainName       string                `json:"MasterChainName"`
	PopularityScore       int                   `json:"PopularityScore"`
	ContentScore          int                   `json:"ContentScore"`
	LicenseNumber         any                   `json:"LicenseNumber"`
	GeoCode               ProviderGeoCode       `json:"GeoCode"`
	ProviderGeocodes      ProviderGeoCode       `json:"ProviderGeocodes"`
	Checkout              ProviderCheckout      `json:"Checkout"`
	Descriptions          []ProviderDescription `json:"Descriptions"`
	Policies              []ProviderPolicy      `json:"Policies"`
	Checkin               ProviderCheckin       `json:"Checkin"`
	Attributes            []ProviderAttribute   `json:"Attributes"`
	Statistics            []ProviderStatistic   `json:"Statistics"`
	Reviews               []ProviderReview      `json:"Reviews"`
	PointOfInterests      []PointOfInterest     `json:"PointOfInterests"`
	Fees                  []ProviderFee         `json:"Fees"`
	Contact               ProviderContact       `json:"Contact"`
	Images                []ProviderImage       `json:"Images"`
}

// ProviderHotelGroup represents a group of provider hotels
type ProviderHotelGroup struct {
	ProviderHotels []ProviderHotel `json:"ProviderHotels"`
	VervotechID    int             `json:"VervotechId"`
}

// GetProviderContentByVervotechIdsResponse represents the response from GetProviderContentByVervotechIds endpoint
type GetProviderContentByVervotechIdsResponse struct {
	Status     bool                 `json:"Status"`
	Message    string               `json:"Message"`
	StatusCode int                  `json:"StatusCode"`
	Hotels     []ProviderHotelGroup `json:"Hotels"`
}
