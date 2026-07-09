package panorama

import (
	"errors"
	"fmt"
)

// --------------------------------- UPDATE CONTACTS REQUEST ---------------------------------

// UpdateContactsRequest is the request structure for UpdateContacts
type UpdateContactsRequest struct {
	Signature string          `json:"signature"`
	Booker    ContactPerson   `json:"booker"`
	Contacts  []ContactPerson `json:"contacts"`
}

func (r *UpdateContactsRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if err := r.Booker.Validate(); err != nil {
		return errors.New("booker: " + err.Error())
	}
	if len(r.Contacts) == 0 {
		return errors.New("contacts is required")
	}
	for i, contact := range r.Contacts {
		if err := contact.Validate(); err != nil {
			return fmt.Errorf("contacts[%d]: %w", i, err)
		}
	}
	return nil
}

// ContactPerson is a booker or contact person for UpdateContacts
type ContactPerson struct {
	Title           string `json:"title"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	CallingCode     string `json:"calling_code"`
	Mobile          string `json:"mobile"`
	NationalityCode string `json:"nationality_code"`
	IsAlsoBooker    bool   `json:"is_also_booker"`
	Gender          string `json:"gender,omitempty"`
	ContactID       string `json:"contact_id,omitempty"`
	Sequence        int    `json:"sequence,omitempty"`
}

func (c *ContactPerson) Validate() error {
	if c.Title == "" {
		return errors.New("title is required")
	}
	if c.FirstName == "" {
		return errors.New("first_name is required")
	}
	if c.LastName == "" {
		return errors.New("last_name is required")
	}
	if c.Email == "" {
		return errors.New("email is required")
	}
	if c.CallingCode == "" {
		return errors.New("calling_code is required")
	}
	if c.Mobile == "" {
		return errors.New("mobile is required")
	}
	if c.NationalityCode == "" {
		return errors.New("nationality_code is required")
	}
	return nil
}

// --------------------------------- UPDATE CONTACTS RESPONSE ---------------------------------

// UpdateContactsResponse is the response structure for UpdateContacts
type UpdateContactsResponse struct {
	Result struct {
		BasicPanoramaResult
		Response UpdateContactsResult `json:"response"`
	} `json:"result"`
}

// UpdateContactsResult is the response payload for UpdateContacts
type UpdateContactsResult struct {
	Booker   ContactPerson   `json:"booker"`
	Contacts []ContactPerson `json:"contacts"`
}
