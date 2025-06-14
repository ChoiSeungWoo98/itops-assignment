package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("issue.db"), &gorm.Config{})
	if err != nil { log.Fatalf("DB 연결 실패: %v", err) }
}