package gocheck

import (
	"fmt"
	"github.com/eyedeekay/httptunnel/multiproxy"
	"github.com/eyedeekay/sam-forwarder/config"
	"github.com/eyedeekay/sam-forwarder/tcp"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c *Check) Parent() {
	err := c.SAMForwarder.Serve()
	if err != nil {
		panic(err)
	}
}

func (c *Check) ParentHTTP() {
	err := c.SAMHTTPProxy.Serve()
	if err != nil {
		panic(err)
	}
}

func (c *Check) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	name := strings.TrimPrefix(rq.URL.Path, "/")
	if name == "style.css" {
		if c.StyleCSS != "" {
			file, err := ioutil.ReadFile(c.StyleCSS)
			if err == nil {
				rw.Header().Set("Content-Type", "text/css")
				fmt.Fprintf(rw, "%s", file)
			} else {
				log.Println(err)
			}
		}
		return
	}
	if name == "script.js" {
		if c.ScriptJS != "" {
			file, err := ioutil.ReadFile(c.ScriptJS)
			if err == nil {
				rw.Header().Set("Content-Type", "text/javascript")
				fmt.Fprintf(rw, "%s", file)
			} else {
				log.Println(err)
			}
		}
		return
	}
	if name == "hosts.txt" {
		fmt.Fprintf(rw, "%s", c.ExportHostsFile())
		return
	}
	if name == "export-sites.json" {
		fmt.Fprintf(rw, "%s", c.ExportJsonArtifact())
		return
	}
	if name == "export-peers.json" {
		fmt.Fprintf(rw, "%s", c.ExportMiniJsonArtifact())
		return
	}
	if name == "" {
		c.DisplayPage(rw, rq, "")
	} else if strings.HasPrefix(name, "web") {
		c.DisplayPage(rw, rq, name)
	} else {
		query := strings.SplitN(name, "/", 1)
		fmt.Fprintf(rw, c.QuerySite(query[0]))
	}
}

func (c *Check) DisplayPage(rw http.ResponseWriter, rq *http.Request, page string) {
	fmt.Fprintf(rw, "<!DOCTYPE html>")
	fmt.Fprintf(rw, "<html>")
	fmt.Fprintf(rw, "<head>")
	fmt.Fprintf(rw, "    <meta charset=\"utf-8\">")
	fmt.Fprintf(rw, "    <link type=\"text/css\" href=\"style.css\" rel=\"stylesheet\">")
	fmt.Fprintf(rw, "    <Title>I2P Site Uptime Checker </Title>")
	fmt.Fprintf(rw, "</head>")
	fmt.Fprintf(rw, "<body>")
	fmt.Fprintf(rw, "<span><a id=\"alive\" href=\"#\">Show Alive Only</a></span>")
	fmt.Fprintf(rw, "<span><a id=\"dead\" href=\"#\">Show Dead Only</a></span>")
	fmt.Fprintf(rw, "<span><a id=\"untested\" href=\"#\">Show Unknown Only</a></span>")
	body := "<h1>Gocheck: Site Uptime Checker</h1>\n"
	body += "<h2>Check a site status from the remote server</h2>\n"
	body += "<pre><code>    http_proxy=http://localhost:4444 curl http://5fma2okrcondmxkf4j2ggwuazaoo5d3z5moyh6wurgob4nthe3oa.b32.i2p/stats.i2p </code></pre>\n"
	body += "<h2>Get an exported history of all the sites this site knows about</h2>\n"
	body += "<pre><code>    http_proxy=http://localhost:4444 curl http://5fma2okrcondmxkf4j2ggwuazaoo5d3z5moyh6wurgob4nthe3oa.b32.i2p/export-sites.json </code></pre>\n"
	body += "<h2>Get an exported history of all the sites this site has history for already</h2>\n"
	body += "<pre><code>    http_proxy=http://localhost:4444 curl http://5fma2okrcondmxkf4j2ggwuazaoo5d3z5moyh6wurgob4nthe3oa.b32.i2p/export-peers.json </code></pre>\n"
	fmt.Fprintf(rw, body)
	log.Println("PAGE IS:", page)
	if page == "" {
		c.AllPages(rw, rq)
	} else {
		c.OnePage(rw, rq, page)
	}
	fmt.Fprintf(rw, "<script src=\"script.js\"></script>")
	fmt.Fprintf(rw, "</body>")
	fmt.Fprintf(rw, "</html>")
}

func (c *Check) OnePage(rw http.ResponseWriter, rq *http.Request, page string) {
	name, err := Validate(strings.TrimPrefix(page, "web"))
	if err != nil {
		log.Println(err)
		return
	}
	c.QuerySite(name)
	for index, site := range c.Sites {
		if site.Url == name {
			fmt.Fprintf(rw, "<div class=\"idnum\" id=\"%v\">%v: %s\n", index, index, site.HTML())
		}
	}
}

func (c *Check) AllPages(rw http.ResponseWriter, rq *http.Request) {
	for index, site := range c.Sites {
		log.Println("SITE LISTING", site.Up(), index, site.SuccessHistory)
		if len(site.SuccessHistory) > 0 {
			fmt.Fprintf(rw, "<div class=\"idnum "+site.Up()+"\" id=\"%v\">%v: %s\n", index, index, site.HTML())
		}
	}
}

func (c *Check) CheckLoop() {
	time.Sleep(time.Second * 10)
	for {
		c.CheckAll()
		time.Sleep(time.Minute * 60)
	}
}

func (c *Check) Serve() error {
	go c.Parent()
	if c.RegularProxy != "" {
		go c.ParentHTTP()
	}
	go c.CheckLoop()
	fmt.Printf("Starting web server", c.Target())
	if err := http.ListenAndServe(c.Target(), c); err != nil {
		return err
	}
	return nil
}

func NewSAMChecker(hostsfile string) (*Check, error) {
	return NewSAMCheckerFromOptions(SetHostsFile(hostsfile))
}

//NewSAMCheckerFromOptions makes a new SAM forwarder with default options, accepts host:port arguments
func NewSAMCheckerFromOptions(opts ...func(*Check) error) (*Check, error) {
	var s Check
	var err error
	s.SAMForwarder = &samforwarder.SAMForwarder{}
	s.SAMHTTPProxy = &i2pbrowserproxy.SAMMultiProxy{
		Conf: &i2ptunconf.Conf{},
	}

	fmt.Println("Initializing eephttpd")
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	s.SAMForwarder.Config().SaveFile = true
	var proxyURL *url.URL
	if s.RegularProxy != "no" {
		proxyURL, err = url.Parse("http://" + s.RegularProxy)
	} else {
		proxyURL, err = url.Parse("http://" + s.SAMHTTPProxy.Target())
	}
	if err != nil {
		return nil, err
	}
	log.Println("ProxyURL", proxyURL)
	proxy := http.ProxyURL(proxyURL)
	s.Transport = &http.Transport{
		Proxy: proxy,
	}
	s.Client = &http.Client{
		Timeout:   time.Duration(time.Minute * 5),
		Transport: s.Transport,
	}
	l, e := s.Load()
	if e != nil {
		return nil, e
	}
	return l.(*Check), nil
}
