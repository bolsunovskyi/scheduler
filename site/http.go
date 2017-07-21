package site

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitHTTP(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "site/index.html", gin.H{})
	})
}
