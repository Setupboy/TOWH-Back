package service

import (
	"strconv"

	"github.com/Setupboy/TOWH-Back/internal/dto"
	"github.com/Setupboy/TOWH-Back/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetProfile(userId uint) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       strconv.Itoa(int(user.ID)),
		Email:    user.Email,
		FirstNme: user.FirstName,
		LastNme:  user.LastName,
		Phone:    user.Phone,
		Role:     string(user.Role),
		IsActive: user.IsActive,
	}, nil
}
func (s *UserService) UpdateProfile(userId uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return s.GetProfile(userId)
}
