package helpers_test

import (
	. "github.com/Wikia/konfigurator/helpers"

	"bufio"
	"bytes"

	"strings"

	"fmt"

	"github.com/mgutz/ansi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", func() {
	var (
		strA string
		strB string
		buff bytes.Buffer
	)

	It("should report no changes when comparing two exact strings", func() {
		strA = "Some simple string\nThat should be the same"
		strB = strA
		sink := bufio.NewWriter(&buff)

		RenderDiff(sink, strA, strB)
		sink.Flush()
		result := strings.Trim(buff.String(), " \n")
		Expect(result).To(Equal(strA))
	})

	It("should report full string changes when comparing some text against empty string", func() {
		strA = "Some simple string\nThat should be the same"
		strB = ""
		expected := fmt.Sprintf("Some simple string\nThat should be the same\n%sSome simple string%s\n%sThat should be the same%s\n%s%s",
			ansi.Red, ansi.Reset, ansi.Red, ansi.Reset, ansi.Green, ansi.Reset)
		sink := bufio.NewWriter(&buff)

		RenderDiff(sink, strA, strB)
		sink.Flush()
		result := strings.Trim(buff.String(), " \n")
		Expect([]byte(result)).To(Equal([]byte(expected)))
	})
})
