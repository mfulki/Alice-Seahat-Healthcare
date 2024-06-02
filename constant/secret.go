package constant

import "time"

const (
	JwtIssuer          = "jwtSeahat"
	JwtDefaultDuration = 1 * time.Hour
)

const (
	TokenTypeConfirm = "confirm"
	TokenTypeReset   = "reset"
	TokenDuration    = 30 * time.Minute
	TokenLength      = 20

	TokenAccess  = "access_token"
	TokenRefresh = "refresh_token"
)

const (
	GeneratedPasswordLength = 15
)
