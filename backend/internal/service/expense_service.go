package service

import (
	"context"
	"fmt"
	"time"

	"github.com/luisfucros/expense-tracker-app/internal/domain/apierror"
	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
	"github.com/luisfucros/expense-tracker-app/internal/repository"
)

const dateLayout = "2006-01-02"

// expenseService implements service.ExpenseService.
type expenseService struct {
	expenseRepo repository.ExpenseRepository
}

// NewExpenseService creates an ExpenseService with the given expense repository.
func NewExpenseService(expenseRepo repository.ExpenseRepository) ExpenseService {
	return &expenseService{expenseRepo: expenseRepo}
}

func (s *expenseService) Create(ctx context.Context, userID uint, input model.CreateExpenseInput) (*model.Expense, error) {
	date, err := parseDate(input.Date)
	if err != nil {
		return nil, apierror.BadRequest(apierror.CodeBadRequest, fmt.Sprintf("invalid date format, expected YYYY-MM-DD: %s", input.Date))
	}

	if !isValidCategory(input.Category) {
		return nil, apierror.BadRequest(apierror.CodeInvalidCategory, "invalid category")
	}

	expense := &model.Expense{
		UserID:   userID,
		Title:    input.Title,
		Amount:   input.Amount,
		Category: input.Category,
		Date:     date,
		Notes:    input.Notes,
	}

	if err := s.expenseRepo.Create(ctx, expense); err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *expenseService) GetByID(ctx context.Context, userID uint, expenseID uint) (*model.Expense, error) {
	return s.expenseRepo.GetByID(ctx, expenseID, userID)
}

func (s *expenseService) Update(ctx context.Context, userID uint, expenseID uint, input model.UpdateExpenseInput) (*model.Expense, error) {
	expense, err := s.expenseRepo.GetByID(ctx, expenseID, userID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		expense.Title = *input.Title
	}
	if input.Amount != nil {
		expense.Amount = *input.Amount
	}
	if input.Category != nil {
		if !isValidCategory(*input.Category) {
			return nil, apierror.BadRequest(apierror.CodeInvalidCategory, "invalid category")
		}
		expense.Category = *input.Category
	}
	if input.Date != nil {
		date, err := parseDate(*input.Date)
		if err != nil {
			return nil, apierror.BadRequest(apierror.CodeBadRequest, fmt.Sprintf("invalid date format, expected YYYY-MM-DD: %s", *input.Date))
		}
		expense.Date = date
	}
	if input.Notes != nil {
		expense.Notes = *input.Notes
	}

	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *expenseService) Delete(ctx context.Context, userID uint, expenseID uint) error {
	return s.expenseRepo.Delete(ctx, expenseID, userID)
}

func (s *expenseService) List(ctx context.Context, userID uint, filter model.ExpenseFilter) (*model.ExpenseListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}

	expenses, total, err := s.expenseRepo.List(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	return &model.ExpenseListResponse{
		Expenses: expenses,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func parseDate(s string) (time.Time, error) {
	return time.Parse(dateLayout, s)
}

func isValidCategory(c model.Category) bool {
	switch c {
	case model.CategoryGroceries,
		model.CategoryLeisure,
		model.CategoryElectronics,
		model.CategoryUtilities,
		model.CategoryClothing,
		model.CategoryHealth,
		model.CategoryOthers:
		return true
	}
	return false
}
