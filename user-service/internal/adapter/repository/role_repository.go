package repository

import (
	"context"
	"errors"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"

	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"
)

type RoleRepositoryI interface {
	GetAll(ctx context.Context, search string) ([]entity.RoleEntity, error)
	GetByID(ctx context.Context, id int64) (*entity.RoleEntity, error)
	Create(ctx context.Context, req entity.RoleEntity) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, req entity.RoleEntity) error
}

type roleRepository struct {
	db *gorm.DB
}

// Create implements RoleRepositoryI.
func (r *roleRepository) Create(ctx context.Context, req entity.RoleEntity) error {
	modelRole := model.Role{
		Name: req.Name,
	}

	if err := r.db.Create(&model.User).Error; err != nil {
		log.Errorf("[RoleRepository-1] Create: %v", err)
		return err
	}
	return nil
}

// Delete implements RoleRepositoryI.
func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	modelRole := model.Role{}

	if err := r.db.Where("id = ?", id).Preload().Error; err != nil {
		log.Errorf("[RoleRepository-1] Delete: %v", err)
		return err
	}

	return nil
}

// GetAll implements RoleRepositoryI.
func (r *roleRepository) GetAll(ctx context.Context, search string) ([]entity.RoleEntity, error) {
	modelRoles := []model.Role{}

	if err := r.db.Where("name ILIKE ?", "%"+search+"%").Find(&modelRoles).Error; err != nil {
		log.Errorf("[RoleRepository-1] GetAll: %v", err)
		return nil, err
	}

	if len(modelRoles) == 0 {
		err := errors.New("404")
		log.Errorf("[RoleRepository-2] GetAll: %v", err)
		return nil, err
	}

	entityRole := []entity.RoleEntity{}
	for _, role := range modelRoles {
		entityRole = append(entityRole, entity.RoleEntity{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	return entityRole, nil
}

// GetByID implements RoleRepositoryI.
func (r *roleRepository) GetByID(ctx context.Context, id int64) (*entity.RoleEntity, error) {
	modelRole := model.Role{}

	if err := r.db.Where("id = ?", id).First(&modelRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New("404")
			log.Errorf("[RoleRepository-1] GetByID: %v", err)
			return nil, err
		}

		log.Errorf("[RoleRepository-2] GetByID: %v", err)
		return nil, err
	}

	return &entity.RoleEntity{
		ID:   modelRole.ID,
		Name: modelRole.Name,
	}, nil

}

// Update implements RoleRepositoryI.
func (r *roleRepository) Update(ctx context.Context, req entity.RoleEntity) error {
	panic("unimplemented")
}

func NewRoleRepository(db *gorm.DB) RoleRepositoryI {
	return &roleRepository{
		db: db,
	}
}
