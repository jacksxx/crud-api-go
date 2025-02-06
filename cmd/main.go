package main

import (
	"crud-api-go/arch/controller"
	"crud-api-go/arch/repository"
	"crud-api-go/arch/router"
	"crud-api-go/arch/service"
	"crud-api-go/config"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadConfig()
	checkDB := config.GetDBConfig()
	fmt.Printf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s \n",
		checkDB.Host, checkDB.Port, checkDB.User, checkDB.DBName, checkDB.Password, checkDB.SSLMode,
	)
	fmt.Println(time.Now())
	dbConnection, err := config.Connect()
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()
	server := echo.New()
	//
	ProductRepository := repository.NewProductRepository(dbConnection)
	ProductService := service.NewProductService(ProductRepository)
	ProductController := controller.NewProductController(ProductService)
	router.ProductsRouter(server, ProductController)

	server.Logger.Fatal(server.Start(":8000"))
}
