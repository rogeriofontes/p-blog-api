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

	filter := bson.M{"post_id": reaction.PostID, "user_id": reaction.UserID}
	count, err := reactionCollection.CountDocuments(context.TODO(), filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar reação"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reação para esse usuário " + reaction.UserID + " já existe"})
		return
	}

	reaction.ID = primitive.NewObjectID()
	reaction.CreatedAt = time.Now()

	_, errCl := reactionCollection.InsertOne(context.TODO(), reaction)
	if errCl != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar reação"})
		return
	}

	c.JSON(http.StatusCreated, reaction)
}

// Buscar todas as reações
// @Summary Buscar todas as reações
// @Description Buscar todas as reações
// @Tags Reactions
// @Accept json
// @Produce json
// @Success 200 {array} models.PostReaction
// @Router /reactions [put]
func UpdateReaction(c *gin.Context) {
	var reaction models.PostReaction
	if err := c.ShouldBindJSON(&reaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"post_id": reaction.PostID, "user_id": reaction.UserID}
	update := bson.M{"$set": bson.M{"reaction": reaction.Reaction}}

	_, err := reactionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar reação"})
		return
	}

	c.JSON(http.StatusOK, reaction)
}

// Buscar todas as reações
// @Summary Buscar todas as reações
// @Description Buscar todas as reações
// @Tags Reactions
// @Accept json
// @Produce json
// @Success 200 {array} models.PostReaction
// @Router /reactions/post/{post_id} [get]
func GetReactionsByPost(c *gin.Context) {

	postID := c.Query("post_id")
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do post inválido"})
		return
	}

	cursor, err := reactionCollection.Find(context.TODO(), bson.M{"post_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar reações"})
		return
	}

	var reactions []models.PostReaction
	if err = cursor.All(context.TODO(), &reactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar reações"})
		return
	}

	c.JSON(http.StatusOK, reactions)
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

// Buscar todas as reações
// @Summary Buscar todas as reações
// @Description Buscar todas as reações
// @Tags Reactions
// @Accept json
// @Produce json
// @Success 200 {array} models.PostReaction
// @Router /reactions [get]
func GetReactions(c *gin.Context) {
	cursor, err := reactionCollection.Find(context.Background(), bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar reações"})
		return
	}

	var reactions []models.PostReaction
	if err = cursor.All(context.Background(), &reactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar reações"})
		return
	}

	c.JSON(http.StatusOK, reactions)
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

// Remover a reação de um usuário em um post
// @Summary Remover a reação de um usuário
// @Description Remove a reação (like/dislike) de um post
// @Tags Reactions
// @Accept json
// @Produce json
// @Param post_id path string true "ID do post"
// @Param user_id query string true "ID do usuário"
// @Success 204 {string} string "Reação removida com sucesso"
// @Router /reactions/{post_id} [delete]
func RemoveReaction(c *gin.Context) {
	reactionID := c.Param("id")

	// Converter para ObjectID
	objID, err := primitive.ObjectIDFromHex(reactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se o post existe
	var existingReaction bson.M
	err = reactionCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existingReaction)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post não encontrado"})
		return
	}

	// Deletar o post
	_, err = reactionCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deletado com sucesso"})
}

// Buscar a reação do usuário para um post específico
// @Summary Buscar a reação do usuário
// @Description Retorna a reação atual do usuário (like/dislike) em um post
// @Tags Reactions
// @Accept json
// @Produce json
// @Param post_id path string true "ID do post"
// @Param user_id query string true "ID do usuário"
// @Success 200 {object} map[string]string "Reação do usuário"
// @Router /reactions/{post_id}/user [get]
func GetUserReaction(c *gin.Context) {
	postID := c.Param("post_id")
	userID := c.Query("user_id")

	if postID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros inválidos"})
		return
	}

	var reaction models.PostReaction
	filter := bson.M{"post_id": postID, "user_id": userID}

	err := reactionCollection.FindOne(context.TODO(), filter).Decode(&reaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{"reaction": nil}) // Sem reação registrada
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar a reação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reaction": reaction.Reaction}) // Retorna "likes" ou "dislikes"
}
