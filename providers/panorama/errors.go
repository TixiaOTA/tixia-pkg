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

	// Access denied.
	AccessDenied PanoramaErrorCode = 4002
)
