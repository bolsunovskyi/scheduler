package main

import (
	"errors"
	"log"
	"net/rpc/jsonrpc"
	"strconv"

	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/natefinch/pie"
)

type SSH struct {
	router *gin.RouterGroup
	db     *gorm.DB
}

func (s SSH) InitDB(_ string, db *gorm.DB) error {
	s.db = db
	return s.migrateDB()
}

func (s SSH) InitRouter(_ string, router *gin.RouterGroup) error {
	s.router = router
	s.initHTTP()
	return nil
}

func (SSH) GetPluginParams(_ string, params *plugins.PluginParams) error {
	*params = plugins.PluginParams{
		Name:        "ssh",
		Description: "Send files or execute commands over SSH",
		Version:     "1.0",
		HasSettings: true,
	}
	return nil
}

func (s SSH) GetBuildParams(_ string, rsp *[]plugins.ItemParam) error {
	if s.db == nil {
		return errors.New("The database was not initialized")
		//s.db = &db
	}

	var servers []Server
	if err := s.db.Find(&servers).Error; err != nil {
		return err
	}
	var paramOptions []plugins.ParamOptions
	for _, s := range servers {
		paramOptions = append(paramOptions, plugins.ParamOptions{
			Name:  s.Name,
			Value: strconv.Itoa(s.ID),
		})
	}

	*rsp = []plugins.ItemParam{
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
	return nil
}

func main() {
	p := pie.NewProvider()
	if err := p.RegisterName("ssh", SSH{}); err != nil {
		log.Fatalln(err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}
