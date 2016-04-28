package headers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/shinji62/route-service-cf/headers"
	"net/http"
)

var routeService *RouteServiceHeaders

var _ = Describe("Headers", func() {
	Describe("Setup properly", func() {
		Context("with valid header", func() {
			headers := &http.Header{}
			headers.Add("X-CF-Proxy-Signature", "ProxySign")
			headers.Add("X-CF-Forwarded-Url", "https://localhost.com")
			headers.Add("X-CF-Proxy-Metadata", "Proxy_Metada")

			routeService = NewRouteServiceHeaders()
			err := routeService.ParseHeadersAndClean(headers)
			It("Should be valid", func() {
				Expect(routeService.IsValidRequest()).To(Equal(true))
				Expect(err).To(BeNil())
			})
			It("Shoud return valid URL", func() {
				Expect(routeService.ParsedUrl.Host).To(Equal("localhost.com"))
			})
			It("Should remove X-CF-Forwarded-Url Header", func() {
				Expect(headers.Get("X-CF-Forwarded-Url")).To(Equal(""))

			})
			It("Should not Remove valid headers", func() {
				Expect(headers.Get("X-CF-Proxy-Signature")).To(Equal("ProxySign"))
				Expect(headers.Get("X-CF-Proxy-Metadata")).To(Equal("Proxy_Metada"))
			})

		})

	})
})
