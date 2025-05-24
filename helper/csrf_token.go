package helper

import (
	"crud-api-go/config"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

// GenerateCSRFToken gera um token CSRF aleatório.
func GenerateCSRFToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

// GenerateSessionID gera um ID de sessão único.
func GenerateSessionID() string {
	sessionID := make([]byte, 16)
	rand.Read(sessionID)
	return hex.EncodeToString(sessionID)
}

// StoreCSRFToken armazena o token CSRF no servidor associado ao ID do usuário no Redis.
func StoreCSRFToken(userID int, token string, sessionID string) error {
	appConfig := config.GetAppConfig()

	// Armazenar o token CSRF no Redis com expiração configurada
	expiration := time.Duration(appConfig.JWTRefreshHours) * time.Hour
	key := fmt.Sprintf("csrf:%d:%s", userID, sessionID) // Formatar chave com ID numérico e sessionID
	err := RedisClient.Set(RedisCtx, key, token, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// InvalidateCSRFToken invalida o token CSRF associado ao ID do usuário no Redis.
func InvalidateCSRFToken(userID int, sessionID string) error {
	// Deletar o token CSRF do Redis
	key := fmt.Sprintf("csrf:%d:%s", userID, sessionID)
	err := RedisClient.Del(RedisCtx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

// RetrieveCSRFToken recupera o token CSRF associado ao ID do usuário do Redis.
func RetrieveCSRFToken(userID int, sessionID string) (string, error) {
	key := fmt.Sprintf("csrf:%d:%s", userID, sessionID)
	token, err := RedisClient.Get(RedisCtx, key).Result()
	if err == redis.Nil {
		return "", errors.New("token CSRF não encontrado")
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

// ValidateCSRFToken verifica se o token fornecido é válido para o usuário.
func ValidateCSRFToken(userID int, token string, sessionID string) bool {
	storedToken, err := RetrieveCSRFToken(userID, sessionID)
	if err != nil {
		return false
	}
	return storedToken == token
}

// StoreRefreshToken armazena o token de refresh como válido no Redis.
func StoreRefreshToken(token string) error {
	appConfig := config.GetAppConfig()
	refreshTokenExpiry := time.Duration(appConfig.JWTRefreshHours) * time.Hour

	// Armazenar o token de refresh no Redis com expiração configurada
	err := RedisClient.Set(RedisCtx, token, "valid", refreshTokenExpiry).Err()
	if err != nil {
		return err
	}
	return nil
}

// InvalidateRefreshToken invalida o token de refresh no Redis.
func InvalidateRefreshToken(token string) error {
	err := RedisClient.Del(RedisCtx, token).Err()
	if err != nil {
		return err
	}
	return nil
}

// SetTokenCookies define os cookies de token de acesso e refresh.
func SetTokenCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	appConfig := config.GetAppConfig()

	accessTokenExpiry := time.Duration(appConfig.JWTExpirationMinutes) * time.Minute
	refreshTokenExpiry := time.Duration(appConfig.JWTRefreshHours) * time.Hour

	now := TimeNow()

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  now.Add(accessTokenExpiry),
		Secure:   appConfig.AppEnv == "production",
		Domain:   appConfig.CookieDomain,
		SameSite: appConfig.SameSite,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  now.Add(refreshTokenExpiry),
		Secure:   appConfig.AppEnv == "production",
		Domain:   appConfig.CookieDomain,
		SameSite: appConfig.SameSite,
	})

}

// RemoveCSRFToken remove o cabeçalho CSRF de uma resposta HTTP.
func RemoveCSRFToken(c echo.Context) {
	c.Response().Header().Del("X-CSRF-Token")
}

// ListUserSessions retorna todas as sessões ativas de um usuário
func ListUserSessions(userID int) ([]string, error) {
	pattern := fmt.Sprintf("csrf:%d:*", userID)
	keys, err := RedisClient.Keys(RedisCtx, pattern).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]string, len(keys))
	for i, key := range keys {
		// Extrai o sessionID da chave (csrf:userID:sessionID)
		parts := strings.Split(key, ":")
		if len(parts) == 3 {
			sessions[i] = parts[2]
		}
	}
	return sessions, nil
}
