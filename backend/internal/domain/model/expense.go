package model

import "time"

// Category represents the type of expense.
type Category string

const (
	CategoryGroceries   Category = "Groceries"
	CategoryLeisure     Category = "Leisure"
	CategoryElectronics Category = "Electronics"
	CategoryUtilities   Category = "Utilities"
	CategoryClothing    Category = "Clothing"
	CategoryHealth      Category = "Health"
	CategoryOthers      Category = "Others"
)

// Expense represents a single expense record in the database.
type Expense struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Category  Category  `gorm:"size:50;not null" json:"category"`
	Date      time.Time `gorm:"not null" json:"date"`
	Notes     string    `gorm:"size:1000" json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateExpenseInput is the request body for creating a new expense.
type CreateExpenseInput struct {
	Title    string   `json:"title" validate:"required,min=1,max=255"`
	Amount   float64  `json:"amount" validate:"required,gt=0"`
	Category Category `json:"category" validate:"required,oneof=Groceries Leisure Electronics Utilities Clothing Health Others"`
	Date     string   `json:"date" validate:"required"` // ISO date string YYYY-MM-DD
	Notes    string   `json:"notes" validate:"max=1000"`
}

// UpdateExpenseInput is the request body for updating an existing expense.
// All fields are optional (pointer types).
type UpdateExpenseInput struct {
	Title    *string   `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Amount   *float64  `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Category *Category `json:"category,omitempty" validate:"omitempty,oneof=Groceries Leisure Electronics Utilities Clothing Health Others"`
	Date     *string   `json:"date,omitempty"`
	Notes    *string   `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

// ExpenseFilter holds query parameters for filtering/paginating expenses.
type ExpenseFilter struct {
	Category  *Category
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	PageSize  int
}

// ExpenseListResponse is the paginated list response for expenses.
type ExpenseListResponse struct {
	Expenses []*Expense `json:"expenses"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}
