package utils

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

func ValidateRuneLength(value string, min, max int, errMsg string) error {
	if min < 0 || max < min {
		return fmt.Errorf("invalid length range: min=%d max=%d", min, max)
	}

	length := utf8.RuneCountInString(value)
	if length < min || length > max {
		if errMsg != "" {
			return errors.New(errMsg)
		}
		return fmt.Errorf("length must be between %d and %d characters", min, max)
	}

	return nil
}
