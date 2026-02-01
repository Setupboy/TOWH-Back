package service

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/Setupboy/TOWH-Back/internal/config"
	"github.com/Setupboy/TOWH-Back/internal/dto"
	"github.com/Setupboy/TOWH-Back/internal/models"
	"github.com/Setupboy/TOWH-Back/internal/repositories"
	"github.com/Setupboy/TOWH-Back/internal/utils"
)

type AuthService struct {
	config   *config.Config
	userRepo repositories.UserRepositoryInterface
}

func NewAuthService(
	config *config.Config,
	userRepo repositories.UserRepositoryInterface,

) *AuthService {
	return &AuthService{
		config:   config,
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// check if exist
	if _, err := s.userRepo.GetByEmail(req.Email); err != nil {
		return nil, errors.New("user not found")
	}

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstNme,
		LastName:  req.LastNme,
		Phone:     req.Phone,
		Role:      models.UserRolePlayer,
	}

	if err = s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmailAndActive(req.Email, true)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// check password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.Jwt.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	refreshToken, err := s.userRepo.GetValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh not found or expired")
	}

	user, err := s.userRepo.GetById(int(claims.UserID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = s.userRepo.DeleteRefreshTokenById(refreshToken.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.userRepo.DeleteRefreshToken(refreshToken)

}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		&s.config.Jwt,
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.Jwt.RefreshTokenExpires),
	}

	if err = s.userRepo.CreateRefreshToken(&refreshTokenModel); err != nil {
		log.Println(err)
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:       strconv.Itoa(int(user.ID)),
			Email:    user.Email,
			FirstNme: user.FirstName,
			LastNme:  user.LastName,
			Phone:    user.Phone,
			Role:     string(user.Role),
			IsActive: user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
