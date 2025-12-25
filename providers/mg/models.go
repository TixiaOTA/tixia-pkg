package mg

import (
	"encoding/json"
	"errors"
	"fmt"
)

// --------------------------------- CONSTANTS ---------------------------------

// MGReservationListDateTypes is the date type of the reservation list
type MGReservationListDateTypes string

// MGReservationListStatuses is the status of the reservation list
type MGReservationListStatuses string

const (
	BookingDateType MGReservationListDateTypes = "1"
	CheckinDateType MGReservationListDateTypes = "2"

	ConfirmedStatus          MGReservationListStatuses = "CONF"
	ReservedStatus           MGReservationListStatuses = "RESERVED"
	NotConfirmedStatus       MGReservationListStatuses = "NOTCONF"
	CancelledStatus          MGReservationListStatuses = "CANCEL"
	CancelledConfirmedStatus MGReservationListStatuses = "CANCELCONF"
	CancelledChargeStatus    MGReservationListStatuses = "CANCELCHARGE"
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

type GetRSVNListRequest struct {
	Login    Credentials                `json:"Login"`
	Date     string                     `json:"Date"`
	DateType MGReservationListDateTypes `json:"DateType"`
	Status   MGReservationListStatuses  `json:"Status,omitempty"`
}

type GetHotelDetailRequest struct {
	Login       Credentials `json:"Login"`
	HotelCode   string      `json:"HotelCode"`
	Language    string      `json:"Language"`
	DetailLevel string      `json:"DetailLevel"`
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

type GetRSVNListResponse struct {
	Status             bool   `json:"status"`
	ResponseTime       int    `json:"responseTime"`
	SessionID          string `json:"sessionID"`
	ReservationDetails []struct {
		MGBookingID     string `json:"mgBookingID"`
		AgencyBookingID string `json:"agencyBookingID"`
		Status          string `json:"status"`
	}
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
	ErrorCode     *string                              `json:"errorCode"`
	ErrorMessage  *string                              `json:"errorMessage"`
	BookingDetail CancellationReservationBookingDetail `json:"bookingDetails"`
}

type GetHotelDetailResponse struct {
	Status       bool             `json:"status"`
	ResponseTime int              `json:"responseTime"`
	SessionID    string           `json:"sessionID"`
	ErrorCode    *string          `json:"errorCode"`
	ErrorMessage *string          `json:"errorMessage"`
	Hotel        HotelDetailHotel `json:"hotel"`
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
	RawBody      *[]byte     `json:"rawBody"`
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

type HotelDetailHotel struct {
	HotelCode        string      `json:"hotelCode"`
	Name             string      `json:"name"`
	Rating           string      `json:"rating"`
	ChainCode        string      `json:"chainCode"`
	ChainName        string      `json:"chainName"`
	BrandCode        string      `json:"brandCode"`
	BrandName        string      `json:"brandName"`
	Type             string      `json:"type"`
	IsMGPreferred    bool        `json:"isMGPreferred"`
	Currency         string      `json:"currency"`
	TopHotelYear     *string     `json:"topHotelYear"`
	Address          Address     `json:"address"`
	Geolocation      Geolocation `json:"geoLocation"`
	ShortDescription Description `json:"shortDescription"`
	LongDescription  Description `json:"longDescription"`
	Reservation      Reservation `json:"reservation"`
	Timezone         string      `json:"timezone"`
	CheckInTime      string      `json:"checkInTime"`
	CheckOutTime     string      `json:"checkOutTime"`
	Website          string      `json:"website"`
	Childpolicy      ChildPolicy `json:"childpolicy"`
	TotalNoOfRooms   string      `json:"totalNoOfRooms"`
	Mealplans        struct {
		Mealplan []Mealplan `json:"mealplan"`
	} `json:"mealplans"`
	Facilities   []Facility `json:"facilities"`
	RoomsDetails struct {
		RoomDetails []HotelDetailRoomDetail `json:"roomDetails"`
	} `json:"roomsDetails"`
	Photo struct {
		Image []Image `json:"image"`
	} `json:"photo"`
}

type Address struct {
	Address1      string `json:"address1"`
	Address2      string `json:"address2"`
	ZipCode       string `json:"zipCode"`
	ContinentCode string `json:"continentCode"`
	ContinentName string `json:"continentName"`
	CountryCode   string `json:"countryCode"`
	CountryName   string `json:"countryName"`
	CityCode      string `json:"cityCode"`
	CityName      string `json:"cityName"`
	Area          string `json:"area"`
	Landmark      string `json:"landmark"`
}

type Geolocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Reservation struct {
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
}

type ChildPolicy struct {
	MinAge int `json:"minAge"`
	MaxAge int `json:"maxAge"`
}

type Mealplan struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Facility struct {
	FacilityGroup []FacilityGroup `json:"facilityGroup"`
}

type FacilityGroup struct {
	Name         string         `json:"name"`
	FacilityType []FacilityType `json:"facilityType"`
}

type FacilityType struct {
	Name     *string `json:"name"`
	Facility []struct {
		Name string `json:"name"`
	} `json:"facility"`
}

type HotelDetailRoomDetail struct {
	RoomCode         string      `json:"roomCode"`
	RoomName         string      `json:"roomName"`
	RoomSize         *string     `json:"roomSize"`
	MaxOccupancy     int         `json:"maxOccupancy"`
	MaxAdults        int         `json:"maxAdults"`
	MaxChildren      int         `json:"maxChildren"`
	IsSmokingAllowed bool        `json:"isSmokingAllowed"`
	RoomDescription  Description `json:"roomDescription"`
	Amenities        struct {
		AmenitiesGroup []AmenitiesGroup `json:"amenitiesGroup"`
	}
	Photo struct {
		Image []Image `json:"image"`
	} `json:"photo"`
}

type AmenitiesGroup struct {
	AmenityGroupName  string              `json:"amenityGroupName"`
	AmenitiesTypeList []AmenitiesTypeList `json:"amenitiesTypeList"`
}

type AmenitiesTypeList struct {
	AmenityTypeName *string   `json:"amenityTypeName"`
	AmenitiesList   []Amenity `json:"amenitiesList"`
}

type Amenity struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Description struct {
	Content      string               `json:"content"`
	CDataContent []DescriptionContent `json:"cDataContent"`
}

type DescriptionContent struct {
	CDataSection string `json:"#cdata-section"`
}

type Image struct {
	URL    string `json:"url"`
	IsMain bool   `json:"isMain"`
	IsCDN  bool   `json:"isCDN"`
}
