package main

import (
	"flag"
	"github.com/eyedeekay/gocheck"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	hostspath := user.HomeDir + "/.i2p/hosts.txt"
	hostsfile := flag.String("hosts", hostspath, "Hosts file to use.")
	flag.Parse()
	check, err := gocheck.NewSAMChecker(*hostsfile)
	if err != nil {
		panic(err)
	}
	check.CheckAll()
	check.Serve()
}
