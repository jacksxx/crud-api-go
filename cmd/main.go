package main

import (
	"crud-api-go/arch/controller"
	"crud-api-go/arch/repository"
	"crud-api-go/arch/router"
	"crud-api-go/arch/service"
	"crud-api-go/config"
	"fmt"
	"time"

	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

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
	//Cria nova instancia de Validator
	validate := validator.New()
	// Cria nova instância do Universal Translator
	ptLocale := pt_BR.New()
	uni := ut.New(ptLocale, ptLocale)
	translator, _ := uni.GetTranslator("pt")

	// Define a base da API
	baseRouter := "/api/v1"

	// Inicializa repositórios, serviços e controllers
	CategoryRepository := repository.NewCategoriasRepository(dbConnection)
	CategoryService := service.NewCategoriaService(CategoryRepository)
	CategoryController := controller.NewCategoriaController(CategoryService, validate, translator)
	router.CategoriesRouter(server, baseRouter+"/categories", CategoryController)

	SubcategoryRepository := repository.NewSubcategoriasRepository(dbConnection)
	SubcategoryService := service.NewSubcategoriaService(SubcategoryRepository)
	SubcategoryController := controller.NewSubcategoriaController(SubcategoryService, validate, translator)
	router.SubcategoriesRouter(server, baseRouter+"/subcategories", SubcategoryController)

	UnitRepository := repository.NewUnidadesRepository(dbConnection)
	UnitService := service.NewUnidadesService(UnitRepository)
	UnitController := controller.NewUnitController(UnitService, validate, translator)
	router.UnitsRouter(server, baseRouter+"/unidades", UnitController)

	ProductRepository := repository.NewProductRepository(dbConnection)
	ProductService := service.NewProductService(ProductRepository)
	ProductController := controller.NewProductController(ProductService, validate, translator)
	router.ProductsRouter(server, baseRouter+"/products", ProductController)

	LstComprasItensRepository := repository.NewLstComprasItensRepository(dbConnection)
	LstComprasItensService := service.NewLstComprasItensService(LstComprasItensRepository)
	LstComprasItensController := controller.NewLstComprasItensController(LstComprasItensService, validate, translator)
	router.LstComprasItensRouter(server, baseRouter+"/lst_compras_itens", LstComprasItensController)

	LstComprasRepository := repository.NewLstComprasRepository(dbConnection)
	LstComprasService := service.NewLstComprasService(LstComprasRepository, LstComprasItensService)
	LstComprasController := controller.NewLstComprasController(LstComprasService, validate, translator)
	router.LstComprasRouter(server, baseRouter+"/lst_compras", LstComprasController)

	// Inicia o servidor na porta 8000
	server.Logger.Fatal(server.Start(":8000"))
}
