package plugins

import (
	"fmt"
	"log"
	"plugin"
)

var loadedItems []*plugin.Plugin

func Load(items []string) {
	log.Println("Loading plugins...")
	for _, v := range items {
		log.Printf("Load plugin %s ...\n", v)
		p, err := plugin.Open(fmt.Sprintf("./plugins/%s/%s.so", v, v))
		if err != nil {
			log.Println(err)
			continue
		}

		s, err := p.Lookup("GetName")
		if err != nil {
			log.Println(err)
			continue
		}
		name := s.(func() string)()
		log.Printf("Plugin name: %s\n", name)

		s, err = p.Lookup("GetDescription")
		if err != nil {
			log.Println(err)
			continue
		}
		description := s.(func() string)()
		log.Printf("Plugin description: %s\n", description)

		s, err = p.Lookup("GetVersion")
		if err != nil {
			log.Println(err)
			continue
		}
		version := s.(func() string)()
		log.Printf("Plugin version: %s\n", version)

		loadedItems = append(loadedItems, p)
		log.Printf("Load plugin %s done.\n", v)
	}
}
