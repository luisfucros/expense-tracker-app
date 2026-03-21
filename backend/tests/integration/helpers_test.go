package integration_test

import (
	"go.uber.org/zap"

	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
)

// noopLogger returns a zap no-op logger suitable for tests.
func noopLogger() *zap.Logger {
	return zap.NewNop()
}

// registerInput builds a model.RegisterInput for integration tests.
func registerInput(email string) model.RegisterInput {
	return model.RegisterInput{
		Name:     "Test User",
		Email:    email,
		Password: "password123",
	}
}
