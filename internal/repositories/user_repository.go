package repositories

import (
	"time"

	"github.com/Setupboy/TOWH-Back/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetById(id int) (*models.User, error) {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetByEmailAndActive(email string, isActive bool) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ? AND is_active = ?", email, isActive).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) Create(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *UserRepository) Update(user *models.User) error {
	return u.db.Save(user).Error
}

func (u *UserRepository) Delete(id uint) error {
	return u.db.Delete(&models.User{}, id).Error
}

func (u *UserRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return u.db.Create(token).Error
}

func (u *UserRepository) GetValidateRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := u.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (u *UserRepository) DeleteRefreshToken(token string) error {
	return u.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

func (u *UserRepository) DeleteRefreshTokenById(id uint) error {
	return u.db.Delete(&models.RefreshToken{}, id).Error
}
