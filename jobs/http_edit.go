package jobs

import (
	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func makeToodleEnabledHandler(db *gorm.DB) {

}

func makeEditGetHandler(db *gorm.DB, loadedPlugins []plugins.Item) gin.HandlerFunc {
	return func(c *gin.Context) {
		j, err := getJobByID(db, c.Param("id"))
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
			"plugins":      loadedPlugins,
			"job":          j,
			"editControls": true,
		})
	}
}
