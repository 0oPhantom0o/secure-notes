package app

import (
	p "github.com/secure-notes/internal/repository/postgres"
	"gorm.io/gorm"
)

func Db() (gorm.DB, error) {

	return p.NewDB()
}
