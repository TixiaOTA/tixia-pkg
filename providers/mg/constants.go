package mg

const (
	// SearchHotelURL is the URL endpoint for searching hotels
	SearchHotelURL = "/SearchHotel"
)

// DefaultHeaders contains the default HTTP headers for MG API requests
var DefaultHeaders = map[string]string{
	"Content-Type": "application/json",
	"Accept":       "application/json",
}
