package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/zuyatna/emoney-microservice/account-service/server/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type AuthInterceptor struct {
	jwtSecret string
	logger    *logrus.Logger
}

func NewAuthInterceptor(jwtSecret string, logger *logrus.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		publicMethods := map[string]bool{
			"/pb.AccountService/CreateAccount": true,
			"/pb.AccountService/Login":         true,
		}

		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		claims, err := i.authorize(ctx)
		if err != nil {
			i.logger.WithError(err).Error("Authorization failed")
			return nil, err
		}
		i.logger.WithField("claims", claims).Info("Authorization successful")

		ctx = context.WithValue(ctx, "claims", claims)
		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context) (*domain.CustomClaim, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, status.Error(codes.Unauthenticated, "invalid token format")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &domain.CustomClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(i.jwtSecret), nil
	})
	if err != nil {
		i.logger.WithError(err).Error("JWT parsing error")
		return nil, status.Error(codes.Unauthenticated, "token is invalid")
	}

	if !token.Valid {
		i.logger.Error("JWT token is not valid")
		return nil, status.Error(codes.Unauthenticated, "token is not valid")
	}

	return claims, nil
}
