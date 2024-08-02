package service

import (
	"auth/genproto"
	"auth/internal/storage"
	"auth/token"
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"
)

type Service struct {
	genproto.UnimplementedAuthServiceServer
	authRepo *storage.AuthRepo
	rds      *redis.Client
}

// NewService создает новый экземпляр Service.
func NewService(authRepo *storage.AuthRepo, rds *redis.Client) *Service {
	return &Service{
		authRepo: authRepo,
		rds:      rds,
	}
}

// SignUp регистрирует нового пользователя.
func (s *Service) SignUp(ctx context.Context, req *genproto.UserReq) (*genproto.UserResp, error) {
	// Подготовка данных пользователя для Redis.
	user := map[string]interface{}{
		"username":   req.Username,
		"email":      req.Email,
		"password":   req.Password,
		"first_name": req.Profile.FirstName,
		"last_name":  req.Profile.LastName,
		"address":    req.Profile.Address,
	}

	// Установка данных пользователя в Redis.
	if err := s.rds.HSet(ctx, req.Email, user).Err(); err != nil {
		return nil, fmt.Errorf("failed to set user data in Redis: %v", err)
	}

	// Установка времени истечения для данных пользователя.
	if err := s.rds.Expire(ctx, req.Email, 90*time.Second).Err(); err != nil {
		return nil, fmt.Errorf("failed to set expiration time for user data: %v", err)
	}

	// Генерация и установка ключа истечения.
	expKey := GenRandomExpiredKey()
	if err := s.rds.Set(ctx, req.Email+"ExpiredKey", expKey, 60*time.Second).Err(); err != nil {
		return nil, fmt.Errorf("failed to set expiration key in Redis: %v", err)
	}

	// Возвращение успешного ответа.
	return &genproto.UserResp{
		Status: "User data stored successfully",
	}, nil
}

// Verify проверяет пользователя по предоставленным данным.
func (s *Service) Verify(ctx context.Context, req *genproto.UserReq) (*genproto.UserResp, error) {
	// Получение ключа истечения для пользователя.
	expiredKey, err := s.rds.Get(ctx, req.Email+"ExpiredKey").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("expired key not found")
		}
		return nil, err
	}

	// Проверка правильности пароля.
	if expiredKey != req.Password {
		return &genproto.UserResp{
			Status: "Verify Password incorrect",
		}, nil
	}

	// Получение данных пользователя из Redis.
	dataUserInRedis, err := s.rds.HGetAll(ctx, req.Email).Result()
	if err != nil {
		return nil, err
	}
	if len(dataUserInRedis) == 0 {
		return nil, fmt.Errorf("data not found in redis hash")
	}

	// Хеширование пароля пользователя.
	hashPassword := base64.StdEncoding.EncodeToString([]byte(dataUserInRedis["password"]))

	dataUser := genproto.UserReq{
		Username:     dataUserInRedis["username"],
		Email:        dataUserInRedis["email"],
		PasswordHash: hashPassword,
		Profile: &genproto.Profile{
			FirstName: dataUserInRedis["first_name"],
			LastName:  dataUserInRedis["last_name"],
			Address:   dataUserInRedis["address"],
		},
	}

	// Создание пользователя в базе данных.
	resp, err := s.authRepo.CreateUser(ctx, &dataUser)
	if err != nil {
		return nil, err
	}

	return &genproto.UserResp{
		UserId: resp.UserId,
		Status: resp.Status,
	}, nil
}

// SignIn выполняет вход пользователя.
func (s *Service) SignIn(ctx context.Context, req *genproto.UserReq) (*genproto.UserResp, error) {
	resp, err := s.authRepo.GetUserByFilter(ctx, req)
	if err != nil {
		return nil, err
	}

	// Проверка правильности пароля.
	passwordHash := base64.StdEncoding.EncodeToString([]byte(req.Password))
	if resp.Users[0].PasswordHash != passwordHash {
		return nil, fmt.Errorf("Login yoki parol xato")
	}

	// Генерация токена для пользователя.
	tkn, err := token.GenerateToken(resp.Users[0].UserId, resp.Users[0].Email, resp.Users[0].PasswordHash)
	if err != nil {
		return nil, err
	}

	// Установка токена в Redis.
	if err := s.rds.Set(ctx, req.Email+"jwtToken", tkn, time.Hour).Err(); err != nil {
		return nil, err
	}

	return &genproto.UserResp{
		UserId: resp.Users[0].UserId,
		Status: "Login bajarildi",
		Email:  resp.Users[0].Email,
		Token:  tkn,
	}, nil
}

// Profile возвращает профиль пользователя.
func (s *Service) Profile(ctx context.Context, req *genproto.UserReq) (*genproto.UserReq, error) {
	result, err := s.rds.Get(ctx, req.Email+"jwtToken").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("Hech qanday malumot topilmadi")
		}
		return nil, err
	}

	// Парсинг токена.
	tokenParse, err := token.ParseToken(result)
	if err != nil {
		return nil, err
	}

	return &genproto.UserReq{
		UserId:       tokenParse.ID,
		Email:        tokenParse.Email,
		PasswordHash: tokenParse.Password,
	}, nil
}

// GetUserById возвращает пользователя по ID.
func (s *Service) GetUserById(ctx context.Context, req *genproto.UserReq) (*genproto.UserReq, error) {
	resp, err := s.authRepo.GetUserById(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GenRandomExpiredKey генерирует случайный ключ истечения.
func GenRandomExpiredKey() int {
	rand.Seed(uint64(time.Now().UnixNano()))
	return rand.Intn(9000) + 1000
}
