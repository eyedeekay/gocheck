package gocheck

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/eyedeekay/httptunnel/multiproxy"

	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"

	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO calculate approx 9's
// TODO save check history(more json?)

type Site struct {
	Url      string   `json:Url,omitempty`
	Dest     []string `json:Dest,omitempty`
	Base32   []string `json:Base32,omitempty`
	Uptime   int      `json:Uptime,omitempty`
	Downtime int      `json:Downtime,omitempty`

	Title          string             `json:Title,omitempty`
	Desc           string             `json:Desc,omitempty`
	SuccessHistory map[time.Time]bool `json:SuccessHistory,omitempty`
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

	Sites       []Site `json:"Sites,omitempty"`
	Peers       []Site `json:"Peers,omitempty"`
	ScriptJS    string `json:"ScriptJS,omitempty"`
	StyleCSS    string `json:"ScriptJS,omitempty"`
	importhosts string `json:"-"`
	importpeers string `json:"-"`
	hostsfile   string `json:"-"`
	peersfile   string `sjon:"-"`
	up          bool   `json:"-"`

	mutex sync.Mutex
}

// TODO: Buffer should be done differently if I want to do requests like this,
// or maybe sites should share a tunnel, or maybe I just use the HTTP proxy
// instead.

//AsyncGet is a misnomer, it has to be done in order for now.
func (c Check) AsyncGet(index int, url string, sh map[time.Time]bool) (success bool, title string, desc string) {
	res, err := c.Client.Get(url)
	if err != nil {
		return
	}
	//	if site.SuccessHistory == nil {
	//		site.SuccessHistory = make(map[time.Time]bool)
	//	}
	log.Println("CHECKING UPNESS")
	defer res.Body.Close()
	if res.StatusCode == 200 {
		//		site.SuccessHistory[time.Now()] = true
		success = true
		//		log.Println(site.SuccessHistory)
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		Title := strings.TrimSpace(doc.Find("Title").Text())
		if Title != "" {
			log.Println(Title)
			title = Title
		}
		Desc := doc.Find("meta[name=description]")
		content, ok := Desc.Attr("content")
		if ok {
			desc = strings.TrimSpace(content)
			log.Println(strings.TrimSpace(content))
		}
		Desc = doc.Find("meta[name=Description]")
		content, ok = Desc.Attr("content")
		if ok {
			desc = strings.TrimSpace(content)
			log.Println(strings.TrimSpace(content))
		}
		return
	} else {
		success = false
		return
	}
}

func (c Check) CheckAll() {
	c.mutex.Lock()
	check := 0
	for index := range c.Sites {
		go func() {
			i := index
			log.Printf("Checking URL: %s", c.Sites[i].Url)
			s, t, d := c.AsyncGet(i, c.Sites[i].Url, c.Sites[i].SuccessHistory)
			c.Sites[i].SuccessHistory[time.Now()] = s
			c.Sites[i].Title = t
			c.Sites[i].Desc = d
			if s {
				fmt.Printf("the eepSite is up: index=%d, Title=\"%s\", Desc=\"%s\", history=%d\n", i, t, d, len(c.Sites[i].SuccessHistory))
			} else {
				fmt.Printf("the eepSite appears to be down: index=%d, url=%s\n", i, c.Sites[i].Url)
			}
			check++
		}()
		time.Sleep(time.Second * 2)
	}
	for {
		if len(c.Sites) == check {
			break
		}
	}
	c.mutex.Unlock()
}

func (c Check) QuerySite(s string) string {
	if u, err := Validate(s); err != nil {
		return "Not a valid URL for checking"
	} else {
		for index, site := range c.Sites {
			u2, err := Validate(site.Url)
			if err != nil {
				return "Invalid URL found in DB"
			}
			if u2 == u {
				s, t, d := c.AsyncGet(index, c.Sites[index].Url, c.Sites[index].SuccessHistory)
				c.Sites[index].SuccessHistory[time.Now()] = s
				c.Sites[index].Title = t
				c.Sites[index].Desc = d
				if s {
					fmt.Printf("the eepSite is up: index=%d, Title=\"%s\", Desc=\"%s\", history=%d\n", index, t, d, len(site.SuccessHistory))
				} else {
					fmt.Printf("the eepSite appears to be down: %v, %s\n", index, site.Url)
				}
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
	s.Sites, e = s.LoadHostsFile(s.hostsfile)
	if e != nil {
		if strings.Contains(e.Error(), "Error hosts file not given") {
			return nil, e
		}
	}
	s.Peers, e = s.LoadHostsFile(s.peersfile)
	if e != nil {
		if strings.Contains(e.Error(), "Error hosts file not given") {
			return nil, e
		}
	}
	s.SAMHTTPProxy = g.(*i2pbrowserproxy.SAMMultiProxy)
	if e := s.ImportSites(s.importhosts); e != nil {
		return nil, e
	}
	if e := s.ImportPeers(s.importpeers); e != nil {
		return nil, e
	}

	s.up = true
	fmt.Printf("Finished putting tunnel up")
	return s, nil
}
