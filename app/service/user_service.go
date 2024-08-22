package service

import (
	"context"
	"errors"

	"github.com/atrariksa/kenalan-user/app/model"
	"github.com/atrariksa/kenalan-user/app/repository"
	"github.com/atrariksa/kenalan-user/app/util"
)

type IUserService interface {
	IsUserExist(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserSubscription(ctx context.Context, email string) (model.User, error)
	GetNextProfileExceptIDs(ctx context.Context, ids []int64, gender string) (model.User, error)
	UpsertSubscription(ctx context.Context, userSubscription model.UserSubscribedProduct) (model.UserSubscribedProduct, error)
}

type UserService struct {
	Repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func (us *UserService) IsUserExist(ctx context.Context, email string) (bool, error) {
	user, err := us.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return false, nil
	}

	return true, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	user, err := us.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}

	if user.ID == 0 {
		return model.User{}, errors.New("user not found")
	}

	return user, nil
}

func (us *UserService) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return model.User{}, errors.New("internal error")
	}
	user.Password = hashedPassword
	err = us.Repo.CreateUser(ctx, &user)
	if err != nil {
		return model.User{}, errors.New("internal error")
	}

	if user.ID == 0 {
		return model.User{}, errors.New("internal error")
	}

	return user, nil
}

func (us *UserService) GetUserSubscription(ctx context.Context, email string) (model.User, error) {
	user, err := us.Repo.GetUserSubscription(ctx, email)
	if err != nil {
		return model.User{}, err
	}

	if user.ID == 0 {
		return model.User{}, errors.New("user not found")
	}

	return user, nil
}

func (us *UserService) GetNextProfileExceptIDs(ctx context.Context, ids []int64, gender string) (model.User, error) {
	return us.Repo.GetNextProfileExceptIDs(ctx, ids, gender)
}

func (us *UserService) UpsertSubscription(ctx context.Context, userSubscription model.UserSubscribedProduct) (model.UserSubscribedProduct, error) {
	err := us.Repo.UpsertSubscription(ctx, &userSubscription)
	if err != nil {
		return model.UserSubscribedProduct{}, errors.New("internal error")
	}

	if userSubscription.ID == 0 {
		return model.UserSubscribedProduct{}, errors.New("internal error")
	}

	return userSubscription, nil
}
