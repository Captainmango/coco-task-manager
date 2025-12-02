package parser

import (
	"errors"
	"fmt"
)

const (
	errInvalidInputFmt         = "%s is not a valid input"
	errMalformedCronFmt        = "malformed cron expression: '%s' invalid character after postion %d"
	errTooManySpacesFmt        = "malformed cron expression: '%s' too many spaces at position %d"
	errUnknownBoundsTypeFmt    = "unknown bounds type %s"
	errFactorsOutsideBoundsFmt = "number outside of %s range %d is not within %d to %d (inclusive)"
	errUnknownFragmentTypeFmt  = "unknown fragment type %s"
	errInvalidFragmentKindFmt  = "invalid fragment kind %s"
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

func ErrInvalidDivisorFragment() error {
	return errors.New("divisor rule only accepts one factor")
}

func ErrInvalidSingleFragment() error {
	return errors.New("single rule only accepts one factor")
}

func ErrInvalidListFragment() error {
	return errors.New("list requires at least 2 factors")
}

func ErrInvalidRangeFragment() error {
	return errors.New("range only accepts 2 factors")
}

func ErrUnknownBoundsType(cft CronFragmentType) error {
	return fmt.Errorf(errUnknownBoundsTypeFmt, cft)
}

func ErrFactorsOutsideBounds(fragmentType, factor, lower, upper any) error {
	return fmt.Errorf(errFactorsOutsideBoundsFmt, fragmentType, factor, lower, upper)
}

func ErrUnknownFragmentType(fragmentType CronFragmentType) error {
	return fmt.Errorf(errUnknownFragmentTypeFmt, fragmentType)
}

func ErrInvalidFragmentKind(fragmentKind OperatorType) error {
	return fmt.Errorf(errInvalidFragmentKindFmt, fragmentKind)
}
