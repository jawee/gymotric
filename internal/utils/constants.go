package utils

const (
	AccessTokenCookieName                 = "X-wt-token"
	RefreshTokenCookieName                = "X-wt-refresh"
	EnvJwtExpireMinutes                   = "JWT_EXPIRE_MINUTES"
	EnvJwtRefreshExpireMinutes            = "JWT_REFRESH_EXPIRE_MINUTES"
	EnvJwtSignKey                         = "JWT_SIGN_KEY"
	EnvJwtRefreshSignKey                  = "JWT_REFRESH_SIGN_KEY"
	EnvSendGridApiKey                     = "SENDGRID_KEY"
	EnvBrevoApiKey                        = "BREVO_KEY"
	EnvApiKey                             = "API_KEY"
	ResetPasswordTokenExpireMinutes       = 10
	EmailConfirmationTokenExpireMinutes   = 10
	AccountConfirmationTokenExpireMinutes = 30
	EnvLogLevel                           = "LOG_LEVEL"
)
