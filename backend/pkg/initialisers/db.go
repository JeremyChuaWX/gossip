package initialisers

import (
	"fmt"
	"gossip/backend/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(config *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
	)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	return DB
}

func MigrateDB(DB *gorm.DB) {
	var err error

	err = DB.SetupJoinTable(&models.User{}, "Subscribed", &models.Subscriptions{})
	if err != nil {
		panic("failed to setup subscriptions table [user model]")
	}

	err = DB.SetupJoinTable(&models.Post{}, "Tags", &models.Taggable{})
	if err != nil {
		panic("failed to setup taggable table [post model]")
	}

	err = DB.SetupJoinTable(&models.Tag{}, "Posts", &models.Taggable{})
	if err != nil {
		panic("failed to setup taggable table [tag model]")
	}

	DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Tag{},
	)
}
