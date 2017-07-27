package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	vd "gopkg.in/go-playground/validator.v9"
)

func (s SSH) initHTTP() {
	s.router.GET("/", func(c *gin.Context) {
		var creds []Credential
		if err := s.db.Find(&creds).Error; err != nil {
			c.HTML(http.StatusOK, "ssh/index.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "ssh/index.html", gin.H{
			"creds": creds,
		})
	})

	s.router.GET("/credential/delete", func(c *gin.Context) {
		if id := c.Query("id"); id != "" {
			if n, err := strconv.Atoi(id); err == nil {
				if err := s.db.Where("id = ?", n).Delete(&Credential{}).Error; err != nil {
					log.Println(err)
				}
			}
		}

		c.Redirect(http.StatusSeeOther, "/a/plugins/ssh")
	})

	s.router.POST("/", func(c *gin.Context) {
		cred := Credential{
			Name:       c.PostForm("name"),
			Username:   c.PostForm("username"),
			Password:   c.PostForm("password"),
			PrivateKey: c.PostForm("key"),
		}

		if err := vd.New().Struct(&cred); err != nil {
			c.HTML(http.StatusOK, "ssh/index.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		if cred.Password == "" && cred.PrivateKey == "" {
			c.HTML(http.StatusOK, "ssh/index.html", gin.H{
				"error": "Password or key must be specified.",
			})
			return
		}

		if err := s.db.Save(&cred).Error; err != nil {
			c.HTML(http.StatusOK, "ssh/index.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Redirect(http.StatusSeeOther, "/a/plugins/ssh")
	})

	s.router.GET("/servers/delete", func(c *gin.Context) {
		if id := c.Query("id"); id != "" {
			if n, err := strconv.Atoi(id); err == nil {
				if err := s.db.Where("id = ?", n).Delete(&Server{}).Error; err != nil {
					log.Println(err)
				}
			}
		}

		c.Redirect(http.StatusSeeOther, "/a/plugins/ssh/servers/")
	})

	s.router.GET("/servers", func(c *gin.Context) {
		var creds []Credential
		if err := s.db.Find(&creds).Error; err != nil {
			c.HTML(http.StatusOK, "ssh/index.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		var servers []Server
		//if err := s.db.Find(&servers).Error; err != nil {
		//	c.HTML(http.StatusOK, "ssh/index.html", gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}

		serverTbl := Server{}.TableName()
		credTbl := Credential{}.TableName()

		rows, err := s.db.Table(serverTbl).
			Select(fmt.Sprintf("%s.*, %s.name", serverTbl, credTbl)).
			Joins(fmt.Sprintf(
				"INNER JOIN %s ON %s.credential_id = %s.id", credTbl, serverTbl, credTbl)).Rows()
		if err != nil {
			c.HTML(http.StatusOK, "ssh/servers.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		for rows.Next() {
			s := Server{}
			if err := rows.Scan(&s.ID, &s.Name, &s.Host, &s.Port, &s.CredentialID, &s.CredentialName); err != nil {
				c.HTML(http.StatusOK, "ssh/servers.html", gin.H{
					"error": err.Error(),
				})
				return
			}
			servers = append(servers, s)
		}

		c.HTML(http.StatusOK, "ssh/servers.html", gin.H{
			"creds":   creds,
			"servers": servers,
		})
	})

	s.router.POST("/servers", func(c *gin.Context) {
		var creds []Credential
		if err := s.db.Find(&creds).Error; err != nil {
			c.HTML(http.StatusOK, "ssh/index.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		port, err := strconv.Atoi(c.PostForm("port"))
		credID, err := strconv.Atoi(c.PostForm("credential"))

		if err != nil {
			c.HTML(http.StatusOK, "ssh/servers.html", gin.H{
				"error": err.Error(),
				"creds": creds,
			})
			return
		}

		srv := Server{
			Name:         c.PostForm("name"),
			Host:         c.PostForm("host"),
			Port:         port,
			CredentialID: credID,
		}

		if err := vd.New().Struct(&srv); err != nil {
			c.HTML(http.StatusOK, "ssh/servers.html", gin.H{
				"error": err.Error(),
				"creds": creds,
			})
			return
		}

		if err := s.db.Where("id = ?", srv.CredentialID).First(&Credential{}).Error; err != nil {
			c.HTML(http.StatusOK, "ssh/servers.html", gin.H{
				"error": err.Error(),
				"creds": creds,
			})
			return
		}

		if err := s.db.Save(&srv).Error; err != nil {
			c.HTML(http.StatusOK, "ssh/servers.html", gin.H{
				"error": err.Error(),
				"creds": creds,
			})
			return
		}

		c.Redirect(http.StatusSeeOther, "/a/plugins/ssh/servers/")
	})
}
