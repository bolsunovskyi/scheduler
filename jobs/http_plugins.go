package jobs

import (
	"net/http"

	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
)

func makePluginSchemaHandler(loadedPlugins []plugins.Params) gin.HandlerFunc {
	return func(c *gin.Context) {
		if pluginName := c.Param("name"); pluginName != "" {
			for _, v := range loadedPlugins {
				if v.Name == pluginName {
					buildSteps, err := v.GetBuildSteps()
					if err == nil {
						c.JSON(http.StatusOK, buildSteps)
						return
					}
				}
			}
		}

		c.JSON(http.StatusNotFound, map[string]string{"message": "Plugin not found"})
	}
}
