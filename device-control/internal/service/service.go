package service  

import (  
	"context"                  
	pb "device-control/genproto" 
	"device-control/internal/storage" 
	"fmt"                        
)  

// Service представляет собой структуру, реализующую интерфейс сервера устройства  
type Service struct {  
	pb.UnimplementedDeviceServiceServer 
	DeviceControlREPO *storage.DeviceControlRepo
}  

// NewService создает новый экземпляр Service с заданным репозиторием  
func NewService(device *storage.DeviceControlRepo) *Service {  
	return &Service{  
		DeviceControlREPO: device, 
	}  
}  

// CreateDevice обрабатывает создание нового устройства  
func (d *Service) CreateDevice(ctx context.Context, req *pb.CreateDeviceRequest) (*pb.Device, error) {  
	return d.DeviceControlREPO.CreateDevice(ctx, req) 
}  

// GetDevice обрабатывает запрос на получение устройства по ID  
func (d *Service) GetDevice(ctx context.Context, req *pb.GetDeviceRequest) (*pb.Device, error) {  
	return d.DeviceControlREPO.GetDevice(ctx, req)
}  

// UpdateDevice обрабатывает обновление данных устройства  
func (d *Service) UpdateDevice(ctx context.Context, req *pb.UpdateDeviceRequest) (*pb.Device, error) {  
	return d.DeviceControlREPO.UpdateDevice(ctx, req) 
}  

// DeleteDevice обрабатывает запрос на удаление устройства  
func (d *Service) DeleteDevice(ctx context.Context, req *pb.DeleteDeviceRequest) (*pb.DeleteDeviceResponse, error) {  
	return d.DeviceControlREPO.DeleteDevice(ctx, req)
}  

// SendCommand обрабатывает команды, отправляемые устройствам  
func (d *Service) SendCommand(ctx context.Context, req *pb.Command) (*pb.CommandResponse, error) {  
	// Получаем устройство по ID  
	device, err := d.DeviceControlREPO.GetDevice(ctx, &pb.GetDeviceRequest{  
		DeviceId: req.GetDeviceId(),  
	})  
	if err != nil {  
		return nil, err
	}  

	// Устанавливаем параметры устройства в зависимости от его типа  
	switch device.DeviceType {  
	case "light":  
		device.ConfigurationSettings["brightness"] = req.CommandPayload["brightness"] 
		device.ConfigurationSettings["color"] = req.CommandPayload["color"]        
	case "tv":  
		device.ConfigurationSettings["volume"] = req.CommandPayload["volume"]       
		device.ConfigurationSettings["channel"] = req.CommandPayload["channel"]     
	case "ac":  
		device.ConfigurationSettings["temperature"] = req.CommandPayload["temperature"]
		device.ConfigurationSettings["mode"] = req.CommandPayload["mode"]               
	default:  
		return nil, fmt.Errorf("unknown device type: %s", device.DeviceType) 
	}  

	// Обновляем устройство в репозитории  
	_, err = d.DeviceControlREPO.UpdateDevice(ctx, &pb.UpdateDeviceRequest{  
		Device: device,  
	})  
	if err != nil {  
		return nil, err 
	}  

	// Отправляем команду устройству и возвращаем ответ  
	return d.DeviceControlREPO.SendMessage(ctx, req)  
}