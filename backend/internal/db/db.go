package db

import (
	"fmt"
	"gossip/backend/internal/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDsn() string {
	godotenv.Load()

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("HOST"),
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DBNAME"),
		os.Getenv("PORT"),
	)
}

func Init() *gorm.DB {
	dsn := getDsn()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
