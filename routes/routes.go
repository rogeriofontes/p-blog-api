package routes

import (
	"github.com/rogeriofontes/p-blog-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/posts", controllers.CreatePost)
		api.GET("/posts", controllers.GetPosts)
		api.GET("/posts/:id", controllers.GetPostByID)

		// Rotas de categorias
		api.POST("/categories", controllers.CreateCategory)
		api.GET("/categories", controllers.GetCategories)
	}
}
