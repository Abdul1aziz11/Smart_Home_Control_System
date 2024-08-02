package pkg  

import (  
	"auth/config"  
	"auth/genproto"  
	"auth/internal/service"  
	"fmt"  
	"net"  

	"google.golang.org/grpc"  
)  

// CopyService определяет структуру для копии сервиса аутентификации.  
type CopyService struct {  
	service service.Service  
}  

// NewCopyService создает новый экземпляр CopyService.  
func NewCopyService(cpyServ service.Service) *CopyService {  
	return &CopyService{  
		service: cpyServ,  
	}  
}  

// RUN запускает gRPC сервер аутентификации.  
func (c *CopyService) RUN(cfg config.Config) error {  
	target := fmt.Sprintf("%s:%s", cfg.AuthServiceHost, cfg.AuthServicePort) // Определение цели для запуска сервера  
	listener, err := net.Listen("tcp", target)  
	if err != nil {  
		return fmt.Errorf("failed to listen on %s: %v", target, err)  
	}  

	newServ := grpc.NewServer()  

	// Регистрация службы аутентификации в gRPC сервере.  
	genproto.RegisterAuthServiceServer(newServ, &c.service)  

	// Запуск сервера и возврат ошибки, если возникла.  
	if err := newServ.Serve(listener); err != nil {  
		return fmt.Errorf("failed to serve: %v", err)  
	}  

	return nil  
}