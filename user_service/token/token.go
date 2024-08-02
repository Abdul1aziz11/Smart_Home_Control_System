package token  

import (  
	"fmt"  
	"time"  

	"github.com/golang-jwt/jwt/v5"  
)  

// MyClaims представляет собой пользовательские данные, которые будут включены в токен.  
type MyClaims struct {  
	ID       string `json:"id"`  
	Email    string `json:"email"`  
	Password string `json:"password"`  
	jwt.RegisteredClaims  
}  

// Секретный ключ для подписи токенов.  
var jwtKey = []byte("your-secret-key")  

// GenerateToken генерирует JWT токен с заданными данными пользователя.  
func GenerateToken(id, email, password string) (string, error) {  
	// Создаем полезную нагрузку (claims) для токена.  
	claims := MyClaims{  
		ID:       id,  
		Email:    email,  
		Password: password,  
		RegisteredClaims: jwt.RegisteredClaims{  
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Срок действия токена  
			Issuer:    "your-app-name",                                    // Эмитент токена  
		},  
	}  

	// Создаем новый токен с заданными claims.  
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)  

	// Подписываем токен и получаем строку токена.  
	tokenString, err := token.SignedString(jwtKey)  
	if err != nil {  
		return "", err  
	}  

	return tokenString, nil  
}  

// ParseToken проверяет и анализирует токен, возвращая данные пользователя.  
func ParseToken(tokenString string) (*MyClaims, error) {  
	// Проверяем и разбираем токен.  
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {  
		return jwtKey, nil  
	})  
	if err != nil {  
		return nil, err  
	}  

	// Проверяем, действителен ли токен и можем ли мы получить claims.  
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {  
		return claims, nil  
	}  

	return nil, fmt.Errorf("invalid token")  
}