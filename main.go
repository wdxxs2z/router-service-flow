package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/shinji62/route-service-cf/proxy"
	"github.com/shinji62/route-service-cf/roundTripper"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
    "encoding/json"
)

var (
	port  = kingpin.Flag("port", "Port to listen").Envar("PORT").Short('p').Required().Int()
	debug = kingpin.Flag("debug", "Port to listen").Envar("DEBUG").Short('d').Bool()
    ratio = kingpin.Flag("ratio", "Http flow ratio").Envar("RATIO").Short('r').Required().String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
        ratiobyte := []byte(*ratio)
	var ratiourl map[string]string
	if err := json.Unmarshal(ratiobyte, &ratiourl); err != nil {
		panic(err)
	}
	httpClient := &http.Client{}
	roundTripper := roundTripper.NewLoggingRoundTripper(*debug)
	proxy := proxy.NewReverseProxy(roundTripper, httpClient, *debug, ratiourl)

	log.Fatal(http.ListenAndServe(":"+fmt.Sprintf("%v", *port), proxy))
}
