package mg

import (
	"encoding/json"
	"time"

	httppkg "github.com/TixiaOTA/tixia-pkg/http/client"
	validator "github.com/go-playground/validator/v10"
)

// MGConfig holds the configuration for MG provider
type MGConfig struct {
	BaseURL   string               `validate:"required,url"`
	AgentID   string               `validate:"required"`
	LoginName string               `validate:"required"`
	Password  string               `validate:"required"`
	Timeout   *time.Duration       // Optional timeout for HTTP requests (default: 2 seconds)
	Retry     *httppkg.RetryConfig // Optional retry configuration (default: MaxRetries=3, InitialBackoff=200ms, MaxBackoff=5s, BackoffMultiplier=2.0)
}

func (c *MGConfig) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// Login represents the login credentials for MG API
type Login struct {
	AgencyCode string `json:"AgencyCode"`
	Username   string `json:"Username"`
	Password   string `json:"Password"`
}

// Hotels represents a list of hotel codes
type Hotels struct {
	Code []string `json:"Code"`
}

// Rooms represents a list of room details
type Rooms struct {
	Room []RoomDetailSearchRequest `json:"Room"`
}

// HotelSearchRequest represents the request structure for searching hotels
type HotelSearchRequest struct {
	Login        *Login  `json:"Login"` // Login is optional, if not provided, it will be set from the config
	Nationality  string  `json:"Nationality" validate:"required"`
	Country      string  `json:"Country"`
	City         string  `json:"City"`
	Hotels       *Hotels `json:"Hotels"` // Hotels is optional
	CheckIn      string  `json:"CheckIn" validate:"required"`
	CheckOut     string  `json:"CheckOut" validate:"required"`
	Rooms        *Rooms  `json:"Rooms"` // Rooms is optional
	Currency     string  `json:"Currency"`
	Language     string  `json:"Language"`
	AvailFlag    bool    `json:"AvailFlag"`
	MaxNoOfHotel int     `json:"MaxNoOfHotel"`
	DetailLevel  string  `json:"DetailLevel"`
}

func (h HotelSearchRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(h)
}

// RoomDetailSearchRequest represents room details in a search request
type RoomDetailSearchRequest struct {
	RoomNo     uint `json:"RoomNo"`
	NoOfAdults uint `json:"NoOfAdults"`
	NoOfChild  uint `json:"NoOfChild"`
	Child1Age  uint `json:"Child1Age"`
	Child2Age  uint `json:"Child2Age"`
	ExtraBed   bool `json:"ExtraBed"`
}

// HotelSearchResponse represents the response structure for hotel search
type HotelSearchResponse struct {
	CheckIn      string     `json:"checkIn"`
	CheckOut     string     `json:"checkOut"`
	Currency     string     `json:"currency"`
	ErrorCode    *string    `json:"errorCode"`
	ErrorMessage *string    `json:"errorMessage"`
	Hotels       HotelsData `json:"hotels"`
	NoOfHotels   int        `json:"noOfHotels"`
	ResponseTime int        `json:"responseTime"`
	SessionID    string     `json:"sessionID"`
	Status       bool       `json:"status"`
	XSD          *string    `json:"xsd"`
	XSI          *string    `json:"xsi"`
}

// HotelsData represents the hotels data in the response
type HotelsData struct {
	Hotel []Hotel `json:"hotel"`
}

// Hotel represents a hotel in the search response
type Hotel struct {
	Code        string                     `json:"code"`
	Latitude    string                     `json:"latitude"`
	Longitude   string                     `json:"longitude"`
	Name        string                     `json:"name"`
	Rating      string                     `json:"rating"`
	RoomDetails []RoomDetailSearchResponse `json:"roomDetails"`
}

// RoomDetailSearchResponse represents room details in a search response
type RoomDetailSearchResponse struct {
	AvailFlag              bool                 `json:"availFlag"`
	AvgNightPrice          float64              `json:"avgNightPrice"`
	B2BMarkup              float64              `json:"b2BMarkup"`
	CanHold                bool                 `json:"canHold"`
	CancellationPolicies   CancellationPolicies `json:"cancellationPolicies"`
	CancellationPolicyType string               `json:"cancellationPolicyType"`
	Code                   string               `json:"code"`
	GrossPrice             float64              `json:"grossPrice"`
	MealPlan               string               `json:"mealPlan"`
	MealPlanName           string               `json:"mealPlanName"`
	Messages               Messages             `json:"messages"`
	MSP                    string               `json:"msp"`
	Name                   string               `json:"name"`
	NetPrice               float64              `json:"netPrice"`
	PackageRate            bool                 `json:"packageRate"`
	PromoCode              string               `json:"promoCode"`
	PromoName              string               `json:"promoName"`
	Rooms                  RoomsInfo            `json:"rooms"`
}

// CancellationPolicies represents cancellation policies
type CancellationPolicies struct {
	Policy []CancellationPolicy `json:"policy"`
}

// CancellationPolicy represents a single cancellation policy
type CancellationPolicy struct {
	Amount      string `json:"amount"`
	B2BMarkup   string `json:"b2BMarkup"`
	FromDate    string `json:"fromDate"`
	GrossAmount string `json:"grossAmount"`
	Nights      string `json:"nights"`
	NoShow      bool   `json:"noShow"`
	Percent     string `json:"percent"`
	ToDate      string `json:"toDate"`
}

// ProviderMessageContent represents a message content
type ProviderMessageContent struct {
	Content string `json:"content"`
}

// Messages represents messages in the response
type Messages struct {
	Message []ProviderMessageContent `json:"message"`
}

// RoomsInfo represents room information
type RoomsInfo struct {
	Room []Room `json:"room"`
}

// Room represents a room in the response
type Room struct {
	AvgNightPrice float64     `json:"avgNightPrice"`
	B2BMarkup     float64     `json:"b2BMarkup"`
	Child1Age     float32     `json:"child1Age"`
	Child2Age     float32     `json:"child2Age"`
	ExtraBed      bool        `json:"extraBed"`
	GrossPrice    float64     `json:"grossPrice"`
	NetPrice      float64     `json:"netPrice"`
	NoOfAdults    int         `json:"noOfAdults"`
	NoOfChild     int         `json:"noOfChild"`
	RateDetails   RateDetails `json:"rateDetails"`
	RateKey       string      `json:"rateKey"`
	RoomNo        int         `json:"roomNo"`
}

// RateDetails represents rate details
type RateDetails struct {
	DiscountDetails *DiscountDetails `json:"discountDetails"`
	NightlyRates    NightlyRates     `json:"nightlyRates"`
}

// DiscountDetails represents discount details
type DiscountDetails struct {
	Total json.Number `json:"total"`
	Promo []Promo     `json:"promo"`
}

// Promo represents a promotion
type Promo struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

// NightlyRates represents nightly rates
type NightlyRates struct {
	NightlyRate []NightlyRate `json:"nightlyRate"`
}

// NightlyRate represents a single nightly rate
type NightlyRate struct {
	SrNo                 int     `json:"srNo"`
	NetPrice             float64 `json:"netPrice"`
	B2BMarkup            float64 `json:"b2BMarkup"`
	GrossPrice           float64 `json:"grossPrice"`
	SupplementaryDetails struct {
		Supplement []Supplement `json:"supplement"`
	} `json:"supplementaryDetails"`
	CompulsoryDetails struct {
		Compulsory []Compulsory `json:"compulsory"`
	} `json:"compulsoryDetails"`
}

// Supplement represents a supplement
type Supplement struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
}

// Compulsory represents a compulsory charge
type Compulsory struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
}
