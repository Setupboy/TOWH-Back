package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FirstNme string `json:"first_nme" binding:"required"`
	LastNme  string `json:"last_nme" binding:"required"`
	Phone    string `json:"phone" binding:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FirstNme string `json:"first_nme"`
	LastNme  string `json:"last_nme"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
}
