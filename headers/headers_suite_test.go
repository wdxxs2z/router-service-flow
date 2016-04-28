package headers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHeaders(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Headers Suite")
}
