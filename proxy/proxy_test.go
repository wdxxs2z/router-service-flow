package proxy_test

import (
	//"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/shinji62/route-service-cf/proxy"
	"github.com/shinji62/route-service-cf/roundTripper"
	"net/http"
	"net/http/httptest"
	//"net/url"
)

var _ = Describe("Proxy", func() {
	var backend *ghttp.Server
	var proxyServer *httptest.Server
	var cookieServer *ghttp.Server
	var statusCode int
	var getReq *http.Request

	BeforeEach(func() {
		roundTripper := roundTripper.NewLoggingRoundTripper(false)
		proxyServ := proxy.NewReverseProxy(roundTripper, &http.Client{}, false)
		proxyServer = httptest.NewServer(proxyServ)
		backend = ghttp.NewServer()
		cookieServer = ghttp.NewServer()
		getReq, _ = http.NewRequest("GET", proxyServer.URL, nil)
		getReq.Host = "some-name"
		getReq.Close = true

	})

	AfterEach(func() {
		backend.Close()
		proxyServer.Close()
		cookieServer.Close()
	})

	Describe("With invalid Header", func() {
		BeforeEach(func() {
			statusCode = http.StatusOK
			var test []byte
			backend.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyBody(test),
				ghttp.RespondWithJSONEncodedPtr(&statusCode, "TEST")))
		})
		It("Should reponse 500 ", func() {
			res, err := http.DefaultClient.Do(getReq)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(res.StatusCode).Should(Equal(500))

		})
	})

	Describe("With Valid Header", func() {
		BeforeEach(func() {
			getReq.Header.Set("X-CF-Proxy-Metadata", "bar")
			getReq.Header.Set("X-CF-Forwarded-Url", backend.URL())
			getReq.Header.Set("X-CF-Proxy-Signature", "foo")
			statusCode = http.StatusOK
			backend.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyBody([]byte{}),
				ghttp.RespondWithJSONEncodedPtr(&statusCode, "TEST")))
		})
		It("Should reponse with 200 to GET Request", func() {
			backend.WrapHandler(0, ghttp.VerifyRequest("GET", "/"))
			res, err := http.DefaultClient.Do(getReq)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(res.StatusCode).Should(Equal(200))
		})

		It("Should reponse with 200 to POST Request", func() {
			backend.WrapHandler(0, ghttp.VerifyRequest("POST", "/"))
			getReq.Method = "POST"
			res, err := http.DefaultClient.Do(getReq)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(res.StatusCode).Should(Equal(200))
		})
		It("Should reponse with 200 to PUT Request", func() {
			backend.WrapHandler(0, ghttp.VerifyRequest("PUT", "/"))
			getReq.Method = "PUT"
			res, err := http.DefaultClient.Do(getReq)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(res.StatusCode).Should(Equal(200))
		})
		It("Should reponse with 404 to GET Request", func() {
			statusCode = http.StatusNotFound
			res, err := http.DefaultClient.Do(getReq)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(res.StatusCode).Should(Equal(404))
			Ω(res.Status).Should(Equal("404 Not Found"))
		})
		/*It("Should send directly status.html", func() {
			backend.WrapHandler(0, ghttp.VerifyRequest("GET", "/status.html"))
			res, err := http.DefaultClient.Do(getReq)
			getReq.URL = url.Parse(fmt.Sprintf("%s/status.html", proxyServer.URL))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(res.StatusCode).Should(Equal(200))
		})*/

	})

})
