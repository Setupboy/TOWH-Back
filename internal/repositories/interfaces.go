package repositories

import "github.com/Setupboy/TOWH-Back/internal/models"

type UserRepositoryInterface interface {
	GetByEmail(email string) (*models.User, error)
	GetById(id int) (*models.User, error)
	GetByEmailAndActive(email string, isActive bool) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error

	CreateRefreshToken(token *models.RefreshToken) error
	GetValidateRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenById(id uint) error
}
