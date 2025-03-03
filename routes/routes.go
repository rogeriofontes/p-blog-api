package routes

import (
	"github.com/rogeriofontes/p-blog-api/controllers"
	"github.com/rogeriofontes/p-blog-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Rota de autenticação
	api.POST("/login", controllers.Login)

	// Rotas protegidas com JWT
	protected := api.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		// Rotas de categorias
		protected.POST("/categories", controllers.CreateCategory)
		protected.GET("/categories", controllers.GetCategories)
		protected.GET("/categories/:id", controllers.GetCategoryByID)
		protected.PUT("/categories/:id", controllers.UpdateCategory)
		protected.DELETE("/categories/:id", controllers.DeleteCategory)

		// Rotas de posts
		protected.POST("/posts", controllers.CreatePost)
		protected.PUT("/posts/:id", controllers.UpdatePost)
		protected.GET("/posts", controllers.GetPosts)
		protected.GET("/posts/:id", controllers.GetPostByID)
		protected.GET("/posts/category/:category_id", controllers.GetPostsByCategory)
		protected.DELETE("/posts/:id", controllers.DeletePost)

		// Rotas de tags
		protected.POST("/tags", controllers.CreateTag)
		protected.GET("/tags", controllers.GetTags)
		protected.GET("/tags/:id", controllers.GetTagByID)

		// Rotas de comentários
		protected.POST("/comments", controllers.CreateComment)
		protected.GET("/comments/post/:post_id", controllers.GetCommentsByPost)
		protected.GET("/comments/:id", controllers.GetCommentByID)
		protected.PUT("/comments/:id", controllers.UpdateComment)
		protected.DELETE("/comments/:id", controllers.DeleteComment)

		// Rotas de reações
		protected.POST("/reactions", controllers.CreateReaction)
		protected.GET("/reactions/:id", controllers.GetReactionByID)
		protected.GET("/reactions/likes/:post_id", controllers.CountLikes)
		protected.GET("/reactions/dislikes/:post_id", controllers.CountDislikes)

		// Rotas de usuários
		protected.POST("/users", controllers.CreateUser)
		protected.GET("/users", controllers.GetAllUsers)
		protected.GET("/users/:id", controllers.GetUserByID)

		protected.POST("/favorites", controllers.AddFavorite)
		protected.DELETE("/favorites/:user_id/:post_id", controllers.RemoveFavorite)
		protected.GET("/favorites/user/:user_id", controllers.GetFavoritesByUser)

		protected.POST("/follow", controllers.FollowUser)
		protected.GET("/followers/user/:user_id", controllers.GetFollowers)
	}
}
