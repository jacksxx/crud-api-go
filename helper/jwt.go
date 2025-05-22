package helper

import (
	"crud-api-go/arch/model"
	"crud-api-go/config"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// GenerateJWT cria um JWT de acesso
func GenerateJWT(username string, codigoUsuario int, sessionID string) (string, error) {
	// Obtenha a configuração do JWT
	appConfig := config.GetAppConfig()

	// Defina o tempo de expiração baseado na configuração
	expirationTime := TimeNow().Add(time.Duration(appConfig.JWTExpirationMinutes) * time.Minute)

	// Criação do JWT com as claims
	claims := &model.Claims{
		Username:      username,
		CodigoUsuario: codigoUsuario,
		SessionID:     sessionID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Criação e assinatura do token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(appConfig.JWTKey))
}

// GenerateRefreshJWT cria um JWT de refresh
func GenerateRefreshJWT(username string, codigoUsuario int, sessionID string) (string, error) {
	appConfig := config.GetAppConfig()

	expirationTime := TimeNow().Add(time.Duration(appConfig.JWTRefreshHours) * time.Hour)

	// Criação do JWT de refresh com as claims, incluindo perfis
	claims := &model.Claims{
		Username:      username,
		CodigoUsuario: codigoUsuario,
		SessionID:     sessionID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Criação e assinatura do token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(appConfig.JWTKey))
}

// GetJWTKey retorna a chave do JWT a partir da configuração
func GetJWTKey() []byte {
	appConfig := config.GetAppConfig()

	return []byte(appConfig.JWTKey)
}

// ExtrairCodigoUsuarioDoToken extrai o código do usuário do token JWT no cookie
func ExtrairCodigoUsuarioDoToken(ctx echo.Context) (int, error) {
	appEnv := config.GetAppConfig().AppEnv

	cookie, err := ctx.Cookie("token")
	if err != nil {
		//Id pra teste do user com id = 1
		if appEnv != "production" {
			return 1, nil
		}
		return 0, fmt.Errorf("token de autenticação ausente ou inválido")
	}

	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJWTKey(), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("token inválido: %v", err)
	}

	return claims.CodigoUsuario, nil
}
