package panorama

import "errors"

// --------------------------------- GET PRICE ITINERARY REQUEST ---------------------------------

// GetPriceItineraryRequest is the request structure for GetPriceItinerary
type GetPriceItineraryRequest struct {
	Signature  string                   `json:"signature"`
	Adult      int                      `json:"adult"`
	Child      int                      `json:"child"`
	Infant     int                      `json:"infant"`
	Schedules  []PriceItinerarySchedule `json:"schedules"`
	PromoCodes []any                    `json:"promo_codes"`
}

// PriceItinerarySchedule is a schedule selection for price itinerary
type PriceItinerarySchedule struct {
	Journeys []PriceItineraryScheduleJourney `json:"journeys"`
	Provider string                          `json:"provider"`
}

// PriceItineraryScheduleJourney is a journey within a schedule selection
type PriceItineraryScheduleJourney struct {
	Segments []PriceItineraryScheduleSegment `json:"segments"`
}

// PriceItineraryScheduleSegment is a segment selection within a journey
type PriceItineraryScheduleSegment struct {
	SegmentCode string `json:"segment_code"`
	FareCode    string `json:"fare_code"`
}

func (r *GetPriceItineraryRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if r.Adult <= 0 {
		return errors.New("adult is required")
	}
	if len(r.Schedules) == 0 {
		return errors.New("schedules is required")
	}
	for _, schedule := range r.Schedules {
		if schedule.Provider == "" {
			return errors.New("provider is required")
		}
		if len(schedule.Journeys) == 0 {
			return errors.New("journeys is required")
		}
		for _, journey := range schedule.Journeys {
			if len(journey.Segments) == 0 {
				return errors.New("segments is required")
			}
			for _, segment := range journey.Segments {
				if segment.SegmentCode == "" {
					return errors.New("segment_code is required")
				}
				if segment.FareCode == "" {
					return errors.New("fare_code is required")
				}
			}
		}
	}
	return nil
}

// --------------------------------- GET PRICE ITINERARY RESPONSE ---------------------------------

// GetPriceItineraryResponse is the response structure for GetPriceItinerary
type GetPriceItineraryResponse struct {
	Result struct {
		BasicPanoramaResult
		Response PriceItineraryResponse `json:"response"`
	} `json:"result"`
}

// PriceItineraryResponse is the response payload for GetPriceItinerary
type PriceItineraryResponse struct {
	TotalPrice             float64                  `json:"total_price"`
	TotalNta               float64                  `json:"total_nta"`
	BalanceDue             float64                  `json:"balance_due"`
	TotalAdminFee          float64                  `json:"total_admin_fee"`
	PriceItineraryProvider []PriceItineraryProvider `json:"price_itinerary_provider"`
	ComplianceText         string                   `json:"compliance_text"`
}

// --------------------------------- GET PRICE ITINERARY RESPONSE TYPES ---------------------------------

// PriceItineraryProvider is a provider block in the GetPriceItinerary response.
type PriceItineraryProvider struct {
	Journeys        []PriceItineraryJourney   `json:"journeys"`
	PromoCodes      []any                     `json:"promo_codes"`
	Passengers      []PriceItineraryPassenger `json:"passengers"`
	Messages        []any                     `json:"messages"`
	HoldDate        string                    `json:"hold_date"`
	CcFeePercentage float64                   `json:"cc_fee_percentage"`
	Rules           []any                     `json:"rules"`
	Sequence        int                       `json:"sequence"`
	Provider        string                    `json:"provider"`
	ScheduleID      int                       `json:"schedule_id"`
}

// PriceItineraryPassenger is a nested type in the GetPriceItinerary response.
type PriceItineraryPassenger struct {
	Title            string                        `json:"title"`
	FirstName        string                        `json:"first_name"`
	LastName         string                        `json:"last_name"`
	PaxType          string                        `json:"pax_type"`
	Gender           string                        `json:"gender"`
	BirthDate        string                        `json:"birth_date"`
	Mobile           string                        `json:"mobile"`
	NationalityCode  string                        `json:"nationality_code"`
	TicketNumber     string                        `json:"ticket_number"`
	TicketURL        string                        `json:"ticket_url"`
	Fees             []any                         `json:"fees"`
	FfNumber         string                        `json:"ff_number"`
	FfCode           string                        `json:"ff_code"`
	RizText          string                        `json:"riz_text"`
	VendorData       PriceItineraryVendorData      `json:"vendor_data"`
	TicketNumberList []any                         `json:"ticket_number_list"`
	TicketURLList    []any                         `json:"ticket_url_list"`
	Fares            []PriceItineraryPassengerFare `json:"fares"`
}

