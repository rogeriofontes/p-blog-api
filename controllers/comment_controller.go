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

var commentCollection *mongo.Collection

// Inicializa a coleção de categorias após conectar ao banco
// @Summary Inicializa a coleção de comentários
// @Description Inicializa a coleção de comentários após conectar ao banco
// @Tags Comentários
// @Accept json
// @Produce json
// @Router /init/comments [get]
func InitCommentController() {
	commentCollection = config.GetCollection("comments")
}

// Criar um comentário
// @Summary Criar um comentário
// @Description Criar um comentário
// @Tags Comentários
// @Accept  json
// @Produce  json
// @Param comment body models.PostComment true "Comentário"
// @Success 201 {object} models.PostComment
// @Router /comments [post]
func CreateComment(c *gin.Context) {
	var comment models.PostComment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.ID = primitive.NewObjectID()
	comment.CreatedAt = time.Now()

	_, err := commentCollection.InsertOne(context.TODO(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar comentário"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// Buscar todos os comentários
// @Summary Buscar todos os comentários
// @Description Buscar todos os comentários
// @Tags Comentários
// @Accept  json
// @Produce  json
// @Success 200 {array} models.PostComment
// @Router /comments [get]
func GetCommentsByPost(c *gin.Context) {
	postID := c.Query("post_id")
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do post inválido"})
		return
	}

	cursor, err := commentCollection.Find(context.TODO(), bson.M{"post_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar comentários"})
		return
	}

	var comments []models.PostComment
	if err = cursor.All(context.TODO(), &comments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar comentários"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// Atualizar um comentário por ID
// @Summary Atualizar um comentário por ID
// @Description Atualizar um comentário por ID
// @Tags Comentários
// @Accept  json
// @Produce  json
// @Param id path string true "ID do comentário"
// @Param comment body models.PostComment true "Comentário"
// @Success 200 {string} message "Comentário atualizado com sucesso"
// @Router /comments/{id} [put]
func UpdateComment(c *gin.Context) {
	commentID := c.Param("id")
	var updatedComment models.PostComment

	objID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = commentCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"content": updatedComment.Content}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar comentário"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comentário atualizado com sucesso"})
}

// Buscar todos os comentários
// @Summary Buscar todos os comentários
// @Description Buscar todos os comentários
// @Tags Comentários
// @Accept  json
// @Produce  json
// @Success 200 {array} models.PostComment
// @Router /comments [get]
func GetAllComments(c *gin.Context) {
	cursor, err := commentCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var comments []models.PostComment
	if err = cursor.All(context.TODO(), &comments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// Buscar comentário por ID
// @Summary Buscar comentário por ID
// @Description Buscar comentário por ID
// @Tags Comentários
// @Accept  json
// @Produce  json
// @Param id path string true "ID do comentário"
// @Success 200 {object} models.PostComment
// @Router /comments/{id} [get]
func GetCommentByID(c *gin.Context) {
	commentID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var comment models.PostComment
	err = commentCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&comment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comentário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// Deletar um comentário por ID
// @Summary Deletar um comentário por ID
// @Description Deletar um comentário por ID
// @Tags Comentários
// @Accept  json
// @Produce  json
// @Param id path string true "ID do comentário"
// @Success 200 {string} message "Comentário deletado com sucesso"
// @Router /comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	_, err = commentCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar comentário"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comentário deletado com sucesso"})
}
