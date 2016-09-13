package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/wdxxs2z/router-service-flow/proxy"
	//"./proxy"
	"github.com/wdxxs2z/router-service-flow/roundTripper"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"encoding/json"
	"github.com/wdxxs2z/router-service-flow/policy"
	//"./policy"
)

var (
	port = kingpin.Flag("port", "Port to listen").Envar("PORT").Short('p').Required().Int()
	debug = kingpin.Flag("debug", "Port to listen").Envar("DEBUG").Short('d').Bool()
	ratio = kingpin.Flag("ratio", "Http flow ratio").Envar("RATIO").Short('r').Required().String()
)

func main() {
	kingpin.Version("0.1.0")
	kingpin.Parse()
	ratiobyte := []byte(*ratio)
	policyRes := policy.PolicyType{}
	// {"typename": "modulo","nodes":[{"index":1,"url":"http://aaa.com","weight":2},{"index":1,"url":"http://bbb.com","weight":3}]}
	if err := json.Unmarshal(ratiobyte, &policyRes); err != nil {
		panic(err)
	}
	httpClient := &http.Client{}
	roundTripper := roundTripper.NewLoggingRoundTripper(*debug)
	proxy := proxy.NewReverseProxy(roundTripper, httpClient, *debug, policyRes)

	log.Fatal(http.ListenAndServe(":" + fmt.Sprintf("%v", *port), proxy))
}