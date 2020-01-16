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
	url := strings.Replace(rq.URL.Path, "/", "", -1)
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
	if strings.Replace(rq.URL.Path, "/", "", -1) != "" {
		query := strings.SplitN(url, "/", 1)
		fmt.Fprintf(rw, c.QuerySite(query[0]))
	} else {
		fmt.Fprintf(rw, "<!DOCTYPE html>")
		fmt.Fprintf(rw, "<html>")
		fmt.Fprintf(rw, "<head>")
		fmt.Fprintf(rw, "    <meta charset=\"utf-8\">")
		fmt.Fprintf(rw, "    <link type=\"text/css\" href=\"style.css\" rel=\"stylesheet\">")
		fmt.Fprintf(rw, "    <title>I2P Site Uptime Checker </title>")
		fmt.Fprintf(rw, "</head>")
		fmt.Fprintf(rw, "<body>")
		for index, site := range c.sites {
			fmt.Fprintf(rw, "<div class=\"idnum\" id=\"%v\">%v: %s\n", index, index, site.HTML())
		}
		fmt.Fprintf(rw, "</body>")
		fmt.Fprintf(rw, "</html>")
	}
}

func (c *Check) CheckLoop() {
	time.Sleep(time.Second * 10)
	for {
		c.CheckAll()
		time.Sleep(time.Minute * 180)
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
