package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

// LoggerMiddleware é um middleware que registra detalhes da requisição.
func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		// Processa a requisição
		err := next(c)

		// Registra a informação da requisição
		log.Printf("Método: %s, URL: %s, Status: %d, Tempo: %s",
			c.Request().Method,
			c.Request().URL.Path,
			c.Response().Status,
			time.Since(start),
		)

		return err
	}
}
