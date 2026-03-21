package unit_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/luisfucros/expense-tracker-app/internal/domain/apierror"
	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
	"github.com/luisfucros/expense-tracker-app/internal/service"
	"github.com/luisfucros/expense-tracker-app/tests/mocks"
)

func TestCreate_Success(t *testing.T) {
	expenseRepo := new(mocks.MockExpenseRepository)
	svc := service.NewExpenseService(expenseRepo)

	input := model.CreateExpenseInput{
		Title:    "Coffee",
		Amount:   4.50,
		Category: model.CategoryLeisure,
		Date:     "2024-01-15",
		Notes:    "Morning coffee",
	}

	expenseRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *model.Expense) bool {
		return e.Title == input.Title &&
			e.Amount == input.Amount &&
			e.Category == input.Category &&
			e.UserID == uint(1)
	})).Run(func(args mock.Arguments) {
		e := args.Get(1).(*model.Expense)
		e.ID = 100
	}).Return(nil)

	expense, err := svc.Create(context.Background(), 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, expense)
	assert.Equal(t, uint(100), expense.ID)
	assert.Equal(t, input.Title, expense.Title)
	assert.Equal(t, input.Amount, expense.Amount)
	assert.Equal(t, input.Category, expense.Category)

	expenseRepo.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	expenseRepo := new(mocks.MockExpenseRepository)
	svc := service.NewExpenseService(expenseRepo)

	expected := &model.Expense{
		ID:       5,
		UserID:   1,
		Title:    "Groceries run",
		Amount:   55.00,
		Category: model.CategoryGroceries,
		Date:     time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}

	expenseRepo.On("GetByID", mock.Anything, uint(5), uint(1)).Return(expected, nil)

	expense, err := svc.GetByID(context.Background(), 1, 5)

	assert.NoError(t, err)
	assert.Equal(t, expected, expense)

	expenseRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	expenseRepo := new(mocks.MockExpenseRepository)
	svc := service.NewExpenseService(expenseRepo)

	expenseRepo.On("GetByID", mock.Anything, uint(999), uint(1)).
		Return(nil, apierror.NotFound("expense not found"))

	expense, err := svc.GetByID(context.Background(), 1, 999)

	assert.Nil(t, expense)
	assert.Error(t, err)

	var apiErr *apierror.APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.Equal(t, apierror.CodeNotFound, apiErr.Code)

	expenseRepo.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	expenseRepo := new(mocks.MockExpenseRepository)
	svc := service.NewExpenseService(expenseRepo)

	existing := &model.Expense{
		ID:       7,
		UserID:   2,
		Title:    "Old title",
		Amount:   10.00,
		Category: model.CategoryOthers,
		Date:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}

	newTitle := "New title"
	newAmount := 25.50
	input := model.UpdateExpenseInput{
		Title:  &newTitle,
		Amount: &newAmount,
	}

	expenseRepo.On("GetByID", mock.Anything, uint(7), uint(2)).Return(existing, nil)
	expenseRepo.On("Update", mock.Anything, mock.MatchedBy(func(e *model.Expense) bool {
		return e.ID == 7 && e.Title == newTitle && e.Amount == newAmount
	})).Return(nil)

	updated, err := svc.Update(context.Background(), 2, 7, input)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, newTitle, updated.Title)
	assert.Equal(t, newAmount, updated.Amount)

	expenseRepo.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	expenseRepo := new(mocks.MockExpenseRepository)
	svc := service.NewExpenseService(expenseRepo)

	expenseRepo.On("Delete", mock.Anything, uint(3), uint(1)).Return(nil)

	err := svc.Delete(context.Background(), 1, 3)

	assert.NoError(t, err)
	expenseRepo.AssertExpectations(t)
}

func TestList_Success(t *testing.T) {
	expenseRepo := new(mocks.MockExpenseRepository)
	svc := service.NewExpenseService(expenseRepo)

	expenses := []*model.Expense{
		{ID: 1, UserID: 1, Title: "Coffee", Amount: 4.50, Category: model.CategoryLeisure},
		{ID: 2, UserID: 1, Title: "Groceries", Amount: 80.00, Category: model.CategoryGroceries},
	}

	filter := model.ExpenseFilter{Page: 1, PageSize: 20}

	expenseRepo.On("List", mock.Anything, uint(1), filter).Return(expenses, int64(2), nil)

	result, err := svc.List(context.Background(), 1, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Expenses, 2)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 20, result.PageSize)

	expenseRepo.AssertExpectations(t)
}
