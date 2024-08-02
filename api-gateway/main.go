package main  

import (  
	"api-gateway/api"  
	"api-gateway/config"  
	authdial "api-gateway/pkg/authDial"  
	devicedial "api-gateway/pkg/deviceDial"  
	"fmt"  
	"log"  
)  

// @title Your API Title  
// @version 1.0  
// @description Your API Description  
// @host localhost:9000  
// @BasePath /  

// @securityDefinitions.apikey BearerAuth  
// @in header  
// @name Authorization  
func main() {  
	// Загрузить конфигурацию  
	cfg := config.LOAD("./")  

	// Вызов службы аутентификации
	authDial, err := authdial.DialWithAuthService(cfg)  
	if err != nil {  
		log.Fatalf("Failed to dial auth service: %v", err)  
	}  

	// Звонок в службу поддержки устройства  
	deviceDial, err := devicedial.DialWithDeviceService(cfg)  
	if err != nil {  
		log.Fatalf("Failed to dial device service: %v", err)  
	}  

	// Создайте API с набираемыми услугами.
	r := api.NewGin(authDial, deviceDial)  

	// Запустите шлюз API 
	log.Printf("Starting API Gateway on port %s...", cfg.ApiGatewayServicePort)  
	if err := r.Run(fmt.Sprintf(":%s", cfg.ApiGatewayServicePort)); err != nil {  
		log.Fatalf("Failed to run API Gateway: %v", err)  
	}  
}