package database

import (
	"log"

	"simple-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupConnection(dbDns string)  {
	connection, err := gorm.Open(mysql.Open(dbDns), &gorm.Config{})
	if err != nil {
		panic("could not connect to the database")
	}

	DB = connection

	if err = connection.AutoMigrate(&models.User{}, &models.Product{}, models.Cart{}, models.CartItem{}); err != nil {
		log.Fatal(1)
	}
}
