package gocheck

import (
	"github.com/eyedeekay/goSam"
	"github.com/eyedeekay/sam-forwarder/tcp"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// TODO calculate approx 9's

type Site struct {
	url  string
	dest string

	title string
	//favicon []byte
	desc           string
	successHistory []bool
}

func (s *Site) JsonString() string {
	var r string
	r = "{" +
		"url:" + s.url +
		"dest:" + s.dest +
		"title:" + s.title +
		"desc:" + s.url +
		"url:" + s.url +
		"}"
	return r
}

type Check struct {
	*samforwarder.SAMForwarder
	*http.Transport
	*http.Client
	i2p   *goSam.Client
	sites []Site
}

func LoadHostsFile(hostsfile string) ([]Site, error) {
	hostbytes, err := ioutil.ReadFile(hostsfile)
	if err != nil {
		return nil, err
	}
	hoststring := string(hostbytes)
	combinedpairs := strings.Split(hoststring, "\n")
	var sites []Site
	for index, pair := range combinedpairs {
		splitpair := strings.SplitN(pair, "=", 2)
		if len(splitpair) == 2 {
			if u, err := Validate(splitpair[0]); err != nil {
				return nil, err
			} else {
				sites = append(
					sites,
					Site{
						url:  u,
						dest: splitpair[1],
					},
				)
				fmt.Printf("LoadHostsFile: (%v)loaded %s\n", index, splitpair[0])
			}
		}
	}

	return sites, nil
}

func NewSAMChecker(hostsfile string) (*Check, error) {
	var c Check
	var err error
	c.i2p, err = goSam.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	c.Transport = &http.Transport{
		Dial: c.i2p.Dial,
	}
	c.Client = &http.Client{
		Transport: c.Transport,
	}
	c.sites, err = LoadHostsFile(hostsfile)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func Validate(u string) (string, error) {
	if !strings.Contains(u, ".i2p") {
		return "", fmt.Errorf("Not an I2P domain")
	}
	if !strings.HasPrefix(u, "http") {
		u = "http://" + u
	}
	if _, err := url.Parse(u); err != nil {
		return "", err
	}
	return u, nil
}

func (c *Check) CheckAll() {
	for index, site := range c.sites {
		fmt.Printf("Checking URL:")
		_, err := c.Client.Get(site.url)
		if err != nil {
			fmt.Printf("the eepSite appears to be down: %v %s\n", index, err)
			site.successHistory = append(site.successHistory, false)
		} else {
			fmt.Printf("the eepSite is up: %v %s\n", index, err)
			site.successHistory = append(site.successHistory, true)
		}
	}
}

func (c *Check) QuerySite(site string) string {
	if u, err := Validate(site); err != nil {
		return "Not a valid URL for checking"
	} else {
		for index, site := range c.sites {
			if site.url == u {
				fmt.Printf("The site was found at %v", index)
				return site.JsonString()
			}
		}
	}
	return "The site was not found"
}
