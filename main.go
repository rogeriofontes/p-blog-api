package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/rogeriofontes/p-blog-api/config"
	"github.com/rogeriofontes/p-blog-api/controllers"
	_ "github.com/rogeriofontes/p-blog-api/docs"
	"github.com/rogeriofontes/p-blog-api/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @title Blog API
// @version 1.0
// @description API para gerenciar posts, comentários, categorias e usuários.
// @termsOfService http://swagger.io/terms/

// @contact.name Rogério Fontes Tomaz
// @contact.url http://www.rogeriofontes.inf.br
// @contact.email rogeriofontes@rogeriofontes.inf.br

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
func main() {
	config.ConnectDatabase()
	controllers.InitPostController()     // Inicializa a coleção de posts
	controllers.InitCategoryController() // Inicializa a coleção de categorias
	controllers.InitTagController()      // Inicializa a coleção de tags
	controllers.InitCommentController()  // Inicializa a coleção de comentários
	controllers.InitReactionController() // Inicializa a coleção de reações
	controllers.InitUserController()     // Inicializa a coleção de usuários
	controllers.InitFavoriteController() // Inicializa a coleção de favoritos
	controllers.InitFollowerController() // Inicializa a coleção de seguidores

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Permitir todos os domínios
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 12 horas
	}))

	// Documentação Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(router)

	router.Run(":8484")
}
