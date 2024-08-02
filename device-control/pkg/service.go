package pkg  

import (  
	"device-control/config"  
	"device-control/genproto"  
	"device-control/internal/service"  
	"fmt"  
	"net"  

	"google.golang.org/grpc"  
)  

// CopyService представляет сервис для обработки копирования устройств.  
type CopyService struct {  
	service service.Service  
}  

// NewCopyService создаёт новый экземпляр CopyService.  
func NewCopyService(service service.Service) *CopyService {  
	return &CopyService{  
		service: service,  
	}  
}  

// RUN запускает GRPC сервер с указанной конфигурацией.  
func (c *CopyService) RUN(cfg config.Config) error {  
	// Создание целевого адреса из конфигурации  
	target := fmt.Sprintf("%s:%s", cfg.DeviceControlHost, cfg.DeviceControlPort)  

	// Создание TCP listener  
	listener, err := net.Listen("tcp", target)  
	if err != nil {  
		return fmt.Errorf("failed to listen on %s: %w", target, err)  
	}  
	defer listener.Close() // Закрытие listener при выходе из функции  

	// Создание нового GRPC сервера  
	newServer := grpc.NewServer()  

	// Регистрация сервиса на GRPC сервере  
	genproto.RegisterDeviceServiceServer(newServer, &c.service)  

	// Запуск GRPC сервера  
	if err := newServer.Serve(listener); err != nil {  
		return fmt.Errorf("failed to serve gRPC: %w", err)  
	}  

	return nil  
}