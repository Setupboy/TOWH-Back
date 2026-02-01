package server

import (
	"net/http"

	"github.com/Setupboy/TOWH-Back/internal/config"
	"github.com/Setupboy/TOWH-Back/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config        *config.Config
	db            *gorm.DB
	logger        *zerolog.Logger
	authService   *service.AuthService
	userService   *service.UserService
	uploadService *service.UploadService
}

func New(
	config *config.Config,
	db *gorm.DB,
	logger *zerolog.Logger,
	authService *service.AuthService,
	userService *service.UserService,
	uploadService *service.UploadService,
) *Server {
	return &Server{
		config:        config,
		db:            db,
		logger:        logger,
		authService:   authService,
		userService:   userService,
		uploadService: uploadService,
	}
}

func (s *Server) SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(s.corsMiddleware())

	router.GET("healthCheck", s.healthCheck)

	api := router.Group("/v1/api")
	{
		auth := api.Group("auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/logout", s.logout)
			auth.POST("/refresh", s.refreshToken)
		}

		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			// user rotes
			user := protected.Group("/users")
			{
				userRote := user
				userRote.GET("/profile", s.profile)
				userRote.PUT("/profile", s.updateProfile)
			}

		}

		// public rotes

	}

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
