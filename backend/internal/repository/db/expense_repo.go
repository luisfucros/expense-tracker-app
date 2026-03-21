package db

import (
	"context"

	"gorm.io/gorm"

	"github.com/luisfucros/expense-tracker-app/internal/domain/apierror"
	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
)

// ExpenseRepo is the GORM-backed implementation of repository.ExpenseRepository.
type ExpenseRepo struct {
	db *gorm.DB
}

// NewExpenseRepository creates a new MySQL expense repository.
func NewExpenseRepository(db *gorm.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) Create(ctx context.Context, e *model.Expense) error {
	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return err
	}
	return nil
}

func (r *ExpenseRepo) GetByID(ctx context.Context, id uint, userID uint) (*model.Expense, error) {
	var e model.Expense
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&e).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.NotFound("expense not found")
		}
		return nil, err
	}
	return &e, nil
}

func (r *ExpenseRepo) Update(ctx context.Context, e *model.Expense) error {
	if err := r.db.WithContext(ctx).Save(e).Error; err != nil {
		return err
	}
	return nil
}

func (r *ExpenseRepo) Delete(ctx context.Context, id uint, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Expense{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apierror.NotFound("expense not found")
	}
	return nil
}

func (r *ExpenseRepo) List(ctx context.Context, userID uint, filter model.ExpenseFilter) ([]*model.Expense, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Expense{}).Where("user_id = ?", userID)

	if filter.Category != nil {
		query = query.Where("category = ?", *filter.Category)
	}
	if filter.StartDate != nil {
		query = query.Where("date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("date <= ?", *filter.EndDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var expenses []*model.Expense
	if err := query.
		Order("date DESC, id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&expenses).Error; err != nil {
		return nil, 0, err
	}

	return expenses, total, nil
}
