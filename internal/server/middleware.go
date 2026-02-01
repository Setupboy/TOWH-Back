package server

import (
	"errors"
	"strings"

	"github.com/Setupboy/TOWH-Back/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header required", errors.New("authorization header required"))
			c.Abort()
			return
		}

		tokenParser := strings.Split(authHeader, " ")
		if len(tokenParser) != 2 || tokenParser[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format", errors.New("invalid authorization header format"))
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenParser[1], s.config.Jwt.Secret)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid token", errors.New("invalid token"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

//func (s *Server) adminMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		role, exist := c.Get("admin_role")
//		if !exist {
//			utils.ForbiddenResponse(c, "Forbidden", nil)
//			c.Abort()
//			return
//		}
//
//		if role != string(models.UserRoleAdmin) {
//			utils.ForbiddenResponse(c, "Forbidden", nil)
//			c.Abort()
//			return
//		}
//	}
//}
