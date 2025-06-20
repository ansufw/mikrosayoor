package repository

import (
	"context"
	"errors"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	modelUser := model.User{}

	if err := u.db.Where("email = ? && is_verified", email, true).
		Preload("Roles").First(&modelUser).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Info("[UserRepository-2] GetUserByEmail: usernotfound")
			return nil, err
		}

		log.Errorf("[UserRepository-1] GetUserByEmail: %v", err)
		return nil, err
	}

	entityUser := entity.UserEntity{
		ID:         modelUser.ID,
		Email:      modelUser.Email,
		Password:   modelUser.Password,
		RoleName:   modelUser.Roles[0].Name,
		Address:    modelUser.Address,
		Lat:        modelUser.Lat,
		Lng:        modelUser.Lng,
		Phone:      modelUser.Phone,
		Photo:      modelUser.Photo,
		IsVerified: modelUser.IsVerified,
	}

	return &entityUser, nil
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}
