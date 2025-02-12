package helper

import (
	"crud-api-go/arch/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
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

func CustomBind(c echo.Context, data interface{}) error {
	if err := c.Bind(data); err != nil {
		// Verifica se o erro é um UnmarshalTypeError
		if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			// Tenta extrair o valor recebido diretamente do JSON
			receivedValue := "valor inválido" // Padrão se não conseguirmos pegar o real

			// Obtém o corpo JSON para verificar o valor recebido
			body, _ := io.ReadAll(c.Request().Body)
			var jsonData map[string]interface{}
			if json.Unmarshal(body, &jsonData) == nil {
				if value, exists := jsonData[unmarshalErr.Field]; exists {
					receivedValue = fmt.Sprintf("%v", value)
				}
			}

			// Mensagem de erro dinâmica
			errorMessage := fmt.Sprintf("Erro ao deserializar o campo '%s': esperava-se um %s, mas recebeu-se '%s'.",
				unmarshalErr.Field,         // Nome do campo com erro
				unmarshalErr.Type.String(), // Tipo esperado
				receivedValue,              // Valor recebido no JSON
			)

			// Retorna erro formatado
			return fmt.Errorf("%s", errorMessage)
		}

		// Se não for um erro de tipo, retorna o erro original
		return err
	}

	// Retorna nil caso o bind tenha sido bem-sucedido
	return nil
}

func BindAndValidate(c echo.Context, validate *validator.Validate, data interface{}, translator ut.Translator) []string {
	if err := CustomBind(c, data); err != nil {
		// Retorna mensagem genérica + erro detalhado
		return []string{fmt.Sprintf("Erro de validação dos dados: %v", err)}
	}

	if err := validate.Struct(data); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Translate(translator))
		}

		return errors
	}

	return nil
}
