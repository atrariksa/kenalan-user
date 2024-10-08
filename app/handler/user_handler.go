package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	pb "github.com/atrariksa/kenalan-user/app/internal/grpc_user_server"
	"github.com/atrariksa/kenalan-user/app/model"
	"github.com/atrariksa/kenalan-user/app/repository"
	"github.com/atrariksa/kenalan-user/app/service"
	"github.com/atrariksa/kenalan-user/app/util"
	"github.com/atrariksa/kenalan-user/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService service.IUserService
}

func GetUserServiceServer(svc service.IUserService) *userServiceServer {
	return &userServiceServer{
		userService: svc,
	}
}

func (s userServiceServer) IsUserExist(ctx context.Context, isUserExistRequest *pb.IsUserExistRequest) (*pb.IsUserExistResponse, error) {
	if isUserExistRequest.Email == "" {
		return nil, errors.New("invalid email")
	}

	isUserExist, err := s.userService.IsUserExist(ctx, isUserExistRequest.Email)
	if err != nil {
		return nil, errors.New("internal error")
	}

	response := pb.IsUserExistResponse{
		Code:        0000,
		IsUserExist: isUserExist,
	}
	return &response, nil
}

func (s userServiceServer) CreateUser(ctx context.Context, createUserRequest *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	dob, err := util.ToDateTimeYYYYMMDD(createUserRequest.User.Dob)
	if err != nil {
		return nil, errors.New("invalid dob")
	}

	timeNow := util.TimeNow()
	_, err = s.userService.CreateUser(ctx, model.User{
		Fullname:  createUserRequest.User.FullName,
		Gender:    createUserRequest.User.Gender,
		DOB:       dob,
		Email:     createUserRequest.User.Email,
		Password:  createUserRequest.User.Password,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	})

	if err != nil {
		return nil, errors.New("internal error")
	}

	return &pb.CreateUserResponse{
		Code:    0000,
		Message: "success",
	}, nil
}

func (s userServiceServer) GetUserByEmail(ctx context.Context, getUserByEmailRequest *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	if getUserByEmailRequest.Email == "" {
		return nil, errors.New("invalid email")
	}

	user, err := s.userService.GetUserByEmail(ctx, getUserByEmailRequest.Email)
	if err != nil {
		return nil, errors.New("internal error")
	}

	response := pb.GetUserByEmailResponse{
		Code: 0000,
		User: &pb.User{
			Id:       user.ID,
			FullName: user.Fullname,
			Gender:   user.Gender,
			Dob:      user.DOB.Format(util.DateFormatYYYYMMDD),
			Email:    user.Email,
			Password: user.Password,
		},
	}
	return &response, nil
}

func (s userServiceServer) GetUserSubscription(ctx context.Context, req *pb.GetUserSubscriptionRequest) (*pb.GetUserSubscriptionResponse, error) {
	user, err := s.userService.GetUserSubscription(ctx, req.Email)
	if err != nil {
		return nil, errors.New("internal error")
	}

	var subscriptions = make([]*pb.UserSubscription, 0)
	if user.UserSubscriptions != nil && len(user.UserSubscriptions) > 0 {
		for i := 0; i < len(user.UserSubscriptions); i++ {
			subscriptions = append(subscriptions, &pb.UserSubscription{
				ExpiredAt:   user.UserSubscriptions[i].ExpiredAt.Format(util.DateFormatYYYYMMDDTHHmmss),
				IsActive:    user.UserSubscriptions[i].IsActive,
				ProductCode: user.UserSubscriptions[i].ProductCode,
				ProductName: user.UserSubscriptions[i].ProductName,
			})
		}
	}

	response := pb.GetUserSubscriptionResponse{
		Code: 0000,
		User: &pb.User{
			Id:       user.ID,
			FullName: user.Fullname,
			Gender:   user.Gender,
			Dob:      user.DOB.Format(util.DateFormatYYYYMMDD),
			Email:    user.Email,
			Password: user.Password,
		},
		Subscriptions: subscriptions,
	}

	return &response, nil
}

func (s userServiceServer) GetNextProfileExceptIDs(ctx context.Context, req *pb.GetNextProfileExceptIDsRequest) (*pb.GetNextProfileExceptIDsResponse, error) {
	user, err := s.userService.GetNextProfileExceptIDs(ctx, req.Ids, req.Gender)
	if err == gorm.ErrRecordNotFound {
		return nil, status.Error(05, "user not found")
	}
	if err != nil {
		return nil, err
	}

	var subscriptions = make([]*pb.UserSubscription, 0)
	if user.UserSubscriptions != nil && len(user.UserSubscriptions) > 0 {
		for i := 0; i < len(user.UserSubscriptions); i++ {
			subscriptions = append(subscriptions, &pb.UserSubscription{
				ExpiredAt:   user.UserSubscriptions[i].ExpiredAt.Format(util.DateFormatYYYYMMDDTHHmmss),
				IsActive:    user.UserSubscriptions[i].IsActive,
				ProductCode: user.UserSubscriptions[i].ProductCode,
				ProductName: user.UserSubscriptions[i].ProductName,
			})
		}
	}

	response := pb.GetNextProfileExceptIDsResponse{
		Code: 0000,
		User: &pb.User{
			Id:       user.ID,
			FullName: user.Fullname,
			Gender:   user.Gender,
			Dob:      user.DOB.Format(util.DateFormatYYYYMMDD),
			Email:    user.Email,
			Password: user.Password,
		},
		Subscriptions: subscriptions,
	}

	return &response, nil
}

func (s userServiceServer) UpsertSubscription(ctx context.Context, upsertSubscriptionRequest *pb.UpsertSubscriptionRequest) (*pb.UpsertSubscriptionResponse, error) {
	expiredAt, err := util.ToDateTimeYYYYMMDDTHHmmss(upsertSubscriptionRequest.ExpiredAt)
	if err != nil {
		return nil, errors.New("invalid expired_at")
	}

	timeNow := util.TimeNow()
	_, err = s.userService.UpsertSubscription(ctx, model.UserSubscribedProduct{
		UserID:      upsertSubscriptionRequest.UserId,
		ExpiredAt:   expiredAt,
		IsActive:    true,
		ProductCode: upsertSubscriptionRequest.ProductCode,
		ProductName: upsertSubscriptionRequest.ProductName,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	})

	if err != nil {
		return nil, errors.New("internal error")
	}

	return &pb.UpsertSubscriptionResponse{
		Code:    0000,
		Message: "success",
	}, nil
}

func SetupServer() {
	fmt.Println("---User Service---")

	cfg := config.GetConfig()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServerConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db := util.GetDB(cfg)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, GetUserServiceServer(userService))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
