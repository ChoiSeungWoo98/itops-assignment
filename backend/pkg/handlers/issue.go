package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"itops-backend/pkg/db"
	"itops-backend/pkg/models"
)

var valid = map[string]bool{"PENDING": true, "IN_PROGRESS": true, "COMPLETED": true, "CANCELLED": true}

// POST /issue
func CreateIssue(c *gin.Context) {
	var req struct {
		Title, Description string
		UserID             *uint `json:"userId"`
	}
	if c.BindJSON(&req) != nil || req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "입력 오류"}); return
	}
	status := "PENDING"
	if req.UserID != nil { status = "IN_PROGRESS" }

	issue := models.Issue{
		Title: req.Title, Description: req.Description, Status: status,
		UserID: req.UserID, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	db.DB.Create(&issue).Preload("User").First(&issue)
	c.JSON(http.StatusCreated, issue)
}

// GET /issues
func ListIssues(c *gin.Context) {
	var list []models.Issue
	if s := c.Query("status"); s != "" {
		db.DB.Preload("User").Where("status = ?", s).Find(&list)
	} else {
		db.DB.Preload("User").Find(&list)
	}
	c.JSON(http.StatusOK, gin.H{"issues": list})
}

// GET /issue/:id
func GetIssue(c *gin.Context) {
	var i models.Issue
	if db.DB.Preload("User").First(&i, c.Param("id")).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "없음"}); return
	}
	c.JSON(http.StatusOK, i)
}

// PATCH /issue/:id
func UpdateIssue(c *gin.Context) {
	var i models.Issue
	if db.DB.First(&i, c.Param("id")).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "없음"}); return
	}
	if i.Status == "COMPLETED" || i.Status == "CANCELLED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "수정 불가"}); return
	}

	var req struct {
		Title, Description, Status *string
		UserID                     *uint `json:"userId"`
	}
	if c.BindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "입력 오류"}); return
	}

	if req.Title != nil { i.Title = *req.Title }
	if req.Description != nil { i.Description = *req.Description }
	if req.UserID != nil {
		if *req.UserID == 0 { i.UserID = nil; i.Status = "PENDING" } else { i.UserID = req.UserID }
	}
	if req.Status != nil {
		if !valid[*req.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 상태"}); return
		}
		if i.UserID == nil && (*req.Status == "IN_PROGRESS" || *req.Status == "COMPLETED") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "담당자 없음"}); return
		}
		i.Status = *req.Status
	}
	i.UpdatedAt = time.Now()
	db.DB.Save(&i).Preload("User").First(&i)
	c.JSON(http.StatusOK, i)
}