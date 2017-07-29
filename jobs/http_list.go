package jobs

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func makeListHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var j []Model
		if err := db.Find(&j).Error; err != nil {
			log.Println(err)
			c.HTML(http.StatusOK, "jobs/index.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "jobs/index.html", gin.H{
			"jobs": j,
		})
	}
}
