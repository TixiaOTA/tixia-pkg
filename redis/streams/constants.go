package streams

// Flight events
const (
	TopicFlightSearchRequested                   string = "flight.search.requested"
	TopicFlightSearchResults                     string = "flight.search.results"
	TopicFlightFareRequested                     string = "flight.fare.requested"
	TopicFlightFareResults                       string = "flight.fare.results"
	TopicFlightBookingRequested                  string = "flight.book.requested"
	TopicFlightBookingResults                    string = "flight.book.results"
	TopicFlightRetrieveBookingRequested          string = "flight.retrievebooking.requested"
	TopicFlightRetrieveBookingResults            string = "flight.retrievebooking.results"
	TopicFlightCancelTicketRequested             string = "flight.cancelticket.requested"
	TopicFlightCancelTicketResults               string = "flight.cancelticket.results"
	TopicFlightIssueTicketRequested              string = "flight.issueticket.requested"
	TopicFlightIssueTicketResults                string = "flight.issueticket.results"
	TopicFlightCheckBalanceRequested             string = "flight.checkbalance.requested"
	TopicFlightCheckBalanceResults               string = "flight.checkbalance.results"
	TopicFlightOrbisSigninRequested              string = "flight.orbis.signin.requested"
	TopicFlightOrbisSigninResults                string = "flight.orbis.signin.results"
	TopicFlightOrbisSearchRequested              string = "flight.orbis.search.requested"
	TopicFlightOrbisSearchResults                string = "flight.orbis.search.results"
	TopicFlightOrbisRulesRequested               string = "flight.orbis.rules.requested"
	TopicFlightOrbisRulesResults                 string = "flight.orbis.rules.results"
	TopicFlightOrbisFareRequested                string = "flight.orbis.fare.requested"
	TopicFlightOrbisFareResults                  string = "flight.orbis.fare.results"
	TopicFlightOrbisSelectJourneyRequested       string = "flight.orbis.selectjourney.requested"
	TopicFlightOrbisSelectJourneyResults         string = "flight.orbis.selectjourney.results"
	TopicFlightOrbisGetFFRequested               string = "flight.orbis.getff.requested"
	TopicFlightOrbisGetFFResults                 string = "flight.orbis.getff.results"
	TopicFlightOrbisSetContactRequested          string = "flight.orbis.setcontact.requested"
	TopicFlightOrbisSetContactResults            string = "flight.orbis.setcontact.results"
	TopicFlightOrbisSetPassengerRequested        string = "flight.orbis.setpassenger.requested"
	TopicFlightOrbisSetPassengerResults          string = "flight.orbis.setpassenger.results"
	TopicFlightOrbisGetSSRRequested              string = "flight.orbis.getssr.requested"
	TopicFlightOrbisGetSSRResults                string = "flight.orbis.getssr.results"
	TopicFlightOrbisSetSSRRequested              string = "flight.orbis.setssr.requested"
	TopicFlightOrbisSetSSRResults                string = "flight.orbis.setssr.results"
	TopicFlightOrbisGetSeatRequested             string = "flight.orbis.getseat.requested"
	TopicFlightOrbisGetSeatResults               string = "flight.orbis.getseat.results"
	TopicFlightOrbisSetSeatRequested             string = "flight.orbis.setseat.requested"
	TopicFlightOrbisSetSeatResults               string = "flight.orbis.setseat.results"
	TopicFlightOrbisCommitBookingRequested       string = "flight.orbis.commitbooking.requested"
	TopicFlightOrbisCommitBookingResults         string = "flight.orbis.commitbooking.results"
	TopicFlightOrbisIssuedRequested              string = "flight.orbis.issued.requested"
	TopicFlightOrbisIssuedResults                string = "flight.orbis.issued.results"
	TopicFlightOrbisGetBookingRequested          string = "flight.orbis.getbooking.requested"
	TopicFlightOrbisGetBookingResults            string = "flight.orbis.getbooking.results"
	TopicFlightOrbisGetRescheduleRequested       string = "flight.orbis.getreschedule.requested"
	TopicFlightOrbisGetRescheduleResults         string = "flight.orbis.getreschedule.results"
	TopicFlightOrbisSetRescheduleRequested       string = "flight.orbis.setreschedule.requested"
	TopicFlightOrbisSetRescheduleResults         string = "flight.orbis.setreschedule.results"
	TopicFlightOrbisSetBookingRequested          string = "flight.orbis.setbooking.requested"
	TopicFlightOrbisSetBookingResults            string = "flight.orbis.setbooking.results"
	TopicFlightOrbisGetPostSsrRequested          string = "flight.orbis.getpostssr.requested"
	TopicFlightOrbisGetPostSsrResults            string = "flight.orbis.getpostssr.results"
	TopicFlightOrbisSetPostSsrsRequested         string = "flight.orbis.setpostssrs.requested"
	TopicFlightOrbisSetPostSsrsResults           string = "flight.orbis.setpostssrs.results"
	TopicFlightOrbisGetPostSeatRequested         string = "flight.orbis.getpostseat.requested"
	TopicFlightOrbisGetPostSeatResults           string = "flight.orbis.getpostseat.results"
	TopicFlightOrbisSetPostSeatsRequested        string = "flight.orbis.setpostseats.requested"
	TopicFlightOrbisSetPostSeatsResults          string = "flight.orbis.setpostseats.results"
	TopicFlightOrbisCancelRequested              string = "flight.orbis.cancel.requested"
	TopicFlightOrbisCancelResults                string = "flight.orbis.cancel.results"
	TopicFlightOrbisPreRefundLoginListRequested  string = "flight.orbis.prerefundloginlist.requested"
	TopicFlightOrbisPreRefundLoginListResults    string = "flight.orbis.prerefundloginlist.results"
	TopicFlightOrbisGetRefundBookingRequested    string = "flight.orbis.getrefundbooking.requested"
	TopicFlightOrbisGetRefundBookingResults      string = "flight.orbis.getrefundbooking.results"
	TopicFlightOrbisGetBalanceRequested          string = "flight.orbis.getbalance.requested"
	TopicFlightOrbisGetBalanceResults            string = "flight.orbis.getbalance.results"
	TopicFlightOrbisGetDestinationListRequested  string = "flight.orbis.getdestinationlist.requested"
	TopicFlightOrbisGetDestinationListResults    string = "flight.orbis.getdestinationlist.results"
	TopicFlightOrbisGetCarriersRequested         string = "flight.orbis.getcarriers.requested"
	TopicFlightOrbisGetCarriersResults           string = "flight.orbis.getcarriers.results"
	TopicFlightOrbisGetCarrierProvidersRequested string = "flight.orbis.getcarrierproviders.requested"
	TopicFlightOrbisGetCarrierProvidersResults   string = "flight.orbis.getcarrierproviders.results"
	TopicFlightOrbisGetProviderListRequested     string = "flight.orbis.getproviderlist.requested"
	TopicFlightOrbisGetProviderListResults       string = "flight.orbis.getproviderlist.results"
	TopicFlightOrbisSignoutRequested             string = "flight.orbis.signout.requested"
	TopicFlightOrbisSignoutResults               string = "flight.orbis.signout.results"
)

