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

var reactionCollection *mongo.Collection

// Inicializa a coleção de categorias após conectar ao banco
// @Summary Inicializa a coleção de reações
// @Description Inicializa a coleção de reações após conectar ao banco
// @Tags Reactions
// @Accept json
// @Produce json
// @Success 200 {string} string "Reações inicializadas"
// @Router /reactions/init [get]
func InitReactionController() {
	reactionCollection = config.GetCollection("reactions")
}

// Buscar todas as reações
// @Summary Buscar todas as reações
// @Description Buscar todas as reações
// @Tags Reactions
// @Accept json
// @Produce json
// @Success 200 {array} models.PostReaction
// @Router /reactions [get]
func CreateReaction(c *gin.Context) {
	var reaction models.PostReaction
	if err := c.ShouldBindJSON(&reaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reaction.ID = primitive.NewObjectID()
	reaction.CreatedAt = time.Now()

	_, err := reactionCollection.InsertOne(context.TODO(), reaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar reação"})
		return
	}

	c.JSON(http.StatusCreated, reaction)
}

// Contar total de likes de um post
// @Summary Contar total de likes de um post
// @Description Contar total de likes de um post
// @Tags Reactions
// @Accept json
// @Produce json
// @Param post_id path string true "ID do post"
// @Success 200 {string} string "Total de likes"
// @Router /reactions/likes/{post_id} [get]
func CountLikes(c *gin.Context) {
	postID := c.Param("post_id")

	// Converter para ObjectID
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do post inválido"})
		return
	}

	count, err := reactionCollection.CountDocuments(context.TODO(), bson.M{"post_id": objID, "reaction": true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar likes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": count})
}

// Buscar reação por ID
// @Summary Buscar reação por ID
// @Description Buscar reação por ID
// @Tags Reactions
// @Accept json
// @Produce json
// @Param id path string true "ID da reação"
// @Success 200 {object} models.PostReaction
// @Router /reactions/{id} [get]
func GetReactionByID(c *gin.Context) {
	reactionID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(reactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var reaction models.PostReaction
	err = reactionCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&reaction)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reação não encontrada"})
		return
	}

	c.JSON(http.StatusOK, reaction)
}

// Contar total de dislikes de um post
// @Summary Contar total de dislikes de um post
// @Description Contar total de dislikes de um post
// @Tags Reactions
// @Accept json
// @Produce json
// @Param post_id path string true "ID do post"
// @Success 200 {string} string "Total de dislikes"
// @Router /reactions/dislikes/{post_id} [get]
func CountDislikes(c *gin.Context) {
	postID := c.Param("post_id")

	// Converter para ObjectID
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do post inválido"})
		return
	}

	count, err := reactionCollection.CountDocuments(context.TODO(), bson.M{"post_id": objID, "reaction": false})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar dislikes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dislikes": count})
}
