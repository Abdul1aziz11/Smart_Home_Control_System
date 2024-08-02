package storage

import (
	"context"
	"device-control/genproto"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *DeviceControlRepo) CreateDevice(ctx context.Context, req *genproto.CreateDeviceRequest) (*genproto.Device, error) {
	device := req.GetDevice()
	now := time.Now().Format(time.RFC3339)

	configMap := make(map[string]string)
	configMap["DeviceType"] = device.DeviceType

	switch device.DeviceType {
	case "light":
		configMap["Brightness"] = "100"
		configMap["Color"] = "white"
	case "tv":
		configMap["Volume"] = "100"
		configMap["Channel"] = "1"
	case "ac":
		configMap["Temperature"] = "30"
		configMap["Mode"] = "cool"
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported device type: %s", device.DeviceType)
	}

	deviceDoc := Device{
		UserID:                device.UserId,
		DeviceType:            device.DeviceType,
		DeviceName:            device.DeviceName,
		DeviceStatus:          device.DeviceStatus,
		ConfigurationSettings: configMap,
		LastUpdated:           now,
		Location:              device.Location,
		FirmwareVersion:       device.FirmwareVersion,
		ConnectivityStatus:    device.ConnectivityStatus,
	}

	result, err := d.coll.Device.InsertOne(ctx, deviceDoc)
	if err != nil {
		return nil, err
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	device.DeviceId = insertedID.Hex()
	device.LastUpdated = now
	device.ConfigurationSettings = configMap

	return device, nil
}

func (d *DeviceControlRepo) UpdateDevice(ctx context.Context, req *genproto.UpdateDeviceRequest) (*genproto.Device, error) {
	device := req.Device
	now := time.Now().Format(time.RFC3339)

	deviceID, err := primitive.ObjectIDFromHex(device.DeviceId)
	if err != nil {
		return nil, fmt.Errorf("invalid device ID: %v", err)
	}

	// Создание карты конфигурации
	configMap := make(map[string]string)
	configMap["DeviceType"] = device.DeviceType

	// Заполнение карты конфигов
	switch device.DeviceType {
	case "light":
		if brightness, ok := req.Device.ConfigurationSettings["brightness"]; ok {
			configMap["Brightness"] = brightness
		}
		if color, ok := req.Device.ConfigurationSettings["color"]; ok {
			configMap["Color"] = color
		}
	case "tv":
		if volume, ok := req.Device.ConfigurationSettings["volume"]; ok {
			configMap["Volume"] = volume
		}
		if channel, ok := req.Device.ConfigurationSettings["channel"]; ok {
			configMap["Channel"] = channel
		}
	case "ac":
		if temperature, ok := req.Device.ConfigurationSettings["temperature"]; ok {
			configMap["Temperature"] = temperature
		}
		if mode, ok := req.Device.ConfigurationSettings["mode"]; ok {
			configMap["Mode"] = mode
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported device type: %s", device.DeviceType)
	}

	// Обновить документ
	updateDoc := bson.M{
		"$set": bson.M{
			"user_id":             device.UserId,
			"device_type":         device.DeviceType,
			"device_name":         device.DeviceName,
			"device_status":       device.DeviceStatus,
			"last_updated":        now,
			"location":            device.Location,
			"firmware_version":    device.FirmwareVersion,
			"connectivity_status": device.ConnectivityStatus,
		},
	}

	// Добавить в updateDoc, если configMap заполнен.
	if len(configMap) > 1 { // Также добавлен DeviceType, поэтому должно быть> 1.
		updateDoc["$set"].(bson.M)["configuration_settings"] = configMap
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedDevice Device
	err = d.coll.Device.FindOneAndUpdate(ctx, bson.M{"_id": deviceID}, updateDoc, opts).Decode(&updatedDevice)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("device with id %s not found", device.DeviceId)
		}
		return nil, err
	}

	return &genproto.Device{
		DeviceId:              updatedDevice.ID.Hex(),
		UserId:                updatedDevice.UserID,
		DeviceType:            updatedDevice.DeviceType,
		DeviceName:            updatedDevice.DeviceName,
		DeviceStatus:          updatedDevice.DeviceStatus,
		ConfigurationSettings: updatedDevice.ConfigurationSettings,
		LastUpdated:           updatedDevice.LastUpdated,
		Location:              updatedDevice.Location,
		FirmwareVersion:       updatedDevice.FirmwareVersion,
		ConnectivityStatus:    updatedDevice.ConnectivityStatus,
	}, nil
}

func (d *DeviceControlRepo) DeleteDevice(ctx context.Context, req *genproto.DeleteDeviceRequest) (*genproto.DeleteDeviceResponse, error) {
	deviceID, err := primitive.ObjectIDFromHex(req.DeviceId)
	if err != nil {
		return nil, fmt.Errorf("invalid device ID: %v", err)
	}

	result, err := d.coll.Device.DeleteOne(ctx, bson.M{"_id": deviceID})
	if err != nil {
		return nil, err
	}
	if result.DeletedCount == 0 {
		return nil, fmt.Errorf("device with id %s not found", req.DeviceId)
	}

	return &genproto.DeleteDeviceResponse{
		Success: true,
	}, nil
}

func (d *DeviceControlRepo) GetDevice(ctx context.Context, req *genproto.GetDeviceRequest) (*genproto.Device, error) {
	deviceID, err := primitive.ObjectIDFromHex(req.DeviceId)
	if err != nil {
		return nil, fmt.Errorf("invalid device ID: %v", err)
	}

	var device Device
	err = d.coll.Device.FindOne(ctx, bson.M{"_id": deviceID}).Decode(&device)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("device with id %s not found", req.DeviceId)
		}
		return nil, fmt.Errorf("error finding device: %v", err)
	}

	return &genproto.Device{
		DeviceId:              device.ID.Hex(),
		UserId:                device.UserID,
		DeviceType:            device.DeviceType,
		DeviceName:            device.DeviceName,
		DeviceStatus:          device.DeviceStatus,
		ConfigurationSettings: device.ConfigurationSettings,
		LastUpdated:           device.LastUpdated,
		Location:              device.Location,
		FirmwareVersion:       device.FirmwareVersion,
		ConnectivityStatus:    device.ConnectivityStatus,
	}, nil
}

func (d *DeviceControlRepo) SendMessage(ctx context.Context, req *genproto.Command) (*genproto.CommandResponse, error) {
	now := time.Now().Format(time.RFC3339)
	_, err := d.coll.Command.InsertOne(ctx, bson.M{
		"device_id":       req.DeviceId,
		"user_id":         req.UserId,
		"command_type":    req.CommandType,
		"command_payload": req.CommandPayload,
		"timestamp":       now,
		"status":          req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &genproto.CommandResponse{
		Status: "Saved",
	}, nil
}
