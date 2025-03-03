package controllers

import (
	"context"

	"github.com/rogeriofontes/p-blog-api/config"
	"github.com/rogeriofontes/p-blog-api/models"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var favoriteCollection *mongo.Collection

// Inicializa a coleção de categorias após conectar ao banco
// @Summary Inicializa a coleção de favoritos
// @Description Inicializa a coleção de favoritos após conectar ao banco
// @Tags Favorites
// @Accept json
// @Produce json
// @Router /init/favorites [get]
func InitFavoriteController() {
	favoriteCollection = config.GetCollection("favorites")
}

// Adicionar um post aos favoritos
// @Summary Adicionar um post aos favoritos
// @Description Adiciona um post aos favoritos de um usuário
// @Tags Favorites
// @Accept json
// @Produce json
// @Param favorite body models.FavoritePost true "Favorito a ser adicionado"
// @Success 201 {object} models.FavoritePost
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /favorites [post]
func AddFavorite(c *gin.Context) {
	var favorite models.FavoritePost

	if err := c.ShouldBindJSON(&favorite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converter IDs para ObjectID
	userID, err := primitive.ObjectIDFromHex(favorite.UserID.Hex())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário inválido"})
		return
	}

	postID, err := primitive.ObjectIDFromHex(favorite.PostID.Hex())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do post inválido"})
		return
	}

	// Verificar se o post já foi favoritado pelo usuário
	count, err := favoriteCollection.CountDocuments(context.TODO(), bson.M{"user_id": userID, "post_id": postID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar favorito"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Post já está nos favoritos"})
		return
	}

	favorite.ID = primitive.NewObjectID()
	favorite.UserID = userID
	favorite.PostID = postID
	favorite.SavedAt = time.Now()

	_, err = favoriteCollection.InsertOne(context.TODO(), favorite)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao favoritar post"})
		return
	}

	c.JSON(http.StatusCreated, favorite)
}

// Remover um post dos favoritos
// @Summary Remover um post dos favoritos
// @Description Remove um post dos favoritos de um usuário
// @Tags Favorites
// @Accept json
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Param post_id path string true "ID do post"
// @Success 200 {string} Message
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /favorites/{user_id}/{post_id} [delete]
func RemoveFavorite(c *gin.Context) {
	userID := c.Param("user_id")
	postID := c.Param("post_id")

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário inválido"})
		return
	}

	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do post inválido"})
		return
	}

	_, err = favoriteCollection.DeleteOne(context.TODO(), bson.M{"user_id": objUserID, "post_id": objPostID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover favorito"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post removido dos favoritos"})
}

// Listar posts favoritados por um usuário
// @Summary Listar posts favoritados por um usuário
// @Description Lista os posts favoritados por um usuário
// @Tags Favorites
// @Accept json
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Success 200 {array} models.FavoritePost
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /favorites/{user_id} [get]
func GetFavoritesByUser(c *gin.Context) {
	userID := c.Param("user_id")

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário inválido"})
		return
	}

	cursor, err := favoriteCollection.Find(context.TODO(), bson.M{"user_id": objUserID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar favoritos"})
		return
	}

	var favorites []models.FavoritePost
	if err = cursor.All(context.TODO(), &favorites); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar favoritos"})
		return
	}

	c.JSON(http.StatusOK, favorites)
}
