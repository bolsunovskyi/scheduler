package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bolsunovskyi/scheduler/site"
	"github.com/gin-gonic/gin"
)

var port int

func init() {
	flag.IntVar(&port, "p", 8080, "http port")
	flag.Parse()
}

func main() {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.Static("/assets", "./_assets")
	router.LoadHTMLGlob("_templates/**/*")

	site.InitHTTP(router)

	log.Printf("HTTP server started on port %d\n", port)
	router.Run(fmt.Sprintf(":%d", port))
}