// PriceItineraryPassengerFare is a nested type in the GetPriceItinerary response.
type PriceItineraryPassengerFare struct {
	CabinClass           string                               `json:"cabin_class"`
	ClassOfService       string                               `json:"class_of_service"`
	TourCode             string                               `json:"tour_code"`
	FareCode             string                               `json:"fare_code"`
	FareBasisCode        string                               `json:"fare_basis_code"`
	FareClass            string                               `json:"fare_class"`
	FareName             string                               `json:"fare_name"`
	FareRules            []any                                `json:"fare_rules"`
	RuleCode             string                               `json:"rule_code"`
	Description          []any                                `json:"description"`
	AvailableCount       int                                  `json:"available_count"`
	ServiceCharges       []PriceItineraryServiceCharge        `json:"service_charges"`
	FareDetails          []PriceItineraryFareDetail           `json:"fare_details"`
	GroupFareID          string                               `json:"group_fare_id"`
	FareRefID            string                               `json:"fare_ref_id"`
	PackageList          []any                                `json:"package_list"`
	EquipmentType        string                               `json:"equipment_type"`
	ValidatingCarrier    string                               `json:"validating_carrier"`
	IsTravelFair         bool                                 `json:"is_travel_fair"`
	TravelFairText       string                               `json:"travel_fair_text"`
	FareRemarks          []any                                `json:"fare_remarks"`
	Sequence             int                                  `json:"sequence"`
	ServiceChargeSummary []PriceItineraryServiceChargeSummary `json:"service_charge_summary"`
	PricingData          PriceItineraryPassengerPricingData   `json:"pricing_data"`
}

// PriceItineraryPassengerPricingData is a nested type in the GetPriceItinerary response.
type PriceItineraryPassengerPricingData struct {
	RuleData            PriceItineraryRuleData            `json:"rule_data"`
	ProviderPricingData PriceItineraryProviderPricingData `json:"provider_pricing_data"`
	AgentPricingData    PriceItineraryAgentPricingData    `json:"agent_pricing_data"`
	CustomerPricingData PriceItineraryCustomerPricingData `json:"customer_pricing_data"`
	AgentCommissionData PriceItineraryAgentCommissionData `json:"agent_commission_data"`
	CalcTextList        []string                          `json:"calc_text_list"`
}

// PriceItineraryVendorData is a nested type in the GetPriceItinerary response.
type PriceItineraryVendorData struct {
}

// PriceItineraryJourney is a nested type in the GetPriceItinerary response.
type PriceItineraryJourney struct {
	Origin                   string                  `json:"origin"`
	OriginTerminal           string                  `json:"origin_terminal"`
	OriginUtc                int                     `json:"origin_utc"`
	OriginCountry            string                  `json:"origin_country"`
	DepartureDate            string                  `json:"departure_date"`
	Destination              string                  `json:"destination"`
	DestinationTerminal      string                  `json:"destination_terminal"`
	DestinationUtc           int                     `json:"destination_utc"`
	DestinationCountry       string                  `json:"destination_country"`
	ArrivalDate              string                  `json:"arrival_date"`
	JourneyCode              string                  `json:"journey_code"`
	RouteCode                string                  `json:"route_code"`
	ElapsedTime              string                  `json:"elapsed_time"`
	TransitCount             int                     `json:"transit_count"`
	Provider                 string                  `json:"provider"`
	CarrierCodeList          []string                `json:"carrier_code_list"`
	OperatingAirlineCodeList []any                   `json:"operating_airline_code_list"`
	Segments                 []PriceItinerarySegment `json:"segments"`
	JourneyRefID             string                  `json:"journey_ref_id"`
	IsVtlFlight              bool                    `json:"is_vtl_flight"`
	Sequence                 int                     `json:"sequence"`
	SearchBanner             []any                   `json:"search_banner"`
}

// PriceItinerarySegment is a nested type in the GetPriceItinerary response.
type PriceItinerarySegment struct {
	Origin               string               `json:"origin"`
	OriginTerminal       string               `json:"origin_terminal"`
	OriginUtc            int                  `json:"origin_utc"`
	OriginCountry        string               `json:"origin_country"`
	DepartureDate        string               `json:"departure_date"`
	Destination          string               `json:"destination"`
	DestinationTerminal  string               `json:"destination_terminal"`
	DestinationUtc       int                  `json:"destination_utc"`
	DestinationCountry   string               `json:"destination_country"`
	ArrivalDate          string               `json:"arrival_date"`
	ElapsedTime          string               `json:"elapsed_time"`
	TransitCount         int                  `json:"transit_count"`
	SegmentCode          string               `json:"segment_code"`
	SegmentCode2         string               `json:"segment_code2"`
	RouteCode            string               `json:"route_code"`
	TransitDuration      string               `json:"transit_duration"`
	CarrierName          string               `json:"carrier_name"`
	CarrierCode          string               `json:"carrier_code"`
	CarrierTypeCode      string               `json:"carrier_type_code"`
	CarrierTypeName      string               `json:"carrier_type_name"`
	CarrierNumber        string               `json:"carrier_number"`
	Provider             string               `json:"provider"`
	OperatingAirlineCode string               `json:"operating_airline_code"`
	Legs                 []PriceItineraryLeg  `json:"legs"`
	Fares                []PriceItineraryFare `json:"fares"`
	EquipmentType        string               `json:"equipment_type"`
	Status               string               `json:"status"`
	SeatCabins           []any                `json:"seat_cabins"`
	Sequence             int                  `json:"sequence"`
}

