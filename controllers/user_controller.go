package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/rogeriofontes/p-blog-api/config"
	"github.com/rogeriofontes/p-blog-api/models"
	"github.com/rogeriofontes/p-blog-api/utils"

	// Importamos o utilitário
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

// Inicializa a coleção de categorias após conectar ao banco
// @Summary Inicializa a coleção de usuários
// @Description Inicializa a coleção de usuários após conectar ao banco
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {string} string "Usuários inicializados"
// @Router /users/init [get]
func InitUserController() {
	userCollection = config.GetCollection("users")
}

// Criar um novo usuário
// @Summary Criar um novo usuário
// @Description Criar um novo usuário
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "Usuário a ser criado"
// @Success 201 {object} models.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash da senha antes de salvar
	user.Password = utils.HashPassword(user.Password)
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()

	// Verificar se o email já existe
	count, err := userCollection.CountDocuments(context.TODO(), bson.M{"email": user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar usuário"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email já cadastrado"})
		return
	}

	// Inserir usuário no banco
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	// Retorna sem expor a senha
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// Buscar um usuário pelo ID
// @Summary Buscar um usuário pelo ID
// @Description Buscar um usuário pelo ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} models.User
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	// Não exibe a senha na resposta
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// Listar todos os usuários (sem mostrar a senha)
// @Summary Listar todos os usuários
// @Description Listar todos os usuários
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} string
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	cursor, err := userCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários"})
		return
	}

	var users []models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar usuários"})
		return
	}

	// Remover as senhas antes de exibir
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

// Atualizar um usuário por ID
// @Summary Atualizar um usuário por ID
// @Description Atualizar um usuário por ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param user body models.User true "Usuário a ser atualizado"
// @Success 200 {string} string "Usuário atualizado com sucesso"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var updatedUser models.User

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Bind dos dados do usuário
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Criar mapa de atualização
	updateData := bson.M{
		"username": updatedUser.Username,
		"email":    updatedUser.Email,
	}

	// Se o usuário enviou uma nova senha, vamos hasheá-la antes de salvar
	if updatedUser.Password != "" {
		updateData["password"] = utils.HashPassword(updatedUser.Password)
	}

	_, err = userCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updateData},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar usuário"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário atualizado com sucesso"})
}

// Deletar um usuário por ID
// @Summary Deletar um usuário por ID
// @Description Deletar um usuário por ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {string} string "Usuário deletado com sucesso"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	_, err = userCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar usuário"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário deletado com sucesso"})
}

// Login autentica um usuário e gera um token JWT
// @Summary Login de usuário
// @Description Autentica um usuário e gera um token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body struct { email string; password string } true "Credenciais de login"
// @Success 200 {string} string "Token
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /login [post]
func Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Buscar usuário pelo email
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": request.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado"})
		return
	}

	// Verificar senha (supondo que esteja em hash)
	if user.Password != utils.HashPassword(request.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Senha incorreta"})
		return
	}

	// Gerar token JWT
	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
