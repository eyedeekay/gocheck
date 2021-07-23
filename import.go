package gocheck

import (
	"encoding/json"
	"fmt"
	"github.com/eyedeekay/i2pasta/convert"
	"io/ioutil"
	"strings"
	"time"
)

func (c *Check) ImportSites(str string) error {
	if str == "" {
		return nil
	}
	bytes, err := ioutil.ReadFile(str)
	if err != nil {
		return fmt.Errorf("Error loading hosts file %s %s", err, str)
	}
	return json.Unmarshal(bytes, &c.Sites)
}

func (c *Check) ImportPeers(str string) error {
	if str == "" {
		return nil
	}
	bytes, err := ioutil.ReadFile(str)
	if err != nil {
		return fmt.Errorf("Error loading peers file %s %s", err, str)
	}
	return json.Unmarshal(bytes, &c.Peers)
}

func (c *Check) LoadHostsFile(hostsfile string) ([]Site, error) {
	fmt.Printf("LoadHostsFile")
	if hostsfile == "" {
		return nil, nil //fmt.Errorf("Error hosts file not given %s", hostsfile)
	}
	hostbytes, err := ioutil.ReadFile(hostsfile)
	if err != nil {
		return nil, fmt.Errorf("Error loading hosts file %s %s", err, hostsfile)
	}
	hoststring := string(hostbytes)
	sites, err := c.LoadHostsLines(hoststring)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func (c *Check) LoadHostsLines(hoststring string) ([]Site, error) {
	combinedpairs := strings.Split(hoststring, "\n")
	var sites []Site

	for index, pair := range combinedpairs {
		if !strings.HasPrefix(strings.Replace(strings.Replace(pair, " ", "", -1), "\t", "", -1), "#") {
			splitpair := strings.SplitN(pair, "=", 2)
			if splitpair[0] != "" {
				site, err := c.LoadHostsLine(splitpair)
				if err != nil {
					return nil, err
				}
				sites = append(sites, site)
				fmt.Printf("LoadHostsFile: (%v)loaded %s\n", index, splitpair[0])
			}
		}
	}
	return sites, nil
}

func (c *Check) LoadHostsLine(splitpair []string) (Site, error) {
	if len(splitpair) == 2 {
		if u, err := Validate(splitpair[0]); err != nil {
			return Site{}, err
		} else {
			b32, err := i2pconv.I2p64to32(splitpair[1])
			if err == nil {
				return Site{
					SuccessHistory: make(map[time.Time]bool),
					Url:            u,
					Dest:           append([]string{}, splitpair[1]),
					Base32:         append([]string{}, b32),
				}, nil
			} else {
				fmt.Printf("%s", err)
				return Site{
					SuccessHistory: make(map[time.Time]bool),
					Url:            u,
					Dest:           append([]string{}, splitpair[1]),
					Base32:         append([]string{}, "fail.b32.i2p"),
				}, nil
			}
		}
	} else if len(splitpair) == 1 {
		if u, err := Validate(splitpair[0]); err != nil {
			return Site{}, err
		} else {
			if err == nil {
				return Site{
					SuccessHistory: make(map[time.Time]bool),
					Url:            u,
				}, nil
			} else {
				fmt.Printf("%s", err)
			}
		}
	}
	return Site{}, nil
}
