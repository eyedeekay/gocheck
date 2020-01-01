package gocheck

import (
	"github.com/eyedeekay/goSam"
	"github.com/eyedeekay/i2pasta/convert"
	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// TODO calculate approx 9's
// TODO save check history(more json?)

type Site struct {
	url      string
	dest     []string
	base32   []string
	uptime   int
	downtime int

	title string
	//favicon []byte
	desc           string
	successHistory []bool
}

func (s *Site) Nines() string {
	trues := 0  //uptime
	falses := 0 //downtime
	for _, v := range s.successHistory {
		if v {
			trues += 1
		} else {
			falses += 1
		}
	}
	s.uptime = (trues + falses) / trues
	s.downtime = (trues + falses) / falses
	r := "<div class=\"" + s.url + "\"" + ">"
	r += "<p>" + strconv.Itoa(s.uptime)
	r += "</p>"
	r += "<p>" + strconv.Itoa(s.downtime)
	r += "</p>"
	r += "</div>"
	return r
}

func (s *Site) HTML() string {
	r := "<div id=\"" + s.url + "\">\n"
	r += "<h3>"
	if s.title == "" {
		s.title = s.url
	}
	r += s.title
	r += "</h3>\n"
	r += "<div class=\"" + s.url + " url\">  URL: " + s.url
	r += "</div>\n"
	r += "<div class=\"" + s.url + " base32\">  Base32: " + s.base32[len(s.base32)-1] + ".b32.i2p"
	r += "</div>\n"
	r += "<div class=\"" + s.url + " desc\">  Description: " + s.desc
	r += "</div>\n"
	r += "<div class=\"" + s.url + " dest\">  Destination: " + s.dest[len(s.dest)-1]
	if len(s.dest) > 1 && s.dest[len(s.dest)-1] != s.dest[len(s.dest)-2] {
		r += "<div class=\"updated\">This URL has been updated since we last checked.</div>"
	}
	r += "</div>\n"
	r += "</div>\n"
	return r
}

func (s *Site) JsonString() string {
	var r string
	changed := "no"
	if len(s.dest) > 1 && s.dest[len(s.dest)-1] != s.dest[len(s.dest)-2] {
		changed = "yes"
	}
	r = "{" +
		"url: \"" + s.url + "\"" +
		"dest: \"" + s.dest[len(s.dest)-1] + "\"" +
		"b32: \"" + s.base32[len(s.base32)-1] + "\"" +
		"title: \"" + s.title + "\"" +
		"desc: \"" + s.desc + "\"" +
		"changed: \"" + changed + "\n" +
		"url: \"" + s.url + "\"" +
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

func (c *Check) LoadHostsFile(hostsfile string) ([]Site, error) {
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
				b32, err := i2pconv.I2p64to32(splitpair[1])
				if err == nil {
					sites = append(
						sites,
						Site{
							url:    u,
							dest:   append([]string{}, splitpair[1]),
							base32: append([]string{}, b32),
						},
					)
				} else {
					fmt.Printf("%s", err)
				}
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

// TODO: Buffer should be done differently if I want to do requests like this,
// or maybe sites should share a tunnel, or maybe I just use the HTTP proxy
// instead.


//AsyncGet is a misnomer, it has to be done in order for now.
func (c *Check) AsyncGet(index int, site Site) {
	var err error
	_, err = c.Client.Get(site.url)
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
		c.AsyncGet(index, site)
		if index != 0 && (index%5) == 0 {
			time.Sleep(time.Minute)
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
