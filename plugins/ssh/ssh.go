package main

import (
	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

type SSH struct {
	router *gin.RouterGroup
	db     *gorm.DB
}

func MakePlugin(params map[string]interface{}) plugins.Item {
	ssh := SSH{
		router: params["router"].(*gin.RouterGroup),
		db:     params["db"].(*gorm.DB),
	}

	ssh.migrateDB()
	ssh.initHTTP()

	return ssh
}

func (SSH) GetName() string {
	return "ssh"
}

func (SSH) GetDescription() string {
	return "Send files or execute commands over SSH"
}

func (SSH) GetVersion() string {
	return "1.0"
}

func (s SSH) GetBuildParams() []plugins.ItemParam {
	var servers []Server
	if err := s.db.Find(&servers).Error; err != nil {
		log.Println(err)
	}
	var paramOptions []plugins.ParamOptions
	for _, s := range servers {
		paramOptions = append(paramOptions, plugins.ParamOptions{
			Name:  s.Name,
			Value: strconv.Itoa(s.ID),
		})
	}

	return []plugins.ItemParam{
		{
			Name:    "server",
			Label:   "Server",
			Type:    plugins.TypeSelect,
			Options: paramOptions,
		},
		{
			Name:  "files",
			Label: "Files to send",
			Type:  plugins.TypeString,
		},
		{
			Name:  "remote_dir",
			Label: "Remote directory",
			Type:  plugins.TypeString,
		},
		{
			Name:  "command",
			Label: "Command",
			Type:  plugins.TypeText,
		},
	}
}

func (SSH) HasSettings() bool {
	return true
}
