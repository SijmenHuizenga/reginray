package decoder

import (
	"github.com/tidwall/match"
	"../model"
	"github.com/vjeantet/grok"
	"log"
)

type Decoder struct {
	grokmatcher     *grok.Grok
	servicePatterns []model.ServicePattern
}

func NewDecoder(sp []model.ServicePattern, grokPattens []model.GrokPatterns) Decoder {
	g, err := grok.NewWithConfig(&grok.Config{
		NamedCapturesOnly:   true,
		RemoveEmptyValues:   false,
		SkipDefaultPatterns: false,
	})

	if err != nil {
		panic("Could not init grok matcher")
	}

	for _, pattern := range grokPattens {
		g.AddPattern(pattern.Key, pattern.Val)
	}

	return Decoder{
		servicePatterns: sp,
		grokmatcher:     g,
	}
}

func (d Decoder) AddFields(logentry *model.LogEntry) {
	pattern := d.findPattern(*logentry)

	if pattern == nil {
		return
	}

	values, err := d.grokmatcher.Parse(pattern.Pattern, logentry.Message)

	if err != nil {
		log.Println("Could not parse log line with type servicepattern %v " + logentry.Message)
		log.Println(err)
		logentry.Fields = map[string]string{}
	}

	logentry.Fields = values

}

func (d Decoder) findPattern(logentry model.LogEntry) *model.ServicePattern {
	for _, servicePattern := range d.servicePatterns {
		if match.Match(logentry.Image.Name, servicePattern.Images) &&
			match.Match(logentry.Container.Name, servicePattern.Containers) {
			return &servicePattern
		}
	}
	return nil
}
