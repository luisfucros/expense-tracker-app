package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
)

// MockAuthService is a testify/mock implementation of service.AuthService.
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuthResponse), args.Error(1)
}

// MockExpenseService is a testify/mock implementation of service.ExpenseService.
type MockExpenseService struct {
	mock.Mock
}

func (m *MockExpenseService) Create(ctx context.Context, userID uint, input model.CreateExpenseInput) (*model.Expense, error) {
	args := m.Called(ctx, userID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Expense), args.Error(1)
}

func (m *MockExpenseService) GetByID(ctx context.Context, userID uint, expenseID uint) (*model.Expense, error) {
	args := m.Called(ctx, userID, expenseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Expense), args.Error(1)
}

func (m *MockExpenseService) Update(ctx context.Context, userID uint, expenseID uint, input model.UpdateExpenseInput) (*model.Expense, error) {
	args := m.Called(ctx, userID, expenseID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Expense), args.Error(1)
}

func (m *MockExpenseService) Delete(ctx context.Context, userID uint, expenseID uint) error {
	args := m.Called(ctx, userID, expenseID)
	return args.Error(0)
}

func (m *MockExpenseService) List(ctx context.Context, userID uint, filter model.ExpenseFilter) (*model.ExpenseListResponse, error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ExpenseListResponse), args.Error(1)
}
