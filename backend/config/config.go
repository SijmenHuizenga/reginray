package config

import (
	"log"
	"gopkg.in/ini.v1"
	"../model"
)

func LoadConfig() ([]model.ServicePattern, []model.GrokPatterns) {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Failed to load config")
	}

	var servicePatterns []model.ServicePattern
	var grokPatterns []model.GrokPatterns

	for _, section := range cfg.Sections() {
		if section.Name() == "DEFAULT" {
			for _, pair := range section.Keys() {
				grokPatterns = append(grokPatterns, model.GrokPatterns{Key: pair.Name(), Val: pair.Value()})
				log.Println("Added pattern " + pair.Name())
			}
			continue
		}
		servicePattern := model.ServicePattern{
			Title:      section.Name(),
			Containers: section.Key("containers").MustString("*"),
			Images:     section.Key("images").MustString("*"),
			Pattern:    section.Key("pattern").Value(),
		}
		servicePatterns = append(servicePatterns, servicePattern)
	}
	return servicePatterns, grokPatterns
}
