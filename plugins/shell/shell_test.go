package main

import (
	"net/http"
	"net/http/httptest"
	"net/rpc/jsonrpc"
	"os"
	"strings"
	"testing"

	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/natefinch/pie"
)

func TestShell_HandleHTTP(t *testing.T) {
	path := "./shell"

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, path)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	r := gin.New()

	r.POST("/", func(c *gin.Context) {
		var rsp plugins.HTTPResponse
		rq := plugins.MakeRequest(c.Request, "")

		if err := client.Call("shell.HandleHTTP", rq, &rsp); err != nil {
			c.String(400, err.Error())
			return
		}

		if rsp.Template != "{}" {
			c.String(400, rsp.Template)
			return
		}

		method, ok := rsp.Data["method"]
		if !ok {
			c.String(400, "Method not found")
			return
		}

		if method.(string) != "POST" {
			c.String(400, "Method not found")
		}

		c.Status(200)
	})

	rq, err := http.NewRequest("POST", "/", strings.NewReader(`{}`))
	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, rq)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
		t.Logf("Body: %s\n", rr.Body.String())
		t.Logf("Path: %s\n", path)
		return
	}
}

func TestShell_GetParams(t *testing.T) {
	path := "./shell"

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, path)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	var params plugins.Params
	if err := client.Call("shell.GetPluginParams", "", &params); err != nil {
		t.Error(err)
		return
	}
	if params.Name != "shell" {
		t.Error("Wrong plugin param name")
		return
	}

	if params.BuildSteps[0].Name != "command" {
		t.Error("Wrong item param")
		return
	}

	if err := client.Call("shell.InitDB", "", nil); err != nil {
		t.Error(err)
		return
	}
}
