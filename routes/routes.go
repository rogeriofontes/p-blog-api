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
	api.POST("/users", controllers.CreateUser)
	api.GET("/posts", controllers.GetPosts)
	api.GET("/posts/:id", controllers.GetPostByID)

	api.GET("/reactions/likes/post/:post_id", controllers.CountLikes)
	api.GET("/reactions/dislikes/post/:post_id", controllers.CountDislikes)

	api.GET("/comments/post", controllers.GetCommentsByPost)

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
		protected.GET("/posts/category/:category_id", controllers.GetPostsByCategory)
		protected.DELETE("/posts/:id", controllers.DeletePost)

		// Rotas de tags
		protected.POST("/tags", controllers.CreateTag)
		protected.GET("/tags", controllers.GetTags)
		protected.GET("/tags/:id", controllers.GetTagByID)
		protected.PUT("/tags/:id", controllers.UpdateTag)
		protected.DELETE("/tags/:id", controllers.DeleteTag)
		protected.GET("/tags/post/:post_id", controllers.GetTagsByPost)

		// Rotas de comentários
		protected.POST("/comments", controllers.CreateComment)
		protected.GET("/comments", controllers.GetAllComments)
		protected.GET("/comments/:id", controllers.GetCommentByID)
		protected.PUT("/comments/:id", controllers.UpdateComment)
		protected.DELETE("/comments/:id", controllers.DeleteComment)

		// Rotas de reações
		protected.GET("/reactions", controllers.GetReactions)
		protected.POST("/reactions", controllers.CreateReaction)
		protected.GET("/reactions/:id", controllers.GetReactionByID)
		protected.GET("/reactions/post", controllers.GetReactionsByPost)
		protected.PUT("/reactions/:id", controllers.UpdateReaction)
		protected.DELETE("/reactions/:id", controllers.RemoveReaction)
		protected.GET("/reactions/user/:post_id", controllers.GetUserReaction)

		// Rotas de usuários
		protected.GET("/users", controllers.GetAllUsers)
		protected.GET("/users/:id", controllers.GetUserByID)
		protected.PUT("/users/:id", controllers.UpdateUser)
		protected.DELETE("/users/:id", controllers.DeleteUser)

		protected.POST("/favorites", controllers.AddFavorite)
		protected.DELETE("/favorites/user/:user_id/post/:post_id", controllers.RemoveFavorite)
		protected.GET("/favorites/user/:user_id", controllers.GetFavoritesByUser)

		protected.POST("/follows", controllers.FollowUser)
		protected.GET("/follows/user/:user_id", controllers.GetFollows)
		protected.DELETE("/follows/user/:user_id/followed/:followed_id", controllers.UnfollowUser)
		protected.GET("/followers/user/:user_id", controllers.GetFollowers)
	}
}
