package jobs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func makeStatusGetHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/a/jobs")
			return
		}

		var j Model
		if err := db.Where("id = ?", id).First(&j).Error; err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/a/jobs")
			return
		}

		err = json.Unmarshal([]byte(j.StepsEncoded), &j.Steps)
		err = json.Unmarshal([]byte(j.PramsEncoded), &j.Params)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/a/jobs")
			return
		}

		c.HTML(http.StatusOK, "jobs/status.html", gin.H{
			"job": j,
		})
	}
}
