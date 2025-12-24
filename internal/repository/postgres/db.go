package postgres

import (
	"context"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(ctx context.Context) (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", host, user, password, name, port)
	fmt.Println(dsn)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to Database successfully!")
	db, err := DB.DB()
	if err != nil {
		return nil, err
	}
	err = db.PingContext(ctx)
	if err != nil {

		return nil, err
	}
	fmt.Println("Database Migrated")
	return DB, nil
}
