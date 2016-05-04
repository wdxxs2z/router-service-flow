package proxy

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/wdxxs2z/router-service-flow/headers"
	"github.com/wdxxs2z/router-service-flow/models"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
	"net/url"
	"strconv"
)

const (
	VcapCookieId    = "__VCAP_ID__"
	StickyCookieKey = "JSESSIONID"
)

func NewReverseProxy(transport http.RoundTripper, httpClient *http.Client, debug bool, ratioMark map[string]string) *httputil.ReverseProxy {

	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			start := time.Now()
			RouterServiceheader := headers.NewRouteServiceHeaders()

			err := RouterServiceheader.ParseHeadersAndClean(&req.Header)

			//session
			if _, err := req.Cookie(StickyCookieKey); err == nil {
				if sticky, err := req.Cookie(VcapCookieId); err == nil {
					log.Println(sticky.Value)
				}
			}

			if RouterServiceheader.IsValidRequest() && err == nil {
				cNormal := models.NewNormal()
				i := 0
				for ratio, ul := range ratioMark {
					if ratioNumber,err := strconv.ParseInt(ratio,10,32) ; err != nil {
						cNormal.AddNode(models.NewNode(i,ul,ratioNumber))
					}
					i++
				}
				winUrl := cNormal.GetWinUrl()
				req.URL, err = url.Parse(winUrl)
				req.Host = req.URL.Host
			} else {
				req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
				req.Host = "No Host"
				log.Print("Header are not Valid")
			}

			if debug {
				dump, err := httputil.DumpRequest(req, true)
				if err != nil {
					log.Fatalln(err.Error())
				}
				log.Printf("%q", dump)
				log.Printf("Time Elapsed header %v ", time.Since(start))

			}

		},
		Transport: transport,
	}
	return reverseProxy
}