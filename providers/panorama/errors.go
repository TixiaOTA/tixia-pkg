package panorama

// Panorama Error Codes
type PanoramaErrorCode int

const (
	// Invalid request.
	// Not a Panorama error code, but used as a fallback.
	// Request body is invalid.
	InvalidRequest PanoramaErrorCode = 400

	// Internal server error.
	// Not a Panorama error code, but used as a fallback.
	// Error likely caused by network issues or server issues.
	InternalServerError PanoramaErrorCode = 500

	// Success.
	Success PanoramaErrorCode = 0

	// Invalid action type.
	// Check headers request.
	InvalidActionType PanoramaErrorCode = 4001

	// Access denied.
	// Error while validating your credential.
	AccessDenied PanoramaErrorCode = 4002

	// Your session has been expired.
	// Sign in retry.
	SessionExpired PanoramaErrorCode = 4003

	// Permission denied.
	// Call panorama admin.
	PermissionDenied PanoramaErrorCode = 4004

	// Provider is unavailable.
	// Check initial data.
	ProviderUnavailable PanoramaErrorCode = 4005

	// Price changed.
	// Must retry search.
	PriceChanged PanoramaErrorCode = 4006

	// SSR not available.
	// Must retry search.
	SsrUnavailable PanoramaErrorCode = 4007

	// Seat not available.
	// Must retry search.
	SeatUnavailable PanoramaErrorCode = 4008

	// Failed to assign SSR.
	// Must retry search.
	FailedSsr PanoramaErrorCode = 4009

	// Failed to assign Seat.
	// Must retry search.
	FailedSeat PanoramaErrorCode = 4010

	// Missing Requirements.
	MissingRequirement PanoramaErrorCode = 4011

	// Duplicate Session, Please create new session.
	// Sign in retry.
	DuplicateSession PanoramaErrorCode = 4012

	// Data Validation Error.
	// Check request data before submit.
	ValidationError PanoramaErrorCode = 4013

	// You had already call this function.
	// Must retry search.
	AlreadyCall PanoramaErrorCode = 4014

	// Destinations are not allowed.
	// Must retry search.
	NotAllowed PanoramaErrorCode = 4015

	// Failed to commit booking.
	// Must retry search.
	FailedCommitBooking PanoramaErrorCode = 4016

	// Failed to get booking.
	// Must retry search.
	FailedGetBooking PanoramaErrorCode = 4017

	// Access restriction.
	AccessRestriction PanoramaErrorCode = 4018

	// Failed to Reissue.
	// Must retry search.
	FailedReissue PanoramaErrorCode = 4019

	// No Schedules.
	// Must retry search.
	NoSchedules PanoramaErrorCode = 4020

	// Price is not available. Please choose another fare class combination.
	// Must retry search.
	PriceUnavailable PanoramaErrorCode = 4021

	// Cancel booking failed.
	// Check get booking.
	FailedCancelBooking PanoramaErrorCode = 4022

	// Reservation not found.
	// Must retry search.
	ReservationNotFound PanoramaErrorCode = 4023

	// Vendor session expired.
	// Call panorama admin.
	VendorSessionExpired PanoramaErrorCode = 4024

	// Flight is not available.
	// Check initial data.
	FlightUnavailable PanoramaErrorCode = 4025

	// Hold booking restricted, You have to Direct Issued this reservation.
	// Must retry search.
	HoldBookingRestricted PanoramaErrorCode = 4026

	// Not enough vendor balance.
	// Call panorama admin.
	NotBalance PanoramaErrorCode = 4027

	// Restricted booking, infant cannot be paired with senior citizen (YCD).
	// Must retry search.
	CannotInfant PanoramaErrorCode = 4028

	// Failed issued.
	// Must retry search.
	FailedIssued PanoramaErrorCode = 4029

	// User is banned.
	// Call panorama admin.
	UserBanned PanoramaErrorCode = 4030

	// Resource not found.
	ResourceNotFound PanoramaErrorCode = 4031

	// Failed Transaction.
	FailedTransaction PanoramaErrorCode = 4032

	// Failed to execute payment.
	FailedPayment PanoramaErrorCode = 4033

	// Transaction has been paid on vendor.
	TransactionPaid PanoramaErrorCode = 4034

	// Failed Transaction, manual resync needed.
	TransactionManual PanoramaErrorCode = 4035

	// Transaction is still in progress, please try again later.
	TransactionProgress PanoramaErrorCode = 4036

	// Frequent Flyer is not available.
	FfUnavailable PanoramaErrorCode = 4037

	// PNR is not found.
	PnrNotFound PanoramaErrorCode = 4038

	// PNR not allowed to process.
	PnrNotAllowed PanoramaErrorCode = 4039

	// Failed Split Booking.
	FailedSplitBooking PanoramaErrorCode = 4040

	// Failed to retrieve vendor balance.
	// Call panorama admin.
	FailedVendorBalance PanoramaErrorCode = 4041

	// No update change detected.
	// Must retry search.
	NoUpdateDetected PanoramaErrorCode = 4042

	// Cancel is not available.
	CancelUnavailable PanoramaErrorCode = 4043

	// Split is not available.
	SplitUnavailable PanoramaErrorCode = 4044
)
