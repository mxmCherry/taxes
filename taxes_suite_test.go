package taxes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTaxes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Taxes Suite")
}
