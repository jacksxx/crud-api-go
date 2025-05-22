package model

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username      string `json:"username"`
	CodigoUsuario int    `json:"codigo_usuario"`
	SessionID     string `json:"session_id"`
	jwt.StandardClaims
}
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
