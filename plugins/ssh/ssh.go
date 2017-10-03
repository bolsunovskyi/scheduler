package main

import (
	"log"
	"strconv"

	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/natefinch/pie"
	"net/rpc/jsonrpc"
)

type SSH struct {
	router *gin.RouterGroup
	db     *gorm.DB
}

func (s SSH) InitDB(dbPath string, _ *string) error {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	s.db = db
	if err := s.migrateDB(); err != nil {
		return err
	}

	return nil
}

func (s SSH) HandleHTTP(rq plugins.HTTPRequest, rsp *plugins.HTTPResponse) error {
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

func (s SSH) GetBuildParams(dbPath string, rsp *[]plugins.ItemParam) error {
	if s.db == nil {
		db, err := gorm.Open("sqlite3", dbPath)
		if err != nil {
			return err
		}
		s.db = db
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
	p.Serve()
}
