package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/siprtcio/balanceservice/modules/log"
	"github.com/siprtcio/balanceservice/router"
)

const (
	DefaultConfFilePath = "configs/config.toml"
)

var (
	confFilePath string
	help         bool
)

// @title "Tiniyo Account Balance API"
// @version 1.0
// @name Tiniyo Account Balance API
// @description This document helps to understand to use account balance apis to manage your tenants balance. For this API's, you can use the your tiniyo `AuthID` as Key and `AuthSecretID` as password for `Basic auth`. If you need any help for integration just reach out to us at Tiniyo at [`support@tiniyo.com`](support@tiniyo.com).
// @termsOfService https://tiniyo.com/legal/tos.html
// @contact.name API Support
// @contact.url http://www.tiniyo.com/support
// @contact.email support@tiniyo.com
// @host api.tiniyo.com
// @BasePath /v1
// @securityDefinitions.basic BasicAuth
func init() {
	flag.StringVar(&confFilePath, "c", DefaultConfFilePath, "Config Path")
	flag.StringVar(&confFilePath, "config", DefaultConfFilePath, "Config Path")
	flag.BoolVar(&help, "h", false, "Show help message")
	flag.BoolVar(&help, "help", false, "Show help message")
	flag.Parse()
	flag.Usage = usage
}

func usage() {
	s := `balanceservice : a Application for Phone Number Management
		Usage: balanceservice [Options...]
		Options:
    		-c,  -config=<path>           Config path of the site. Default is configs/config.toml.
		Other options:
    		-h,  -help                  Show help message.
		`
	fmt.Printf(s)
	os.Exit(0)
}

func main() {

	log.Info("Application for Managing Users")

	if help {
		usage()
		return
	}
	log.Debugf("run with conf:%s", confFilePath)
	router.RunSubdomains(confFilePath)
}
