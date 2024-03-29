package main

import (
	db "example/service/api/db"
	"example/service/api/docs"

	"example/service/api/config"
	notifier "example/service/api/notifier"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

func main() {
	config.InitConfig()

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Group("/api/v1")

	db.AddApiRoutes(v1)
	go notifier.CreateObjectCreationNotifierFunc()(notifier.ObjectCreationNotificationChannel)
	// go prod.CreateConsumerFunc()()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run() // listen and serve on 0.0.0.0:8080
}
