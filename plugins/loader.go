package plugins

import (
	"html/template"
	"log"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/natefinch/pie"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"os"
	"runtime"
)

const (
	TypeString = "string"
	TypeText   = "text"
	TypeSelect = "select"
)

type Item interface {
	GetName(_ string, rsp *string) error
	GetDescription(_ string, rsp *string) error
	GetVersion(_ string, rsp *string) error
	GetBuildParams(_ string, rsp *[]ItemParam) error
	HasSettings(_ string, rsp *bool) error
	InitPlugin(params map[string]interface{}, rsp Item) error
}

type ItemParam struct {
	Name        string         `json:"name"`
	Label       string         `json:"label"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Options     []ParamOptions `json:"options"`
	Value       string         `json:"value"`
}

type PluginParams struct {
	Name        string
	Description string
	Version     string
	HasSettings bool
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
}

type ParamOptions struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func makePluginHTTPHandler(pluginName string, group *gin.RouterGroup, dbPath string) {
	group.Group(fmt.Sprintf("/%s", pluginName), func(c *gin.Context) {
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

		c.HTML(rsp.StatusCode, rsp.Template, rsp.Data)
	})
}

func loadPlugin(pluginName string, baseTemplate *template.Template, group *gin.RouterGroup, dbPath string) error {
	log.Printf("Load plugin %s ...\n", pluginName)

	pluginExec := "./plugins/" + pluginName + "/" + pluginName
	if runtime.GOOS == "windows" {
		pluginExec = pluginExec + ".exe"
	}

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stdout, pluginExec)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Call(pluginName+".InitDB", dbPath, nil); err != nil {
		return err
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

	var params PluginParams
	if err := client.Call(pluginName+".GetPluginParams", "", &params); err != nil {
		return err
	}

	log.Printf("Plugin name: %s\n", params.Name)
	log.Printf("Plugin description: %s\n", params.Description)
	log.Printf("Plugin version: %s\n", params.Version)

	return nil
}

func Load(baseTemplate *template.Template, group *gin.RouterGroup, items []string, dbPath string) []string {
	var loadedItems []string

	log.Println("Loading plugins...")
	for _, item := range items {

		if err := loadPlugin(item, baseTemplate, group, dbPath); err != nil {
			log.Println(err)
			continue
		}

		loadedItems = append(loadedItems, item)
		log.Printf("Load plugin [%s] done.\n", item)
	}

	return loadedItems
}
