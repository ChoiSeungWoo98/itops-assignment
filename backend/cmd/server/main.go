package main

import (
	"itops-backend/pkg/db"
	"itops-backend/pkg/models"
	"itops-backend/pkg/routes"
)

func main() {
	db.Init()
	db.DB.AutoMigrate(&models.User{}, &models.Issue{})

	// 초기 사용자 삽입
	seed := []models.User{{1, "김개발"}, {2, "이디자인"}, {3, "박기획"}}
	for _, u := range seed { db.DB.FirstOrCreate(&u, u) }

	r := routes.Setup()
	r.Run(":8080")
}