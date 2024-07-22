package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sandipbera35/jwt_authservice/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Connect *gorm.DB

func DbConnect() error {
	errr := godotenv.Load(".env")
	if errr != nil {
		log.Fatalf("Error loading .env file")
	}

	dialect := os.Getenv("DIALECT")
	_ = dialect
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")

	port, _ := strconv.Atoi(dbPort)

	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	dburi := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		host, user, password, dbName, port)
	var err error
	if strings.ToUpper(os.Getenv("DBLOGTYPE")) == "INFO" {
		Connect, err = gorm.Open(postgres.Open(dburi), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else if strings.ToUpper(os.Getenv("DBLOGTYPE")) == "WARNING" {
		Connect, err = gorm.Open(postgres.Open(dburi), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})

	} else {
		Connect, err = gorm.Open(postgres.Open(dburi), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
	}

	if err != nil {
		return err
	}

	if os.Getenv("MIGRATION") == "true" {
		Migrate(models.User{})
	}
	log.Println("Database connection was successful!!")
	return nil

}

func Migrate(User models.User) {
	Connect.Exec("ALTER DATABASE " + os.Getenv("NAME") + " SET timezone = 'UTC'")
	Connect.AutoMigrate(&User)
	log.Println("DataBase migration was successfull ")
}
