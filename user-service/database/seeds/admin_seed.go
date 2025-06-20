package seeds

import (
	"log"
	"user-service/internal/core/domain/model"
	"user-service/utils/conv"

	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	bytes, err := conv.HashPassword("admin123")
	if err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	}

	modelRole := model.Role{}
	if err := db.Where("name = ?", "Super Admin").First(&modelRole).Error; err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	}

	admin := model.User{
		Name:       "super admin",
		Email:      "superadmin@mail.com",
		Password:   string(bytes),
		IsVerified: true,
		Roles:      []*model.Role{&modelRole},
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: admin.Email}).Error; err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	} else {
		log.Printf("Admin %s created successfully", admin.Name)
	}
}
