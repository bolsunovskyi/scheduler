package plugins

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func InitHTTP(r *gin.RouterGroup, db *gorm.DB, loadedPlugins []Params) {
	r.GET("/plugins", func(c *gin.Context) {
		c.HTML(http.StatusOK, "plugins/index.html", gin.H{
			"plugins": loadedPlugins,
		})
	})
}

func InitTemplate() (*template.Template, error) {
	tpl := template.New("app").Delims("[[", "]]").Funcs(map[string]interface{}{
		"split": func(s string) []string {
			return strings.Split(s, "\n")
		},
		"time": func(t time.Time) string {
			return t.Format("2.01.2006 15:04:05")
		},
	})

	if _, err := tpl.ParseGlob("./_templates/**/*"); err != nil {
		return nil, err
	}

	return tpl, nil
}
