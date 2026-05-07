package derr

var (
	BadRequestError     = NewBadRequestError("Bad Request")
	UnauthorizedError   = NewUnauthorizedError("Unauthorized")
	NotFoundError       = NewNotFoundError("Not Found")
	InternalServerError = NewInternalError("Internal Server Error")
)
var (
	InvalidCredentialsError = NewClientError("INVALID_CREDENTIALS", "Invalid credentials")
	EmailRequired           = NewClientError("EMAIL_REQUIRED", "Email required")
	InvalidURLError         = NewClientError("INVALID_URL", "Invalid url")
	UserAlreadyExists       = NewClientError("USER_ALREADY_EXISTS", "User already exists")
	PasswordRequired        = NewClientError("PASSWORD_REQUIRED", "Password required")
	InvalidCredentials      = NewClientError("INVALID_CREDENTIALS", "Invalid credentials")
	InvalidEmail            = NewClientError("INVALID_EMAIL", "Invalid email")
)

var (
	InvalidUsername    = NewClientError("INVALID_USERNAME", "Invalid username")
	BiographyIsTooLong = NewClientError("BIOGRAPHY_IS_TOO_LONG", "Biography is too long")
	NameIsTooShort     = NewClientError("NAME_IS_TOO_SHORT", "Name is too short")
	NameIsTooLong      = NewClientError("NAME_IS_TOO_LONG", "Name is too long")
	WeakPassword       = NewClientError("WEAK_PASSWORD", "Weak password")
)
