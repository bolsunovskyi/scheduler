package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/rpc/jsonrpc"
	"os"
	"strconv"
	"strings"

	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/natefinch/pie"
)

type SSH struct {
	router *gin.RouterGroup
	db     *gorm.DB
}

func (s SSH) InitDB(dbPath string, _ *string) error {
	if s.db == nil {
		db, err := gorm.Open("sqlite3", dbPath)
		if err != nil {
			return err
		}

		s.db = db
	}

	if err := s.migrateDB(); err != nil {
		return err
	}

	return nil
}

func (s SSH) HandleHTTP(rq plugins.HTTPRequest, rsp *plugins.HTTPResponse) error {
	if s.db == nil {
		db, err := gorm.Open("sqlite3", rq.DBPath)
		if err != nil {
			return err
		}

		s.db = db
	}

	hrq, err := http.NewRequest(rq.Method, rq.URL.Path+"?"+rq.URL.RawQuery, strings.NewReader(rq.BodyStr))
	if err != nil {
		return err
	}
	hrq.Header = rq.Header

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.NoRoute(func(c *gin.Context) {
		c.String(404, fmt.Sprintf("Path [%s] not found", c.Request.URL.Path))
	})

	tpl, err := plugins.InitTemplate()
	if err != nil {
		return err
	}

	if _, err := tpl.ParseGlob("./_templates/**/*"); err != nil {
		log.Fatalln(err)
	}

	if _, err := os.Stat("./plugins/ssh/_templates"); err == nil {
		if _, err := tpl.ParseGlob("./plugins/ssh/_templates/*"); err != nil {
			log.Fatalln(err)
		}
	}

	//load assets
	if _, err := os.Stat("./plugins/ssh/_assets"); err == nil {
		r.Static("ssh/assets/", "./plugins/ssh/_assets")
	}

	r.SetHTMLTemplate(tpl)

	s.router = r.Group("/a/plugins/ssh")
	s.initHTTP()
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, hrq)

	rsp.StatusCode = rr.Code
	rsp.RawBody = rr.Body.String()
	rsp.Header = rr.HeaderMap
	rsp.Raw = true

	return nil
}

func (s SSH) GetPluginParams(dbPath string, params *plugins.Params) error {
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
	var paramOptions []plugins.BuildStepOptions
	for _, s := range servers {
		paramOptions = append(paramOptions, plugins.BuildStepOptions{
			Name:  s.Name,
			Value: strconv.Itoa(s.ID),
		})
	}

	*params = plugins.Params{
		Name:        "ssh",
		Description: "Send files or execute commands over SSH",
		Version:     "1.0",
		HasSettings: true,
		BuildSteps: []plugins.BuildStep{
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
