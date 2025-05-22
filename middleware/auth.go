package middleware

import (
	"crud-api-go/arch/model"
	"crud-api-go/config"
	"crud-api-go/helper"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			env := config.GetAppConfig().AppEnv
			if env == "developmentWithoutAuth" {
				return next(c)
			}
			// Obtenção do token do cookie
			cookie, err := c.Cookie("token")
			if err != nil {
				return helper.BuildResponse(c, http.StatusUnauthorized, nil, []string{"Token de acesso não encontrado"})
			}

			// Parsing do token JWT
			claims := &model.Claims{}

			appConfig := config.GetAppConfig()
			jwtKey := appConfig.JWTKey

			if jwtKey == "" {
				return echo.ErrInternalServerError
			}

			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			})
			if err != nil || !token.Valid {
				return helper.BuildResponse(c, http.StatusUnauthorized, nil, []string{"Token de acesso inválido"})
			}
			// Verifica a expiração do token
			if claims.ExpiresAt < helper.TimeNow().Unix() {
				return helper.BuildResponse(c, http.StatusUnauthorized, nil, []string{"Token de acesso expirado"})
			}

			// Verificação de CSRF para métodos não-GET
			if c.Request().Method != http.MethodGet {
				// Obtém o token CSRF do cabeçalho
				csrfToken := c.Request().Header.Get("X-CSRF-Token")
				if csrfToken == "" {
					return helper.BuildResponse(c, http.StatusForbidden, nil, []string{"Token CSRF não fornecido"})
				}

				// Valida o token CSRF
				if !helper.ValidateCSRFToken(claims.CodigoUsuario, csrfToken, claims.SessionID) {
					return helper.BuildResponse(c, http.StatusForbidden, nil, []string{"Token CSRF inválido"})
				}
			}

			// Salva as claims no contexto para uso nos handlers seguintes
			c.Set("claims", claims)
			return next(c)
		}
	}
}
