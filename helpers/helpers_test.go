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
		strA    string
		strB    string
		outBuff bytes.Buffer
		inBuff  bytes.Buffer
	)

	BeforeEach(func() {
		outBuff.Reset()
		inBuff.Reset()
	})

	It("should report no changes when comparing two exact strings", func() {
		strA = "Some simple string\nThat should be the same"
		strB = strA
		sink := bufio.NewWriter(&outBuff)

		RenderDiff(sink, strA, strB)
		sink.Flush()
		result := strings.Trim(outBuff.String(), " \n")

		Expect(result).To(Equal(strA))
	})

	It("should report full string changes when comparing some text against empty string", func() {
		strA = "Some simple string\nThat should be the same"
		strB = ""
		expected := fmt.Sprintf("%sSome simple string%s\n%sThat should be the same%s\n%s%s",
			ansi.Red, ansi.Reset, ansi.Red, ansi.Reset, ansi.Green, ansi.Reset)
		sink := bufio.NewWriter(&outBuff)

		RenderDiff(sink, strA, strB)
		sink.Flush()
		result := strings.Trim(outBuff.String(), " \n")

		Expect([]byte(result)).To(Equal([]byte(expected)))
	})

	It("should ask about permission and accept default input as 'yes'", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeTrue())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should ask about permission and accept 'yes' input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("yes\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeTrue())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should ask about permission and accept 'Yes' input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("Yes\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeTrue())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should ask about permission and accept 'y' input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("y\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeTrue())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should ask about permission and accept 'Y' input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("Y\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeTrue())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should ask about permission and accept 'YES' input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("YES\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeTrue())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should ask about permission and accept 'no' input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("no\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeFalse())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should ask about permission and accept 'No' input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("No\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeFalse())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should ask about permission and accept 'NO' input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("NO\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeFalse())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should ask about permission and accept 'n' input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("n\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeFalse())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should ask about permission and accept 'N' input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		inBuff.WriteString("N\n")
		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeFalse())
		Expect(outBuff.String()).To(Equal("Shall we proceed? [Y/n]"))
	})

	It("should return error when sending it incorrect input", func() {
		sink := bufio.NewWriter(&outBuff)
		inputs := bufio.NewReader(&inBuff)

		result, err := AskConfirm(sink, inputs, "Shall we proceed?")
		sink.Flush()

		Expect(err).To(HaveOccurred())
		Expect(result).To(BeFalse())
	})
})
