package bankgovua_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBankgovua(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bankgovua Suite")
}
