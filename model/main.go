package model

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LineUser struct {
	gorm.Model
	UserId string
	InstallationId int
	DefaultScheduleTitle string
}

func Migrate() {
	db, err := DbOpen()
	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&LineUser{})
}

func DbOpen() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		return db, err
	}
	return db, err
}
