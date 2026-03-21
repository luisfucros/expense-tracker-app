package unit_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/luisfucros/expense-tracker-app/internal/config"
	"github.com/luisfucros/expense-tracker-app/internal/domain/apierror"
	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
	"github.com/luisfucros/expense-tracker-app/internal/service"
	"github.com/luisfucros/expense-tracker-app/tests/mocks"
)

func newTestConfig() *config.Config {
	return &config.Config{
		JWTSecret:      "test-secret-key-long-enough",
		JWTExpiryHours: 24,
	}
}

func TestRegister_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(userRepo, cfg)

	input := model.RegisterInput{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	}

	// Expect Create to be called and succeed; it will set ID on the passed user.
	userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.Email == input.Email && u.Name == input.Name
	})).Run(func(args mock.Arguments) {
		u := args.Get(1).(*model.User)
		u.ID = 1
	}).Return(nil)

	resp, err := svc.Register(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, uint(1), resp.User.ID)
	assert.Equal(t, input.Name, resp.User.Name)
	assert.Equal(t, input.Email, resp.User.Email)

	userRepo.AssertExpectations(t)
}

func TestRegister_EmailConflict(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(userRepo, cfg)

	input := model.RegisterInput{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "password123",
	}

	conflictErr := apierror.Conflict("email already registered")
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(conflictErr)

	resp, err := svc.Register(context.Background(), input)

	assert.Nil(t, resp)
	assert.Error(t, err)

	var apiErr *apierror.APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.Equal(t, apierror.CodeConflict, apiErr.Code)

	userRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(userRepo, cfg)

	// Pre-hash a known password
	// bcrypt hash of "password123" with cost 12 (pre-computed for test speed)
	// We register first to get a real hash, then use it for login test.
	// Instead, we call Register then Login to keep it integration-style within unit.
	registerInput := model.RegisterInput{
		Name:     "Carol",
		Email:    "carol@example.com",
		Password: "password123",
	}

	var createdUser *model.User
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).
		Run(func(args mock.Arguments) {
			u := args.Get(1).(*model.User)
			u.ID = 42
			createdUser = &model.User{
				ID:           u.ID,
				Name:         u.Name,
				Email:        u.Email,
				PasswordHash: u.PasswordHash,
			}
		}).Return(nil)

	_, err := svc.Register(context.Background(), registerInput)
	assert.NoError(t, err)

	// Now test Login
	userRepo.On("GetByEmail", mock.Anything, registerInput.Email).Return(createdUser, nil)

	loginInput := model.LoginInput{
		Email:    registerInput.Email,
		Password: registerInput.Password,
	}

	resp, err := svc.Login(context.Background(), loginInput)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, uint(42), resp.User.ID)

	userRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(userRepo, cfg)

	// Build a user with a known bcrypt hash for "correctpassword"
	registerInput := model.RegisterInput{
		Name:     "Dave",
		Email:    "dave@example.com",
		Password: "correctpassword",
	}
	var storedUser *model.User
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).
		Run(func(args mock.Arguments) {
			u := args.Get(1).(*model.User)
			u.ID = 10
			storedUser = &model.User{
				ID:           u.ID,
				Name:         u.Name,
				Email:        u.Email,
				PasswordHash: u.PasswordHash,
			}
		}).Return(nil)
	_, _ = svc.Register(context.Background(), registerInput)

	userRepo.On("GetByEmail", mock.Anything, registerInput.Email).Return(storedUser, nil)

	loginInput := model.LoginInput{
		Email:    registerInput.Email,
		Password: "wrongpassword",
	}

	resp, err := svc.Login(context.Background(), loginInput)

	assert.Nil(t, resp)
	assert.Error(t, err)

	var apiErr *apierror.APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.Equal(t, apierror.CodeUnauthorized, apiErr.Code)

	userRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(userRepo, cfg)

	userRepo.On("GetByEmail", mock.Anything, "nobody@example.com").
		Return(nil, apierror.NotFound("user not found"))

	resp, err := svc.Login(context.Background(), model.LoginInput{
		Email:    "nobody@example.com",
		Password: "whatever",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)

	var apiErr *apierror.APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.Equal(t, apierror.CodeUnauthorized, apiErr.Code)

	userRepo.AssertExpectations(t)
}
