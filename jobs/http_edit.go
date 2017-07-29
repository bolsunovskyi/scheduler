package jobs

import (
	"encoding/json"
	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func makeToodleEnabledHandler(db *gorm.DB) {

}

func makeEditGetHandler(db *gorm.DB, loadedPlugins []plugins.Item) gin.HandlerFunc {
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

		if c.Request.Header.Get("Accept") == "application/json" {
			for i, step := range j.Steps {
				for _, p := range loadedPlugins {
					if step.PluginName == p.GetName() {
						j.Steps[i].Description = p.GetDescription()
					}
				}
			}
			c.JSON(http.StatusOK, j)
			return
		}

		c.HTML(http.StatusOK, "jobs/create.html", gin.H{
			"plugins": loadedPlugins,
			"job":     j,
		})
	}
}
