package plugins

import (
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

type ParamOptions struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func Load(db *gorm.DB, baseTemplate *template.Template, group *gin.RouterGroup, items []string) []Item {
	var loadedItems []Item

	//params := map[string]interface{}{
	//	"db": db,
	//}

	log.Println("Loading plugins...")
	for _, item := range items {
		log.Printf("Load plugin %s ...\n", item)

		//p, err := plugin.Open(fmt.Sprintf("./plugins/%s/%s.so", item, item))
		//if err != nil {
		//	log.Println(err)
		//	continue
		//}
		//
		//s, err := p.Lookup("MakePlugin")
		//if err != nil {
		//	log.Println(err)
		//	continue
		//}
		//
		//switch v := s.(type) {
		//default:
		//	log.Printf("unexpected type %T\n", v)
		//	continue
		//case func(params map[string]interface{}) Item:
		//	pluginParams := params
		//	pluginParams["router"] = group.Group(fmt.Sprintf("/%s", item))
		//
		//	//load plugin templates
		//	if _, err := os.Stat(fmt.Sprintf("./plugins/%s/_templates", item)); err == nil {
		//		if _, err := baseTemplate.ParseGlob(fmt.Sprintf("./plugins/%s/_templates/*", item)); err != nil {
		//			log.Printf("%s plugin, %s", item, err.Error())
		//		}
		//	}
		//
		//	//load assets
		//	if _, err := os.Stat(fmt.Sprintf("./plugins/%s/_assets", item)); err == nil {
		//		group.Static(fmt.Sprintf("%s/assets/", item), fmt.Sprintf("./plugins/%s/_assets", item))
		//	}
		//
		//	p := s.(func(params map[string]interface{}) Item)(params)
		//	log.Printf("Plugin name: %s\n", p.GetName())
		//	log.Printf("Plugin description: %s\n", p.GetDescription())
		//	log.Printf("Plugin version: %s\n", p.GetVersion())
		//	loadedItems = append(loadedItems, p)
		//	log.Printf("Load plugin [%s] done.\n", item)
		//}
	}

	return loadedItems
}
