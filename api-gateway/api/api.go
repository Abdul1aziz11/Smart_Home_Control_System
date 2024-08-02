package api

import (
	"api-gateway/api/handlers"
	"api-gateway/genproto"
	"api-gateway/middleware"
	"api-gateway/rule"

	_ "api-gateway/docs"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9000
// @BasePath /
func NewGin(
	authServ genproto.AuthServiceClient,
	deviceServ genproto.DeviceServiceClient,
) *gin.Engine {
	// Создаем новый экземпляр Gin роутера
	r := gin.Default()

	// Инициализация Casbin enforcer
	enforcer, err := casbin.NewEnforcer("./rule/model.conf", "./rule/policy.csv")
	if err != nil {
		panic(err)
	}

	// Инициализация обработчиков с инъекцией зависимостей
	hnd := handlers.Handlers{
		Auth:          authServ,
		DeviceControl: deviceServ,
	}

	// Маршрут для Swagger документации
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Публичные маршруты
	auth := r.Group("/users")
	{
		auth.POST("/register", hnd.Register)
		auth.POST("/verify", hnd.Verify)
		auth.POST("/login", hnd.SignIn)
	}

	// Применение middleware аутентификации ко всем последующим маршрутам
	r.Use(middleware.TokenAuthMiddleware(authServ))

	// Маршруты, требующие аутентификации
	authenticated := r.Group("/")
	{
		authenticated.GET("/profile", rule.CasbinMiddleware(enforcer), hnd.Profile)

		devices := authenticated.Group("/devices").Use(rule.CasbinMiddleware(enforcer))
		{
			devices.POST("/", hnd.CreateDevice)
			devices.PUT("/:id", hnd.UpdateDevice)
			devices.DELETE("/:id", hnd.DeleteDevice)
			devices.GET("/:id", hnd.GetDevice)
		}

		controls := authenticated.Group("/control").Use(rule.CasbinMiddleware(enforcer))
		{
			controls.POST("/", hnd.ControlDevice)
		}
	}

	return r
}
