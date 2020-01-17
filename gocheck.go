package gocheck

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/eyedeekay/httptunnel/multiproxy"
	"github.com/eyedeekay/i2pasta/convert"
	"github.com/eyedeekay/sam-forwarder/config"
	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"
	"github.com/eyedeekay/sam3/i2pkeys"

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

func (s *Site) HostPair() string {
	return s.url + "=" + s.dest[len(s.dest)-1]
}

func (s *Site) Up() string {
	//log.Println("Length of history", len(s.successHistory))
	if len(s.successHistory) < 1 {
		return "unknown"
	}

	if s.successHistory[len(s.successHistory)-1] {
		return "true"
	}
	return "false"
}

func (s *Site) JSONUp() string {
	if len(s.successHistory) < 1 {
		return "[]"
	}
	a := "[\n"
	for _, s := range s.successHistory {
		if s {
			a += "    true,\n"
		} else {
			a += "    false,\n"
		}
	}
	a += "]"
	return strings.TrimSuffix(a, ",\n]") + "\n  ]"
}

func (s *Site) HTML() string {
	r := "<span id=\"" + s.url + "\" class=\"site " + s.Up() + "\">\n"
	r += "<h3>"
	if s.title == "" {
		s.title = s.url
	}
	r += s.title
	r += "</h3>\n"
	r += "<div><span class=\"" + s.url + " label url\">  URL: </span> <span class=\"field\"><a href=\"" + s.url + "\">" + s.url + "</a></span></div>\n"
	r += "<div><span class=\"" + s.url + " label base32\">  Base32: </span> <span class=\"field\"><a href=\"http://" + s.base32[len(s.base32)-1] + ".b32.i2p\">" + s.base32[len(s.base32)-1] + "</a>" + "</span></div>\n"
	r += "<div><span class=\"" + s.url + " label desc\">  Description: </span> <span class=\"field\">" + s.desc + "</span></div>\n"
	r += "<div><span class=\"" + s.url + " label stat\">  Alive: </span> <span class=\"field\">" + s.Up() + "</span></div>\n"
	r += "<div><span class=\"" + s.url + " label dest\">  Destination: </span> <span class=\"field\">" + s.dest[len(s.dest)-1] + "</span></div>\n"
	if len(s.dest) > 1 && s.dest[len(s.dest)-1] != s.dest[len(s.dest)-2] {
		r += "<div><span class=\" label updated\">Changed?:</span> <span class=\"field\">This URL has been updated since we last checked.</span></div>"
	}
	r += "</span>\n"
	r += "</div>\n"
	return r
}

func (s *Site) JsonString() string {
	var r string
	changed := "no"
	if len(s.dest) > 1 && s.dest[len(s.dest)-1] != s.dest[len(s.dest)-2] {
		changed = "yes"
	}
	r = "{\n" +
		"  \"url\": \"" + s.url + "\",\n" +
		"  \"dest\": \"" + s.dest[len(s.dest)-1] + "\",\n" +
		"  \"b32\": \"" + s.base32[len(s.base32)-1] + "\",\n" +
		"  \"title\": \"" + s.title + "\",\n" +
		"  \"desc\": \"" + s.desc + "\",\n" +
		"  \"up\": " + s.JSONUp() + ",\n" +
		"  \"changed\": \"" + changed + "\",\n" +
		"  \"url: \"" + s.url + "\"\n" +
		"}"
	return r
}

type Check struct {
	*samforwarder.SAMForwarder
	SAMHTTPProxy *i2pbrowserproxy.SAMMultiProxy
	*http.Transport
	*http.Client
	RegularProxy string

	sites     []Site
	hostsfile string
	up        bool
}

func (c *Check) Base32() string {
	return c.SAMForwarder.Base32()
}

func (c *Check) Base32Readable() string {
	return c.SAMForwarder.Base32Readable()
}

func (c *Check) Base64() string {
	return c.SAMForwarder.Base64()
}

