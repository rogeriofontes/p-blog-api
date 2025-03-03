package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rogeriofontes/p-blog-api/config"
	"github.com/rogeriofontes/p-blog-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection

// Inicializa a coleção de categorias após conectar ao banco
// @Summary Inicializa a coleção de categorias
// @Description Inicializa a coleção de categorias após conectar ao banco
// @Tags categories
// @Accept json
// @Produce json
// @Router /init/categories [get]
func InitCategoryController() {
	categoryCollection = config.GetCollection("categories")
}

// Criar uma nova categoria
// @Summary Criar uma nova categoria
// @Description Criar uma nova categoria
// @Tags categories
// @Accept  json
// @Produce  json
// @Param category body models.PostCategory true "Category"
// @Success 201 {object} models.PostCategory
// @Router /categories [post]
func CreateCategory(c *gin.Context) {
	var category models.PostCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()

	_, err := categoryCollection.InsertOne(context.TODO(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar no banco"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// Listar todas as categorias
// @Summary Listar todas as categorias
// @Description Listar todas as categorias
// @Tags categories
// @Accept  json
// @Produce  json
// @Success 200 {array} models.PostCategory
// @Router /categories [get]
func GetCategories(c *gin.Context) {
	cursor, err := categoryCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var categories []models.PostCategory
	if err = cursor.All(context.TODO(), &categories); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, categories)
}

// Atualizar uma categoria por ID
// @Summary Atualizar uma categoria por ID
// @Description Atualizar uma categoria por ID
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path string true "ID da categoria"
// @Param category body models.PostCategory true "Category"
// @Success 200 {object} models.PostCategory
// @Router /categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	var updatedCategory models.PostCategory

	objID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = categoryCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"name": updatedCategory.Name}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar categoria"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categoria atualizada com sucesso"})
}

// Buscar categoria por ID
// @Summary Buscar categoria por ID
// @Description Buscar categoria por ID
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path string true "ID da categoria"
// @Success 200 {object} models.PostCategory
// @Router /categories/{id} [get]
func GetCategoryByID(c *gin.Context) {
	categoryID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var category models.PostCategory
	err = categoryCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&category)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoria não encontrada"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Deletar uma categoria por ID
// @Summary Deletar uma categoria por ID
// @Description Deletar uma categoria por ID
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path string true "ID da categoria"
// @Success 200 {string} string
// @Router /categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	_, err = categoryCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar categoria"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categoria deletada com sucesso"})
}
