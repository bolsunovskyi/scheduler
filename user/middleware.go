package user

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Middleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie("session")
		if err != nil {
			log.Println(err)
			c.Next()
			return
		}

		var s Session
		if err := db.Where("value = ?", session).First(&s).Error; err != nil {
			log.Println(err)
			c.Next()
			return
		}

		var u Model
		if err := db.Where("id = ?", s.UserID).First(&u).Error; err != nil {
			log.Println(err)
			c.Next()
			return
		}

		s.CreatedAt = time.Now()
		if err := db.Save(s).Error; err != nil {
			log.Println(err)
			return
		}

		c.Set("user", &u)
		c.Next()
	}
}
