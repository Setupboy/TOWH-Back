package server

import (
	"github.com/Setupboy/TOWH-Back/internal/dto"
	"github.com/Setupboy/TOWH-Back/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	response, err := s.authService.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Registration failed", err)
		return
	}

	utils.CreatedResponse(c, "Registration successfully", response)
}

func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	response, err := s.authService.Login(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Login failed", err)
		return
	}
	utils.SuccessResponse(c, "Login successfully", response)
}

func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	response, err := s.authService.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Refresh token failed", err)
	}
	utils.SuccessResponse(c, "Refresh token successfully", response)
}

func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	if err := s.authService.Logout(req.RefreshToken); err != nil {
		utils.InternalServerErrorResponse(c, "Logout failed", err)
	}

	utils.SuccessResponse(c, "Logout successfully", nil)
}
