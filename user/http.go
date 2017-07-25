package user

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
)

type loginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func makeLoginFunc(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rq := loginRequest{
			Email:    c.PostForm("email"),
			Password: c.PostForm("password"),
		}

		if err := validator.New().Struct(&rq); err != nil {
			c.HTML(http.StatusOK, "user/index.html", gin.H{"error": err.Error(), "email": rq.Email})
			return
		}

		var u Model
		if err := db.Where("email = ?", rq.Email).First(&u).Error; err != nil {
			c.HTML(http.StatusOK, "user/index.html", gin.H{
				"error": "Email or password is wrong", "email": rq.Email})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rq.Password)); err != nil {
			c.HTML(http.StatusOK, "user/index.html", gin.H{
				"error": "Email or password is wrong", "email": rq.Email})
			return
		}

		node, _ := snowflake.NewNode(time.Now().Unix() % 1022)

		s := Session{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IP:        c.Request.RemoteAddr,
			UserID:    u.ID,
			Value:     node.Generate().Base36(),
		}
		if err := db.Save(&s).Error; err != nil {
			c.HTML(http.StatusOK, "user/index.html", gin.H{
				"error": err.Error(), "email": rq.Email})
			return
		}

		c.SetCookie("session", s.Value, int(24*3600),
			"/", "", false, false)
		c.Redirect(http.StatusSeeOther, "/a/jobs")
	}
}

func InitHTTP(r *gin.Engine, db *gorm.DB) {
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404", gin.H{})
	})
	r.NoMethod(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404", gin.H{})
	})
	r.GET("/", func(c *gin.Context) {
		if _, e := c.Get("user"); e {
			c.Redirect(http.StatusSeeOther, "/a/jobs")
			return
		}
		c.HTML(http.StatusOK, "user/index.html", gin.H{})
	})
	r.POST("/", makeLoginFunc(db))

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("session", "", -10, "/", "", false, false)
		c.Redirect(http.StatusSeeOther, "/")
	})
}
