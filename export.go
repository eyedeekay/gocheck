package gocheck

import (
	"encoding/json"
	"log"
)

func (s *Site) JsonString() string {
	r, e := json.Marshal(s)
	if e != nil {
		log.Println("JSON ERROR", e.Error())
		return e.Error()
	}
	return string(r)
}

func (c *Check) ExportJsonArtifact() string {
	r, e := json.Marshal(c)
	if e != nil {
		log.Println("JSON ERROR", e.Error())
		return e.Error()
	}
	return string(r)
}

func (c *Check) ExportMiniJsonArtifact() string {
	var temp []Site
	for _, site := range c.Sites {
		if len(site.SuccessHistory) > 0 {
			temp = append(temp, site)
		}
	}
	r, e := json.Marshal(temp)
	if e != nil {
		log.Println("JSON ERROR", e.Error())
		return e.Error()
	}
	return string(r)
}

func (c *Check) ExportHostsFile() string {
	export := ""
	for _, site := range c.Sites {
		export += site.HostPair() + "\n"
	}
	return export
}
