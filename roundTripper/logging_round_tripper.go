package roundTripper

import (
	"crypto/tls"
	"errors"
	"github.com/wdxxs2z/router-service-flow/headers"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type LoggingRoundTripper struct {
	transport http.RoundTripper
	debug     bool
}

func NewLoggingRoundTripper(debug bool) *LoggingRoundTripper {
	return &LoggingRoundTripper{
		transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		debug: debug,
	}
}

// forward to the URL
// Send response back to the Router

func (lrt *LoggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	var err error
	var res *http.Response
	start := time.Now()
	if request.Host == "No Host" {
		return nil, errors.New("Invalid Header")
	}
	res, err = lrt.transport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if lrt.debug {
		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			log.Fatalln(err.Error())
		}

		log.Printf("%q", dump)
		log.Printf("Time Elapsed RoundTrip %v", time.Since(start))

	}
	//Adding CF headers
	res.Header.Add(headers.RouteServiceMetadata, request.Header.Get(headers.RouteServiceMetadata))
	res.Header.Add(headers.RouteServiceSignature, request.Header.Get(headers.RouteServiceSignature))
	return res, err
}
