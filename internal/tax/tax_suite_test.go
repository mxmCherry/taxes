package tax_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTax(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tax Suite")
}
