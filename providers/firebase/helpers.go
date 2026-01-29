package firebase

import (
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
)

// mapAuthError maps Firebase Auth errors to custom FirebaseError types
func mapAuthError(err error) *FirebaseError {
	if err == nil {
		return nil
	}

	// Check for specific Firebase Auth errors
	if auth.IsUserNotFound(err) {
		return &FirebaseError{
			Code:    UserNotFound,
			Message: "User not found",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	if auth.IsEmailAlreadyExists(err) {
		return &FirebaseError{
			Code:    EmailAlreadyExists,
			Message: "Email already exists",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	if auth.IsPhoneNumberAlreadyExists(err) {
		return &FirebaseError{
			Code:    PhoneNumberAlreadyExists,
			Message: "Phone number already exists",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	if auth.IsInvalidEmail(err) {
		return &FirebaseError{
			Code:    InvalidEmail,
			Message: "Invalid email address",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	// Check for invalid phone number by string matching
	if strings.Contains(err.Error(), "invalid phone number") ||
		strings.Contains(err.Error(), "PHONE_NUMBER") {
		return &FirebaseError{
			Code:    InvalidPhoneNumber,
			Message: "Invalid phone number",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	// Check for weak password by string matching
	if strings.Contains(err.Error(), "weak password") ||
		strings.Contains(err.Error(), "WEAK_PASSWORD") {
		return &FirebaseError{
			Code:    WeakPassword,
			Message: "Password is too weak",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	// Check for invalid ID token by string matching
	if strings.Contains(err.Error(), "invalid id token") ||
		strings.Contains(err.Error(), "INVALID_ID_TOKEN") {
		return &FirebaseError{
			Code:    InvalidIDToken,
			Message: "Invalid ID token",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	// Check for insufficient permission errors
	if strings.Contains(err.Error(), "insufficient_permission") ||
		strings.Contains(err.Error(), "permission-denied") {
		return &FirebaseError{
			Code:    InsufficientPermission,
			Message: "Insufficient permissions to perform this operation",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	// Check for invalid argument errors
	if strings.Contains(err.Error(), "invalid_argument") ||
		strings.Contains(err.Error(), "INVALID_ARGUMENT") {
		return &FirebaseError{
			Code:    InvalidArgument,
			Message: "Invalid argument provided",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	// Default to internal error for unknown auth errors
	return &FirebaseError{
		Code:    InternalError,
		Message: "Internal authentication error",
		Details: map[string]any{
			"originalError": err.Error(),
		},
	}
}

// mapMessagingError maps Firebase Messaging errors to custom FirebaseError types
func mapMessagingError(err error) *FirebaseError {
	if err == nil {
		return nil
	}

	// Check for specific Firebase Messaging errors
	if messaging.IsInvalidArgument(err) {
		return &FirebaseError{
			Code:    InvalidArgument,
			Message: "Invalid message argument",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	// Check for invalid registration by string matching
	if strings.Contains(err.Error(), "invalid registration") ||
		strings.Contains(err.Error(), "INVALID_REGISTRATION") {
		return &FirebaseError{
			Code:    InvalidRegistration,
			Message: "Invalid registration token",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	if messaging.IsUnregistered(err) {
		return &FirebaseError{
			Code:    Unregistered,
			Message: "Registration token is unregistered",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	// Check for invalid package name by string matching
	if strings.Contains(err.Error(), "invalid package name") ||
		strings.Contains(err.Error(), "INVALID_PACKAGE_NAME") {
		return &FirebaseError{
			Code:    InvalidPackageName,
			Message: "Invalid package name",
			Details: map[string]any{
				"originalError": err.Error(),
			},
		}
	}

	// Default to messaging error for unknown messaging errors
	return &FirebaseError{
		Code:    MessagingError,
		Message: "Messaging error",
		Details: map[string]any{
			"originalError": err.Error(),
		},
	}
}

// convertUserRecord converts Firebase SDK UserRecord to our custom UserRecord
func convertUserRecord(user *auth.UserRecord) *UserRecord {
	if user == nil {
		return nil
	}

	var lastLoginAt *time.Time
	if user.UserMetadata.LastLogInTimestamp != 0 {
		t := time.Unix(user.UserMetadata.LastLogInTimestamp, 0)
		lastLoginAt = &t
	}

	return &UserRecord{
		UID:           user.UID,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		DisplayName:   user.DisplayName,
		PhotoURL:      user.PhotoURL,
		Disabled:      user.Disabled,
		EmailVerified: user.EmailVerified,
		CustomClaims:  user.CustomClaims,
		CreatedAt:     time.Unix(user.UserMetadata.CreationTimestamp, 0),
		LastLoginAt:   lastLoginAt,
	}
}

// convertToken converts Firebase SDK Token to our custom Token
func convertToken(token *auth.Token) *Token {
	if token == nil {
		return nil
	}

	return &Token{
		Issuer:    token.Issuer,
		Audience:  token.Audience,
		Subject:   token.Subject,
		IssuedAt:  token.IssuedAt,
		ExpiresAt: token.Expires,
		Claims:    token.Claims,
	}
}