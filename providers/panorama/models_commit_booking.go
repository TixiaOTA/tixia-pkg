package panorama

import (
	"encoding/json"
	"errors"
)

// --------------------------------- COMMIT BOOKING REQUEST ---------------------------------

// CommitBookingRequest is the request structure for CommitBooking
type CommitBookingRequest struct {
	Signature   string `json:"signature"`
	ForceIssued bool   `json:"force_issued"`
}

func (r *CommitBookingRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

// --------------------------------- COMMIT BOOKING RESPONSE ---------------------------------

// CommitBookingResponse is the response structure for CommitBooking
type CommitBookingResponse struct {
	Result struct {
		BasicPanoramaResult
		Response BookingResponse `json:"response"`
	} `json:"result"`
}

// --------------------------------- BOOKING RESPONSE ---------------------------------

// BookingResponse is the booking payload returned by commit_booking / get_booking / cancel
type BookingResponse struct {
	OrderNumber            string             `json:"order_number"`
	HoID                   int                `json:"ho_id"`
	AgentName              string             `json:"agent_name"`
	AgentSeqID             string             `json:"agent_seq_id"`
	AgentTypeName          string             `json:"agent_type_name"`
	CustomerParentID       any                `json:"customer_parent_id"`
	CustomerParentName     string             `json:"customer_parent_name"`
	CustomerParentSeqID    string             `json:"customer_parent_seq_id"`
	CustomerParentTypeName string             `json:"customer_parent_type_name"`
	PNR                    string             `json:"pnr"`
	State                  string             `json:"state"`
	StateDescription       string             `json:"state_description"`
	ReconcileState         string             `json:"reconcile_state"`
	HoldDate               string             `json:"hold_date"`
	PreviewImageURL        any                `json:"preview_image_url"`
	YCD                    int                `json:"YCD"`
	ADT                    int                `json:"ADT"`
	CHD                    int                `json:"CHD"`
	INF                    int                `json:"INF"`
	IsAgent                bool               `json:"is_agent"`
	Contact                BookingContact     `json:"contact"`
	VoucherReference       any                `json:"voucher_reference"`
	VoucherDiscount        float64            `json:"voucher_discount"`
	Booker                 BookingBooker      `json:"booker"`
	DepartureDate          string             `json:"departure_date"`
	ArrivalDate            string             `json:"arrival_date"`
	ProviderType           string             `json:"provider_type"`
	PaymentAcquirerNumber  json.RawMessage    `json:"payment_acquirer_number"`
	IssuedRequestStatus    string             `json:"issued_request_status"`
	AgentNta               float64            `json:"agent_nta"`
	CreateDate             string             `json:"create_date"`
	BookedDate             string             `json:"booked_date"`
	BookedBy               string             `json:"booked_by"`
	IssuedBy               any                `json:"issued_by"`
	IssuedDate             string             `json:"issued_date"`
	UsePoint               bool               `json:"use_point"`
	SignatureBooked        string             `json:"signature_booked"`
	Currency               string             `json:"currency"`
	EstimatedCurrency      EstimatedCurrency  `json:"estimated_currency"`
	Total                  float64            `json:"total"`
	PointReward            float64            `json:"point_reward"`
	Invoices               []string           `json:"invoices"`
	InvoiceLines           []string           `json:"invoice_lines"`
	IsHoldDateSync         bool               `json:"is_hold_date_sync"`
	Direction              string             `json:"direction"`
	Origin                 string             `json:"origin"`
	OriginDisplayName      string             `json:"origin_display_name"`
	Destination            string             `json:"destination"`
	DestinationDisplayName string             `json:"destination_display_name"`
	SectorType             string             `json:"sector_type"`
	Passengers             []BookingPassenger `json:"passengers"`
	ProviderBookings       []ProviderBooking  `json:"provider_bookings"`
	RefundList             []any              `json:"refund_list"`
	RescheduleList         []any              `json:"reschedule_list"`
	ExpiredDate            string             `json:"expired_date"`
}

// EstimatedCurrency is currency conversion info on a booking
type EstimatedCurrency struct {
	TotalPrice    float64                  `json:"total_price"`
	Currency      string                   `json:"currency"`
	OtherCurrency []EstimatedOtherCurrency `json:"other_currency"`
}

// EstimatedOtherCurrency is an alternate currency amount on a booking
type EstimatedOtherCurrency struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
}

// BookingContact is the contact summary on a booking
type BookingContact struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// BookingBooker is the booker details on a booking
type BookingBooker struct {
	Name            string          `json:"name"`
	FirstName       string          `json:"first_name"`
	FaceImage       []any           `json:"face_image"`
	LastName        string          `json:"last_name"`
	Gender          string          `json:"gender"`
	BirthDate       string          `json:"birth_date"`
	NationalityCode string          `json:"nationality_code"`
	NationalityName string          `json:"nationality_name"`
	MaritalStatus   string          `json:"marital_status"`
	Phones          []BookingPhone  `json:"phones"`
	Email           string          `json:"email"`
	SeqID           string          `json:"seq_id"`
	Identities      json.RawMessage `json:"identities"`
	Behaviors       json.RawMessage `json:"behaviors"`
	OriginalAgent   string          `json:"original_agent"`
	FrequentFlyers  []any           `json:"frequent_flyers"`
	IsIDAltName     bool            `json:"is_id_alt_name"`
	RizText         string          `json:"riz_text"`
}

