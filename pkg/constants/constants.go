package constants

//CLI colors
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

// Errors.
const (
	ErrMalformedRequest   string = "Malformed"
	ErrUnknownResource    string = "Resource Not Found"
	ErrMissingProperties  string = "Missing Values"
	ErrRequestDecoding    string = "JSON Decoding Error"
	ErrRequestBody        string = "Invalid Request Body"
	ErrAccessTokenInvalid string = "Invalid Access Token"
	ErrDatabaseOperation  string = "Database Operation Error"
	ErrInternal           string = "Internal Error"
	ErrBasicAuth          string = "Invalid Auth"
	ErrMissingStartTime   string = "Missing 'start' Query Parameters"
	ErrMissingEndTime     string = "Missing 'end' Query Parameters"
	ErrWrongURLParamType  string = "URL 'id' Parameter must be a number"
)
