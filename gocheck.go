package gocheck

import (
	"github.com/eyedeekay/goSam"
	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"

	"fmt"
	"io/ioutil"
	"log"
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

func (s *Site) HTML() string {
	r := "<div id=\"" + s.url + "\">\n"
	r += "<h3>"
	if s.title == "" {
		s.title = s.url
	}
	r += s.title
	r += "</h3>\n"
	r += "<p>  URL: " + s.url
	r += "</p>\n"
	r += "<p>  Description: " + s.desc
	r += "</p>\n"
	r += "<p>  Destination: " + s.dest
	r += "</p>\n"
	r += "</div>\n"
	return r
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
	i2p       *goSam.Client
	sites     []Site
	hostsfile string
	up        bool
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
	return NewSAMCheckerFromOptions(SetHostsFile(hostsfile))
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

func (c *Check) AsyncGet(index int, site Site) {
	_, err := c.Client.Get(site.url)
	if err != nil {
		fmt.Printf("the eepSite appears to be down: %v %s\n", index, err)
		site.successHistory = append(site.successHistory, false)
	} else {
		fmt.Printf("the eepSite is up: %v %s\n", index, err)
		site.successHistory = append(site.successHistory, true)
	}
}

func (c *Check) CheckAll() {
	for index, site := range c.sites {
		log.Printf("Checking URL: %s", site.url)
		go c.AsyncGet(index, site)
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

func (s *Check) Load() (samtunnel.SAMTunnel, error) {
	if !s.up {
		fmt.Printf("Started putting tunnel up")
	}
	f, e := s.SAMForwarder.Load()
	if e != nil {
		return nil, e
	}
	s.SAMForwarder = f.(*samforwarder.SAMForwarder)
	//s.mark = markdown.New(markdown.XHTMLOutput(true))
	s.up = true
	fmt.Printf("Finished putting tunnel up")
	return s, nil
}
