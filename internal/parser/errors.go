package parser

import (
	"errors"
	"fmt"
)

const (
	ErrInvalidInputFmt = "%s is not a valid input"
	ErrMalformedCronFmt = "malformed cron expression: '%s' invalid character after postion %d"
	ErrTooManySpacesFmt = "malformed cron expression: '%s' too many spaces at position %d"
)

func ErrInvalidInput(val any) error {
	return fmt.Errorf(ErrInvalidInputFmt, val)
}

func ErrUnableToPullNextFragment() error {
	return errors.New("next fragment not available")
}

func ErrMalformedCron(input any, position uint8) error {
	return fmt.Errorf(ErrMalformedCronFmt, input, position)
}

func ErrTooManySpaces(input any, position uint8) error {
	return fmt.Errorf(ErrTooManySpacesFmt, input, position)
}


