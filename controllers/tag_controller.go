package controllers

import (
	"context"
	"net/http"

	"github.com/rogeriofontes/p-blog-api/config"
	"github.com/rogeriofontes/p-blog-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var tagCollection *mongo.Collection

// Inicializa a coleção de categorias após conectar ao banco
// @Summary Inicializa a coleção de tags
// @Description Inicializa a coleção de tags após conectar ao banco
// @Tags Tags
// @Accept json
// @Produce json
// @Success 200 {string} string "Tags inicializadas"
// @Router /tags/init [get]
func InitTagController() {
	tagCollection = config.GetCollection("tags")
}

func CreateTag(c *gin.Context) {
	if tagCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}

	var tag models.PostTag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := tagCollection.InsertOne(context.TODO(), tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// Buscar todas as tags
// @Summary Buscar todas as tags
// @Description Buscar todas as tags
// @Tags Tags
// @Accept json
// @Produce json
// @Success 200 {array} models.PostTag
// @Router /tags [get]
func GetTags(c *gin.Context) {
	if tagCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}

	cursor, err := tagCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var tags []models.PostTag
	if err = cursor.All(context.TODO(), &tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// Atualizar uma tag por ID
// @Summary Atualizar uma tag por ID
// @Description Atualizar uma tag por ID
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "ID da tag"
// @Success 200 {string} string "Tag atualizada com sucesso"
// @Router /tags/{id} [put]
func UpdateTag(c *gin.Context) {
	tagID := c.Param("id")
	var updatedTag models.PostTag

	objID, err := primitive.ObjectIDFromHex(tagID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.ShouldBindJSON(&updatedTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = tagCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"name": updatedTag.Name}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag atualizada com sucesso"})
}

// Buscar tag por ID
// @Summary Buscar tag por ID
// @Description Buscar tag por ID
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "ID da tag"
// @Success 200 {object} models.PostTag
// @Router /tags/{id} [get]
func GetTagByID(c *gin.Context) {
	tagID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(tagID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var tag models.PostTag
	err = tagCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&tag)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag não encontrada"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// Deletar uma tag por ID
// @Summary Deletar uma tag por ID
// @Description Deletar uma tag por ID
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "ID da tag"
// @Success 200 {string} string "Tag deletada com sucesso"
// @Router /tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	tagID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(tagID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	_, err = tagCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deletada com sucesso"})
}
