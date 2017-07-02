package bankgovua_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBankgovua(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bankgovua Suite")
}
