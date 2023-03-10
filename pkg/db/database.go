package db

import (
	"fmt"
	"os"
	"time"

	"github.com/nikhilnarayanan623/webApp/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

// to delete balck if the token time is expired
func DeleteBlackListToken() {

	DB.Where("end_time < ?", float64(time.Now().Unix())).Delete(&models.JwtBlackList{})

	fmt.Println("delted black listed token from database")
}

// connect to database
func ConnnectToDb() {

	dsn := os.Getenv("DATABASE")

	if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("Faild to Connect Database")
		return
	}

	fmt.Println("Successfully Connected to database")
}

// migrate table struct if there is table is available
func MigrateToDB() {

	if DB.AutoMigrate(&models.Admin{}, &models.User{}, &models.Product{}, &models.JwtBlackList{}); err != nil {
		fmt.Println("faild to sync database")
		return
	}

	fmt.Println("Successfully synced to database")
}
