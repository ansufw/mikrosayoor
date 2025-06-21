package seeds

import (
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"

	"user-service/internal/core/domain/model"
)

func SeedRole(db *gorm.DB) {
	roles := []model.Role{
		{Name: "Super Admin"},
		{Name: "Customer"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Role{Name: role.Name}).Error; err != nil {
			log.Fatalf("%s: %v", "SeedRole", err)
		} else {
			log.Printf("Role %s created", role.Name)
		}
	}

}
