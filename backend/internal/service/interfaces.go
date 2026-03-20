package service

import (
	"context"

	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
)

// AuthService handles user authentication operations.
type AuthService interface {
	Register(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error)
	Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error)
}

// ExpenseService handles business logic for expense management.
// type ExpenseService interface {
// 	Create(ctx context.Context, userID uint, input model.CreateExpenseInput) (*model.Expense, error)
// 	GetByID(ctx context.Context, userID uint, expenseID uint) (*model.Expense, error)
// 	Update(ctx context.Context, userID uint, expenseID uint, input model.UpdateExpenseInput) (*model.Expense, error)
// 	Delete(ctx context.Context, userID uint, expenseID uint) error
// 	List(ctx context.Context, userID uint, filter model.ExpenseFilter) (*model.ExpenseListResponse, error)
// }
