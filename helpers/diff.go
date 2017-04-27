package helpers

import (
	"fmt"
	"io"
	"strings"

	"github.com/aryann/difflib"
	"github.com/mgutz/ansi"
)

func RenderDiff(w io.Writer, first string, second string) {
	diffs := difflib.Diff(strings.Split(first, "\n"), strings.Split(second, "\n"))

	for _, diff := range diffs {
		text := diff.Payload

		switch diff.Delta {
		case difflib.RightOnly:
			fmt.Fprintf(w, "%s\n", ansi.Color(text, "green"))
		case difflib.LeftOnly:
			fmt.Fprintf(w, "%s\n", ansi.Color(text, "red"))
		case difflib.Common:
			fmt.Fprintf(w, "%s\n", text)
		}
	}
}
