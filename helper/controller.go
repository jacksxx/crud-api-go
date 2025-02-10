package helper

import (
	"crud-api-go/arch/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

// BuildResponse cria um objeto WebResponse e envia a resposta para o cliente
func BuildResponse(c echo.Context, code int, data interface{}, errors []string) error {
	// Cria um objeto do tipo WebResponse com os dados fornecidos
	response := model.WebResponse{
		Code:   code,                  // Código de status HTTP da resposta
		Status: http.StatusText(code), // Converte o código HTTP em um status textual (ex: 200 -> "OK")
		Data:   data,                  // Os dados da resposta (payload)
		Errors: errors,                // Lista de erros, caso existam
	}

	// Chama WriteResponseBody para enviar a resposta ao cliente
	return WriteResponseBody(c, response)
}

func SuccessResponse(c echo.Context, data interface{}) error {
	response := model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	return WriteResponseBody(c, response)
}