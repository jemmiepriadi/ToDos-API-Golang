package model

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	errDb error
)

func ConnectDataBase() {
	dsn := "host=localhost user=postgres password=password123 dbname=todo port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Profile{})

	DB = database
}

type Profile struct {
	ID          int
	Nama        string    `json:"Name" example:"Jemmi"`
	Date        time.Time `example:"0001-01-01T07:00:00+07:00"`
	DateInput   string    `json:"DateInput" example:"2022-07-05"`
	Description string    `json:"Description" example:"Berenang"`
}
