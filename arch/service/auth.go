package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"crud-api-go/helper"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type AuthService interface {
	Authenticate(c echo.Context, credentials model.LoginCredentials) (int, error)
	Refresh(c echo.Context) (*model.AuthTokens, int, error)
	Logout(c echo.Context) error
}

type authService struct {
	userRepository repository.UsuarioRepository
}

func NewAuthService(userRepository repository.UsuarioRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func (s *authService) Authenticate(c echo.Context, credentials model.LoginCredentials) (int, error) {
	usernameLower := strings.ToLower(credentials.Username)
	usuario, err := s.userRepository.GetUsuariosByUsername(usernameLower)
	if usuario.Id == 0 {
		return http.StatusNotFound, fmt.Errorf("usuário não encontrado")
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valid := helper.CheckPassword(*usuario.Senha, *usuario.Salt, credentials.Password)
	if valid {
		// Gera sessionID único para esta sessão
		sessionID := helper.GenerateSessionID()

		// Gera o token CSRF
		csrfToken, err := helper.GenerateCSRFToken()
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// Armazena o token CSRF associado ao usuário e sessionID
		err = helper.StoreCSRFToken(usuario.Id, csrfToken, sessionID)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// Modifica a geração do JWT para incluir o sessionID
		accessToken, err := helper.GenerateJWT(usernameLower, usuario.Id, sessionID)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		refreshToken, err := helper.GenerateRefreshJWT(usernameLower, usuario.Id, sessionID)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// Define o token CSRF no cabeçalho de resposta
		c.Response().Header().Set("X-CSRF-Token", csrfToken)

		// Define os cookies para os tokens
		helper.SetTokenCookies(c.Response().Writer, accessToken, refreshToken)

		// Retorna status OK e nil para indicar sucesso
		return http.StatusOK, nil
	}
	return http.StatusUnauthorized, fmt.Errorf("usuário ou senha inválidos")
}

func (s *authService) Refresh(c echo.Context) (*model.AuthTokens, int, error) {

	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("token de atualização não encontrado")
	}
	claims := &model.Claims{}
	refreshToken, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return helper.GetJWTKey(), nil
	})
	if err != nil || !refreshToken.Valid {
		return nil, http.StatusForbidden, fmt.Errorf("sessão expirada. Por favor, faça login novamente")
	}

	// Extrai o sessionID do token atual
	sessionID := claims.SessionID

	// Gera novos tokens mantendo o mesmo sessionID
	accessToken, err := helper.GenerateJWT(claims.Username, claims.CodigoUsuario, sessionID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	newRefreshToken, err := helper.GenerateRefreshJWT(claims.Username, claims.CodigoUsuario, sessionID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Define os novos cookies
	helper.SetTokenCookies(c.Response().Writer, accessToken, newRefreshToken)

	// Gera e define um novo token CSRF no cabeçalho da resposta
	newCsrfToken, err := helper.GenerateCSRFToken()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	c.Response().Header().Set("X-CSRF-Token", newCsrfToken)

	// Armazena o novo token CSRF associado ao usuário
	helper.StoreCSRFToken(claims.CodigoUsuario, newCsrfToken, sessionID)

	// Retorna os tokens junto com o código de status HTTP e sem erro
	return &model.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, http.StatusOK, nil
}

func (s *authService) Logout(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err == nil {
		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return helper.GetJWTKey(), nil
		})

		if err == nil && token.Valid {
			// Invalida o token CSRF usando o sessionID
			helper.InvalidateCSRFToken(claims.CodigoUsuario, claims.SessionID)
		}
	}

	helper.SetTokenCookies(c.Response().Writer, "", "")
	helper.RemoveCSRFToken(c)

	return nil
}
