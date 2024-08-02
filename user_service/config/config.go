package config  

import (  
	"github.com/joho/godotenv" 
	"github.com/spf13/viper"   
)  

// Структура для хранения настроек подключения к MongoDB  
type Mongo struct {  
	MongoUser       string   
	MongoPassword   string 
	MongoDatabase   string  
	MongoCollection string  
}  

// Структура для хранения настроек подключения к Redis  
type Redis struct {  
	RedisHost string 
	RedisPort string 
}  

// Основная структура конфигурации  
type Config struct {  
	Mongo           Mongo  
	Redis           Redis   
	AuthServiceHost string
	AuthServicePort string 
}  

// Функция для загрузки конфигурации из файла окружения  
func Load(path string) Config {  
	// Загружает переменные окружения из файла .env  
	godotenv.Load(path + "/.env")  
	  
	// Создаёт новый экземпляр viper для работы с конфигурацией  
	cfg := viper.New()  
	  
	// Автоматически загружает значения из переменных окружения  
	cfg.AutomaticEnv()  

	// Создаёт конфигурацию, заполняя её значениями из окружения  
	conf := Config{  
		Mongo: Mongo{  
			MongoUser:       cfg.GetString("MONGODB_USERNAME"), 
			MongoPassword:   cfg.GetString("MONGODB_PASSWORD"), 
			MongoDatabase:   cfg.GetString("MONGODB_DATABASE"), 
			MongoCollection: cfg.GetString("MONGODB_COLLECTION"), 
		},  
		Redis: Redis{  
			RedisHost: cfg.GetString("REDIS_HOST"), 
			RedisPort: cfg.GetString("REDIS_PORT"),
		},  
		AuthServiceHost: cfg.GetString("AUTH_SERVICE_HOST"), 
		AuthServicePort: cfg.GetString("AUTH_SERVICE_PORT"), 
	}  

	return conf // Возвращает заполненную структуру конфигурации  
}