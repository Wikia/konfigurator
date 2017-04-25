package helpers

import (
	"fmt"
	"strings"
)

func AskConfirm(prompt string) (bool, error) {
	var response string

	for {
		fmt.Print(prompt + " [Y/n]")

		_, err := fmt.Scanln(&response)
		if err != nil {
			return false, err
		}

		okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
		nayResponses := []string{"n", "N", "no", "No", "NO"}

		response = strings.Trim(response, " \n\r")
		for _, yes := range okayResponses {
			if yes == response {
				return true, nil
			}
		}

		for _, nay := range nayResponses {
			if nay == response {
				return false, nil
			}
		}
	}
}
