package services

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	appConfig "github.com/vinicius-lima-barbosa/boilerplate-golang/internal/config"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/dto"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/models"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyInUse  = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInactiveUser       = errors.New("inactive user")
)

type AuthService interface {
	Signup(name, email, password string) (*dto.AuthResponse, error)
	Login(email, password string) (*dto.AuthResponse, error)
	Refresh(refreshToken string) (*dto.AuthTokensResponse, error)
	ValidateAccessToken(token string) (*dto.AuthTokenClaims, error)
	WithTrx(trx *gorm.DB) AuthService
}

type authService struct {
	userRepo         repositories.UserRepository
	jwtAccessSecret  []byte
	jwtRefreshSecret []byte
	accessTTLMinutes int
	refreshTTLHours  int
}

type tokenClaims struct {
	UserID    uint        `json:"user_id"`
	UserType  models.Type `json:"user_type"`
	TokenType string      `json:"token_type"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	env := appConfig.LoadEnv()

	return &authService{
		userRepo:         userRepo,
		jwtAccessSecret:  []byte(env.JWTAccessSecret),
		jwtRefreshSecret: []byte(env.JWTRefreshSecret),
		accessTTLMinutes: env.JWTAccessTTLMinutes,
		refreshTTLHours:  env.JWTRefreshTTLHours,
	}
}

func (s *authService) WithTrx(trx *gorm.DB) AuthService {
	if trx == nil {
		return s
	}

	return &authService{
		userRepo:         s.userRepo.WithTrx(trx),
		jwtAccessSecret:  s.jwtAccessSecret,
		jwtRefreshSecret: s.jwtRefreshSecret,
		accessTTLMinutes: s.accessTTLMinutes,
		refreshTTLHours:  s.refreshTTLHours,
	}
}

func (s *authService) Signup(name, email, password string) (*dto.AuthResponse, error) {
	existingUser, err := s.userRepo.GetByEmail(context.Background(), email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyInUse
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := hashToken(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		Type:         models.USER,
		IsActive:     true,
	}

	createdUser, err := s.userRepo.Create(context.Background(), user)
	if err != nil {
		return nil, err
	}

	tokens, err := s.generateTokenPair(createdUser)
	if err != nil {
		return nil, err
	}

	refreshTokenHash := hashRefreshToken(tokens.RefreshToken)

	if err := s.userRepo.UpdateRefreshTokenHash(context.Background(), createdUser.ID, &refreshTokenHash); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User:   dto.ToAuthUserResponse(createdUser),
		Tokens: *tokens,
	}, nil
}

func (s *authService) Login(email, password string) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if user.PasswordHash == "" || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrInactiveUser
	}

	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, err
	}

	refreshTokenHash := hashRefreshToken(tokens.RefreshToken)

	if err := s.userRepo.UpdateRefreshTokenHash(context.Background(), user.ID, &refreshTokenHash); err != nil {
		return nil, err
	}

	if err := s.userRepo.UpdateLastLoginAt(context.Background(), user.ID, time.Now().UTC()); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User:   dto.ToAuthUserResponse(user),
		Tokens: *tokens,
	}, nil
}

func (s *authService) Refresh(refreshToken string) (*dto.AuthTokensResponse, error) {
	claims, err := s.parseToken(refreshToken, s.jwtRefreshSecret, "refresh")
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(context.Background(), claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrInactiveUser
	}

	if user.RefreshTokenHash == nil || !compareRefreshTokenHash(*user.RefreshTokenHash, refreshToken) {
		return nil, ErrInvalidToken
	}

	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, err
	}

	newRefreshTokenHash := hashRefreshToken(tokens.RefreshToken)

	if err := s.userRepo.UpdateRefreshTokenHash(context.Background(), user.ID, &newRefreshTokenHash); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *authService) ValidateAccessToken(token string) (*dto.AuthTokenClaims, error) {
	claims, err := s.parseToken(token, s.jwtAccessSecret, "access")
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &dto.AuthTokenClaims{
		UserID:    claims.UserID,
		UserType:  claims.UserType,
		TokenType: claims.TokenType,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}

func (s *authService) generateTokenPair(user *models.User) (*dto.AuthTokensResponse, error) {
	accessExpiresAt := time.Now().UTC().Add(time.Duration(s.accessTTLMinutes) * time.Minute)
	refreshExpiresAt := time.Now().UTC().Add(time.Duration(s.refreshTTLHours) * time.Hour)

	accessToken, err := s.generateToken(user, "access", accessExpiresAt, s.jwtAccessSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user, "refresh", refreshExpiresAt, s.jwtRefreshSecret)
	if err != nil {
		return nil, err
	}

	return &dto.AuthTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTTLMinutes * 60),
	}, nil
}

func (s *authService) generateToken(user *models.User, tokenType string, expiresAt time.Time, secret []byte) (string, error) {
	claims := tokenClaims{
		UserID:    user.ID,
		UserType:  user.Type,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (s *authService) parseToken(rawToken string, secret []byte, expectedTokenType string) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(rawToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != expectedTokenType {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func hashToken(value string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func hashRefreshToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func compareRefreshTokenHash(expectedHash, token string) bool {
	providedHash := hashRefreshToken(token)
	return subtle.ConstantTimeCompare([]byte(expectedHash), []byte(providedHash)) == 1
}