// Payment events
const (
	TopicPaymentCreateRequested string = "payment.create.requested"
	TopicPaymentCreateResults   string = "payment.create.results"
	TopicPaymentCallback        string = "payment.callback"
	TopicPaymentPostProcessing  string = "payment.postprocessing"
)

// PDF events
const (
	TopicInvoiceGenerationRequested string = "invoice.generation.requested"
	TopicInvoiceGenerationResults   string = "invoice.generation.results"
	TopicTicketGenerationRequested  string = "ticket.generation.requested"
	TopicTicketGenerationResults    string = "ticket.generation.results"
	TopicVisaGenerationRequested    string = "visa.pdf.generation.requested"
	TopicVisaGenerationResults      string = "visa.pdf.generation.results"
	TopicTourGenerationRequested    string = "tour.pdf.generation.requested"
	TopicTourGenerationResults      string = "tour.pdf.generation.results"
)

// Email events
const (
	TopicMailSendRequested string = "mail.send.requested"
)

// OTP events
const (
	TopicOTPSendRequested string = "otp.send.requested"
)

// Hotel events
const (
	TopicHotelSearchRequested             string = "hotel.search.requested"
	TopicHotelSearchResults               string = "hotel.search.results"
	TopicHotelRecheckRateRequested        string = "hotel.recheckrate.requested"
	TopicHotelRecheckRateResults          string = "hotel.recheckrate.results"
	TopicHotelBookRequested               string = "hotel.book.requested"
	TopicHotelBookResults                 string = "hotel.book.results"
	TopicHotelReservationDetailsRequested string = "hotel.reservationdetails.requested"
	TopicHotelReservationDetailsResults   string = "hotel.reservationdetails.results"
)

// Population events
const (
	TopicHotelPopulationRequested     string = "hotel.population.requested"
	TopicHotelCityPopulationRequested string = "hotel.city.population.requested"
)

// Booking events
const (
	// This event is emitted when a booking is created and the schedule is set to expire after a certain time.
	TopicBookingScheduleExpiredRequested string = "booking.schedule.expired.requested"
	TopicBookingScheduleExpiredResults   string = "booking.schedule.expired.results"

	// This event is emitted when a booking is cancelled and the schedule is to be removed from the system.
	TopicBookingScheduleCancelled string = "booking.schedule.cancelled"

	TopicBookingQueueRequested string = "booking.queue.requested"
	TopicBookingQueueResults   string = "booking.queue.results"
)

var CompressedEventTopics = map[string]bool{
	TopicFlightSearchResults: true,
	TopicHotelSearchResults:  true,
}
