package mg

import (
	"encoding/json"
	"errors"
	"fmt"
)

// --------------------------------- REQUESTS ---------------------------------

type SearchHotelRequest struct {
	Login       Credentials `json:"Login"`
	Nationality string      `json:"Nationality"`
	Country     string      `json:"Country"`
	City        string      `json:"City"`
	Hotels      *struct {
		Code []string `json:"Code"`
	} `json:"Hotels,omitempty"`
	CheckIn  string `json:"CheckIn"`
	CheckOut string `json:"CheckOut"`
	Rooms    *struct {
		Room []RoomRequest `json:"Room"`
	} `json:"Rooms"`
	Currency     string `json:"Currency"`
	Language     string `json:"Language"`
	AvailFlag    bool   `json:"AvailFlag"`
	MaxNoOfHotel int    `json:"MaxNoOfHotel"`
	DetailLevel  string `json:"DetailLevel"`
}

type RecheckHotelRequest struct {
	Login       Credentials         `json:"Login"`
	SessionID   string              `json:"SessionId"`
	Nationality string              `json:"Nationality"`
	Country     string              `json:"Country"`
	City        string              `json:"City"`
	HotelCode   string              `json:"HotelCode,omitempty"`
	CheckIn     string              `json:"CheckIn"`  // Format: YYYY-MM-DD
	CheckOut    string              `json:"CheckOut"` // Format: YYYY-MM-DD
	RoomDetails *RoomDetailsRequest `json:"RoomDetails"`
	Rooms       struct {
		Room []RoomRequest `json:"Room"`
	} `json:"Rooms"`
	Currency    string `json:"Currency"`
	Language    string `json:"Language,omitempty"`
	AvailFlag   bool   `json:"AvailFlag,omitempty"`
	DetailLevel string `json:"DetailLevel,omitempty"`
}

// BookHotelRequest represents the request structure for booking a hotel
type BookHotelRequest struct {
	Login           Credentials        `json:"Login"`
	SessionID       string             `json:"SessionId"`
	AgencyBookingID string             `json:"AgencyBookingID"` // Tixia generated booking ID
	Nationality     string             `json:"Nationality"`
	CheckIn         string             `json:"CheckIn"`  // Format: YYYY-MM-DD
	CheckOut        string             `json:"CheckOut"` // Format: YYYY-MM-DD
	HotelCode       string             `json:"HotelCode"`
	RoomDetails     RoomDetailsRequest `json:"RoomDetails"`
	Rooms           struct {
		Room []RoomBookingRequest `json:"Room"`
	} `json:"Rooms"`
	Currency    string `json:"Currency"`
	Language    string `json:"Language,omitempty"`
	AvailFlag   bool   `json:"AvailFlag,omitempty"`
	OnHold      bool   `json:"OnHold,omitempty"`
	SpecialReq  string `json:"SpecialReq,omitempty"`
	DetailLevel string `json:"DetailLevel,omitempty"`
}

// GetRSVNDetailsRequest represents the request structure for getting the RSVN details
type GetRSVNDetailsRequest struct {
	Login           Credentials `json:"Login"`
	MGBookingID     string      `json:"MGBookingID"`
	AgencyBookingID string      `json:"AgencyBookingID"` // Tixia generated booking ID
}

type CancelReservationRequest struct {
	Login           Credentials `json:"Login"`
	AgencyBookingID string      `json:"AgencyBookingID"` // Tixia generated booking ID
}

// --------------------------------- RESPONSES ---------------------------------

type SearchHotelResponse struct {
	CheckIn      string  `json:"checkIn"`
	CheckOut     string  `json:"checkOut"`
	Currency     string  `json:"currency"`
	ErrorCode    *string `json:"errorCode"`
	ErrorMessage *string `json:"errorMessage"`
	Hotels       struct {
		Hotel []Hotel `json:"Hotel"`
	} `json:"hotels"`
	NoOfHotels   int     `json:"noOfHotels"`
	ResponseTime int     `json:"responseTime"`
	SessionID    string  `json:"sessionID"`
	Status       bool    `json:"status"`
	XSD          *string `json:"xsd"`
	XSI          *string `json:"xsi"`
}

