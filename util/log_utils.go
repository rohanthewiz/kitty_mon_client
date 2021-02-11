package util

import (
	"crypto/sha1"
	"fmt"
	"kitty_mon_client/config"
	"log"
	"net"
	"strings"
	"time"
)

var Fpf = fmt.Printf
var Fpl = fmt.Println
var Lpl = log.Println
var Lf = log.Fatal

func Pd(params ...interface{}) {
	if config.Opts.Debug {
		log.Println(params...)
	}
}

func Pl(params ...interface{}) {
	if config.Opts.Verbose {
		fmt.Println(params...)
	}
}

func Pf(msg string, params ...interface{}) {
	if config.Opts.Verbose {
		fmt.Printf(msg, params...)
	}
}

func Random_sha1() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte("%$"+time.Now().String()+"e{")))
}

func Short_sha(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}

func Trim_whitespace(in_str string) string {
	return strings.Trim(in_str, " \n\r\t")
}

// Return all IPs on Geoforce subnets
func IPs(class_c_only bool) string {
	var ret_addr []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		addr_str := addr.String()
		if class_c_only && !strings.Contains(addr_str, "192.") {
			continue
		}
		if strings.Contains(addr_str, ".") {
			ret := strings.Split(addr_str, "/")[0]
			ret_addr = append(ret_addr, ret)
		}
	}
	return strings.Join(ret_addr, ", ")
}
