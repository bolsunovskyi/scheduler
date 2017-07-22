package jobs

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func InitHTTP(r *gin.RouterGroup, db *gorm.DB) {
	r.GET("/jobs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "jobs/index.html", gin.H{})
	})

	r.GET("/jobs/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "jobs/create.html", gin.H{})
	})
}
