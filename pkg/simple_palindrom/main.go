package simple_palindrom

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

func IsSimplePalindrom(s string) (bool, error) {
	trimmedString := strings.Trim(s, " ")

	if trimmedString == "" {
		return false, errors.New("empty string")
	}

	trimmedStringLen := len(strings.Trim(trimmedString, " "))
	for i := range trimmedStringLen / 2 {
		if strings.Trim(s, " ")[i] != strings.Trim(s, " ")[trimmedStringLen-i-1] {
			return false, nil
		}
	}

	return true, nil
}

func IsSimplePalindromList(stringsForPalindrom []string) (map[string]bool, error) {
	palindromListResult := make(map[string]bool)

	for _, s := range stringsForPalindrom {
		isPalindrom, err := IsSimplePalindrom(s)

		if err != nil {
			slog.Error("error in func IsSimplePalindrom", s, fmt.Errorf("%w", err))
			palindromListResult[s] = false
			return nil, err
		}

		if isPalindrom {
			fmt.Printf("%q is palindrom\n", s)
			palindromListResult[s] = true
			continue
		}

		fmt.Printf("%q is NOT palindrom\n", s)
		palindromListResult[s] = false
	}

	return palindromListResult, nil
}
