package firebase

import "fmt"

// FirebaseErrorCode represents Firebase-specific error codes
type FirebaseErrorCode string

const (
	// Authentication errors
	UserNotFound             FirebaseErrorCode = "USER_NOT_FOUND"
	InvalidIDToken           FirebaseErrorCode = "INVALID_ID_TOKEN"
	InsufficientPermission   FirebaseErrorCode = "INSUFFICIENT_PERMISSION"
	InvalidArgument          FirebaseErrorCode = "INVALID_ARGUMENT"
	EmailAlreadyExists       FirebaseErrorCode = "EMAIL_ALREADY_EXISTS"
	PhoneNumberAlreadyExists FirebaseErrorCode = "PHONE_NUMBER_ALREADY_EXISTS"
	WeakPassword             FirebaseErrorCode = "WEAK_PASSWORD"
	InvalidEmail             FirebaseErrorCode = "INVALID_EMAIL"
	InvalidPhoneNumber       FirebaseErrorCode = "INVALID_PHONE_NUMBER"

	// Messaging errors
	MessagingError      FirebaseErrorCode = "MESSAGING_ERROR"
	InvalidRegistration FirebaseErrorCode = "INVALID_REGISTRATION"
	Unregistered        FirebaseErrorCode = "UNREGISTERED"
	InvalidPackageName  FirebaseErrorCode = "INVALID_PACKAGE_NAME"

	// Generic errors
	InternalError FirebaseErrorCode = "INTERNAL_ERROR"
	UnknownError  FirebaseErrorCode = "UNKNOWN_ERROR"
)

// FirebaseError represents a Firebase-specific error
type FirebaseError struct {
	Code    FirebaseErrorCode `json:"code"`
	Message string            `json:"message"`
	Details map[string]any    `json:"details,omitempty"`
}

func (e *FirebaseError) Error() string {
	return fmt.Sprintf("Firebase error (%s): %s", e.Code, e.Message)
}
