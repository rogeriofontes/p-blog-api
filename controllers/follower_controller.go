package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/rogeriofontes/p-blog-api/config"
	"github.com/rogeriofontes/p-blog-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var followerCollection *mongo.Collection

// Inicializa a coleção de categorias após conectar ao banco
// @Summary Inicializa a coleção de seguidores
// @Description Inicializa a coleção de seguidores após conectar ao banco
// @Tags Followers
// @Accept json
// @Produce json
// @Router /init/followers [get]
func InitFollowerController() {
	followerCollection = config.GetCollection("followers")
}

// Seguir um usuário
// @Summary Seguir um usuário
// @Description Seguir um usuário
// @Tags Followers
// @Accept json
// @Produce json
// @Param follower body models.Follower true "Seguidor a ser adicionado"
// @Success 201 {object} models.Follower
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /followers [post]
func FollowUser(c *gin.Context) {
	var follower models.Follower

	if err := c.ShouldBindJSON(&follower); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converter IDs para ObjectID
	userID, err := primitive.ObjectIDFromHex(follower.UserID.Hex())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário inválido"})
		return
	}

	followID, err := primitive.ObjectIDFromHex(follower.FollowID.Hex())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário seguido inválido"})
		return
	}

	// Verificar se o usuário já segue o outro
	count, err := followerCollection.CountDocuments(context.TODO(), bson.M{"user_id": userID, "follow_id": followID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar seguidor"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Usuário já está sendo seguido"})
		return
	}

	follower.ID = primitive.NewObjectID()
	follower.UserID = userID
	follower.FollowID = followID
	follower.FollowedAt = time.Now()

	_, err = followerCollection.InsertOne(context.TODO(), follower)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao seguir usuário"})
		return
	}

	c.JSON(http.StatusCreated, follower)
}

// Listar seguidores de um usuário
// @Summary Listar seguidores de um usuário
// @Description Listar seguidores de um usuário
// @Tags Followers
// @Accept json
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Success 200 {object} models.Follower
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /followers/{user_id} [get]
func GetFollowers(c *gin.Context) {
	userID := c.Param("user_id")

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário inválido"})
		return
	}

	cursor, err := followerCollection.Find(context.TODO(), bson.M{"follow_id": objUserID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar seguidores"})
		return
	}

	var followers []models.Follower
	if err = cursor.All(context.TODO(), &followers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar seguidores"})
		return
	}

	c.JSON(http.StatusOK, followers)
}
