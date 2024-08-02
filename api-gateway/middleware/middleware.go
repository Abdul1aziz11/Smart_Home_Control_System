package middleware

import (
	"api-gateway/genproto"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("your-secret-key")

// TokenAuthMiddleware проверка токена с использованием middleware
func TokenAuthMiddleware(authServ genproto.AuthServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Получение заголовка Authorization
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			ctx.Abort()
			return
		}

		// Удаление префикса Bearer
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Проверка валидности токена
		if !isValidToken(tokenString) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// Разбор токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error parsing token"})
			ctx.Abort()
			return
		}

		// Получение claims из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		// Получение User ID из claims
		userID, ok := claims["id"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}
		email, ok := claims["email"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}
		password, ok := claims["password"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		// Получение информации о пользователе из базы данных
		user, err := authServ.GetUserById(context.Background(), &genproto.UserReq{
			UserId: userID,
		})
		if err != nil || user == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			ctx.Abort()
			return
		}

		// Проверка email и пароля пользователя
		if user.Email != email || user.PasswordHash != password {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			ctx.Abort()
			return
		}

		// Добавление информации о пользователе в контекст (если необходимо)
		ctx.Set("user", user)

		// Продолжение выполнения запроса
		ctx.Next()
	}
}

// isValidToken проверяет валидность токена
func isValidToken(tokenString string) bool {
	// Разбор токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	// Возвращает false, если ошибка при разборе токена или токен недействителен
	if err != nil || !token.Valid {
		return false
	}

	// Проверка срока действия токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				return false
			}
		}
		return true
	}
	return false
}
