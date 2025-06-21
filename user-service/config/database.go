package config

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"user-service/database/seeds"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Psql.User,
		cfg.Psql.Password,
		cfg.Psql.Host,
		cfg.Psql.Port,
		cfg.Psql.DBName,
	)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Errorf("ConnectionPostgres-1 Failed to connect to database " + cfg.Psql.Host)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("ConnectionPostgres-2 Failed to get database connection")
		return nil, err
	}

	seeds.SeedRole(db)
	seeds.SeedAdmin(db)

	sqlDB.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Psql.DBMaxIdle)

	return &Postgres{DB: db}, nil
}
