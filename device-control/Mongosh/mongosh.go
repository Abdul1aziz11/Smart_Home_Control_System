package mongosh  

import (  
	"context"  
	"device-control/config"  
	"fmt"  
	"log"  

	"go.mongodb.org/mongo-driver/mongo"  
	"go.mongodb.org/mongo-driver/mongo/options"  
)  

// Collections содержит коллекции базы данных  
type Collections struct {  
	Device  *mongo.Collection  
	Command *mongo.Collection  
}  

// InitMongo инициализирует подключение к MongoDB и возвращает коллекции  
func InitMongo(cfg config.Config) (*Collections, error) {  
	ctx := context.Background()  

	// Формируем URI подключения  
	uri := fmt.Sprintf("mongodb://%s:%d", "mongo", 27017)  

	// Настройки клиента MongoDB с аутентификацией  
	clientOptions := options.Client().ApplyURI(uri).  
		SetAuth(options.Credential{  
			Username: "Bek10022006",  
			Password: "Bek10022006",  
		})  

	client, err := mongo.Connect(ctx, clientOptions)  
	if err != nil {  
		log.Fatalf("Failed to connect to MongoDB: %v", err)  
	}  

	// Проверяем подключение к базе данных  
	if err := client.Ping(ctx, nil); err != nil {  
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)  
	}  

	// Устанавливаем коллекции  
	collectionDevice := client.Database(cfg.Mongo.MongoDB).Collection(cfg.Mongo.MongoCollectionDevice)  
	collectionCommand := client.Database(cfg.Mongo.MongoDB).Collection(cfg.Mongo.MongoCollectionCommand)  

	return &Collections{  
		Device:  collectionDevice,  
		Command: collectionCommand,  
	}, nil  
}