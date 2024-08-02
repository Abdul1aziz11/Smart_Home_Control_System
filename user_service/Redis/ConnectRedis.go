package redis  

import (  
	"auth/config"  
	"context"  
	"fmt"  

	"github.com/redis/go-redis/v9"  
)  

// ConnectRedis устанавливает соединение с Redis-сервером на основе конфигурации.  
func ConnectRedis(cfg config.Config) (*redis.Client, error) {  
	// Создаем новый клиент Redis с использованием параметров из конфигурации.  
	options := &redis.Options{  
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.RedisHost, cfg.Redis.RedisPort), // Форматируем адрес  
	}  
	client := redis.NewClient(options)  

	// Пингуем сервер Redis для проверки соединения.  
	if err := client.Ping(context.Background()).Err(); err != nil {  
		return nil, fmt.Errorf("could not connect to redis: %w", err) // Возвращаем ошибку с уточнением  
	}  

	return client, nil // Возвращаем клиент, если все прошло успешно  
}