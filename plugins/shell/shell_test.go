package main

import (
	"net/rpc/jsonrpc"
	"os"
	"testing"

	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/natefinch/pie"
)

func TestShell_GetParams(t *testing.T) {
	path := "./shell"

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, path)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	var param string
	if err := client.Call("shell.GetName", "", &param); err != nil {
		t.Error(err)
		return
	}

	if err := client.Call("shell.GetDescription", "", &param); err != nil {
		t.Error(err)
		return
	}

	if err := client.Call("shell.GetVersion", "", &param); err != nil {
		t.Error(err)
		return
	}

	var bParam bool
	if err := client.Call("shell.HasSettings", "", &bParam); err != nil {
		t.Error(err)
		return
	}

	var itemParam []plugins.ItemParam
	if err := client.Call("shell.GetBuildParams", "", &itemParam); err != nil {
		t.Error(err)
		return
	}

	initPlugin := map[string]interface{}{
		"foo": 125,
	}

	if err := client.Call("shell.InitPlugin", initPlugin, &bParam); err != nil {
		t.Error(err)
		return
	}
}
