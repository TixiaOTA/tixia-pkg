package panorama

import (
	"context"
)

// panoramaClient defines the interface for client methods used by the usecase
type panoramaClient interface {
	SignIn(ctx context.Context) (*SignInResponse, *PanoramaError)
	Search(ctx context.Context, req SearchRequest) (*SearchResponse, *PanoramaError)
	GetBooking(ctx context.Context, req GetBookingRequest) (*GetBookingResponse, *PanoramaError)
	CancelVoid(ctx context.Context, req CancelVoidRequest) (*GetBookingResponse, *PanoramaError)
}

type Usecase struct {
	client panoramaClient
}

func NewUsecase(client *Client) *Usecase {
	return &Usecase{
		client: client,
	}
}

func (u *Usecase) Search(ctx context.Context, req SearchBookingRequest) (*SearchBookingResponse, *PanoramaError) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, buildPanoramaError(InvalidRequest, err.Error())
	}

	// Sign in to get signature
	signInResp, err := u.client.SignIn(ctx)
	if err != nil {
		return nil, err
	}

	signature := signInResp.Result.Response.Signature

	var journeyLists []JourneyListRequest
	for _, v := range req.JourneyList {
		journeyList := JourneyListRequest{
			Origin:        v.Origin,
			Destination:   v.Destination,
			DepartureDate: v.DepartureDate,
		}
		journeyLists = append(journeyLists, journeyList)
	}

	// Search to get flight
	searchResp, err := u.client.Search(ctx, SearchRequest{
		Signature:    signature,
		JourneyList:  journeyLists,
		Adult:        req.Adult,
		Child:        req.Child,
		Infant:       req.Infant,
		CabinClass:   req.CabinClass,
		CarrierCodes: req.CarrierCodes,
		IsComboPrice: req.IsComboPrice,
		Provider:     req.Provider,
		PromoCodes:   req.PromoCodes,
	})
	if err != nil {
		return nil, err
	}

	return &SearchBookingResponse{
		Signature: signature,
		Results:   searchResp.Result.Response,
	}, nil
}

func (u *Usecase) CancelBooking(ctx context.Context, req CancelBookingRequest) (*CancelBookingResponse, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Sign in to get signature
	signInResp, err := u.client.SignIn(ctx)
	if err != nil {
		return nil, err
	}

	signature := signInResp.Result.Response.Signature

	var results []CancelBookingResult

	// Loop over order numbers
	for _, orderNumber := range req.OrderNumbers {
		result := CancelBookingResult{
			OrderNumber: orderNumber,
			Success:     false,
		}

		// Get booking details
		getBookingResp, err := u.client.GetBooking(ctx, GetBookingRequest{
			Signature:   signature,
			OrderNumber: orderNumber,
		})
		if err != nil {
			result.Message = err.Error()
			results = append(results, result)
			continue
		}

		// Check if booking state is "booked"
		if getBookingResp.Result.Response.State != "booked" {
			result.Message = "Booking state is not 'booked', skipping cancel"
			result.State = getBookingResp.Result.Response.State
			results = append(results, result)
			continue
		}

		// Cancel the booking
		cancelResp, err := u.client.CancelVoid(ctx, CancelVoidRequest{
			Signature:   signature,
			OrderNumber: orderNumber,
		})
		if err != nil {
			result.Message = err.Error()
			results = append(results, result)
			continue
		}

		// Check if state is now "void"
		if cancelResp.Result.Response.State == "void" {
			result.Success = true
			result.Message = "Booking successfully canceled"
			result.State = cancelResp.Result.Response.State
		} else {
			result.Message = "Cancel operation completed but state is not 'void'"
			result.State = cancelResp.Result.Response.State
		}

		results = append(results, result)
	}

	return &CancelBookingResponse{
		Results: results,
	}, nil
}
