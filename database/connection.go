package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"simple-api/models"
)

var DB *gorm.DB

func SetupConnection()  {
	connection, err := gorm.Open(mysql.Open("sammy:Password@123@/simple_api"), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database")
	}

	DB = connection

	if err = connection.AutoMigrate(&models.User{}, &models.Product{}, models.Cart{}, models.CartItem{}); err != nil {
		log.Fatal(1)
	}
}
