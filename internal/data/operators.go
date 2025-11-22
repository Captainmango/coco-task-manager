package data

import "fmt"

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
	switch cf.FragmentType {
	case MINUTE:
		for _, num := range cf.factors {
			if num > 60 {
				return fmt.Errorf("number outside of minute range %d is not within 0 to 60 (inclusive)", num)
			}
		}
		return nil
	case HOUR:
		for _, num := range cf.factors {
			if num > 23 {
				return fmt.Errorf("number outside of hour range %d is not within 0 to 23 (inclusive)", num)
			}
		}
		return nil
	case DAY:
		for _, num := range cf.factors {
			if num < 1 || num > 31 {
				return fmt.Errorf("number outside of day range %d is not within 1 to 7 (inclusive)", num)
			}
		}
		return nil
	case MONTH:
		for _, num := range cf.factors {
			if num < 1 || num > 12 {
				return fmt.Errorf("number outside of month range %d is not within 1 to 12 (inclusive)", num)
			}
		}
		return nil
	case WEEKDAY:
		for _, num := range cf.factors {
			if num < 1 || num > 12 {
				return fmt.Errorf("number outside of month range %d is not within 1 to 12 (inclusive)", num)
			}
		}
		return nil
	default:
		return nil
	}
}

func wildcard(cf CronFragment) ([]uint8, error) {
	cf.validate()
	return []uint8{}, nil
}

func rangeOp(cf CronFragment) ([]uint8, error) {
	cf.validate()
	return []uint8{}, nil
}

func list(cf CronFragment) ([]uint8, error) {
	cf.validate()
	return []uint8{}, nil
}

func divisor(cf CronFragment) ([]uint8, error) {
	cf.validate()
	return []uint8{}, nil
}

func single(cf CronFragment) ([]uint8, error) {
	cf.validate()
	return cf.factors, nil
}
