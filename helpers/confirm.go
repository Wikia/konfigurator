package helpers

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func AskConfirm(w io.Writer, r io.Reader, prompt string) (bool, error) {
	reader := bufio.NewReader(r)

	for {
		_, _ = fmt.Fprint(w, prompt+" [Y/n]")

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		// default action
		if len(response) == 0 {
			return true, nil
		}

		if response == "yes" || response == "y" {
			return true, nil
		} else {
			return false, nil
		}
	}
}
