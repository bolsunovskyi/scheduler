package plugins

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/pie"
)

const (
	TypeString = "string"
	TypeText   = "text"
	TypeSelect = "select"
)

type BuildStep struct {
	Name        string             `json:"name"`
	Label       string             `json:"label"`
	Type        string             `json:"type"`
	Description string             `json:"description"`
	Options     []BuildStepOptions `json:"options"`
	Value       string             `json:"value"`
}

type BuildStepOptions struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Params struct {
	Name        string
	Description string
	Version     string
	HasSettings bool
	BuildSteps  []BuildStep
}

type HTTPRequest struct {
	DBPath           string
	Method           string
	URL              *url.URL
	Proto            string // "HTTP/1.0"
	ProtoMajor       int    // 1
	ProtoMinor       int    // 0
	Header           http.Header
	ContentLength    int64
	TransferEncoding []string
	Host             string
	Form             url.Values
	PostForm         url.Values
	MultipartForm    *multipart.Form
	RemoteAddr       string
	RequestURI       string
	Body             []byte
	BodyStr          string
}

func MakeRequest(rq *http.Request, dbPath string) HTTPRequest {
	rq.ParseForm()
	rq.ParseMultipartForm(32 << 20) //

	bts, _ := ioutil.ReadAll(rq.Body)

	return HTTPRequest{
		Method:           rq.Method,
		PostForm:         rq.PostForm,
		DBPath:           dbPath,
		Form:             rq.Form,
		Header:           rq.Header,
		RemoteAddr:       rq.RemoteAddr,
		Host:             rq.Host,
		URL:              rq.URL,
		ContentLength:    rq.ContentLength,
		MultipartForm:    rq.MultipartForm,
		Proto:            rq.Proto,
		ProtoMajor:       rq.ProtoMajor,
		ProtoMinor:       rq.ProtoMinor,
		RequestURI:       rq.RequestURI,
		TransferEncoding: rq.TransferEncoding,
		Body:             bts,
		BodyStr:          string(bts),
	}
}

type HTTPResponse struct {
	StatusCode int
	Data       map[string]interface{}
	Template   string
	Json       bool
	Raw        bool
	RawBody    string
}

func makePluginHTTPHandler(pluginName string, group *gin.RouterGroup, dbPath string) {
	pluginGroup := group.Group(fmt.Sprintf("/%s", pluginName))
	pluginGroup.Any("", func(c *gin.Context) {
		//Call plugin http handler
		pluginExec := "./plugins/" + pluginName + "/" + pluginName
		if runtime.GOOS == "windows" {
			pluginExec = pluginExec + ".exe"
		}

		client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, c.Writer, pluginExec)
		if err != nil {
			c.String(400, err.Error())
			return
		}
		defer client.Close()

		var rsp HTTPResponse
		if err := client.Call(pluginName+".HandleHTTP", MakeRequest(c.Request, dbPath), &rsp); err != nil {
			c.String(400, err.Error())
			return
		}

		if rsp.Json {
			c.JSON(rsp.StatusCode, rsp.Data)
			return
		}

		if rsp.Template != "" {
			c.HTML(rsp.StatusCode, rsp.Template, rsp.Data)
			return
		}

		c.String(rsp.StatusCode, rsp.RawBody)
	})
}

func loadPlugin(pluginName string, baseTemplate *template.Template, group *gin.RouterGroup, dbPath string) (*Params, error) {
	log.Printf("Load plugin %s ...\n", pluginName)

	pluginExec := "./plugins/" + pluginName + "/" + pluginName
	if runtime.GOOS == "windows" {
		pluginExec = pluginExec + ".exe"
	}

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stdout, pluginExec)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	if err := client.Call(pluginName+".InitDB", dbPath, nil); err != nil {
		return nil, err
	}

	//make http handler
	makePluginHTTPHandler(pluginName, group, dbPath)

	//load templates
	if _, err := os.Stat(fmt.Sprintf("./plugins/%s/_templates", pluginName)); err == nil {
		if _, err := baseTemplate.ParseGlob(fmt.Sprintf("./plugins/%s/_templates/*", pluginName)); err != nil {
			log.Printf("%s plugin, %s", pluginName, err.Error())
		}
	}

	//load assets
	if _, err := os.Stat(fmt.Sprintf("./plugins/%s/_assets", pluginName)); err == nil {
		group.Static(fmt.Sprintf("%s/assets/", pluginName), fmt.Sprintf("./plugins/%s/_assets", pluginName))
	}

	var params Params
	if err := client.Call(pluginName+".GetPluginParams", dbPath, &params); err != nil {
		return nil, err
	}

	log.Printf("Plugin name: %s\n", params.Name)
	log.Printf("Plugin description: %s\n", params.Description)
	log.Printf("Plugin version: %s\n", params.Version)

	return &params, nil
}

func Load(baseTemplate *template.Template, group *gin.RouterGroup, items []string, dbPath string) []Params {
	var loadedItems []Params

	log.Println("Loading plugins...")
	for _, item := range items {

		params, err := loadPlugin(item, baseTemplate, group, dbPath)
		if err != nil {
			log.Println(err)
			continue
		}

		loadedItems = append(loadedItems, *params)
		log.Printf("Load plugin [%s] done.\n", item)
	}

	return loadedItems
}
