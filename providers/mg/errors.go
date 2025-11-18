package mg

// MG Error Codes
type MGErrorCode string

const (
	// Internal server error.
	// Not an MG error code, but used as a fallback.
	// Error likely caused by network issues or server issues.
	InternalServerError MGErrorCode = "INTERNAL_SERVER_ERROR"

	// Intermittent server error.
	// Must retry
	JRVXML161 MGErrorCode = "JRVXML161"

	// Invalid request.
	// Request body is invalid.
	JRVXML013 MGErrorCode = "JRVXML013"
)
