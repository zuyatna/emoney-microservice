package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zuyatna/emoney-microservice/account-service/server/domain"
	"github.com/zuyatna/emoney-microservice/account-service/server/pb"
	"github.com/zuyatna/emoney-microservice/account-service/server/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountHandler struct {
	pb.UnimplementedAccountServiceServer
	usecase usecase.AccountUseCase
	logger  *logrus.Entry
}

func NewAccountHandler(usecase usecase.AccountUseCase, logger *logrus.Entry) *AccountHandler {
	return &AccountHandler{usecase: usecase, logger: logger}
}

func (h *AccountHandler) Create(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	id, err := h.usecase.CreateAccount(ctx, req.GetName(), req.GetEmail(), req.GetPassword())
	if err != nil {
		h.logger.WithError(err).Error("Error creating account")
		return nil, status.Error(codes.Internal, "Error creating account")
	}

	h.logger.WithField("account_id", id).Info("Account created successfully")
	return &pb.CreateAccountResponse{Id: id, Message: "Account created successfully"}, nil
}

func (h *AccountHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := h.usecase.LoginAccount(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		h.logger.WithError(err).Error("Error logging in")
		return nil, status.Error(codes.Internal, "Error logging in")
	}

	h.logger.WithField("token", token).Info("Account logged in")
	return &pb.LoginResponse{AccessToken: token}, nil
}

func (h *AccountHandler) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
	claims, ok := ctx.Value("claims").(*domain.CustomClaim)
	if !ok || claims.ID != req.GetAccountId() {
		return nil, status.Error(codes.PermissionDenied, "You can only view your own account")
	}

	account, err := h.usecase.GetAccountByID(ctx, req.GetAccountId())
	if err != nil {
		h.logger.WithError(err).WithField("account_id", req.GetAccountId()).Error("Error getting account")
		return nil, status.Errorf(codes.NotFound, "Could not find account: %v", err)
	}

	return &pb.Account{
		Id:        account.ID,
		Name:      account.Name,
		Email:     account.Email,
		Balance:   account.Balance,
		CreatedAt: timestamppb.New(account.CreatedAt),
		UpdatedAt: timestamppb.New(account.UpdatedAt),
	}, nil
}
