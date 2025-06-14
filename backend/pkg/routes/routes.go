package routes

import (
	"github.com/gin-gonic/gin"
	"itops-backend/pkg/handlers"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.POST("/issue",  handlers.CreateIssue)
	r.GET("/issues", handlers.ListIssues)
	r.GET("/issue/:id", handlers.GetIssue)
	r.PATCH("/issue/:id", handlers.UpdateIssue)
	return r
}