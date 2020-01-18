package gocheck

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func Validate(u string) (string, error) {
	if !strings.Contains(u, ".i2p") {
		return "", fmt.Errorf("Not an I2P domain:%s", u)
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
