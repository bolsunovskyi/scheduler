package main

import (
	"net/rpc/jsonrpc"
	"os"
	"testing"

	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/natefinch/pie"
)

func TestShell_GetParams(t *testing.T) {
	path := "./ssh"

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, path)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	var itemParam []plugins.ItemParam
	if err := client.Call("ssh.GetBuildParams", "", &itemParam); err != nil {
		t.Error(err)
		return
	}

	dbPath := "/tmp/test.sqlite"
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(path)

	if err := client.Call("ssh.InitDB", "", db); err != nil {
		t.Error(err)
		return
	}

	group := gin.New().Group("/test")
	if err := client.Call("ssh.InitRouter", "", group); err != nil {
		t.Error(err)
		return
	}

	var params plugins.PluginParams
	if err := client.Call("ssh.GetPluginParams", "", &params); err != nil {
		t.Error(err)
		return
	}
}
