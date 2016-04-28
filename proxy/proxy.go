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
    "math/rand"
	"strconv"
)

func NewReverseProxy(transport http.RoundTripper, httpClient *http.Client, debug bool, radioMark map[string]string) *httputil.ReverseProxy {

	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			start := time.Now()
			RouterServiceheader := headers.NewRouteServiceHeaders()

			err := RouterServiceheader.ParseHeadersAndClean(&req.Header)
			if RouterServiceheader.IsValidRequest() && err == nil {
                randomNum := rand.Intn(501)
				log.Println("The randomNum is: %v" , randomNum)
				sum := int64(0)
				for ratio,_ := range radioMark {
					ratioNum,err := strconv.ParseInt(ratio,10,32)
					if err == nil {
						sum += ratioNum
					}else{
						panic(err)
					}
				}
				mod := int64(randomNum) % sum
				log.Println("The mod is : %v", mod)
				for radio,ul := range radioMark {
					intratio,err := strconv.ParseInt(radio,10,32)
					if err == nil {
						intratio = int64(intratio)
					}else {
						panic(err)
					}
					// case 1
					if intratio > (sum - intratio) {
						if mod >= intratio {
							req.URL, err = url.Parse(ul)
							req.Host = req.URL.Host
							break
						}else{
							index := strconv.Itoa(int(sum)-int(intratio))
							req.URL, err = url.Parse(radioMark[index])
							req.Host = req.URL.Host
							break
						}
					}
					// case 2
					if intratio < (sum - intratio) {
						if mod < (sum -intratio){
							req.URL, err = url.Parse(ul)
							req.Host = req.URL.Host
							break
						}else{
							index := strconv.Itoa(int(sum)-int(intratio))
							req.URL, err = url.Parse(radioMark[index])
							req.Host = req.URL.Host
							break
						}
					}
				}
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
