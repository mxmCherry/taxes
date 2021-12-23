package table_test

import (
	"bytes"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mxmCherry/taxes/v2/internal/format/table"
	"github.com/mxmCherry/taxes/v2/internal/tax"
)

var _ = Describe("Format", func() {
	It("formats", func() {
		var calc *tax.CalcRun
		Expect(yaml.Unmarshal(read("testdata/data.yaml"), &calc)).To(Succeed())

		buf := bytes.NewBuffer(nil)
		subject := table.New(buf)
		Expect(subject.Format(calc)).To(Succeed())
		Expect(subject.Close()).To(Succeed())
		Expect(buf.String()).To(Equal(string(read("testdata/table.txt"))))
	})
})

func read(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return b
}
