package data

import (
	"fmt"
	"slices"
)

var (
	WILDCARD OperatorType = "WILDCARD"
	DIVISOR  OperatorType = "DIVISOR"
	LIST     OperatorType = "LIST"
	RANGE    OperatorType = "RANGE"
	SINGLE   OperatorType = "SINGLE"
)

type OperatorType string

func (cf CronFragment) GetPossibleValues() ([]uint8, error) {
	switch cf.Kind {
	case WILDCARD:
		return wildcard(cf)
	case DIVISOR:
		return divisor(cf)
	case LIST:
		return list(cf)
	case RANGE:
		return rangeOp(cf)
	case SINGLE:
		return single(cf)
	default:
		return nil, ErrInvalidFragmentKind(cf.Kind)
	}
}

func (cf CronFragment) validate() error {
	bounds, err := getBounds(cf.FragmentType)

	if err != nil {
		return err
	}

	switch cf.FragmentType {
	case MINUTE, HOUR:
		for _, num := range cf.Factors {
			if num > bounds.upper {
				return ErrFactorsOutsideBounds(cf.FragmentType, num, bounds.lower, bounds.upper)
			}
		}
		return nil
	case DAY, MONTH, WEEKDAY:
		for _, num := range cf.Factors {
			if num < bounds.lower || num > bounds.upper {
				return ErrFactorsOutsideBounds(cf.FragmentType, num, bounds.lower, bounds.upper)
			}
		}
		return nil
	default:
		return ErrUnknownFragmentType(cf.FragmentType)
	}
}

func wildcard(cf CronFragment) ([]uint8, error) {
	if err := cf.validate(); err != nil {
		return nil, err
	}

	var output []uint8
	bounds, _ := getBounds(cf.FragmentType)

	for i := bounds.lower; i <= bounds.upper; i++ {
		output = append(output, i)
	}

	return output, nil
}

func rangeOp(cf CronFragment) ([]uint8, error) {
	if err := cf.validate(); err != nil {
		return nil, err
	}

	var output []uint8
	slices.Sort(cf.Factors)

	if len(cf.Factors) != 2 {
		return nil, fmt.Errorf("incorrect factors passed: %v", cf.Factors)
	}

	for i := cf.Factors[0]; i <= cf.Factors[1]; i++ {
		output = append(output, i)
	}

	return output, nil
}

func list(cf CronFragment) ([]uint8, error) {
	if err := cf.validate(); err != nil {
		return nil, err
	}

	slices.Sort(cf.Factors)

	return cf.Factors, nil
}

func divisor(cf CronFragment) ([]uint8, error) {
	if err := cf.validate(); err != nil {
		return nil, err
	}

	var output []uint8
	bounds, _ := getBounds(cf.FragmentType)

	for i := 0; i <= int(bounds.upper); i += int(cf.Factors[0]) {
		if i < int(bounds.lower) || i > int(bounds.upper) {
			continue
		}

		output = append(output, uint8(i))
	}

	return output, nil
}

func single(cf CronFragment) ([]uint8, error) {
	if err := cf.validate(); err != nil {
		return nil, err
	}
	return cf.Factors, nil
}

func getBounds(cft CronFragmentType) (FragmentBounds, error) {
	bounds, ok := cronOutputBounds[cft]

	if !ok {
		return FragmentBounds{}, ErrUnknownBoundsType(cft)
	}

	return bounds, nil
}
