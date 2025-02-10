package helper

import (
	"crud-api-go/arch/model"
	"net/http"
	
	"github.com/labstack/echo/v4"
)

// WriteResponseBody envia a resposta JSON para o cliente
func WriteResponseBody(c echo.Context, response model.WebResponse) error {
	// Define o cabeçalho da resposta para indicar que o conteúdo é JSON
	c.Response().Header().Set("Content-Type", "application/json")

	// Envia a resposta JSON com o código de status HTTP e o corpo da resposta
	err := c.JSON(response.Code, response)

	// Caso ocorra um erro ao enviar a resposta JSON, trata o erro enviando um JSON de erro genérico
	if err != nil {
		// Tenta enviar uma nova resposta JSON com um erro interno do servidor (500)
		err := c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		
		// Caso ocorra outro erro ao tentar enviar a resposta de erro, retorna esse erro
		if err != nil {
			return err
		}
	}

	// Retorna o erro original, caso tenha ocorrido, ou nil se tudo ocorreu corretamente
	return err
}
