package devicedial

import (
	"api-gateway/config"
	"api-gateway/genproto"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DeviceServiceClient обертка для клиента DeviceService
type DeviceServiceClient struct {
	genproto.DeviceServiceClient
}

// DialWithDeviceService создает соединение с сервисом DeviceService
func DialWithDeviceService(cfg config.Config) (*DeviceServiceClient, error) {
	// Формирование адреса для подключения
	target := fmt.Sprintf("%s:%d", cfg.DeviceServiceHost, cfg.DeviceServicePort)

	// Создание gRPC соединения с использованием незащищенных (insecure) учетных данных
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Создание клиента DeviceService
	deviceService := genproto.NewDeviceServiceClient(conn)

	return &DeviceServiceClient{DeviceServiceClient: deviceService}, nil
}
