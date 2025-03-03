package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDatabase() {
	mongoURI := os.Getenv("MONGO_URI")                   // Pegando do environment
	clientOptions := options.Client().ApplyURI(mongoURI) //mongodb://admin:admin@localhost:27017

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verifica se o MongoDB está acessível
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Erro ao pingar o MongoDB: %v", err)
	}

	fmt.Println("🔥 Conectado ao MongoDB!")
	DB = client
}

// Retorna uma coleção do banco de dados
func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		log.Println("A conexão com o banco de dados não foi inicializada. Tentando conectar...")
		ConnectDatabase()
	}
	return DB.Database("blog").Collection(collectionName)
}
