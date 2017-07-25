package jobs

import (
	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func makeCreatePostHandler(db *gorm.DB, loadedPlugins []plugins.Item) gin.HandlerFunc {
	return func(c *gin.Context) {
		//name := c.PostForm("name")
		//params := c.PostFormArray("params[]")
		//fmt.Println(name)
		//fmt.Println(params)
		//fmt.Printf("%+v\n", c.Request.PostForm)

		c.HTML(http.StatusOK, "jobs/create.html", gin.H{
			"plugins": loadedPlugins,
		})
	}
}

func makeCreateGetHandler(db *gorm.DB, loadedPlugins []plugins.Item) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "jobs/create.html", gin.H{
			"plugins": loadedPlugins,
		})
	}
}
