package data

import (
	"errors"
	"fmt"
)

const (
	ErrUnknownBoundsTypeFmt = "unknown bounds type %s"
	ErrFactorsOutsideBoundsFmt = "number outside of %s range %d is not within %d to %d (inclusive)"
	ErrUnknownFragmentTypeFmt = "unknown fragment type %s"
	ErrInvalidFragmentKindFmt = "invalid fragment kind %s"
)

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
	return fmt.Errorf(ErrUnknownBoundsTypeFmt, cft)
}

func ErrFactorsOutsideBounds(fragmentType, factor, lower, upper any) error {
	return fmt.Errorf(ErrFactorsOutsideBoundsFmt, fragmentType, factor, lower, upper)
}

func ErrUnknownFragmentType(fragmentType CronFragmentType) error {
	return fmt.Errorf(ErrUnknownFragmentTypeFmt, fragmentType)
}

func ErrInvalidFragmentKind(fragmentKind OperatorType) error {
	return fmt.Errorf(ErrInvalidFragmentKindFmt, fragmentKind)
}
