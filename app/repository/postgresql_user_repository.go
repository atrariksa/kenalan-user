package repository

import (
	"context"

	"github.com/atrariksa/kenalan-user/app/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	GetUserSubscription(ctx context.Context, email string) (model.User, error)
	GetNextProfileExceptIDs(ctx context.Context, ids []int64) (model.User, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	err := ur.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := ur.DB.Where(&model.User{Email: email}).Take(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return model.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserSubscription(ctx context.Context, email string) (model.User, error) {
	var user model.User
	db := ur.DB
	db = db.Preload("UserSubscriptions")
	db = db.Joins("JOIN user_subscribed_products on users.id = user_subscribed_products.user_id")
	err := db.Where(&model.User{Email: email}).Take(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return model.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) GetNextProfileExceptIDs(ctx context.Context, ids []int64) (model.User, error) {
	var user model.User
	db := ur.DB
	db = db.Preload("UserSubscriptions")
	db = db.Joins("JOIN user_subscribed_products on users.id = user_subscribed_products.user_id")
	err := db.Where("id NOT IN (?)", ids).Take(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return model.User{}, err
	}

	return user, nil
}
