package server

import (
	"github.com/Setupboy/TOWH-Back/internal/dto"
	"github.com/Setupboy/TOWH-Back/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) profile(c *gin.Context) {
	userId := c.GetUint("user_id")

	profile, err := s.userService.GetProfile(userId)
	if err != nil {
		utils.NotFoundResponse(c, "User not found", err)
		return
	}

	utils.SuccessResponse(c, "profile retried successfully", profile)
}

func (s *Server) updateProfile(c *gin.Context) {
	userId := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid request", err)
		return
	}

	profile, err := s.userService.UpdateProfile(userId, &req)
	if err != nil {
		utils.NotFoundResponse(c, "Failed to update profile", err)
	}

	utils.SuccessResponse(c, "profile updated successfully", profile)
}
