package main

import (
	"github.com/rogeriofontes/p-blog-api/config"
	"github.com/rogeriofontes/p-blog-api/controllers"
	"github.com/rogeriofontes/p-blog-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	controllers.InitPostController()     // Inicializa a coleção de posts
	controllers.InitCategoryController() // Inicializa a coleção de categorias

	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run(":8080")
}
