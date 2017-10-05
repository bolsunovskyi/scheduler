package jobs

import (
	"encoding/json"
	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func makeCreatePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var j Model
		if err := json.NewDecoder(c.Request.Body).Decode(&j); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
			return
		}

		bts, err := json.Marshal(j.Params)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
			return
		}
		j.PramsEncoded = string(bts)

		bts, err = json.Marshal(j.Steps)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
			return
		}
		j.StepsEncoded = string(bts)

		if err := db.Save(&j).Error; err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, j)
	}
}

func makeCreateGetHandler(loadedPlugins []plugins.Params) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "jobs/create.html", gin.H{
			"plugins": loadedPlugins,
		})
	}
}