// PriceItineraryFare is a nested type in the GetPriceItinerary response.
type PriceItineraryFare struct {
	CabinClass           string                               `json:"cabin_class"`
	ClassOfService       string                               `json:"class_of_service"`
	TourCode             string                               `json:"tour_code"`
	FareCode             string                               `json:"fare_code"`
	FareBasisCode        string                               `json:"fare_basis_code"`
	FareClass            string                               `json:"fare_class"`
	FareName             string                               `json:"fare_name"`
	FareRules            []any                                `json:"fare_rules"`
	RuleCode             string                               `json:"rule_code"`
	Description          []any                                `json:"description"`
	AvailableCount       int                                  `json:"available_count"`
	FareDetails          []PriceItineraryFareDetail           `json:"fare_details"`
	GroupFareID          string                               `json:"group_fare_id"`
	FareRefID            string                               `json:"fare_ref_id"`
	PackageList          []any                                `json:"package_list"`
	EquipmentType        string                               `json:"equipment_type"`
	ValidatingCarrier    string                               `json:"validating_carrier"`
	IsTravelFair         bool                                 `json:"is_travel_fair"`
	TravelFairText       string                               `json:"travel_fair_text"`
	FareRemarks          []any                                `json:"fare_remarks"`
	Sequence             int                                  `json:"sequence"`
	ServiceChargeSummary []PriceItineraryServiceChargeSummary `json:"service_charge_summary"`
	PricingDetails       []PriceItineraryPricingDetail        `json:"pricing_details"`
}

// PriceItineraryPricingDetail is a nested type in the GetPriceItinerary response.
type PriceItineraryPricingDetail struct {
	RuleData            PriceItineraryRuleData            `json:"rule_data"`
	ProviderPricingData PriceItineraryProviderPricingData `json:"provider_pricing_data"`
	AgentPricingData    PriceItineraryAgentPricingData    `json:"agent_pricing_data"`
	CustomerPricingData PriceItineraryCustomerPricingData `json:"customer_pricing_data"`
	AgentCommissionData PriceItineraryAgentCommissionData `json:"agent_commission_data"`
	CalcTextList        []string                          `json:"calc_text_list"`
	PaxRefs             []string                          `json:"pax_refs"`
	IsReservation       bool                              `json:"is_reservation"`
}

