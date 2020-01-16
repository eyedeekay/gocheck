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
	url := strings.TrimPrefix(rq.URL.Path, "/")
	if url == "style.css" {
		file, err := ioutil.ReadFile("style.css")
		if err == nil {
			rw.Header().Set("Content-Type", "text/css")
			fmt.Fprintf(rw, "%s", file)
		} else {
			log.Println(err)
		}
		return
	}
	if url == "script.js" {
		file, err := ioutil.ReadFile("script.js")
		if err == nil {
			rw.Header().Set("Content-Type", "text/javascript")
			fmt.Fprintf(rw, "%s", file)
		} else {
			log.Println(err)
		}
		return
	}
	if url == "hosts.txt" {
		fmt.Fprintf(rw, "%s", c.ExportHostsFile())
		return
	}
	if url == "export.json" {
		fmt.Fprintf(rw, "%s", c.ExportJsonArtifact())
		return
	}
	if url == "export-mini.json" {
		fmt.Fprintf(rw, "%s", c.ExportMiniJsonArtifact())
		return
	}
	log.Println("URL", url)
	if url != "" {
		query := strings.SplitN(url, "/", 1)
		fmt.Fprintf(rw, c.QuerySite(query[0]))
	} else if strings.HasPrefix(url, "web") {
		c.DisplayPage(rw, rq, url)
	} else {
		c.DisplayPage(rw, rq, "")
	}
}

func (c *Check) DisplayPage(rw http.ResponseWriter, rq *http.Request, page string) {
	fmt.Fprintf(rw, "<!DOCTYPE html>")
	fmt.Fprintf(rw, "<html>")
	fmt.Fprintf(rw, "<head>")
	fmt.Fprintf(rw, "    <meta charset=\"utf-8\">")
	fmt.Fprintf(rw, "    <link type=\"text/css\" href=\"style.css\" rel=\"stylesheet\">")
	fmt.Fprintf(rw, "    <title>I2P Site Uptime Checker </title>")
	fmt.Fprintf(rw, "</head>")
	fmt.Fprintf(rw, "<body>")
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
	for index, site := range c.sites {
		if site.url == name {
			fmt.Fprintf(rw, "<div class=\"idnum\" id=\"%v\">%v: %s\n", index, index, site.HTML())
		}
	}
}

func (c *Check) AllPages(rw http.ResponseWriter, rq *http.Request) {
	for index, site := range c.sites {
		if len(site.successHistory) > 0 {
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

	s.sites, err = s.LoadHostsFile(s.hostsfile)
	if err != nil {
		return nil, err
	}
	l, e := s.Load()
	if e != nil {
		return nil, e
	}
	return l.(*Check), nil
}
