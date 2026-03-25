package dto

import "github.com/vinicius-lima-barbosa/boilerplate-golang/internal/models"

type AuthUserResponse struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	Type        models.Type `json:"type"`
	Description string      `json:"description"`
}

type AuthTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthResponse struct {
	User   AuthUserResponse   `json:"user"`
	Tokens AuthTokensResponse `json:"tokens"`
}

type AuthTokenClaims struct {
	UserID    uint
	UserType  models.Type
	TokenType string
	ExpiresAt int64
}

func ToAuthUserResponse(user *models.User) AuthUserResponse {
	return AuthUserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Type:        user.Type,
		Description: user.Description,
	}
}
