package storage  

import (  
	mongosh "device-control/Mongosh"  
	"go.mongodb.org/mongo-driver/bson/primitive"  
)  

// DeviceControlRepo управляет взаимодействием с коллекцией устройств в базе данных  
type DeviceControlRepo struct {  
	coll mongosh.Collections  
}  

// Device представляет устройство с его атрибутами  
type Device struct {  
	ID                    primitive.ObjectID `bson:"_id,omitempty"`  
	UserID                string             `bson:"user_id"`  
	DeviceType            string             `bson:"device_type"`  
	DeviceName            string             `bson:"device_name"`  
	DeviceStatus          string             `bson:"device_status"`  
	ConfigurationSettings map[string]string  `bson:"configuration_settings"`  
	LastUpdated           string             `bson:"last_updated"`  
	Location              string             `bson:"location"`  
	FirmwareVersion       string             `bson:"firmware_version"`  
	ConnectivityStatus    string             `bson:"connectivity_status"`  
}  

// DeviceTypes определяет тип устройства  
type DeviceTypes struct {  
	DeviceType interface{}  
}  

// Light представляет настройки для светильника  
type Light struct {  
	Brightness string  
	Color      string  
}  

// TV представляет настройки для телевизора  
type TV struct {  
	Volume  string  
	Channel string  
}  

// AC представляет настройки для кондиционера  
type AC struct {  
	Temperature string  
	Mode        string  
}  

// NewDeviceControlRepo создает новый экземпляр DeviceControlRepo с заданной коллекцией  
func NewDeviceControlRepo(coll mongosh.Collections) *DeviceControlRepo {  
	return &DeviceControlRepo{coll: coll}  
}