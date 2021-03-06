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
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

var (
	port int
	conf config
)

type database struct {
	Type     string
	Path     string
	Name     string
	Host     string
	Port     int
	User     string
	Password string
}

type admin struct {
	Email    string
	Password string
}

type config struct {
	DB               database
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

	//m, err := migrate.New(
	//	"file://_migrations",
	//	fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
	//		conf.DB.User, conf.DB.Password, conf.DB.Host, conf.DB.Port, conf.DB.Name))
	//
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	//	log.Fatalln(err)
	//}
}

func main() {
	//TODO: refactor this
	var db *gorm.DB
	var err error
	if conf.DB.Type == "postgres" {
		db, err = gorm.Open("postgres",
			fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
				conf.DB.User, conf.DB.Password, conf.DB.Host, conf.DB.Port, conf.DB.Name))
		if err != nil {
			return
		}
	} else if conf.DB.Type == "sqlite3" {
		db, err = gorm.Open("sqlite3", conf.DB.Path)
	}

	router := gin.New()
	router.Use( /*gin.Logger(), */ gin.Recovery(), user.Middleware(db))

	router.Static("/assets", "./_assets")
	tpl := template.New("app").Delims("[[", "]]").Funcs(map[string]interface{}{
		"split": func(s string) []string {
			return strings.Split(s, "\n")
		},
		"time": func(t time.Time) string {
			return t.Format("2.01.2006 15:04:05")
		},
	})
	if _, err := tpl.ParseGlob("_templates/**/*"); err != nil {
		log.Fatalln(err)
	}

	if err := user.Init(router, db, conf.Admin.Email, conf.Admin.Password); err != nil {
		log.Fatalln(err)
	}

	auth := router.Group("/a")
	auth.Use(user.AbortUnAuth())

	loadedPlugins := plugins.Load(db, tpl, auth.Group("/plugins"), conf.Plugins)

	jobs.Init(auth, db, loadedPlugins)
	plugins.InitHTTP(auth, db, loadedPlugins)

	router.SetHTMLTemplate(tpl)

	log.Printf("HTTP server started on port %d\n", port)
	router.Run(fmt.Sprintf(":%d", port))
}
