package db

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"github.com/luisfucros/expense-tracker-app/internal/domain/apierror"
	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
)

// UserRepo is the GORM-backed implementation of repository.UserRepository.
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepository creates a new MySQL user repository.
func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *model.User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		if isDuplicateEntry(err) {
			return apierror.Conflict("email already registered")
		}
		return err
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.NotFound("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.NotFound("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Update(ctx context.Context, u *model.User) error {
	if err := r.db.WithContext(ctx).Save(u).Error; err != nil {
		if isDuplicateEntry(err) {
			return apierror.Conflict("email already in use")
		}
		return err
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apierror.NotFound("user not found")
	}
	return nil
}

// isDuplicateEntry checks if a MySQL error is a duplicate entry (error 1062).
func isDuplicateEntry(err error) bool {
	return strings.Contains(err.Error(), "1062") ||
		strings.Contains(err.Error(), "Duplicate entry") ||
		strings.Contains(err.Error(), "duplicate")
}