// BookingPhone is a phone entry on booker/contact
type BookingPhone struct {
	Type          string `json:"type"`
	CallingCode   string `json:"calling_code"`
	CallingNumber string `json:"calling_number"`
}

// BookingPassenger is a passenger on a booking
type BookingPassenger struct {
	Name                        string                              `json:"name"`
	FirstName                   string                              `json:"first_name"`
	LastName                    string                              `json:"last_name"`
	CustFirstName               string                              `json:"cust_first_name"`
	CustLastName                string                              `json:"cust_last_name"`
	Gender                      string                              `json:"gender"`
	Title                       string                              `json:"title"`
	BirthDate                   string                              `json:"birth_date"`
	NationalityCode             string                              `json:"nationality_code"`
	IdentityCountryOfIssuedCode string                              `json:"identity_country_of_issued_code"`
	IdentityType                string                              `json:"identity_type"`
	IdentityNumber              string                              `json:"identity_number"`
	IdentityExpdate             string                              `json:"identity_expdate"`
	Sequence                    int                                 `json:"sequence"`
	IsValidIdentity             bool                                `json:"is_valid_identity"`
	DescriptionFileURL          any                                 `json:"description_file_url"`
	Description                 any                                 `json:"description"`
	SaleServiceCharges          map[string]map[string]ServiceCharge `json:"sale_service_charges"`
	ServiceChargeDetails        []ServiceChargeDetail               `json:"service_charge_details"`
	Fees                        []any                               `json:"fees"`
	Behaviors                   json.RawMessage                     `json:"behaviors"`
	SeqID                       string                              `json:"seq_id"`
	PaxType                     string                              `json:"pax_type"`
	RizText                     string                              `json:"riz_text"`
	FrequentFlyers              []any                               `json:"frequent_flyers"`
	PassengerNumber             int                                 `json:"passenger_number"`
}

// ServiceCharge is a single charge line under sale_service_charges
type ServiceCharge struct {
	Amount          float64 `json:"amount"`
	ForeignAmount   float64 `json:"foreign_amount"`
	ChargeCode      string  `json:"charge_code"`
	Currency        string  `json:"currency"`
	ForeignCurrency string  `json:"foreign_currency"`
	PaxType         string  `json:"pax_type"`
}

// ServiceChargeDetail is a service charge detail block for a PNR
type ServiceChargeDetail struct {
	PNR            string                     `json:"pnr"`
	ServiceCharges map[string][]ServiceCharge `json:"service_charges"`
	PaxType        string                     `json:"pax_type"`
	PaxCount       int                        `json:"pax_count"`
	BaseFare       float64                    `json:"base_fare"`
	BaseTax        float64                    `json:"base_tax"`
	BaseNta        float64                    `json:"base_nta"`
	BasePrice      float64                    `json:"base_price"`
	BaseCommission float64                    `json:"base_commission"`
}

// ProviderBooking is a provider-level booking under the order
type ProviderBooking struct {
	PNR                    string           `json:"pnr"`
	PNR2                   string           `json:"pnr2"`
	Provider               string           `json:"provider"`
	ProviderName           string           `json:"provider_name"`
	ProviderAlias          string           `json:"provider_alias"`
	AgentID                int              `json:"agent_id"`
	State                  string           `json:"state"`
	StateDescription       string           `json:"state_description"`
	Sequence               int              `json:"sequence"`
	Origin                 string           `json:"origin"`
	OriginDisplayName      string           `json:"origin_display_name"`
	Destination            string           `json:"destination"`
	DestinationDisplayName string           `json:"destination_display_name"`
	DepartureDate          string           `json:"departure_date"`
	ArrivalDate            string           `json:"arrival_date"`
	Journeys               []BookingJourney `json:"journeys"`
	Currency               string           `json:"currency"`
	HoldDate               string           `json:"hold_date"`
	Tickets                []BookingTicket  `json:"tickets"`
	SectorType             string           `json:"sector_type"`
	ErrorMsg               string           `json:"error_msg"`
	Reference              string           `json:"reference"`
	TotalPriceCost         float64          `json:"total_price_cost"`
	PenaltyAmount          float64          `json:"penalty_amount"`
	PenaltyCurrency        string           `json:"penalty_currency"`
	IsAdvancePurchase      bool             `json:"is_advance_purchase"`
	Rules                  []BookingRule    `json:"rules"`
	PNRReferences          []any            `json:"pnr_references"`
	Duplicates             []any            `json:"duplicates"`
	DuplicatesBackend      string           `json:"duplicates_backend"`
	VendorRemarks          []any            `json:"vendor_remarks"`
	IsHaltIssued           bool             `json:"is_halt_issued"`
	Fees                   []any            `json:"fees"`
	PromoCodes             []any            `json:"promo_codes"`
}

