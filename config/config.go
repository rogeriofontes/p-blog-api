package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

const dbName = "blog"

// Nome das coleções que queremos garantir que existem
var collections = []string{"users", "categories", "posts", "comments", "tags", "reactions", "favorites", "followers"}

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

func EnsuseDatabaseAndCollections() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, collection := range collections {
		col := DB.Database(dbName).Collection(collection)

		// Insere um documento temporário para garantir a criação
		_, err := col.InsertOne(ctx, bson.M{"init": true})
		if err != nil {
			log.Printf("Erro ao inicializar a coleção %s: %v", collection, err)
		} else {
			log.Printf("Coleção %s inicializada com sucesso!", collection)
		}

		// Remove o documento temporário
		_, _ = col.DeleteOne(ctx, bson.M{"init": true})
	}
}

// Retorna uma coleção do banco de dados
func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		log.Println("A conexão com o banco de dados não foi inicializada. Tentando conectar...")
		ConnectDatabase()
	}
	return DB.Database(dbName).Collection(collectionName)
}
