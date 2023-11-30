package database

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB() (*gorm.DB, error) {

	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	dbName := viper.GetString("database.dbname")
	dbUser := viper.GetString("database.user")
	password := viper.GetString("database.password")

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

	// db.AutoMigrate(&model.NFAD{})

	return db, nil
}
