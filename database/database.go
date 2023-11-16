package database

import (
	"fmt"
	"log"

	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB() (*gorm.DB, error) {
	host := "localhost"
	port := "5432"
	dbName := "lolkek"
	dbUser := "postgres"
	password := "Primlab2020"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
		return nil, err
	}

	fmt.Println("Database connection successful...")

	db.AutoMigrate(&model.NFAD{})

	return db, nil
}
