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
func InitCategoryController() {
	categoryCollection = config.GetCollection("categories")
}

// Criar uma nova categoria
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
