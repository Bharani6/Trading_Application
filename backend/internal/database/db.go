package database

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	lHost := viper.GetString("db.host")
	lPort := viper.GetString("db.port")
	lUser := viper.GetString("db.user")
	lPassword := viper.GetString("db.password")
	lDbname := viper.GetString("db.name")

	lDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		lHost, lUser, lPassword, lDbname, lPort)

	var lErr error
	DB, lErr = gorm.Open(postgres.Open(lDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if lErr != nil {
		log.Fatalf("Failed to connect to Postgres database: %v", lErr)
	}

	log.Println("Successfully connected to Postgres database")
}
