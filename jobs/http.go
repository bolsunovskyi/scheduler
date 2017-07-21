package jobs

import (
	"github.com/bolsunovskyi/scheduler/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func InitHTTP(r *gin.Engine, db *gorm.DB) {
	auth := r.Group("/a")
	auth.Use(user.AbortUnAuth())

	auth.GET("/jobs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "jobs/index.html", gin.H{})
	})
}
