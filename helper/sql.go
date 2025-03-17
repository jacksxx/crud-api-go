package helper

import (
	"database/sql"
	"fmt"
	"net/http"
)

func CheckRowsAffected(queryResult sql.Result, responseText string) (int, error) {
	rowsAffected, err := queryResult.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao obter linhas afetadas: %v", err.Error())
	}
	if rowsAffected == 0 {
		return http.StatusNotFound, fmt.Errorf("%s", responseText)
	}
	return http.StatusOK, nil
}