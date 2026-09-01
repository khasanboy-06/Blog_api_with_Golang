package database


import (
	"fmt"
	"log"
	"os"

	"Blog_project_with_Go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect (){
	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tashkent",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil{
		log.Fatal("Databasega ulanishda xatolik: ", err)
	}

	err = DB.AutoMigrate(&models.Post{}, &models.User{})

	if err != nil {
		log.Fatal("Migrate xatosi: ", err)
	}

	fmt.Println("Databasega muvaffaqqiyatli ulandi")
}