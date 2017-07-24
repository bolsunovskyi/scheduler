package plugins

import (
	"fmt"
	"log"
	"plugin"

	"os"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	TypeString = "string"
	TypeText   = "text"
	TypeSelect = "select"
)

type JobStep interface {
	GetName() string
	GetDescription() string
	GetVersion() string
	GetBuildParams() []StepParam
	HasSettings() bool
}

type StepParam struct {
	Name    string
	Type    string
	Options string
	Value   string
}

var loadedItems []JobStep

func GetPlugins() []JobStep {
	return loadedItems
}

func Load(db *gorm.DB, baseTemplate *template.Template, group *gin.RouterGroup, items []string) {
	params := map[string]interface{}{
		"db": db,
	}

	log.Println("Loading plugins...")
	for _, item := range items {
		log.Printf("Load plugin %s ...\n", item)
		p, err := plugin.Open(fmt.Sprintf("./plugins/%s/%s.so", item, item))
		if err != nil {
			log.Println(err)
			continue
		}

		s, err := p.Lookup("MakePlugin")
		if err != nil {
			log.Println(err)
			continue
		}

		switch v := s.(type) {
		default:
			log.Printf("unexpected type %T\n", v)
			continue
		case func(params map[string]interface{}) JobStep:
			pluginParams := params
			pluginParams["router"] = group.Group(fmt.Sprintf("/%s", item))

			//load plugin templates
			if _, err := os.Stat(fmt.Sprintf("./plugins/%s/_templates", item)); err == nil {
				if _, err := baseTemplate.ParseGlob(fmt.Sprintf("./plugins/%s/_templates/*", item)); err != nil {
					log.Printf("%s plugin, %s", item, err.Error())
				}
			}

			//load assets
			if _, err := os.Stat(fmt.Sprintf("./plugins/%s/_assets", item)); err == nil {
				group.Static(fmt.Sprintf("%s/assets/", item), fmt.Sprintf("./plugins/%s/_assets", item))
			}

			p := s.(func(params map[string]interface{}) JobStep)(params)
			log.Printf("Plugin name: %s\n", p.GetName())
			log.Printf("Plugin description: %s\n", p.GetDescription())
			log.Printf("Plugin version: %s\n", p.GetVersion())
			loadedItems = append(loadedItems, p)
			log.Printf("Load plugin [%s] done.\n", item)
		}
	}
}
