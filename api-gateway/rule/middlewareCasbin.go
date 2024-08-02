package rule

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware - Middleware для проверки разрешений с использованием Casbin
func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получение роли пользователя из заголовка
		role := c.GetHeader("X-User-Role")
		if role == "" {
			role = "user" // или другая роль по умолчанию, если не указана
		}

		// Проверка разрешений Casbin
		ok, _ := enforcer.Enforce(role, c.Request.URL.Path, c.Request.Method)
		if !ok {
			// Если доступ запрещен, отправить ответ с ошибкой 403
			c.JSON(http.StatusForbidden, gin.H{"message": "Access is denied"})
			c.Abort()
			return
		}

		// Если доступ разрешен, продолжить обработку запроса
		c.Next()
	}
}
