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

var postCollection *mongo.Collection

// Inicializa a coleção depois que o banco de dados estiver conectado
func InitPostController() {
	postCollection = config.GetCollection("posts")
}

// Criar um novo post
func CreatePost(c *gin.Context) {
	if postCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Banco de dados não conectado!"})
		return
	}

	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now()

	_, err := postCollection.InsertOne(context.TODO(), post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar no banco"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// Listar todos os posts
func GetPosts(c *gin.Context) {
	if postCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Banco de dados não conectado!"})
		return
	}

	cursor, err := postCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var posts []models.Post
	if err = cursor.All(context.TODO(), &posts); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, posts)
}

// Buscar post por ID
func GetPostByID(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	err := postCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&post)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post não encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar post"})
		return
	}

	c.JSON(http.StatusOK, post)
}
