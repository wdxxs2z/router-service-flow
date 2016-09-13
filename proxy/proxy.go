package proxy

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/wdxxs2z/router-service-flow/headers"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
	"net/url"
	"github.com/wdxxs2z/router-service-flow/policy"
	//"../policy"
)

const (
	VcapCookieId    = "__VCAP_ID__"
	StickyCookieKey = "JSESSIONID"
)

func NewReverseProxy(transport http.RoundTripper, httpClient *http.Client, debug bool, policyType policy.PolicyType) *httputil.ReverseProxy {

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
				var winUrl string
				switch policyType.TypeName {
				case policy.POLICY_MODULO:
					policyModulo := policy.NewModulo(policyType.TypeName,policyType.Nodes)
					winUrl = policyModulo.WinUrl()
				case policy.POLICY_ROUNDROBIN:
					policyRoundrobin := policy.NewRoundRobin(policyType.TypeName,policyType.Nodes)
					winUrl = policyRoundrobin.WinUrl()
				case policy.POLICY_SOURCEHASH:
					policySourceHash := policy.NewSourceHash(policyType.TypeName,policyType.Nodes,req.RemoteAddr)
					winUrl = policySourceHash.WinUrl()
				case policy.POLICY_ROBIN_WEIGHT:
					policyRoundrobinweight := policy.NewRoundRobinWeight(policyType.TypeName,policyType.Nodes)
					winUrl = policyRoundrobinweight.WinUrl()
				}
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