// BookingJourney is a journey under a provider booking
type BookingJourney struct {
	Sequence               int              `json:"sequence"`
	Origin                 string           `json:"origin"`
	OriginDisplayName      string           `json:"origin_display_name"`
	Destination            string           `json:"destination"`
	DestinationDisplayName string           `json:"destination_display_name"`
	JourneyKey             string           `json:"journey_key"`
	DepartureDate          string           `json:"departure_date"`
	ArrivalDate            string           `json:"arrival_date"`
	Segments               []BookingSegment `json:"segments"`
	SearchBanner           []any            `json:"search_banner"`
	IsVtlFlight            bool             `json:"is_vtl_flight"`
	JourneyCode            string           `json:"journey_code"`
	SectorType             string           `json:"sector_type"`
}

// BookingSegment is a segment under a booking journey
type BookingSegment struct {
	SegmentCode            string              `json:"segment_code"`
	FareCode               string              `json:"fare_code"`
	PNR                    string              `json:"pnr"`
	CarrierName            string              `json:"carrier_name"`
	CarrierCode            string              `json:"carrier_code"`
	CarrierTypeCode        string              `json:"carrier_type_code"`
	CarrierTypeName        string              `json:"carrier_type_name"`
	CarrierNumber          string              `json:"carrier_number"`
	Provider               string              `json:"provider"`
	Origin                 string              `json:"origin"`
	OriginDisplayName      string              `json:"origin_display_name"`
	OriginTerminal         string              `json:"origin_terminal"`
	Destination            string              `json:"destination"`
	DestinationDisplayName string              `json:"destination_display_name"`
	DestinationTerminal    string              `json:"destination_terminal"`
	DepartureDate          string              `json:"departure_date"`
	ArrivalDate            string              `json:"arrival_date"`
	ElapsedTime            string              `json:"elapsed_time"`
	ClassOfService         string              `json:"class_of_service"`
	CabinClass             string              `json:"cabin_class"`
	GroupCode              string              `json:"group_code"`
	Sequence               int                 `json:"sequence"`
	Seats                  []any               `json:"seats"`
	Legs                   []BookingLeg        `json:"legs"`
	FareDetails            []BookingFareDetail `json:"fare_details"`
	FareBasisCode          string              `json:"fare_basis_code"`
	FareClass              string              `json:"fare_class"`
	FareName               string              `json:"fare_name"`
	OperatingAirlineCode   string              `json:"operating_airline_code"`
	Description            []any               `json:"description"`
	FareRemarks            []any               `json:"fare_remarks"`
	Status                 string              `json:"status"`
	ValidatingCarrier      string              `json:"validating_carrier"`
}

// BookingLeg is a leg under a booking segment
type BookingLeg struct {
	LegCode                string `json:"leg_code"`
	Origin                 string `json:"origin"`
	OriginDisplayName      string `json:"origin_display_name"`
	Destination            string `json:"destination"`
	DestinationDisplayName string `json:"destination_display_name"`
	OriginTerminal         string `json:"origin_terminal"`
	DestinationTerminal    string `json:"destination_terminal"`
	DepartureDate          string `json:"departure_date"`
	ArrivalDate            string `json:"arrival_date"`
	ElapsedTime            string `json:"elapsed_time"`
	Sequence               int    `json:"sequence"`
}

// BookingFareDetail is a fare detail under a booking segment
type BookingFareDetail struct {
	DetailCode  string  `json:"detail_code"`
	DetailType  string  `json:"detail_type"`
	DetailName  string  `json:"detail_name"`
	PaxTypeJSON string  `json:"pax_type_json"`
	Amount      float64 `json:"amount"`
	Unit        string  `json:"unit"`
	Description string  `json:"description"`
	Sequence    int     `json:"sequence"`
}

// BookingTicket is a ticket under a provider booking
type BookingTicket struct {
	TicketID                    string `json:"ticket_id"`
	Passenger                   string `json:"passenger"`
	PaxType                     string `json:"pax_type"`
	TicketNumber                string `json:"ticket_number"`
	FfNumber                    string `json:"ff_number"`
	FfCode                      string `json:"ff_code"`
	FfName                      string `json:"ff_name"`
	Fees                        []any  `json:"fees"`
	Sequence                    any    `json:"sequence"`
	PassengerNumber             int    `json:"passenger_number"`
	Title                       string `json:"title"`
	FirstName                   string `json:"first_name"`
	LastName                    string `json:"last_name"`
	BirthDate                   string `json:"birth_date"`
	IdentityType                string `json:"identity_type"`
	IdentityNumber              string `json:"identity_number"`
	IdentityExpdate             string `json:"identity_expdate"`
	IdentityCountryOfIssuedCode string `json:"identity_country_of_issued_code"`
	TicketNumberList            []any  `json:"ticket_number_list"`
	TicketURLList               []any  `json:"ticket_url_list"`
	RizText                     string `json:"riz_text"`
	PassengerDesc               string `json:"passenger_desc"`
}

// BookingRule is a rule entry under a provider booking
type BookingRule struct {
	Name        string   `json:"name"`
	Description []string `json:"description"`
}
