package controllers

import (
	"context"
	"fmt"
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
// @Summary Inicializa a coleção de posts
// @Description Inicializa a coleção de posts após conectar ao banco
// @Tags posts
// @Accept json
// @Produce json
// @Router /init/posts [get]
func InitPostController() {
	postCollection = config.GetCollection("posts")
}

// Criar um novo post
// CreatePost cria um novo post
// @Summary Criar um post
// @Description Cria um novo post e salva no banco de dados
// @Tags posts
// @Accept  json
// @Produce  json
// @Param post body models.Post true "Dados do post"
// @Success 201 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts [post]
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

// Editar um post
// UpdatePost atualiza um post existente
// @Summary Atualizar um post
// @Description Atualiza um post existente no banco de dados
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path string true "ID do post"
// @Param post body models.Post true "Dados do post"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [put]
func UpdatePost(c *gin.Context) {
	postID := c.Param("id")
	var updatedPost models.Post

	// Buscar post original
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se o post existe
	var existingPost models.Post
	err = postCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existingPost)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post não encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualizar no banco
	_, err = postCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updatedPost},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post atualizado com sucesso"})
}

// Listar todos os posts
// GetPosts lista todos os posts
// @Summary Listar posts
// @Description Retorna uma lista de posts do banco de dados
// @Tags posts
// @Produce  json
// @Success 200 {array} models.Post
// @Failure 500 {object} models.ErrorResponse
// @Router /posts [get]
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
// GetPostByID busca um post por ID
// @Summary Buscar post por ID
// @Description Retorna um post específico do banco de dados
// @Tags posts
// @Produce  json
// @Param id path string true "ID do post"
// @Success 200 {object} models.Post
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [get]
func GetPostByID(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("ID recebido:", id) // Log para debug

	objID, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var post models.Post
	err := postCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&post)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post não encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// Buscar posts por categoria
// GetPostsByCategory busca posts por categoria
// @Summary Buscar posts por categoria
// @Description Retorna uma lista de posts de uma categoria específica
// @Tags posts
// @Produce  json
// @Param category_id path string true "ID da categoria"
// @Success 200 {array} models.Post
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/category/{category_id} [get]
func GetPostsByCategory(c *gin.Context) {
	categoryID := c.Param("category_id")

	// Converter para ObjectID
	objID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	cursor, err := postCollection.Find(context.TODO(), bson.M{"category_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar posts"})
		return
	}

	var posts []models.Post
	if err = cursor.All(context.TODO(), &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// Deletar um post por ID
// DeletePost deleta um post por ID
// @Summary Deletar post por ID
// @Description Deleta um post específico do banco de dados
// @Tags posts
// @Param id path string true "ID do post"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [delete]
func DeletePost(c *gin.Context) {
	postID := c.Param("id")

	// Converter para ObjectID
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se o post existe
	var existingPost bson.M
	err = postCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existingPost)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post não encontrado"})
		return
	}

	// Deletar o post
	_, err = postCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deletado com sucesso"})
}
