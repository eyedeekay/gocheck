package gocheck

import (
	"encoding/json"
	"io/ioutil"
)

func (c *Check) Import(str string) error {
	bytes, err := ioutil.ReadFile(str)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, c)
}

func (c *Check) ImportSites(str string) error {
	bytes, err := ioutil.ReadFile(str)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, c.Sites)
}

func (c *Check) ImportPeers(str string) error {
	bytes, err := ioutil.ReadFile(str)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, c.Peers)
}