// PriceItineraryAgentCommissionData is a nested type in the GetPriceItinerary response.
type PriceItineraryAgentCommissionData struct {
	ID                int    `json:"id"`
	Sequence          int    `json:"sequence"`
	Name              string `json:"name"`
	CurrencyCode      string `json:"currency_code"`
	SetExpirationDate bool   `json:"set_expiration_date"`
	DateFrom          string `json:"date_from"`
	DateTo            string `json:"date_to"`
	Route             struct {
		Origin struct {
			AccessType          string `json:"access_type"`
			DestinationCodeList []any  `json:"destination_code_list"`
			CityCodeList        []any  `json:"city_code_list"`
			CountryCodeList     []any  `json:"country_code_list"`
		} `json:"origin"`
		Destination struct {
			AccessType          string `json:"access_type"`
			DestinationCodeList []any  `json:"destination_code_list"`
			CityCodeList        []any  `json:"city_code_list"`
			CountryCodeList     []any  `json:"country_code_list"`
		} `json:"destination"`
		Carrier struct {
			AccessType  string `json:"access_type"`
			CarrierList []any  `json:"carrier_list"`
		} `json:"carrier"`
		Operated struct {
			AccessType   string `json:"access_type"`
			OperatedList []any  `json:"operated_list"`
		} `json:"operated"`
		ClassOfService struct {
			AccessType         string `json:"access_type"`
			ClassOfServiceList []any  `json:"class_of_service_list"`
		} `json:"class_of_service"`
		ChargeCode struct {
			AccessType     string `json:"access_type"`
			ChargeCodeList []any  `json:"charge_code_list"`
		} `json:"charge_code"`
		TourCode struct {
			AccessType   string `json:"access_type"`
			TourCodeList []any  `json:"tour_code_list"`
		} `json:"tour_code"`
		FareRemark struct {
			AccessType     string `json:"access_type"`
			FareRemarkList []any  `json:"fare_remark_list"`
		} `json:"fare_remark"`
		DateOfTravel struct {
			AccessType string `json:"access_type"`
			StartDate  string `json:"start_date"`
			EndDate    string `json:"end_date"`
		} `json:"date_of_travel"`
		Total struct {
			AccessType     string  `json:"access_type"`
			IsLessEqual    bool    `json:"is_less_equal"`
			LessAmount     float64 `json:"less_amount"`
			IsGreaterEqual bool    `json:"is_greater_equal"`
			GreaterAmount  float64 `json:"greater_amount"`
		} `json:"total"`
	} `json:"route"`
	Commission struct {
		Agent struct {
			CommissionByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"commission_by_percentage"`
			CommissionByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsPax     bool    `json:"is_pax"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"commission_by_amount"`
			DiscountByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"discount_by_percentage"`
			DiscountByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsPax     bool    `json:"is_pax"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"discount_by_amount"`
		} `json:"agent"`
		Parent struct {
			ChargeByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"charge_by_percentage"`
			ChargeByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsPax     bool    `json:"is_pax"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"charge_by_amount"`
		} `json:"parent"`
		Ho struct {
			ChargeByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"charge_by_percentage"`
			ChargeByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsPax     bool    `json:"is_pax"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"charge_by_amount"`
		} `json:"ho"`
		ResidualAmountTo string `json:"residual_amount_to"`
		UplineList       []struct {
			ID                     int    `json:"id"`
			Sequence               int    `json:"sequence"`
			AgentTypeCode          string `json:"agent_type_code"`
			CommissionByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"commission_by_percentage"`
			CommissionByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsPax     bool    `json:"is_pax"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"commission_by_amount"`
			State string `json:"state"`
		} `json:"upline_list"`
	} `json:"commission"`
	PricingBreakdown bool   `json:"pricing_breakdown"`
	State            string `json:"state"`
	PricingID        int    `json:"pricing_id"`
	ParentData       struct {
		ID        int    `json:"id"`
		Sequence  int    `json:"sequence"`
		Name      string `json:"name"`
		AgentType struct {
			AccessType        string `json:"access_type"`
			AgentTypeCodeList []any  `json:"agent_type_code_list"`
		} `json:"agent_type"`
		Provider struct {
			AccessType       string `json:"access_type"`
			ProviderCodeList []any  `json:"provider_code_list"`
		} `json:"provider"`
		Carrier struct {
			AccessType      string `json:"access_type"`
			CarrierCodeList []any  `json:"carrier_code_list"`
		} `json:"carrier"`
		ProviderType struct {
			AccessType           string `json:"access_type"`
			ProviderTypeCodeList []any  `json:"provider_type_code_list"`
		} `json:"provider_type"`
		Agent struct {
			AccessType  string `json:"access_type"`
			AgentIDList []any  `json:"agent_id_list"`
		} `json:"agent"`
		State string `json:"state"`
	} `json:"parent_data"`
	RuleID       int `json:"rule_id"`
	PricingIndex int `json:"pricing_index"`
	RuleIndex    int `json:"rule_index"`
}

// PriceItineraryCustomerPricingData is a nested type in the GetPriceItinerary response.
type PriceItineraryCustomerPricingData struct {
}

// PriceItineraryAgentPricingData is a nested type in the GetPriceItinerary response.
type PriceItineraryAgentPricingData struct {
	ID                   int    `json:"id"`
	Sequence             int    `json:"sequence"`
	Name                 string `json:"name"`
	CurrencyCode         string `json:"currency_code"`
	SetExpirationDate    bool   `json:"set_expiration_date"`
	DateFrom             string `json:"date_from"`
	DateTo               string `json:"date_to"`
	IsOrderNumberPricing bool   `json:"is_order_number_pricing"`
	Route                struct {
		Origin struct {
			AccessType          string `json:"access_type"`
			DestinationCodeList []any  `json:"destination_code_list"`
			CityCodeList        []any  `json:"city_code_list"`
			CountryCodeList     []any  `json:"country_code_list"`
		} `json:"origin"`
		Destination struct {
			AccessType          string `json:"access_type"`
			DestinationCodeList []any  `json:"destination_code_list"`
			CityCodeList        []any  `json:"city_code_list"`
			CountryCodeList     []any  `json:"country_code_list"`
		} `json:"destination"`
		Carrier struct {
			AccessType  string `json:"access_type"`
			CarrierList []any  `json:"carrier_list"`
		} `json:"carrier"`
		Operated struct {
			AccessType   string `json:"access_type"`
			OperatedList []any  `json:"operated_list"`
		} `json:"operated"`
		ClassOfService struct {
			AccessType         string `json:"access_type"`
			ClassOfServiceList []any  `json:"class_of_service_list"`
		} `json:"class_of_service"`
		ChargeCode struct {
			AccessType     string `json:"access_type"`
			ChargeCodeList []any  `json:"charge_code_list"`
		} `json:"charge_code"`
		TourCode struct {
			AccessType   string `json:"access_type"`
			TourCodeList []any  `json:"tour_code_list"`
		} `json:"tour_code"`
		FareRemark struct {
			AccessType     string `json:"access_type"`
			FareRemarkList []any  `json:"fare_remark_list"`
		} `json:"fare_remark"`
		DateOfTravel struct {
			AccessType string `json:"access_type"`
			StartDate  string `json:"start_date"`
			EndDate    string `json:"end_date"`
		} `json:"date_of_travel"`
		Total struct {
			AccessType     string  `json:"access_type"`
			IsLessEqual    bool    `json:"is_less_equal"`
			LessAmount     float64 `json:"less_amount"`
			IsGreaterEqual bool    `json:"is_greater_equal"`
			GreaterAmount  float64 `json:"greater_amount"`
		} `json:"total"`
	} `json:"route"`
	Commission struct {
		Agent struct {
			CommissionByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"commission_by_percentage"`
			CommissionByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsPax     bool    `json:"is_pax"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"commission_by_amount"`
		} `json:"agent"`
		Parent struct {
			ChargeByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"charge_by_percentage"`
			ChargeByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsPax     bool    `json:"is_pax"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"charge_by_amount"`
		} `json:"parent"`
		Ho struct {
			ChargeByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"charge_by_percentage"`
			ChargeByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsPax     bool    `json:"is_pax"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"charge_by_amount"`
		} `json:"ho"`
		ResidualAmountTo string `json:"residual_amount_to"`
		UplineList       []any  `json:"upline_list"`
	} `json:"commission"`
	Ticketing struct {
		Sales struct {
			Fare struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"fare"`
			Tax struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"tax"`
			Total struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"total"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"upsell_by_percentage"`
			UpsellByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"upsell_by_amount"`
		} `json:"sales"`
		NtaAgent struct {
			Fare struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"fare"`
			Tax struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"tax"`
			Total struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"total"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"upsell_by_percentage"`
			UpsellByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"upsell_by_amount"`
		} `json:"nta_agent"`
		Commission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"commission"`
		HoCommission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"ho_commission"`
	} `json:"ticketing"`
	Ancillary struct {
		Sales struct {
			Fare struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"fare"`
			Tax struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"tax"`
			Total struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"total"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"upsell_by_percentage"`
			UpsellByAmount struct {
				Amount float64 `json:"amount"`
			} `json:"upsell_by_amount"`
		} `json:"sales"`
		NtaAgent struct {
			Fare struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"fare"`
			Tax struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"tax"`
			Total struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"total"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"upsell_by_percentage"`
			UpsellByAmount struct {
				Amount float64 `json:"amount"`
			} `json:"upsell_by_amount"`
		} `json:"nta_agent"`
		Commission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"commission"`
		HoCommission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"ho_commission"`
	} `json:"ancillary"`
	Reservation struct {
		Sales struct {
			UpsellByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
			} `json:"upsell_by_amount"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"upsell_by_percentage"`
		} `json:"sales"`
		NtaAgent struct {
			UpsellByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
			} `json:"upsell_by_amount"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"upsell_by_percentage"`
		} `json:"nta_agent"`
		Commission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"commission"`
		HoCommission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"ho_commission"`
	} `json:"reservation"`
	PricingBreakdown bool   `json:"pricing_breakdown"`
	State            string `json:"state"`
	PricingType      string `json:"pricing_type"`
	PricingID        int    `json:"pricing_id"`
	ParentData       struct {
		ID        int    `json:"id"`
		Sequence  int    `json:"sequence"`
		Name      string `json:"name"`
		AgentType struct {
			AccessType        string `json:"access_type"`
			AgentTypeCodeList []any  `json:"agent_type_code_list"`
		} `json:"agent_type"`
		Provider struct {
			AccessType       string `json:"access_type"`
			ProviderCodeList []any  `json:"provider_code_list"`
		} `json:"provider"`
		Carrier struct {
			AccessType      string `json:"access_type"`
			CarrierCodeList []any  `json:"carrier_code_list"`
		} `json:"carrier"`
		ProviderType struct {
			AccessType           string `json:"access_type"`
			ProviderTypeCodeList []any  `json:"provider_type_code_list"`
		} `json:"provider_type"`
		Agent struct {
			AccessType  string `json:"access_type"`
			AgentIDList []any  `json:"agent_id_list"`
		} `json:"agent"`
		State string `json:"state"`
	} `json:"parent_data"`
	RuleID       int `json:"rule_id"`
	PricingIndex int `json:"pricing_index"`
	RuleIndex    int `json:"rule_index"`
}

