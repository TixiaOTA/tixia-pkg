package panorama

// PanoramaErrorCode is a provider-specific error code returned in result.error_code.
type PanoramaErrorCode int

const (
	// InvalidRequest is not a Panorama error code; used as a client-side fallback.
	InvalidRequest PanoramaErrorCode = 400

	// InternalServerError is not a Panorama error code; used as a client-side fallback.
	InternalServerError PanoramaErrorCode = 500

	// Success indicates the provider call succeeded.
	Success PanoramaErrorCode = 0

	InvalidActionType                   PanoramaErrorCode = 4001
	ErrorValidatingCredential           PanoramaErrorCode = 4002
	SessionExpired                      PanoramaErrorCode = 4003
	PermissionDenied                    PanoramaErrorCode = 4004
	ProviderUnavailable                 PanoramaErrorCode = 4005
	PriceChanged                        PanoramaErrorCode = 4006
	SSRNotAvailable                     PanoramaErrorCode = 4007
	SeatNotAvailable                    PanoramaErrorCode = 4008
	FailedToAssignSSR                   PanoramaErrorCode = 4009
	FailedToAssignSeat                  PanoramaErrorCode = 4010
	MissingRequirements                 PanoramaErrorCode = 4011
	DuplicateSession                    PanoramaErrorCode = 4012
	DataValidationError                 PanoramaErrorCode = 4013
	AlreadyCalledFunction               PanoramaErrorCode = 4014
	DestinationsNotAllowed              PanoramaErrorCode = 4015
	FailedToCommitBooking               PanoramaErrorCode = 4016
	FailedToGetBooking                  PanoramaErrorCode = 4017
	AccessRestriction                   PanoramaErrorCode = 4018
	FailedToReissue                     PanoramaErrorCode = 4019
	NoSchedules                         PanoramaErrorCode = 4020
	PriceNotAvailable                   PanoramaErrorCode = 4021
	CancelBookingFailed                 PanoramaErrorCode = 4022
	ReservationNotFound                 PanoramaErrorCode = 4023
	VendorSessionExpired                PanoramaErrorCode = 4024
	FlightNotAvailable                  PanoramaErrorCode = 4025
	HoldBookingRestricted               PanoramaErrorCode = 4026
	NotEnoughVendorBalance              PanoramaErrorCode = 4027
	RestrictedBooking                   PanoramaErrorCode = 4028
	FailedIssued                        PanoramaErrorCode = 4029
	UserBanned                          PanoramaErrorCode = 4030
	ResourceNotFound                    PanoramaErrorCode = 4031
	FailedTransaction                   PanoramaErrorCode = 4032
	FailedToExecutePayment              PanoramaErrorCode = 4033
	TransactionPaidOnVendor             PanoramaErrorCode = 4034
	FailedTransactionManualResyncNeeded PanoramaErrorCode = 4035
	TransactionInProgress               PanoramaErrorCode = 4036
	FrequentFlyerNotAvailable           PanoramaErrorCode = 4037
	PNRNotFound                         PanoramaErrorCode = 4038
	PNRNotAllowedToProcess              PanoramaErrorCode = 4039
	FailedSplitBooking                  PanoramaErrorCode = 4040
	FailedToRetrieveVendorBalance       PanoramaErrorCode = 4041
	NoUpdateChangeDetected              PanoramaErrorCode = 4042
	CancelNotAvailable                  PanoramaErrorCode = 4043
	SplitNotAvailable                   PanoramaErrorCode = 4044

	// AccessDenied is kept as an alias for credential validation failures.
	AccessDenied = ErrorValidatingCredential
)
