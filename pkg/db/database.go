package db

import (
	"fmt"
	"os"
	"webApp/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConnnectToDb() {

	dsn := os.Getenv("DATABASE")

	if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("Faild to Connect Database")
		return
	}

	fmt.Println("Successfully Connected to database")
}

func MigrateToDB() {

	if DB.AutoMigrate(&models.Admin{}, &models.User{}, &models.JwtBlackList{}); err != nil {
		fmt.Println("faild to sync database")
		return
	}

	fmt.Println("Successfully synced to database")
}