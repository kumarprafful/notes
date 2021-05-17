package controllers

import (
	"notes/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
}

func (s *Server) DBConnect() {
	var err error
	// dsn := "host=" + os.Getenv("HOST") + " user=" + os.Getenv("POSTGRES_USER") + " dbname=" + os.Getenv("POSTGRES_DB") + " port=5432 sslmode=disable password=" + os.Getenv("POSTGRES_PASSWORD")
	dsn := "host=localhost user=prafful password=password dbname=notes port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	// DBURL:= fmt.Sprintf("host=%s port%s user=%s dbname=%s sslmode=disable password=%s", )
	s.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to DB")
	}

	s.DB.AutoMigrate(&models.User{})
}