type RecheckHotelResponse struct {
	Status       bool    `json:"status"`
	ResponseTime int     `json:"responseTime"`
	SessionID    string  `json:"sessionID"`
	CheckIn      string  `json:"checkIn"`
	CheckOut     string  `json:"checkOut"`
	Currency     string  `json:"currency"`
	ErrorCode    *string `json:"errorCode"`
	ErrorMessage *string `json:"errorMessage"`
	Nationality  *string `json:"nationality"`

	Hotels struct {
		Hotel HotelRecheckResponse `json:"hotel"`
	} `json:"hotels,omitempty"`
}

type BookHotelResponse struct {
	Status        bool           `json:"status"`
	ResponseTime  int            `json:"responseTime"`
	SessionID     string         `json:"sessionID"`
	BookingDetail BookingDetails `json:"bookingDetails"`
}

// GetRSVNDetailsResponse represents the response structure for getting the RSVN details
type GetRSVNDetailsResponse struct {
	Status        bool           `json:"status"`
	ResponseTime  int            `json:"responseTime"`
	SessionID     string         `json:"sessionID"`
	BookingDetail BookingDetails `json:"bookingDetails"`
}

// CancelReservationResponse represents the response structure for canceling a reservation
type CancelReservationResponse struct {
	Status        bool                                 `json:"status"`
	ResponseTime  int                                  `json:"responseTime"`
	SessionID     string                               `json:"sessionID"`
	BookingDetail CancellationReservationBookingDetail `json:"bookingDetails"`
}

// --------------------------------- SHARED TYPES ---------------------------------

type Credentials struct {
	AgencyCode string `json:"AgencyCode" log:"redact"`
	Username   string `json:"Username" log:"redact"`
	Password   string `json:"Password" log:"redact"`
}

func (c *Credentials) Validate() error {
	if c.AgencyCode == "" {
		return errors.New("agency code is required")
	}
	if c.Username == "" {
		return errors.New("username is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (code %d): %s", e.Code, e.Message)
}

type MGError struct {
	Status       bool        `json:"status"`
	ErrorCode    MGErrorCode `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
}

func (e *MGError) Error() string {
	return fmt.Sprintf("MG error (status %t): %s - %s", e.Status, e.ErrorCode, e.ErrorMessage)
}

type RoomRequest struct {
	RoomNo     uint   `json:"RoomNo"`
	NoOfAdults uint   `json:"NoOfAdults"`
	NoOfChild  uint   `json:"NoOfChild"`
	Child1Age  uint   `json:"Child1Age"`
	Child2Age  uint   `json:"Child2Age"`
	ExtraBed   bool   `json:"ExtraBed"`
	RateKey    string `json:"RateKey,omitempty"`
}

type Hotel struct {
	Code        string       `json:"code"`
	Latitude    string       `json:"latitude"`
	Longitude   string       `json:"longitude"`
	Name        string       `json:"name"`
	Rating      string       `json:"rating"`
	RoomDetails []RoomDetail `json:"roomDetails"`
}

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

type RoomDetail struct {
	AvailFlag            bool    `json:"availFlag"`
	AvgNightPrice        float64 `json:"avgNightPrice"`
	B2BMarkup            float64 `json:"b2BMarkup"`
	CanHold              bool    `json:"canHold"`
	CancellationPolicies struct {
		CancellationPolicy []CancellationPolicy `json:"policy"`
	} `json:"cancellationPolicies"`
	CancellationPolicyType string  `json:"cancellationPolicyType"`
	Code                   string  `json:"code"`
	GrossPrice             float64 `json:"grossPrice"`
	MealPlan               string  `json:"mealPlan"`
	MealPlanName           string  `json:"mealPlanName"`
	Messages               struct {
		Message []struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"messages"`
	MSP         string  `json:"msp"`
	Name        string  `json:"name"`
	NetPrice    float64 `json:"netPrice"`
	PackageRate bool    `json:"packageRate"`
	PromoCode   string  `json:"promoCode"`
	PromoName   string  `json:"promoName"`
	Rooms       struct {
		Room []Room `json:"Room"`
	} `json:"rooms"`
}

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

type RateDetails struct {
	DiscountDetails *DiscountDetails `json:"discountDetails"` // assumed to be nullable
	NightlyRates    NightlyRates     `json:"nightlyRates"`
}

type DiscountDetails struct {
	Total json.Number `json:"total"`
	Promo []Promo     `json:"promo"`
}

type Promo struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

type NightlyRates struct {
	NightlyRate []NightlyRate `json:"nightlyRate"`
}

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

type Supplement struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
}

