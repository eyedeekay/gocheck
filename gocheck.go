package gocheck

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/eyedeekay/httptunnel/multiproxy"
	"github.com/eyedeekay/i2pasta/convert"
	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TODO calculate approx 9's
// TODO save check history(more json?)

type Site struct {
	Url      string   `json:url,omitempty`
	Dest     []string `json:dest,omitempty`
	Base32   []string `json:base32,omitempty`
	Uptime   int      `json:uptime,omitempty`
	Downtime int      `json:downtime,omitempty`

	Title          string             `json:title,omitempty`
	Desc           string             `json:desc,omitempty`
	SuccessHistory map[time.Time]bool `json:successHistory,omitempty`
}

func (s *Site) Nines() string {
	trues := 0  //Uptime
	falses := 0 //Downtime
	for _, v := range s.SuccessHistory {
		if v {
			trues += 1
		} else {
			falses += 1
		}
	}
	s.Uptime = (trues + falses) / trues
	s.Downtime = (trues + falses) / falses
	r := "<div class=\"" + s.Url + "\"" + ">"
	r += "<p>" + strconv.Itoa(s.Uptime)
	r += "</p>"
	r += "<p>" + strconv.Itoa(s.Downtime)
	r += "</p>"
	r += "</div>"
	return r
}

func (s *Site) HostPair() string {
	return s.Url + "=" + s.Dest[len(s.Dest)-1]
}

func (s *Site) sort() []time.Time {
	var keys []time.Time
	for k := range s.SuccessHistory {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})
	return keys
}

func (s *Site) SortedSuccesses() []bool {
	var vals []bool
	for _, k := range s.sort() {
		vals = append(vals, s.SuccessHistory[k])
	}
	return vals
}

func (s *Site) Up() string {
	//log.Println("Length of history", len(s.SuccessHistory))
	if len(s.SuccessHistory) < 1 {
		return "unknown"
	}

	if s.SortedSuccesses()[len(s.SortedSuccesses())-1] {
		return "true"
	}
	return "false"
}

func (s *Site) HTML() string {
	r := "<span id=\"" + s.Url + "\" class=\"site " + s.Up() + "\">\n"
	r += "<h3>"
	if s.Title == "" {
		s.Title = s.Url
	}
	r += s.Title
	r += "</h3>\n"
	r += "<div><span class=\"" + s.Url + " label Url\">  URL: </span> <span class=\"field\"><a href=\"" + s.Url + "\">" + s.Url + "</a></span></div>\n"
	r += "<div><span class=\"" + s.Url + " label Base32\">  Base32: </span> <span class=\"field\"><a href=\"http://" + s.Base32[len(s.Base32)-1] + ".b32.i2p\">" + s.Base32[len(s.Base32)-1] + "</a>" + "</span></div>\n"
	r += "<div><span class=\"" + s.Url + " label Desc\">  Description: </span> <span class=\"field\">" + s.Desc + "</span></div>\n"
	r += "<div><span class=\"" + s.Url + " label stat\">  Alive: </span> <span class=\"field\">" + s.Up() + "</span></div>\n"
	r += "<div><span class=\"" + s.Url + " label Dest\">  Destination: </span> <span class=\"field\">" + s.Dest[len(s.Dest)-1] + "</span></div>\n"
	if len(s.Dest) > 1 && s.Dest[len(s.Dest)-1] != s.Dest[len(s.Dest)-2] {
		r += "<div><span class=\" label updated\">Changed?:</span> <span class=\"field\">This URL has been updated since we last checked.</span></div>"
	}
	r += "</span>\n"
	r += "</div>\n"
	return r
}

// MarshalJSON returns *m as the JSON encoding of m.
/*func (s *Site) MarshalJSON() ([]byte, error) {
	return []byte(s.JsonString()), nil
}*/

// UnmarshalJSON sets *m to a copy of data.
/*func (s *Site) UnmarshalJSON(data []byte) error {
	if s == nil {
		return fmt.Errorf("RawString: UnmarshalJSON on nil pointer")
	}
	*s += Site(data)
	return nil
}*/

type Check struct {
	*samforwarder.SAMForwarder `json:"-"`
	SAMHTTPProxy               *i2pbrowserproxy.SAMMultiProxy `json:"-"`
	*http.Transport            `json:"-"`
	*http.Client               `json:"-"`
	RegularProxy               string `json:"-"`

	Sites     []Site `json:"Sites,omitempty"`
	Peers     []Site `json:"Peers,omitempty"`
	hostsfile string `json:"-"`
	up        bool   `json:"-"`
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
							Url:    u,
							Dest:   append([]string{}, splitpair[1]),
							Base32: append([]string{}, b32),
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

// TODO: Buffer should be done differently if I want to do requests like this,
// or maybe sites should share a tunnel, or maybe I just use the HTTP proxy
// instead.

//AsyncGet is a misnomer, it has to be done in order for now.
func (c *Check) AsyncGet(index int, site *Site) {
	res, err := c.Client.Get(site.Url)
	if site.SuccessHistory == nil {
		site.SuccessHistory = make(map[time.Time]bool)
	}
	log.Println("CHECKING UPNESS")
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		site.SuccessHistory[time.Now()] = true
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		Title := strings.TrimSpace(doc.Find("Title").Text())
		if Title != "" {
			log.Println(Title)
			site.Title = Title
		}
		Desc := doc.Find("meta[name=Description]")
		content, ok := Desc.Attr("content")
		if ok {
			site.Desc = strings.TrimSpace(content)
			log.Println(strings.TrimSpace(content))
		}
		fmt.Printf("the eepSite is up: index=%d, Title=\"%s\", Desc=\"%s\", history=%d\n", index, site.Title, site.Desc, len(site.SuccessHistory))
	} else {
		fmt.Printf("the eepSite appears to be down: %v, %v\n", index, res.StatusCode)
		site.SuccessHistory[time.Now()] = false
	}
}

func (c *Check) CheckAll() {
	for index, site := range c.Sites {
		log.Printf("Checking URL: %s", site.Url)
		go c.AsyncGet(index, &site)
		time.Sleep(time.Second * 20)
	}
}

func (c *Check) QuerySite(s string) string {
	if u, err := Validate(s); err != nil {
		return "Not a valid URL for checking"
	} else {
		for index, site := range c.Sites {
			u2, err := Validate(site.Url)
			if err != nil {
				return "Invalid URL found in DB"
			}
			if u2 == u {
				c.AsyncGet(index, &c.Sites[index])
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
