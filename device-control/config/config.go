package config  

import (  
	"github.com/joho/godotenv" 
	"github.com/spf13/viper"   
)  

// Структура для хранения конфигурации MongoDB  
type Mongo struct {  
	MongoUser              string 
	MongoPassword          string
	MongoDB                string  
	MongoCollectionDevice  string 
	MongoCollectionCommand string 
}  

// Главная структура конфигурации, содержащая все необходимые параметры  
type Config struct {  
	Mongo             Mongo  
	DeviceControlHost string   
	DeviceControlPort string   
}  

// Функция для загрузки конфигурации из .env файла и возврата структуры Config  
func LOAD(path string) Config {  
	godotenv.Load(path + "/.env") 
	cfg := viper.New()              
	cfg.AutomaticEnv()              

	// Создаем экземпляр Config и заполняем его значениями из окружения  
	conf := Config{  
		Mongo: Mongo{  
			MongoUser:              cfg.GetString("MONGODB_USER"),             
			MongoPassword:          cfg.GetString("MONGODB_PASSWORD"),          
			MongoDB:                cfg.GetString("MONGODB_DATABASE"),         
			MongoCollectionDevice:  cfg.GetString("MONGODB_COLLECTION_DEVICE"),  
			MongoCollectionCommand: cfg.GetString("MONGODB_COLLECTION_COMMAND"), 
		},  

		DeviceControlHost: cfg.GetString("DEVICE_SERVICE_HOST"), 
		DeviceControlPort: cfg.GetString("DEVICE_SERVICE_PORT"), 
	}  

	return conf 
}