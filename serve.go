package gocheck

import (
	"fmt"
	"github.com/eyedeekay/goSam"
	"github.com/eyedeekay/sam-forwarder/tcp"
	"net/http"
	"strings"
)

func (c *Check) Parent() {
	err := c.SAMForwarder.Serve()
	if err != nil {
		panic(err)
	}
}

func (c *Check) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	if strings.Replace(rq.URL.Path, "/", "", -1) != "" {
		query := strings.SplitN(rq.URL.Path, "/", 1)
		fmt.Fprintf(rw, c.QuerySite(query[0]))
	} else {
		for index, site := range c.sites {
			fmt.Fprintf(rw, "%v: %s\n", index, site.HTML())
		}
	}
}

func (c *Check) Serve() error {
	go c.Parent()
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
	fmt.Println("Initializing eephttpd")
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	s.SAMForwarder.Config().SaveFile = true
	s.i2p, err = goSam.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	s.Transport = &http.Transport{
		Dial: s.i2p.Dial,
	}
	s.Client = &http.Client{
		Transport: s.Transport,
	}
	s.sites, err = LoadHostsFile(s.hostsfile)
	if err != nil {
		return nil, err
	}
	l, e := s.Load()
	if e != nil {
		return nil, e
	}
	return l.(*Check), nil
}
