package train

import (
	"fmt"
	"strings"
)

// CommitBookingRequest matches the Orbis train commit_booking payload.
// A commit creates a provider booking; it does not issue tickets.
type CommitBookingRequest struct {
	Contacts    []BookingContact    `json:"contacts"`
	Passengers  []BookingPassenger  `json:"passengers"`
	Schedules   []BookingSchedule   `json:"schedules"`
	Booker      BookingBooker       `json:"booker"`
	ForceIssued bool                `json:"force_issued"`
}

type BookingBooker struct {
	Title           string `json:"title"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	CallingCode     string `json:"calling_code"`
	Mobile          string `json:"mobile"`
	NationalityCode string `json:"nationality_code"`
}

type BookingContact struct {
	BookingBooker
	IsAlsoBooker bool `json:"is_also_booker"`
}

type BookingPassenger struct {
	PaxType         string          `json:"pax_type"`
	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	Title           string          `json:"title"`
	BirthDate       string          `json:"birth_date"`
	NationalityCode string          `json:"nationality_code"`
	Identity        BookingIdentity `json:"identity"`
}

type BookingIdentity struct {
	IdentityCountryOfIssuedCode string `json:"identity_country_of_issued_code"`
	IdentityExpDate             string `json:"identity_expdate"`
	IdentityNumber              string `json:"identity_number"`
	IdentityType                string `json:"identity_type"`
}

type BookingSchedule struct {
	Journeys []BookingJourney `json:"journeys"`
	Provider string           `json:"provider"`
}

type BookingJourney struct {
	JourneyCode string `json:"journey_code"`
	FareCode    string `json:"fare_code"`
}

func (r CommitBookingRequest) Validate() error {
	if r.ForceIssued {
		return fmt.Errorf("force_issued must be false; ticket issuance is a separate operation")
	}
	if len(r.Contacts) == 0 {
		return fmt.Errorf("contacts is required")
	}
	if len(r.Passengers) == 0 {
		return fmt.Errorf("passengers is required")
	}
	if len(r.Schedules) < 1 || len(r.Schedules) > 2 {
		return fmt.Errorf("schedules must contain one or two journeys")
	}
	if err := validateBookingBooker(r.Booker); err != nil {
		return fmt.Errorf("booker: %w", err)
	}
	for i, contact := range r.Contacts {
		if err := validateBookingBooker(contact.BookingBooker); err != nil {
			return fmt.Errorf("contacts[%d]: %w", i, err)
		}
	}
	for i, passenger := range r.Passengers {
		if strings.TrimSpace(passenger.PaxType) == "" || strings.TrimSpace(passenger.FirstName) == "" ||
			strings.TrimSpace(passenger.LastName) == "" || strings.TrimSpace(passenger.Title) == "" ||
			strings.TrimSpace(passenger.BirthDate) == "" || strings.TrimSpace(passenger.NationalityCode) == "" ||
			strings.TrimSpace(passenger.Identity.IdentityType) == "" || strings.TrimSpace(passenger.Identity.IdentityNumber) == "" {
			return fmt.Errorf("passengers[%d] has missing required fields", i)
		}
	}
	for i, schedule := range r.Schedules {
		if schedule.Provider != ProviderKAI {
			return fmt.Errorf("schedules[%d].provider must be %q", i, ProviderKAI)
		}
		if len(schedule.Journeys) != 1 {
			return fmt.Errorf("schedules[%d].journeys must contain one journey", i)
		}
		for j, journey := range schedule.Journeys {
			if strings.TrimSpace(journey.JourneyCode) == "" || strings.TrimSpace(journey.FareCode) == "" {
				return fmt.Errorf("schedules[%d].journeys[%d] requires journey_code and fare_code", i, j)
			}
		}
	}
	return nil
}

func validateBookingBooker(b BookingBooker) error {
	if strings.TrimSpace(b.Title) == "" || strings.TrimSpace(b.FirstName) == "" ||
		strings.TrimSpace(b.LastName) == "" || strings.TrimSpace(b.Email) == "" ||
		strings.TrimSpace(b.CallingCode) == "" || strings.TrimSpace(b.Mobile) == "" ||
		strings.TrimSpace(b.NationalityCode) == "" {
		return fmt.Errorf("missing required contact details")
	}
	return nil
}

// CommitBookingResponse retains provider booking and seat assignment details
// for downstream persistence and eventual issuing / ticket printing.
type CommitBookingResponse struct {
	JSONRPC string              `json:"jsonrpc"`
	Result  CommitBookingResult `json:"result"`
}

type CommitBookingResult struct {
	BasicResult
	Response CommitBookingData `json:"response"`
}

type CommitBookingData struct {
	OrderNumber      string                    `json:"order_number"`
	PNR              string                    `json:"pnr"`
	State            string                    `json:"state"`
	StateDescription string                    `json:"state_description"`
	HoldDate         string                    `json:"hold_date"`
	Direction        string                    `json:"direction"`
	Origin           string                    `json:"origin"`
	Destination      string                    `json:"destination"`
	DepartureDate    string                    `json:"departure_date"`
	ArrivalDate      string                    `json:"arrival_date"`
	ProviderBookings []CommittedProviderBooking `json:"provider_bookings"`
}

type CommittedProviderBooking struct {
	PNR              string             `json:"pnr"`
	Provider         string             `json:"provider"`
	State            string             `json:"state"`
	StateDescription string             `json:"state_description"`
	Sequence         int                `json:"sequence"`
	Origin           string             `json:"origin"`
	Destination      string             `json:"destination"`
	DepartureDate    string             `json:"departure_date"`
	ArrivalDate      string             `json:"arrival_date"`
	Journeys         []CommittedJourney `json:"journeys"`
	Currency         string             `json:"currency"`
	HoldDate         string             `json:"hold_date"`
	Tickets          []CommittedTicket  `json:"tickets"`
	ErrorMsg         string             `json:"error_msg"`
}

type CommittedJourney struct {
	Sequence      int             `json:"sequence"`
	Origin        string          `json:"origin"`
	Destination   string          `json:"destination"`
	DepartureDate string          `json:"departure_date"`
	ArrivalDate   string          `json:"arrival_date"`
	PNR           string          `json:"pnr"`
	CarrierCode   string          `json:"carrier_code"`
	CarrierNumber string          `json:"carrier_number"`
	CarrierName   string          `json:"carrier_name"`
	JourneyCode   string          `json:"journey_code"`
	FareCode      string          `json:"fare_code"`
	Seats         []CommittedSeat `json:"seats"`
	CabinClass    string          `json:"cabin_class"`
	ClassOfService string         `json:"class_of_service"`
}

type CommittedSeat struct {
	Seat              string `json:"seat"`
	SeatCode          string `json:"seat_code"`
	Passenger         string `json:"passenger"`
	PassengerSequence int    `json:"passenger_sequence"`
}

type CommittedTicket struct {
	Passenger         string `json:"passenger"`
	PaxType           string `json:"pax_type"`
	TicketNumber      string `json:"ticket_number"`
	PassengerSequence int    `json:"passenger_sequence"`
}