// PriceItineraryProviderPricingData is a nested type in the GetPriceItinerary response.
type PriceItineraryProviderPricingData struct {
	ID                   int    `json:"id"`
	Sequence             int    `json:"sequence"`
	Name                 string `json:"name"`
	PricingType          string `json:"pricing_type"`
	CurrencyCode         string `json:"currency_code"`
	SetExpirationDate    bool   `json:"set_expiration_date"`
	DateFrom             string `json:"date_from"`
	DateTo               string `json:"date_to"`
	IsOrderNumberPricing bool   `json:"is_order_number_pricing"`
	Route                struct {
		Origin struct {
			AccessType          string `json:"access_type"`
			DestinationCodeList []any  `json:"destination_code_list"`
			CityCodeList        []any  `json:"city_code_list"`
			CountryCodeList     []any  `json:"country_code_list"`
		} `json:"origin"`
		Destination struct {
			AccessType          string `json:"access_type"`
			DestinationCodeList []any  `json:"destination_code_list"`
			CityCodeList        []any  `json:"city_code_list"`
			CountryCodeList     []any  `json:"country_code_list"`
		} `json:"destination"`
		Carrier struct {
			AccessType  string `json:"access_type"`
			CarrierList []any  `json:"carrier_list"`
		} `json:"carrier"`
		Operated struct {
			AccessType   string `json:"access_type"`
			OperatedList []any  `json:"operated_list"`
		} `json:"operated"`
		ClassOfService struct {
			AccessType         string `json:"access_type"`
			ClassOfServiceList []any  `json:"class_of_service_list"`
		} `json:"class_of_service"`
		ChargeCode struct {
			AccessType     string `json:"access_type"`
			ChargeCodeList []any  `json:"charge_code_list"`
		} `json:"charge_code"`
		TourCode struct {
			AccessType   string `json:"access_type"`
			TourCodeList []any  `json:"tour_code_list"`
		} `json:"tour_code"`
		FareRemark struct {
			AccessType     string `json:"access_type"`
			FareRemarkList []any  `json:"fare_remark_list"`
		} `json:"fare_remark"`
		DateOfTravel struct {
			AccessType string `json:"access_type"`
			StartDate  string `json:"start_date"`
			EndDate    string `json:"end_date"`
		} `json:"date_of_travel"`
		Total struct {
			AccessType     string  `json:"access_type"`
			IsLessEqual    bool    `json:"is_less_equal"`
			LessAmount     float64 `json:"less_amount"`
			IsGreaterEqual bool    `json:"is_greater_equal"`
			GreaterAmount  float64 `json:"greater_amount"`
		} `json:"total"`
	} `json:"route"`
	Less struct {
		Percentage float64 `json:"percentage"`
		IsInfant   bool    `json:"is_infant"`
		TourCode   string  `json:"tour_code"`
	} `json:"less"`
	Ticketing struct {
		Nta struct {
			Fare struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"fare"`
			Tax struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"tax"`
			Total struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"total"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"upsell_by_percentage"`
			UpsellByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"upsell_by_amount"`
		} `json:"nta"`
		NtaAgent struct {
			Fare struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"fare"`
			Tax struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"tax"`
			Total struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"total"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"upsell_by_percentage"`
			UpsellByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"upsell_by_amount"`
		} `json:"nta_agent"`
		Sales struct {
			Fare struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"fare"`
			Tax struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"tax"`
			Total struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"total"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
				IsInfant   bool    `json:"is_infant"`
			} `json:"upsell_by_percentage"`
			UpsellByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
				IsInfant  bool    `json:"is_infant"`
			} `json:"upsell_by_amount"`
		} `json:"sales"`
		Commission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"commission"`
		HoCommission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"ho_commission"`
	} `json:"ticketing"`
	Ancillary struct {
		Nta struct {
			Fare struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"fare"`
			Tax struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"tax"`
			Total struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"total"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"upsell_by_percentage"`
			UpsellByAmount struct {
				Amount float64 `json:"amount"`
			} `json:"upsell_by_amount"`
		} `json:"nta"`
		NtaAgent struct {
			Fare struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"fare"`
			Tax struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"tax"`
			Total struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"total"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"upsell_by_percentage"`
			UpsellByAmount struct {
				Amount float64 `json:"amount"`
			} `json:"upsell_by_amount"`
		} `json:"nta_agent"`
		Sales struct {
			Fare struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"fare"`
			Tax struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"tax"`
			Total struct {
				Percentage float64 `json:"percentage"`
				Amount     float64 `json:"amount"`
			} `json:"total"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"upsell_by_percentage"`
			UpsellByAmount struct {
				Amount float64 `json:"amount"`
			} `json:"upsell_by_amount"`
		} `json:"sales"`
		Commission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"commission"`
		HoCommission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"ho_commission"`
	} `json:"ancillary"`
	Reservation struct {
		Nta struct {
			UpsellByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
			} `json:"upsell_by_amount"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"upsell_by_percentage"`
		} `json:"nta"`
		NtaAgent struct {
			UpsellByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
			} `json:"upsell_by_amount"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"upsell_by_percentage"`
		} `json:"nta_agent"`
		Sales struct {
			UpsellByAmount struct {
				Amount    float64 `json:"amount"`
				IsRoute   bool    `json:"is_route"`
				IsSegment bool    `json:"is_segment"`
			} `json:"upsell_by_amount"`
			UpsellByPercentage struct {
				Percentage float64 `json:"percentage"`
				Minimum    float64 `json:"minimum"`
				HasMinimum bool    `json:"has_minimum"`
				Maximum    float64 `json:"maximum"`
				HasMaximum bool    `json:"has_maximum"`
			} `json:"upsell_by_percentage"`
		} `json:"sales"`
		Commission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"commission"`
		HoCommission struct {
			TaxType       string  `json:"tax_type"`
			TaxAmount     float64 `json:"tax_amount"`
			TaxPercentage float64 `json:"tax_percentage"`
			Rounding      int     `json:"rounding"`
		} `json:"ho_commission"`
	} `json:"reservation"`
	PricingBreakdown bool   `json:"pricing_breakdown"`
	State            string `json:"state"`
	PricingID        int    `json:"pricing_id"`
	ParentData       struct {
		ID               int    `json:"id"`
		Sequence         int    `json:"sequence"`
		Name             string `json:"name"`
		ProviderTypeCode string `json:"provider_type_code"`
		Provider         struct {
			AccessType       string `json:"access_type"`
			ProviderCodeList []any  `json:"provider_code_list"`
		} `json:"provider"`
		Carrier struct {
			AccessType      string `json:"access_type"`
			CarrierCodeList []any  `json:"carrier_code_list"`
		} `json:"carrier"`
		Agent struct {
			AccessType  string `json:"access_type"`
			AgentIDList []any  `json:"agent_id_list"`
		} `json:"agent"`
		AgentType struct {
			AccessType        string `json:"access_type"`
			AgentTypeCodeList []any  `json:"agent_type_code_list"`
		} `json:"agent_type"`
		State string `json:"state"`
	} `json:"parent_data"`
	RuleID       int `json:"rule_id"`
	PricingIndex int `json:"pricing_index"`
	RuleIndex    int `json:"rule_index"`
}

// PriceItineraryRuleData is a nested type in the GetPriceItinerary response.
type PriceItineraryRuleData struct {
	AgentID                  int      `json:"agent_id"`
	AgentTypeCode            string   `json:"agent_type_code"`
	ProviderTypeCode         string   `json:"provider_type_code"`
	CustomerParentTypeCode   string   `json:"customer_parent_type_code"`
	CustomerParentID         string   `json:"customer_parent_id"`
	ProviderCode             string   `json:"provider_code"`
	CarrierCode              string   `json:"carrier_code"`
	OriginCode               string   `json:"origin_code"`
	OriginCity               string   `json:"origin_city"`
	OriginCountry            string   `json:"origin_country"`
	DestinationCode          string   `json:"destination_code"`
	DestinationCity          string   `json:"destination_city"`
	DestinationCountry       string   `json:"destination_country"`
	ClassOfServiceList       []string `json:"class_of_service_list"`
	ChargeCodeList           []string `json:"charge_code_list"`
	TourCodeList             []any    `json:"tour_code_list"`
	PricingDatetime          string   `json:"pricing_datetime"`
	DepartureDateList        []string `json:"departure_date_list"`
	CarrierCodeList          []string `json:"carrier_code_list"`
	OperatingAirlineCodeList []any    `json:"operating_airline_code_list"`
	FareRemarkList           []any    `json:"fare_remark_list"`
	CurrencyCodeList         []string `json:"currency_code_list"`
	UpsellByAmountCharge     bool     `json:"upsell_by_amount_charge"`
	Context                  struct {
		UID                   int    `json:"uid"`
		UserName              string `json:"user_name"`
		UserLogin             string `json:"user_login"`
		AgentSeqID            string `json:"agent_seq_id"`
		HoID                  int    `json:"ho_id"`
		AgentID               int    `json:"agent_id"`
		AgentName             string `json:"agent_name"`
		AgentTypeID           int    `json:"agent_type_id"`
		AgentTypeName         string `json:"agent_type_name"`
		AgentTypeCode         string `json:"agent_type_code"`
		AgentFrontendSecurity []any  `json:"agent_frontend_security"`
		HoSeqID               string `json:"ho_seq_id"`
		HoName                string `json:"ho_name"`
		HoCurrency            string `json:"ho_currency"`
		BreakdownPricing      struct {
			All string `json:"all"`
		} `json:"breakdown_pricing"`
		ShowFareOnly bool   `json:"show_fare_only"`
		APIRole      string `json:"api_role"`
		DeviceType   string `json:"device_type"`
		HostIps      []any  `json:"host_ips"`
		Configs      struct {
			Airline struct {
				ProviderAccess string `json:"provider_access"`
				Providers      struct {
				} `json:"providers"`
				Blackouts []any `json:"blackouts"`
			} `json:"airline"`
			Hotel struct {
				ProviderAccess string `json:"provider_access"`
				Providers      struct {
				} `json:"providers"`
				Blackouts []any `json:"blackouts"`
			} `json:"hotel"`
			Train struct {
				ProviderAccess string `json:"provider_access"`
				Providers      struct {
				} `json:"providers"`
				Blackouts []any `json:"blackouts"`
			} `json:"train"`
		} `json:"configs"`
		CoUID                   int    `json:"co_uid"`
		CoUserName              string `json:"co_user_name"`
		CoUserLogin             string `json:"co_user_login"`
		CoAgentSeqID            string `json:"co_agent_seq_id"`
		CoHoID                  int    `json:"co_ho_id"`
		CoAgentID               int    `json:"co_agent_id"`
		CoAgentName             string `json:"co_agent_name"`
		CoAgentTypeID           int    `json:"co_agent_type_id"`
		CoAgentTypeName         string `json:"co_agent_type_name"`
		CoAgentTypeCode         string `json:"co_agent_type_code"`
		CoAgentFrontendSecurity []any  `json:"co_agent_frontend_security"`
		PricingBreakdown        bool   `json:"pricing_breakdown"`
		CoHoBaseCountry         string `json:"co_ho_base_country"`
		CoHoSeqID               string `json:"co_ho_seq_id"`
		CoHoName                string `json:"co_ho_name"`
		CoHoCurrency            string `json:"co_ho_currency"`
		CoBreakdownPricing      struct {
			All string `json:"all"`
		} `json:"co_breakdown_pricing"`
		CoShowFareOnly bool   `json:"co_show_fare_only"`
		Sid            string `json:"sid"`
		Signature      string `json:"signature"`
		CoUserInfo     []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			AgentTypeID struct {
				ID                 int    `json:"id"`
				Name               string `json:"name"`
				Code               string `json:"code"`
				RoundingAmountType string `json:"rounding_amount_type"`
				RoundingPlaces     int    `json:"rounding_places"`
			} `json:"agent_type_id"`
			AgentLevel int `json:"agent_level"`
		} `json:"co_user_info"`
		APIKey      string `json:"api_key"`
		ExpiredDate string `json:"expired_date"`
	} `json:"context"`
	TotalAmount float64 `json:"total_amount"`
}

// PriceItineraryServiceChargeSummary is a nested type in the GetPriceItinerary response.
type PriceItineraryServiceChargeSummary struct {
	ServiceCharges            []PriceItineraryServiceCharge `json:"service_charges"`
	PaxType                   string                        `json:"pax_type"`
	PaxCount                  int                           `json:"pax_count"`
	BaseFare                  float64                       `json:"base_fare"`
	BaseTax                   float64                       `json:"base_tax"`
	BaseAncillary             float64                       `json:"base_ancillary"`
	BaseSeat                  float64                       `json:"base_seat"`
	BasePriceOtt              float64                       `json:"base_price_ott"`
	BaseUpsell                float64                       `json:"base_upsell"`
	BaseUpsellVendor          float64                       `json:"base_upsell_vendor"`
	BaseUpsellAgent           float64                       `json:"base_upsell_agent"`
	BaseDiscount              float64                       `json:"base_discount"`
	BasePrice                 float64                       `json:"base_price"`
	BaseCommission            float64                       `json:"base_commission"`
	BaseCommissionUpline      float64                       `json:"base_commission_upline"`
	BaseCommissionCharge      float64                       `json:"base_commission_charge"`
	BaseCommissionVendor      float64                       `json:"base_commission_vendor"`
	BaseFeeHo                 float64                       `json:"base_fee_ho"`
	BaseVatHo                 float64                       `json:"base_vat_ho"`
	BaseFeeCustomer           float64                       `json:"base_fee_customer"`
	BaseVatCustomer           float64                       `json:"base_vat_customer"`
	BaseNta                   float64                       `json:"base_nta"`
	BaseNtaVendor             float64                       `json:"base_nta_vendor"`
	BaseHiddenFeeHo           float64                       `json:"base_hidden_fee_ho"`
	BaseHiddenVatHo           float64                       `json:"base_hidden_vat_ho"`
	BaseHiddenVendorVat       float64                       `json:"base_hidden_vendor_vat"`
	BaseHiddenCommissionHo    float64                       `json:"base_hidden_commission_ho"`
	BaseNoHiddenCommissionHo  float64                       `json:"base_no_hidden_commission_ho"`
	BaseCommissionHo          float64                       `json:"base_commission_ho"`
	TotalFare                 float64                       `json:"total_fare"`
	TotalTax                  float64                       `json:"total_tax"`
	TotalAncillary            float64                       `json:"total_ancillary"`
	TotalSeat                 float64                       `json:"total_seat"`
	TotalPriceOtt             float64                       `json:"total_price_ott"`
	TotalUpsell               float64                       `json:"total_upsell"`
	TotalUpsellVendor         float64                       `json:"total_upsell_vendor"`
	TotalUpsellAgent          float64                       `json:"total_upsell_agent"`
	TotalDiscount             float64                       `json:"total_discount"`
	TotalPrice                float64                       `json:"total_price"`
	TotalCommission           float64                       `json:"total_commission"`
	TotalCommissionUpline     float64                       `json:"total_commission_upline"`
	TotalCommissionCharge     float64                       `json:"total_commission_charge"`
	TotalCommissionVendor     float64                       `json:"total_commission_vendor"`
	TotalFeeHo                float64                       `json:"total_fee_ho"`
	TotalVatHo                float64                       `json:"total_vat_ho"`
	TotalFeeCustomer          float64                       `json:"total_fee_customer"`
	TotalVatCustomer          float64                       `json:"total_vat_customer"`
	TotalNta                  float64                       `json:"total_nta"`
	TotalNtaVendor            float64                       `json:"total_nta_vendor"`
	TotalHiddenFeeHo          float64                       `json:"total_hidden_fee_ho"`
	TotalHiddenVatHo          float64                       `json:"total_hidden_vat_ho"`
	TotalHiddenVendorVat      float64                       `json:"total_hidden_vendor_vat"`
	TotalHiddenCommissionHo   float64                       `json:"total_hidden_commission_ho"`
	TotalNoHiddenCommissionHo float64                       `json:"total_no_hidden_commission_ho"`
	TotalCommissionHo         float64                       `json:"total_commission_ho"`
}

// PriceItineraryServiceCharge is a nested type in the GetPriceItinerary response.
type PriceItineraryServiceCharge struct {
	ChargeCode        string  `json:"charge_code"`
	ChargeType        string  `json:"charge_type"`
	Currency          string  `json:"currency"`
	PaxType           string  `json:"pax_type"`
	PaxCount          int     `json:"pax_count"`
	Amount            float64 `json:"amount"`
	ForeignCurrency   string  `json:"foreign_currency"`
	ForeignAmount     float64 `json:"foreign_amount"`
	Total             float64 `json:"total"`
	Sequence          int     `json:"sequence"`
	CommissionAgentID any     `json:"commission_agent_id"`
	IsNew             bool    `json:"is_new,omitempty"`
	IsOriginal        bool    `json:"is_original,omitempty"`
}

// PriceItineraryFareDetail is a nested type in the GetPriceItinerary response.
type PriceItineraryFareDetail struct {
	DetailCode  string  `json:"detail_code"`
	DetailType  string  `json:"detail_type"`
	Amount      float64 `json:"amount"`
	Unit        string  `json:"unit"`
	DetailName  string  `json:"detail_name"`
	Description []any   `json:"description"`
	Sequence    int     `json:"sequence"`
}

// PriceItineraryLeg is a nested type in the GetPriceItinerary response.
type PriceItineraryLeg struct {
	Origin               string `json:"origin"`
	OriginTerminal       string `json:"origin_terminal"`
	OriginUtc            int    `json:"origin_utc"`
	OriginCountry        string `json:"origin_country"`
	DepartureDate        string `json:"departure_date"`
	Destination          string `json:"destination"`
	DestinationTerminal  string `json:"destination_terminal"`
	DestinationUtc       int    `json:"destination_utc"`
	DestinationCountry   string `json:"destination_country"`
	ArrivalDate          string `json:"arrival_date"`
	ElapsedTime          string `json:"elapsed_time"`
	TransitDuration      string `json:"transit_duration"`
	CarrierName          string `json:"carrier_name"`
	CarrierCode          string `json:"carrier_code"`
	CarrierTypeCode      string `json:"carrier_type_code"`
	CarrierTypeName      string `json:"carrier_type_name"`
	CarrierNumber        string `json:"carrier_number"`
	LegCode              FlexibleString `json:"leg_code"`
	RouteCode            string         `json:"route_code"`
	Provider             string         `json:"provider"`
	OperatingAirlineCode string         `json:"operating_airline_code"`
	LegCode2             FlexibleString `json:"leg_code2"`
	EquipmentType        string `json:"equipment_type"`
	Sequence             int    `json:"sequence"`
}
