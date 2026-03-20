package repository

import (
	"context"

	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
)

// UserRepository defines persistence operations for User entities.
type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id uint) error
}

// ExpenseRepository defines persistence operations for Expense entities.
// type ExpenseRepository interface {
// 	Create(ctx context.Context, e *model.Expense) error
// 	GetByID(ctx context.Context, id uint, userID uint) (*model.Expense, error)
// 	Update(ctx context.Context, e *model.Expense) error
// 	Delete(ctx context.Context, id uint, userID uint) error
// 	List(ctx context.Context, userID uint, filter model.ExpenseFilter) ([]*model.Expense, int64, error)
// }
