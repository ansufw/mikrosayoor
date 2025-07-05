package repository

import (
	"context"
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"

	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"
)

type UserRepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	CreateUserAccount(ctx context.Context, req entity.UserEntity) error
	UpdateUserVerified(ctx context.Context, userID int64) (*entity.UserEntity, error)
	UpdatePasswordByID(ctx context.Context, req entity.UserEntity) error
	GetUserByID(ctx context.Context, userID int64) (*entity.UserEntity, error)
	UpdateDataUser(ctx context.Context, req entity.UserEntity) error
}

type userRepository struct {
	db *gorm.DB
}

// UpdateDataUser implements UserRepositoryInterface.
func (u *userRepository) UpdateDataUser(ctx context.Context, req entity.UserEntity) error {
	modelUser := model.User{}
	modelUser.Name = req.Name
	modelUser.Email = req.Email
	modelUser.Address = req.Address
	modelUser.Lat = req.Lat
	modelUser.Lng = req.Lng
	modelUser.Phone = req.Phone
	modelUser.Photo = req.Photo

	if err := u.db.Where("id = ? AND is_verified = true", req.ID).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] UpdateDataUser: %v", err)
			return err
		}

		log.Errorf("[UserRepository-2] UpdateDataUser: %v", err)
		return err
	}

	if err := u.db.Save(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-3] UpdateDataUser: %v", err)
		return err
	}

	return nil
}

// GetUserByID implements UserRepositoryInterface.
func (u *userRepository) GetUserByID(ctx context.Context, userID int64) (*entity.UserEntity, error) {
	modelUser := model.User{}

	if err := u.db.Where("id = ? AND is_verified = ?", userID, true).Preload("Roles").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] GetUserByID: %v", err)
			return nil, err
		}

		log.Errorf("[UserRepository-2] GetUserByID: %v", err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:       modelUser.ID,
		Email:    modelUser.Email,
		Name:     modelUser.Name,
		RoleName: modelUser.Roles[0].Name,
		Address:  modelUser.Address,
		Lat:      modelUser.Lat,
		Lng:      modelUser.Lng,
		Phone:    modelUser.Phone,
		Photo:    modelUser.Photo,
	}, nil

}

// UpdatePasswordByID implements UserRepositoryInterface.
func (u *userRepository) UpdatePasswordByID(ctx context.Context, req entity.UserEntity) error {
	modelUser := model.User{}

	if err := u.db.Where("id = ?", req.ID).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] UpdatePasswordByID: %v", err)
			return err
		}

		log.Errorf("[UserRepository-2] UpdatePasswordByID: %v", err)
		return err
	}

	modelUser.Password = req.Password
	if err := u.db.Save(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-3] UpdatePasswordByID: %v", err)
		return err
	}
	return nil
}

// UpdateUserVerified implements UserRepositoryInterface.
func (u *userRepository) UpdateUserVerified(ctx context.Context, userID int64) (*entity.UserEntity, error) {
	modelUser := model.User{}
	if err := u.db.Where("id = ?", userID).Preload("Roles").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] UpdateUserVerified: %v", err)
			return nil, err
		}

		log.Errorf("[UserRepository-2] UpdateUserVerified: %v", err)
		return nil, err
	}

	modelUser.IsVerified = true
	if err := u.db.Save(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-3] UpdateUserVerified: %v", err)
		return nil, err
	}

	return &entity.UserEntity{
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
	}, nil
}

func (u *userRepository) CreateUserAccount(ctx context.Context, req entity.UserEntity) error {
	modelRole := &model.Role{}
	if err := u.db.Where("name = ?", "Customer").First(&modelRole).Error; err != nil {
		log.Errorf("[UserRepository-1] CreateUserAccount: %v", err)
		return err
	}

	modelUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Roles:    []*model.Role{modelRole},
	}

	if err := u.db.Create(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-2] CreateUserAccount: %v", err)
		return err
	}

	currentTime := time.Now()
	modelVerify := model.VerificationToken{
		UserID:    modelUser.ID,
		Token:     req.Token,
		TokenType: "email_verification",
		ExpiresAt: currentTime.Add(1 * time.Hour),
		User:      modelUser,
	}

	if err := u.db.Create(&modelVerify).Error; err != nil {
		log.Errorf("[UserRepository-3] CreateUserAccount: %v", err)
		return err
	}

	return nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	modelUser := model.User{}

	if err := u.db.Where("email = ? AND is_verified = ?", email, true).
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
