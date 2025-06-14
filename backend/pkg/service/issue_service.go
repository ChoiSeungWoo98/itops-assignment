package service

import (
	"fmt"
	"time"
	"github.com/ChoiSeungWoo98/itops-assignment/internal/models"
	"github.com/ChoiSeungWoo98/itops-assignment/internal/repository"
)

var valid = map[string]bool{"PENDING": true,"IN_PROGRESS": true,"COMPLETED": true,"CANCELLED": true}
var repo = repository.IssueRepo{}

func CreateIssue(req models.Issue) (models.Issue, error) {
	if req.UserID == nil { req.Status = "PENDING" } else { req.Status = "IN_PROGRESS" }
	req.CreatedAt, req.UpdatedAt = time.Now(), time.Now()
	repo.Create(&req); return req, nil
}

func UpdateIssue(src *models.Issue, req models.Issue) error {
	if src.Status == "COMPLETED" || src.Status == "CANCELLED" { return  fmt.Errorf("수정 불가") }

	if req.Title != ""       { src.Title = req.Title }
	if req.Description != "" { src.Description = req.Description }
	if req.UserID != nil {
		if *req.UserID == 0 { src.UserID = nil; src.Status = "PENDING" } else { src.UserID = req.UserID }
	}
	if req.Status != "" {
		if !valid[req.Status] { return fmt.Errorf("잘못된 상태") }
		if src.UserID == nil && (req.Status == "IN_PROGRESS" || req.Status == "COMPLETED") {
			return fmt.Errorf("담당자 없음")
		}
		src.Status = req.Status
	}
	src.UpdatedAt = time.Now()
	repo.Save(src); return nil
}