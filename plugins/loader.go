package plugins

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"plugin"
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
}

type StepParam struct {
	Name    string
	Type    string
	Options string
	Value   string
}

var loadedItems []JobStep

func Load(db *gorm.DB, items []string) {
	params := map[string]interface{}{
		"db": db,
	}

	log.Println("Loading plugins...")
	for _, v := range items {
		log.Printf("Load plugin %s ...\n", v)
		p, err := plugin.Open(fmt.Sprintf("./plugins/%s/%s.so", v, v))
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
			p := s.(func(params map[string]interface{}) JobStep)(params)
			log.Printf("Plugin name: %s\n", p.GetName())
			log.Printf("Plugin description: %s\n", p.GetDescription())
			log.Printf("Plugin version: %s\n", p.GetVersion())
			loadedItems = append(loadedItems, p)
			log.Printf("Load plugin %s done.\n", v)
		}
	}
}
