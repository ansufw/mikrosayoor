package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"

	"user-service/config"
	"user-service/internal/adapter/message"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/domain/entity"
	"user-service/utils/conv"
)

type UserServiceInterface interface {
	SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error)
	CreateUserAccount(ctx context.Context, req entity.UserEntity) error
	ForgotPassword(ctx context.Context, req entity.UserEntity) error
}

type userService struct {
	repo       repository.UserRepositoryInterface
	cfg        *config.Config
	jwtService JwtServiceInterface
	repoToken  repository.VerificationTokenRepositoryInterface
}

// ForgotPassword implements UserServiceInterface.
func (u *userService) ForgotPassword(ctx context.Context, req entity.UserEntity) error {
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Errorf("[UserService-1] ForgotPasswd: %v", err)
		return err
	}

	token := uuid.New().String()

	reqEntity := entity.VerificationTokenEntity{
		UserID:    user.ID,
		Token:     token,
		TokenType: "forgot_password",
	}

	err = u.repoToken.CreateVerificationToken(ctx, reqEntity)
	if err != nil {
		log.Errorf("[UserService-2] ForgotPasswd: %v", err)
		return err
	}

	urlForgot := fmt.Sprintf("%s/forgot-password?token=%s", u.cfg.App.UrlForgotPassword, token)

	messageParam := fmt.Sprintf("please click link below for reset password: %v", urlForgot)
	err = message.PublishMessage(req.Email, messageParam, "forgot-password")
	if err != nil {
		log.Errorf("[UserService-3] ForgotPasswd: %v", err)
		return err
	}
	return nil
}

// CreateUserAccount implements UserServiceInterface.
func (u *userService) CreateUserAccount(ctx context.Context, req entity.UserEntity) error {
	passwd, err := conv.HashPassword(req.Password)
	if err != nil {
		log.Errorf("[UserService-1] CreateUserAccount: %v", err)
		return err
	}

	req.Password = passwd
	token := uuid.New().String()
	req.Token = token

	err = u.repo.CreateUserAccount(ctx, req)
	if err != nil {
		log.Errorf("[UserService-2] CreateUserAccount: %v", err)
		return err
	}

	urlVerify := fmt.Sprintf("http://localhost:8080/verify?token=%v", req.Token)
	messageParam := fmt.Sprintf("please verify your account with click link below: %v", req.Token)
	err = message.PublishMessage(req.Email, messageParam, urlVerify)
	if err != nil {
		log.Errorf("[UserService-3] CreateUserAccount: %v", err)
		return err
	}

	return nil
}

func (u *userService) SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error) {
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Errorf("[UserService-1] SignIn: %v", err)
		return nil, "", err
	}

	if checkPass := conv.CheckPasswordHash(req.Password, user.Password); !checkPass {
		err = errors.New("invalid password")
		log.Errorf("[UserService-2] SignIn: %v", err)
		return nil, "", err
	}

	token, err := u.jwtService.GenerateToken(user.ID)
	if err != nil {
		log.Errorf("[UserService-3] SignIn: %v", err)
		return nil, "", err
	}

	sessionData := map[string]interface{}{
		"user_id":    user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"logged_in":  true,
		"created_at": time.Now().String(),
		"token":      token,
	}

	redisConn := config.NewRedisClient()
	err = redisConn.HSet(ctx, token, sessionData).Err()
	if err != nil {
		log.Errorf("[UserService-4] SignIn: %v", err)
		return nil, "", err
	}

	return user, token, nil
}

func NewUserService(repo repository.UserRepositoryInterface, cfg *config.Config, jwtService JwtServiceInterface, repoToken repository.VerificationTokenRepositoryInterface) UserServiceInterface {
	return &userService{
		repo:       repo,
		cfg:        cfg,
		jwtService: jwtService,
		repoToken:  repoToken,
	}
}
