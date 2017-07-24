package jobs

import (
	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"net/http"
)

func makePluginSchemaHandler(loadedPlugins []plugins.Item) gin.HandlerFunc {
	return func(c *gin.Context) {
		if pluginName := c.Query("name"); pluginName != "" {
			for _, v := range loadedPlugins {
				if v.GetName() == pluginName {
					c.JSON(http.StatusOK, v.GetBuildParams())
					return
				}
			}
		}

		c.JSON(http.StatusNotFound, map[string]string{"message": "Plugin not found"})
	}
}
