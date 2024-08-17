package repository

import (
	"context"

	"github.com/atrariksa/kenalan-user/app/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	GetUserSubscription(ctx context.Context, email string) (model.User, error)
	GetNextProfileExceptIDs(ctx context.Context, ids []int64, gender string) (model.User, error)
	UpsertSubscription(ctx context.Context, userSubscription *model.UserSubscribedProduct) error
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
	err := db.Where(&model.User{Email: email}).Take(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return model.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) GetNextProfileExceptIDs(ctx context.Context, ids []int64, gender string) (model.User, error) {
	var user model.User
	db := ur.DB
	db = db.Preload("UserSubscriptions")
	err := db.Where("users.id NOT IN (?) AND users.gender=?", ids, gender).Take(&user).Error
	return user, err
}

func (ur *UserRepository) UpsertSubscription(ctx context.Context, userSubscription *model.UserSubscribedProduct) error {
	var existingData model.UserSubscribedProduct
	db := ur.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Where(
		"user_id = ? AND product_code = ?",
		userSubscription.UserID,
		userSubscription.ProductCode).Take(&existingData)
	if existingData.ID == 0 {
		return ur.DB.Create(userSubscription).Error
	}
	userSubscription.ID = existingData.ID
	return db.Exec(
		`UPDATE "user_subscribed_products" 
		 SET "expired_at"=?,"updated_at"=?, "is_active"=?
		 WHERE "id" = ?`,
		userSubscription.ExpiredAt,
		userSubscription.UpdatedAt,
		userSubscription.IsActive,
		existingData.ID).Error
}
