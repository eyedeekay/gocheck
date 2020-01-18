package gocheck

import (
	"encoding/json"
	"log"
)

func (s *Site) JsonString() string {
	r, e := json.MarshalIndent(s, "", "  ")
	if e != nil {
		log.Println("JSON ERROR", e.Error())
		return e.Error()
	}
	return string(r)
}

func (c *Check) ExportJsonArtifact() string {
	var Sites []Site
	for _, site := range c.Sites {
		Sites = append(Sites, site)
	}
	r, e := json.MarshalIndent(Sites, "", "  ")
	if e != nil {
		log.Println("JSON ERROR", e.Error())
		return e.Error()
	}
	return string(r)
}

func (c *Check) ExportMiniJsonArtifact() string {
	var Peers []Site
	for _, site := range c.Peers {
		Peers = append(Peers, site)
	}
	r, e := json.MarshalIndent(Peers, "", "  ")
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
