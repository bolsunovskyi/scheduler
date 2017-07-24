package plugins

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func InitHTTP(r *gin.RouterGroup, db *gorm.DB, loadedPlugins []Item) {
	r.GET("/plugins", func(c *gin.Context) {
		c.HTML(http.StatusOK, "plugins/index.html", gin.H{
			"plugins": loadedPlugins,
		})
	})
}
