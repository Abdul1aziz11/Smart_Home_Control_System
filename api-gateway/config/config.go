package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config структура для хранения конфигурационных данных
type Config struct {
	ApiGatewayServiceHost string
	ApiGatewayServicePort string

	AuthServiceHost string
	AuthServicePort string

	DeviceServiceHost string
	DeviceServicePort string
}

// LOAD загружает конфигурацию из .env файла и среды выполнения
func LOAD(path string) Config {
	// Загружаем переменные окружения из .env файла
	godotenv.Load(path + "/.env")
	
	// Создаем новый экземпляр viper
	cfg := viper.New()
	
	// Автоматически считываем переменные окружения
	cfg.AutomaticEnv()

	// Инициализируем конфигурацию из переменных окружения
	conf := Config{
		ApiGatewayServiceHost: cfg.GetString("API_GATEWAY_SERVICE_HOST"),
		ApiGatewayServicePort: cfg.GetString("API_GATEWAY_SERVICE_PORT"),

		AuthServiceHost: cfg.GetString("AUTH_SERVICE_HOST"),
		AuthServicePort: cfg.GetString("AUTH_SERVICE_PORT"),

		DeviceServiceHost: cfg.GetString("DEVICE_SERVICE_HOST"),
		DeviceServicePort: cfg.GetString("DEVICE_SERVICE_PORT"),
	}

	return conf
}
