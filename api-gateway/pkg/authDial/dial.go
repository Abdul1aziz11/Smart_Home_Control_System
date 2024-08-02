package authdial

import (
	"api-gateway/config"
	"api-gateway/genproto"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthServiceClient обертка для клиента AuthService
type AuthServiceClient struct {
	genproto.AuthServiceClient
}

// DialWithAuthService создает соединение с сервисом AuthService
func DialWithAuthService(cfg config.Config) (*AuthServiceClient, error) {
	// Формирование адреса для подключения
	target := fmt.Sprintf("%s:%s", cfg.AuthServiceHost, cfg.AuthServicePort)

	// Создание gRPC соединения с использованием незащищенных (insecure) учетных данных
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Создание клиента AuthService
	authService := genproto.NewAuthServiceClient(conn)

	return &AuthServiceClient{AuthServiceClient: authService}, nil
}
