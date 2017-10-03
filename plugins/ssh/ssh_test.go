package main

import (
	"os"
	"testing"

	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/jinzhu/gorm"
	"github.com/natefinch/pie"
)

func TestShell_GetParams(t *testing.T) {
	path := "./ssh"

	dbPath := "/tmp/test.sqlite"
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(path)
	defer db.Close()

	client, err := pie.StartProvider(os.Stderr, path)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	var itemParam []plugins.ItemParam
	if err := client.Call("ssh.GetBuildParams", dbPath, &itemParam); err != nil {
		t.Error(err)
		return
	}

	if err := client.Call("ssh.InitDB", dbPath, nil); err != nil {
		t.Error(err)
		return
	}

	//group := gin.New().Group("/test")
	//if err := client.Call("ssh.InitRouter", "", group); err != nil {
	//	t.Error(err)
	//	return
	//}

	var params plugins.PluginParams
	if err := client.Call("ssh.GetPluginParams", "", &params); err != nil {
		t.Error(err)
		return
	}
}