func (c *Check) Cleanup() {
	c.SAMForwarder.Cleanup()
}

func (c *Check) Close() error {
	return c.SAMForwarder.Close()
}

func (c *Check) Config() *i2ptunconf.Conf {
	return c.SAMForwarder.Config()
}

func (c *Check) GetType() string {
	return "uptimer"
}

func (c *Check) ID() string {
	return c.SAMForwarder.ID()
}

func (c *Check) Print() string {
	return c.SAMForwarder.Print()
}

func (c *Check) Props() map[string]string {
	return c.SAMForwarder.Props()
}

func (c *Check) Keys() i2pkeys.I2PKeys {
	return c.SAMForwarder.Keys()
}

func (c *Check) Search(s string) string {
	return c.SAMForwarder.Search(s)
}

func (c *Check) Target() string {
	return c.SAMForwarder.Target()
}

func (c *Check) Up() bool {
	return c.SAMForwarder.Up()
}

func (c *Check) ExportJsonArtifact() string {
	export := "{"
	for _, site := range c.sites {
		export += site.JsonString() + ",\n"
	}
	export += "}"
	return strings.TrimSuffix(export, "},\n}") + "}\n}\n"
}

func (c *Check) ExportMiniJsonArtifact() string {
	export := "{\n"
	for _, site := range c.sites {
		if len(site.successHistory) > 0 {
			export += site.JsonString() + ",\n"
		}
	}
	export += "}"
	return strings.TrimSuffix(export, "},\n}") + "}\n}\n"
}

func (c *Check) ExportHostsFile() string {
	export := ""
	for _, site := range c.sites {
		export += site.HostPair() + "\n"
	}
	return export
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
		u = strings.Replace("http://"+u, "///", "//", -1)
	}
	if _, err := url.Parse(u); err != nil {
		log.Println("ERR", err)
		return "", err
	}
	return u, nil
}

// TODO: Buffer should be done differently if I want to do requests like this,
// or maybe sites should share a tunnel, or maybe I just use the HTTP proxy
// instead.

//AsyncGet is a misnomer, it has to be done in order for now.
func (c *Check) AsyncGet(index int, site *Site) {
	res, err := c.Client.Get(site.url)
	log.Println("CHECKING UPNESS")
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		site.successHistory = append(site.successHistory, true)
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		title := strings.TrimSpace(doc.Find("title").Text())
		if title != "" {
			log.Println(title)
			site.title = title
		}
		desc := doc.Find("meta[name=description]")
		content, ok := desc.Attr("content")
		if ok {
			site.desc = strings.TrimSpace(content)
			log.Println(strings.TrimSpace(content))
		}
		fmt.Printf("the eepSite is up: index=%d, title=\"%s\", desc=\"%s\", history=%d\n", index, site.title, site.desc, len(site.successHistory))
	} else {
		fmt.Printf("the eepSite appears to be down: %v, %v\n", index, res.StatusCode)
		site.successHistory = append(site.successHistory, false)
	}
}

func (c *Check) CheckAll() {
	for index, site := range c.sites {
		log.Printf("Checking URL: %s", site.url)
		go c.AsyncGet(index, &site)
		time.Sleep(time.Second * 20)
	}
}

func (c *Check) QuerySite(s string) string {
	if u, err := Validate(s); err != nil {
		return "Not a valid URL for checking"
	} else {
		for index, site := range c.sites {
			u2, err := Validate(site.url)
			if err != nil {
				return "Invalid URL found in DB"
			}
			if u2 == u {
				c.AsyncGet(index, &c.sites[index])
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
	g, e := s.SAMHTTPProxy.Load()
	if e != nil {
		return nil, e
	}
	s.SAMHTTPProxy = g.(*i2pbrowserproxy.SAMMultiProxy)
	s.up = true
	fmt.Printf("Finished putting tunnel up")
	return s, nil
}
