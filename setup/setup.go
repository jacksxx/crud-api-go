package setup

import (
	"crud-api-go/arch/controller"
	"crud-api-go/arch/repository"
	"crud-api-go/arch/router"
	"crud-api-go/arch/service"
	"database/sql"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func InitApp(server *echo.Echo, dbConnection *sql.DB) {
	// Cria nova instância do Validador
	validate := validator.New()
	// Cria nova instância do Universal Translator
	ptLocale := pt_BR.New()
	uni := ut.New(ptLocale, ptLocale)
	translator, _ := uni.GetTranslator("pt")

	// Define a base da API
	baseRouter := "/api/v1"

	// Inicializa repositórios, serviços e controllers
	UserRepository := repository.NewUsuarioRepository(dbConnection)
	UserService := service.NewUsuarioService(UserRepository)
	UserController := controller.NewUsuarioController(UserService, validate, translator)
	router.UsuarioRouter(server, baseRouter+"/users", UserController)

	
	AuthService := service.NewAuthService(UserRepository)
	AuthController := controller.NewAuthController(AuthService, validate, translator)
	router.AuthRoutes(server, baseRouter+"/auth", AuthController)

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
}
