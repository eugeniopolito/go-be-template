package api

import (
	"time"

	db "github.com/eugeniopolito/gobetemplate/db/sqlc"
)

type PaginationRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"size" binding:"required,min=1,max=100"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	// The Role of the user
	// example: 1 for admin, 2 for user
	Role     int    `json:"role" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	// The username of a thing
	// example: joedoe
	Username string `json:"username"`
	// The Name of he user
	// example: Some name
	Name string `json:"name"`
	// The Surname of he user
	// example: Some name
	Surname string `json:"surname"`
	// The enabled/disabeld flag
	// example: 0 for disabled, 1 for enabled
	Enabled bool `json:"enabled"`
	// The Email of the user
	// example: joe.doe@email.com
	Email string `json:"email"`
	// The Role of the user
	// example: 1 for admin, 2 for user
	Role             int       `json:"role"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func CreateUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:         user.Username,
		Name:             user.Name,
		Surname:          user.Surname,
		Email:            user.Email,
		Role:             int(user.Role.Int32),
		Enabled:          user.Enabled,
		CreatedAt:        user.CreatedAt,
		PasswordChangeAt: user.PasswordChangeAt,
	}
}

type VerifyEmailRequest struct {
	EmailId    int    `form:"email_id" binding:"required"`
	SecretCode string `form:"secret_code" binding:"required"`
}

type VerifyEmailResponse struct {
	IsVerified bool `json:"is_verified"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

type GetUserRequest struct {
	Username string `uri:"username" binding:"required"`
}

type CountUsersResponse struct {
	Count int `json:"count"`
}
