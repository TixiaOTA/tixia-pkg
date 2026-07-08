package panorama

// --------------------------------- SELL JOURNEYS REQUEST ---------------------------------

// SellJourneysRequest is the request structure for SellJourneys
type SellJourneysRequest = GetPriceItineraryRequest

// --------------------------------- SELL JOURNEYS RESPONSE ---------------------------------

// SellJourneysResponse is the response structure for SellJourneys
type SellJourneysResponse struct {
	Result struct {
		BasicPanoramaResult
		Response SellJourneyResponse `json:"response"`
	} `json:"result"`
}

// SellJourneyResponse is the response payload for SellJourneys
type SellJourneyResponse struct {
	TotalPrice          float64                  `json:"total_price"`
	TotalNta            float64                  `json:"total_nta"`
	ComplianceText      string                   `json:"compliance_text"`
	SellJourneyProvider []PriceItineraryProvider `json:"sell_journey_provider"`
}
