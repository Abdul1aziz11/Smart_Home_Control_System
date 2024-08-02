package mongosh  

import (  
	"auth/config"  
	"context"  
	"fmt"  
	"log"  

	"go.mongodb.org/mongo-driver/mongo"  
	"go.mongodb.org/mongo-driver/mongo/options"  
)  

// InitMongo инициализирует соединение с MongoDB и возвращает коллекцию.  
func InitMongo(cfg config.Config) (*mongo.Collection, error) {  
	// Создание контекста для работы с MongoDB.  
	ctx := context.Background()  

	// Формирование URI для подключения к MongoDB.  
	uri := fmt.Sprintf("mongodb://%s:%d", "mongo", 27017)  

	clientOptions := options.Client().  
		ApplyURI(uri).  
		SetAuth(options.Credential{  
			Username: "Bek10022006",  
			Password: "Bek10022006",  
		})  

	// Подключение к клиенту MongoDB.  
	client, err := mongo.Connect(ctx, clientOptions)  
	if err != nil {  
		log.Fatalf("Failed to connect to MongoDB: %v", err) // Логирование фатальной ошибки подключения  
	}  

	// Проверка подключения к MongoDB.  
	if err := client.Ping(ctx, nil); err != nil {  
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err) // Возврат ошибки при неудаче  
	}  

	// Возвращаем коллекцию из указанной базы данных.  
	collection := client.Database(cfg.Mongo.MongoDatabase).Collection(cfg.Mongo.MongoCollection)  

	return collection, nil  
}