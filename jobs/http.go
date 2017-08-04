package jobs

import (
	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Init(r *gin.RouterGroup, db *gorm.DB, loadedPlugins []plugins.Item) error {
	if err := db.AutoMigrate(&Model{}, &Build{}, &Tab{}).Error; err != nil {
		return err
	}

	r.GET("/jobs", makeListHandler(db))
	r.GET("/jobs/plugins/schema/:name", makePluginSchemaHandler(loadedPlugins))
	r.GET("/jobs/create", makeCreateGetHandler(loadedPlugins))
	r.POST("/jobs/create", makeCreatePostHandler(db))
	r.GET("/jobs/edit/:id", makeEditGetHandler(db, loadedPlugins))
	r.GET("/jobs/status/:id", makeStatusGetHandler(db))
	r.GET("/jobs/build/:id", makeBuildGetHandler(db))
	r.POST("/jobs/build/:id", makeBuildPostHandler(db))

	go parseBuildQueue(db)

	return nil
}
