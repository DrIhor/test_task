package jwt

import "github.com/golang-jwt/jwt"

// JWT tocken info
type GoogleModel struct {
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}
