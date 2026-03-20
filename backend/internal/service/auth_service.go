package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/luisfucros/expense-tracker-app/internal/config"
	"github.com/luisfucros/expense-tracker-app/internal/domain/apierror"
	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
	"github.com/luisfucros/expense-tracker-app/internal/repository"
)

// authService implements service.AuthService.
type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewAuthService creates an AuthService with the given user repository and config.
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, apierror.Internal("failed to hash password")
	}

	user := &model.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, apierror.Internal("failed to generate token")
	}

	return &model.AuthResponse{
		Token: token,
		User: model.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		// Return unauthorized instead of not found to avoid user enumeration
		return nil, apierror.Unauthorized("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, apierror.Unauthorized("invalid email or password")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, apierror.Internal("failed to generate token")
	}

	return &model.AuthResponse{
		Token: token,
		User: model.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

// generateToken creates a signed JWT for the given user ID.
func (s *authService) generateToken(userID uint) (string, error) {
	expiry := time.Duration(s.cfg.JWTExpiryHours) * time.Hour

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
