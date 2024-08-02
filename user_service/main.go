package main  

import (  
	mongosh "auth/Mongosh"      
	redis "auth/Redis"          
	"auth/config"                
	"auth/internal/service"     
	"auth/internal/storage"       
	"auth/pkg"                  
	"log"                      
)  

func main() {  
	// Загружаем конфигурацию из файла .env  
	conf := config.Load("./")  

	// Инициализируем соединение с MongoDB  
	collDB, err := mongosh.InitMongo(conf)  
	if err != nil {  
		log.Fatal(err)
	}  

	// Создаем новый репозиторий для аутентификации на основе коллекции MongoDB  
	authRepo := storage.NewAuthRepo(collDB)  

	// Подключаемся к Redis  
	rds, err := redis.ConnectRedis(conf)  
	if err != nil {  
		log.Fatal(err) // Завершаем программу, если не удается подключиться к Redis  
	}  

	// Создаем новый сервис аутентификации, используя созданный репозиторий и подключение к Redis  
	service := service.NewService(authRepo, rds)  

	// Создаем сервис для копирования, передавая в него экземпляр сервиса аутентификации  
	cpyServ := pkg.NewCopyService(*service)  

	// Логируем информацию о запуске сервиса аутентификации  
	log.Printf("Auth Service Running on :%s port", conf.AuthServicePort)  

	// Запускаем сервис копирования и обрабатываем возможные ошибки  
	if err := cpyServ.RUN(conf); err != nil {  
		log.Fatal(err.Error()) 
	}  
}