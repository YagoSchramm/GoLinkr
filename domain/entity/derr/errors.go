package derr

var (
	BadRequestError     = NewBadRequestError("Bad Request")
	UnauthorizedError   = NewUnauthorizedError("Unauthorized")
	NotFoundError       = NewNotFoundError("Not Found")
	InternalServerError = NewInternalError("Internal Server Error")
)
var (
	InvalidCredentials = NewClientError("INVALID_CREDENTIALS", "Invalid credentials")
	InvalidEmail       = NewClientError("INVALID_EMAIL", "Invalid email")
	InvalidURL         = NewClientError("INVALID_URL", "Invalid url")
	UserAlreadyExists  = NewClientError("USER_ALREADY_EXISTS", "User already exists")
)
