package db

import (
	"gossip/backend/pkg/config"
	"gossip/backend/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	dbUrl := config.GetDbUrl()

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.SetupJoinTable(&models.User{}, "Subscribed", &models.Subscriptions{})
	if err != nil {
		panic("failed to setup subscriptions table [user model]")
	}

	err = db.SetupJoinTable(&models.Post{}, "Tags", &models.Taggable{})
	if err != nil {
		panic("failed to setup taggable table [post model]")
	}

	err = db.SetupJoinTable(&models.Tag{}, "Posts", &models.Taggable{})
	if err != nil {
		panic("failed to setup taggable table [tag model]")
	}

	db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Tag{},
	)

	return db
}
