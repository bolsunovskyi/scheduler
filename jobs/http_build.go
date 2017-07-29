package jobs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func makeBuildGetHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		j, err := getJobByID(db, c.Param("id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/a/jobs")
			return
		}

		if len(j.Params) > 0 {
			c.HTML(http.StatusOK, "jobs/build.html", gin.H{
				"job": j,
			})
			return
		}

		u, _ := c.Get("userID")

		if err := j.build(db, map[string]string{}, u.(int)); err != nil {
			log.Println(err)
		}

		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/a/jobs/status/%d", j.ID))
	}
}

func makeBuildPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		j, err := getJobByID(db, c.Param("id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/a/jobs")
			return
		}

		c.Request.ParseForm()
		params := make(map[string]string)
		for k, v := range c.Request.Form {
			params[k] = v[0]
		}

		u, _ := c.Get("userID")

		if err := j.build(db, params, u.(int)); err != nil {
			log.Println(err)
		}

		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/a/jobs/status/%d", j.ID))
	}
}
