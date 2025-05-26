package main

import (
	"crud-api-go/arch/model"
	"crud-api-go/config"
	_ "crud-api-go/docs/swagger"
	"crud-api-go/helper"
	"crud-api-go/middleware"
	"crud-api-go/setup"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title CRUD-API-GO
// @version 0.1
// @description Uma api go criada apenas para teste e estudo
// @host localhost:8000
// @BasePath /api/v1
func main() {
	// Carrega a configuração do banco de dados
	config.LoadConfig()
	fmt.Println(time.Now())
	// Conecta ao banco de dados
	dbConnection, err := config.Connect()
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()

	// Conecta ao Redis com tratamento de erro
	if err := helper.InitRedis(); err != nil {
		panic(err)
	}
	// Finaliza o redis e o context
	defer helper.RedisClient.Close()
	defer helper.ShutdownRedis()

	// Carrega a configuração do app
	appConfig := config.GetAppConfig()
	// Cria a instância do Echo
	server := echo.New()

	// Custom HTTP error handler
	server.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusNotFound {
			response := model.WebResponse{
				Code:   http.StatusNotFound,
				Status: http.StatusText(http.StatusNotFound),
				Errors: []string{"Rota não encontrada"},
			}
			c.JSON(http.StatusNotFound, response)
			return
		}
		// Para outros erros, mantém o comportamento padrão do Echo
		server.DefaultHTTPErrorHandler(err, c)
	}

	var corsOrigins []string
	if appConfig.AppEnv == "production" {
		corsOrigins = []string{appConfig.CorsOrigin}
	} else {
		corsOrigins = []string{"http://localhost:3000"}
	}

	server.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PATCH, echo.PUT, echo.DELETE},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Csrf-Token"},
	}))

	server.Use(middleware.LoggerMiddleware)
	//swagger
	server.GET("/swagger/*", echoSwagger.WrapHandler)
	//Carrega app setup
	setup.InitApp(server, dbConnection)

	// Inicia o servidor na porta 8000
	server.Logger.Fatal(server.Start(":8000"))
}
