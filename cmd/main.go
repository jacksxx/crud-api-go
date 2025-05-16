package main

import (
	"crud-api-go/config"
	_ "crud-api-go/docs/swagger"
	"crud-api-go/setup"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"time"
)

// @title CRUD-API-GO
// @version 0.1
// @description Uma api go criada apenas para teste e estudo
// @host localhost:8000
// @BasePath /api/v1
func main() {
	// Carrega a configuração do banco de dados
	config.LoadConfig()
	checkDB := config.GetDBConfig()
	fmt.Printf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s \n",
		checkDB.Host, checkDB.Port, checkDB.User, checkDB.DBName, checkDB.Password, checkDB.SSLMode,
	)
	fmt.Println(time.Now())

	// Conecta ao banco de dados
	dbConnection, err := config.Connect()
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()

	// Cria a instância do Echo
	server := echo.New()
	//swagger
	server.GET("/swagger/*", echoSwagger.WrapHandler)
	//Cria nova instancia de Validator
	setup.InitApp(server, dbConnection)

	// Inicia o servidor na porta 8000
	server.Logger.Fatal(server.Start(":8000"))
}
