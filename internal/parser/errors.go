package parser

import (
	"errors"
	"fmt"
)

const (
	errInvalidInputFmt = "%s is not a valid input"
	errMalformedCronFmt = "malformed cron expression: '%s' invalid character after postion %d"
	errTooManySpacesFmt = "malformed cron expression: '%s' too many spaces at position %d"
)

func ErrInvalidInput(val any) error {
	return fmt.Errorf(errInvalidInputFmt, val)
}

func ErrUnableToPullNextFragment() error {
	return errors.New("next fragment not available")
}

func ErrMalformedCron(input any, position uint8) error {
	return fmt.Errorf(errMalformedCronFmt, input, position)
}

func ErrTooManySpaces(input any, position uint8) error {
	return fmt.Errorf(errTooManySpacesFmt, input, position)
}

func invalidCronTabEntry(input string) error {
	return fmt.Errorf("%s is not a valid crontab entry", input)
}

