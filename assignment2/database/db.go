package database

import (
	"assignment2/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=postgres dbname=assignment2 port=5432 sslmode=disable"
	connections, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = connections

	err = connections.AutoMigrate(&models.Items{}, &models.Order{})
	if err != nil {
		return
	}
}
