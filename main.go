package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bolsunovskyi/scheduler/jobs"
	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/bolsunovskyi/scheduler/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	port int
	conf config
)

type admin struct {
	Email    string
	Password string
}

type config struct {
	DBPath           string
	Admin            admin
	Plugins          []string
	DefaultBuildPath string
}

func init() {
	flag.IntVar(&port, "p", 8080, "http port")
	flag.Parse()

	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Println("Unable to read config file")
		log.Fatalln(err)
	}
}

func initTemplate() (*template.Template, error) {
	tpl := template.New("app").Delims("[[", "]]").Funcs(map[string]interface{}{
		"split": func(s string) []string {
			return strings.Split(s, "\n")
		},
		"time": func(t time.Time) string {
			return t.Format("2.01.2006 15:04:05")
		},
	})
	if _, err := tpl.ParseGlob("_templates/**/*"); err != nil {
		return nil, err
	}

	return tpl, nil
}

func main() {
	db, err := gorm.Open("sqlite3", conf.DBPath)
	if err != nil {
		log.Fatalln(err)
	}

	tpl, err := initTemplate()
	if err != nil {
		log.Fatalln(err)
	}

	router := gin.New()
	router.Use( /*gin.Logger(), */ gin.Recovery(), user.Middleware(db))
	router.Static("/assets", "./_assets")
	router.SetHTMLTemplate(tpl) //maybe move down

	//init user package
	if err := user.Init(router, db, conf.Admin.Email, conf.Admin.Password); err != nil {
		log.Fatalln(err)
	}

	//create auth router group
	auth := router.Group("/a")
	auth.Use(user.AbortUnAuth())

	//load plugins
	loadedPlugins := plugins.Load(db, tpl, auth.Group("/plugins"), conf.Plugins, conf.DBPath)

	//init jobs package
	jobs.Init(auth, db, loadedPlugins)

	//init plugins http transport
	plugins.InitHTTP(auth, db, loadedPlugins)

	//start router
	log.Printf("HTTP server started on port %d\n", port)
	router.Run(fmt.Sprintf(":%d", port))
}
