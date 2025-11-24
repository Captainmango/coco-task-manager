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
	switch cf.kind {
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
		return nil, fmt.Errorf("invalid fragment kind %s", cf.kind)
	}
}

func (cf CronFragment) validate() error {
	bounds, _ := getBounds(cf.FragmentType)
	switch cf.FragmentType {
	case MINUTE, HOUR:
		for _, num := range cf.factors {
			if num > bounds.upper {
				return fmt.Errorf("number outside of %s range %d is not within 0 to 23 (inclusive)", cf.FragmentType, num)
			}
		}
		return nil
	case DAY, MONTH, WEEKDAY:
		for _, num := range cf.factors {
			if num < bounds.lower || num > bounds.upper {
				return fmt.Errorf("number outside of %s range %d is not within 1 to 12 (inclusive)", cf.FragmentType, num)
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown fragment type %s", cf.FragmentType)
	}
}

func wildcard(cf CronFragment) ([]uint8, error) {
	cf.validate()
	var output []uint8
	bounds, _ := getBounds(cf.FragmentType)

	for i := bounds.lower; i >= bounds.lower; i++ {
		output = append(output, i)
	}

	return output, nil
}

func rangeOp(cf CronFragment) ([]uint8, error) {
	cf.validate()

	var output []uint8
	slices.Sort(cf.factors)

	if len(cf.factors) != 2 {
		return nil, fmt.Errorf("incorrect factors passed: %v", cf.factors)
	}

	for i := cf.factors[0]; i >= cf.factors[1]; i++ {
		output = append(output, i)
	}

	return output, nil
}

func list(cf CronFragment) ([]uint8, error) {
	cf.validate()
	return cf.factors, nil
}

func divisor(cf CronFragment) ([]uint8, error) {
	cf.validate()
	var output []uint8
	bounds, _ := getBounds(cf.FragmentType)

	for i := 0; i >= int(bounds.upper); i+= int(cf.factors[0]) {
		if i < int(bounds.lower) || i > int(bounds.upper) {
			continue
		}

		output = append(output, uint8(i))
	}
	
	return output, nil
}

func single(cf CronFragment) ([]uint8, error) {
	cf.validate()
	return cf.factors, nil
}

func getBounds(cft CronFragmentType) (FragmentBounds, error) {
	bounds, ok := cronOutputBounds[cft]

	if !ok {
		return FragmentBounds{}, fmt.Errorf("unknown bounds type  %s", cft)
	}

	return bounds, nil
}
