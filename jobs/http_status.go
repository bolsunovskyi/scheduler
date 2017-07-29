package jobs

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"math"
	"net/http"
	"strconv"
)

const buildsPerPage = 25

func makeStatusGetHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		j, err := getJobByID(db, c.Param("id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/a/jobs")
			return
		}

		page := 1
		if pageStr, ok := c.GetQuery("page"); ok {
			if pageN, err := strconv.Atoi(pageStr); err == nil {
				page = pageN
			}
		}

		offset := (page - 1) * buildsPerPage
		var builds []Build
		if err := db.Where("job_id = ?", j.ID).Order("number DESC").
			Limit(buildsPerPage).Offset(offset).Find(&builds).Error; err != nil {
			log.Println(err)
		}
		var count int
		if err := db.Table(Build{}.TableName()).Where("job_id = ?", j.ID).Count(&count).Error; err != nil {
			log.Println(err)
		}

		count = int(math.Ceil(float64(count) / float64(buildsPerPage)))
		var pages []int
		for i := 1; i <= count; i++ {
			pages = append(pages, i)
		}

		params := gin.H{
			"job":    j,
			"builds": builds,
			"pages":  pages,
			"page":   page,
			"count":  count,
		}

		if page > 1 {
			params["prevPage"] = page - 1
		}

		if page < count {
			params["nextPage"] = page + 1
		}

		c.HTML(http.StatusOK, "jobs/status.html", params)
	}
}
