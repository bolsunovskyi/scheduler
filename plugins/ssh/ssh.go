package main

import (
	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type SSH struct {
	router *gin.RouterGroup
	db     *gorm.DB
}

func MakePlugin(params map[string]interface{}) plugins.JobStep {
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

func (SSH) GetBuildParams() []plugins.StepParam {
	return []plugins.StepParam{
		{
			Name: "Command",
			Type: plugins.TypeString,
		},
	}
}

func (SSH) HasSettings() bool {
	return true
}
