package panorama

import (
	"errors"
	"fmt"
)

// --------------------------------- UPDATE PASSENGERS REQUEST ---------------------------------

// UpdatePassengersRequest is the request structure for UpdatePassengers
type UpdatePassengersRequest struct {
	Signature  string            `json:"signature"`
	Passengers []UpdatePassenger `json:"passengers"`
}

func (r *UpdatePassengersRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if len(r.Passengers) == 0 {
		return errors.New("passengers is required")
	}
	for i, passenger := range r.Passengers {
		if err := passenger.Validate(); err != nil {
			return fmt.Errorf("passengers[%d]: %w", i, err)
		}
	}
	return nil
}

// UpdatePassenger is a passenger in an UpdatePassengers request
type UpdatePassenger struct {
	Title           string              `json:"title"`
	FirstName       string              `json:"first_name"`
	LastName        string              `json:"last_name"`
	BirthDate       string              `json:"birth_date"`
	PaxType         string              `json:"pax_type"`
	NationalityCode string              `json:"nationality_code"`
	Identity        *PassengerIdentity  `json:"identity,omitempty"`
	FfNumbers       []PassengerFfNumber `json:"ff_numbers,omitempty"`
}

func (p *UpdatePassenger) Validate() error {
	if p.Title == "" {
		return errors.New("title is required")
	}
	if p.FirstName == "" {
		return errors.New("first_name is required")
	}
	if p.LastName == "" {
		return errors.New("last_name is required")
	}
	if p.BirthDate == "" {
		return errors.New("birth_date is required")
	}
	if p.PaxType == "" {
		return errors.New("pax_type is required")
	}
	if p.NationalityCode == "" {
		return errors.New("nationality_code is required")
	}
	if p.Identity != nil {
		if err := p.Identity.Validate(); err != nil {
			return errors.New("identity: " + err.Error())
		}
	}
	return nil
}

// PassengerIdentity is optional identity document data for a passenger
type PassengerIdentity struct {
	IdentityNumber              string `json:"identity_number"`
	IdentityExpdate             string `json:"identity_expdate"`
	IdentityCountryOfIssuedCode string `json:"identity_country_of_issued_code"`
	IdentityType                string `json:"identity_type"`
}

func (i *PassengerIdentity) Validate() error {
	if i.IdentityNumber == "" {
		return errors.New("identity_number is required")
	}
	if i.IdentityExpdate == "" {
		return errors.New("identity_expdate is required")
	}
	if i.IdentityCountryOfIssuedCode == "" {
		return errors.New("identity_country_of_issued_code is required")
	}
	if i.IdentityType == "" {
		return errors.New("identity_type is required")
	}
	return nil
}

// PassengerFfNumber is an optional frequent flyer number for a passenger
type PassengerFfNumber struct {
	ScheduleID int    `json:"schedule_id"`
	FfNumber   string `json:"ff_number"`
	FfCode     string `json:"ff_code"`
}

// --------------------------------- UPDATE PASSENGERS RESPONSE ---------------------------------

// UpdatePassengersResponse is the response structure for UpdatePassengers
type UpdatePassengersResponse struct {
	Result struct {
		BasicPanoramaResult
		Response UpdatePassengersResult `json:"response"`
	} `json:"result"`
}

// UpdatePassengersResult is the response payload for UpdatePassengers
type UpdatePassengersResult struct {
	Passengers []UpdatePassengerResult `json:"passengers"`
}

// UpdatePassengerResult is a passenger in an UpdatePassengers response
type UpdatePassengerResult struct {
	Title           string                    `json:"title"`
	FirstName       string                    `json:"first_name"`
	LastName        string                    `json:"last_name"`
	BirthDate       string                    `json:"birth_date"`
	PaxType         string                    `json:"pax_type"`
	NationalityCode string                    `json:"nationality_code"`
	IdentityFromAPI *PassengerIdentityFromAPI `json:"identity_from_api,omitempty"`
	PassengerNumber int                       `json:"passenger_number"`
	Sequence        int                       `json:"sequence"`
	PassengerID     string                    `json:"passenger_id"`
	Gender          string                    `json:"gender"`
	IsAlsoBooker    bool                      `json:"is_also_booker"`
	IsAlsoContact   bool                      `json:"is_also_contact"`
}

// PassengerIdentityFromAPI is identity data returned by the API
type PassengerIdentityFromAPI struct {
	PassengerIdentity
	IsValidIdentity bool `json:"is_valid_identity"`
}
