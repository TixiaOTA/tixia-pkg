package mg

// --------------------------------- REQUESTS ---------------------------------

type SearchHotelRequestDTO struct {
	Nationality string `json:"Nationality"`
	Country     string `json:"Country"`
	City        string `json:"City"`
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

func (dto *SearchHotelRequestDTO) ToInternalModel(credentials Credentials) SearchHotelRequest {
	return SearchHotelRequest{
		Login:        credentials,
		Nationality:  dto.Nationality,
		Country:      dto.Country,
		City:         dto.City,
		Hotels:       dto.Hotels,
		CheckIn:      dto.CheckIn,
		CheckOut:     dto.CheckOut,
		Rooms:        dto.Rooms,
		Currency:     dto.Currency,
		Language:     dto.Language,
		AvailFlag:    dto.AvailFlag,
		MaxNoOfHotel: dto.MaxNoOfHotel,
		DetailLevel:  dto.DetailLevel,
	}
}

type RecheckHotelRequestDTO struct {
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

func (dto *RecheckHotelRequestDTO) ToInternalModel(credentials Credentials) RecheckHotelRequest {
	return RecheckHotelRequest{
		Login:       credentials,
		SessionID:   dto.SessionID,
		Nationality: dto.Nationality,
		Country:     dto.Country,
		City:        dto.City,
		HotelCode:   dto.HotelCode,
		CheckIn:     dto.CheckIn,
		CheckOut:    dto.CheckOut,
		RoomDetails: dto.RoomDetails,
		Rooms:       dto.Rooms,
		Currency:    dto.Currency,
		Language:    dto.Language,
		AvailFlag:   dto.AvailFlag,
		DetailLevel: dto.DetailLevel,
	}
}

type BookHotelRequestDTO struct {
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

func (dto *BookHotelRequestDTO) ToInternalModel(credentials Credentials) BookHotelRequest {
	return BookHotelRequest{
		Login:           credentials,
		SessionID:       dto.SessionID,
		AgencyBookingID: dto.AgencyBookingID,
		Nationality:     dto.Nationality,
		CheckIn:         dto.CheckIn,
		CheckOut:        dto.CheckOut,
		HotelCode:       dto.HotelCode,
		RoomDetails:     dto.RoomDetails,
		Rooms:           dto.Rooms,
		Currency:        dto.Currency,
		Language:        dto.Language,
		AvailFlag:       dto.AvailFlag,
		OnHold:          dto.OnHold,
		SpecialReq:      dto.SpecialReq,
		DetailLevel:     dto.DetailLevel,
	}
}

type GetRSVNDetailsRequestDTO struct {
	MGBookingID     string `json:"MGBookingID"`
	AgencyBookingID string `json:"AgencyBookingID"` // Tixia generated booking ID
}

func (dto *GetRSVNDetailsRequestDTO) ToInternalModel(credentials Credentials) GetRSVNDetailsRequest {
	return GetRSVNDetailsRequest{
		Login:           credentials,
		MGBookingID:     dto.MGBookingID,
		AgencyBookingID: dto.AgencyBookingID,
	}
}

type CancelReservationRequestDTO struct {
	AgencyBookingID string `json:"AgencyBookingID"` // Tixia generated booking ID
}

func (dto *CancelReservationRequestDTO) ToInternalModel(credentials Credentials) CancelReservationRequest {
	return CancelReservationRequest{
		Login:           credentials,
		AgencyBookingID: dto.AgencyBookingID,
	}
}
