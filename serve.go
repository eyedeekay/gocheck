package gocheck

import (
	"fmt"
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
	if rq.URL.Path != "" {
		query := strings.SplitN(rq.URL.Path, "/", 1)
		fmt.Fprintf(rw, c.QuerySite(query[0]))
	} else {
		for index, site := range c.sites {
			fmt.Fprintf(rw, "%v: %s\n", index, site.JsonString())
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
	s.SAMForwarder = &samforwarder.SAMForwarder{}
	fmt.Println("Initializing eephttpd")
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	s.SAMForwarder.Config().SaveFile = true
	l, e := s.Load()
	if e != nil {
		return nil, e
	}
	return l.(*Check), nil
}