type Compulsory struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
}

type RoomDetailsRequest struct {
	Code                   string `json:"Code"`
	MealPlan               string `json:"MealPlan"`
	CancellationPolicyType string `json:"CancellationPolicyType"`
	PackageRate            bool   `json:"PackageRate"`
}

type HotelRecheckResponse struct {
	Code        string     `json:"code"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	Name        string     `json:"name"`
	Rating      string     `json:"rating"`
	RoomDetails RoomDetail `json:"roomDetails"`
}

type RoomBookingRequest struct {
	PaxDetails struct {
		Pax []Pax `json:"Pax"`
	} `json:"PaxDetails"`
	RoomNo     uint   `json:"RoomNo"`
	NoOfAdults uint   `json:"NoOfAdults"`
	NoOfChild  uint   `json:"NoOfChild"`
	ExtraBed   bool   `json:"ExtraBed"`
	RateKey    string `json:"RateKey"`
}

type BookingDetails struct {
	Nationality       string       `json:"nationality"`
	CheckIn           string       `json:"checkIn"`
	CheckOut          string       `json:"checkOut"`
	Currency          string       `json:"currency"`
	MgBookingID       string       `json:"mgBookingID"`
	AgencyBookingID   string       `json:"agencyBookingID"`
	MgBookingVerID    string       `json:"mgBookingVersionID"`
	AgencyVoucherNo   string       `json:"agencyVoucherNo"`
	AgencyVoucherDate string       `json:"agencyVoucherDate"`
	Status            string       `json:"status"`
	SpecialRequest    CDataWrapper `json:"specialRequest"`
	Hotels            struct {
		Hotel HotelRecheckResponse `json:"hotel"`
	} `json:"hotels"`
}

type CDataWrapper struct {
	CDataContent []CData `json:"cDataContent"`
}

type CData struct {
	CDataSection string `json:"#cdata-section"`
}

type Pax struct {
	Salutation string `json:"Salutation"`
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	Age        uint   `json:"Age"`
}

type CancellationReservationCancellationCharges struct {
	Amount      string `json:"amount"`
	B2BMarkup   string `json:"b2BMarkup"`
	GrossAmount string `json:"grossAmount"`
	Currency    string `json:"currency"`
}

type CancellationReservationCancellationPolicies struct {
	CancellationCharges CancellationReservationCancellationCharges `json:"cancellationCharges"`
}

type CancellationReservationBookingDetail struct {
	MgBookingID          string                                      `json:"mgBookingID"`
	AgencyBookingID      string                                      `json:"agencyBookingID"`
	MgBookingVerID       string                                      `json:"mgBookingVersionID"`
	AgencyVoucherNo      string                                      `json:"agencyVoucherNo"`
	AgencyVoucherDate    string                                      `json:"agencyVoucherDate"`
	HotelBookingID       string                                      `json:"hotelBookingID"`
	Status               string                                      `json:"status"`
	CancellationPolicies CancellationReservationCancellationPolicies `json:"cancellationPolicies"`
}
