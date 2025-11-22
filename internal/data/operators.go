package data

var (
	WILDCARD OperatorType = "WILDCARD"
	DIVISOR  OperatorType = "DIVISOR"
	LIST     OperatorType = "LIST"
	RANGE    OperatorType = "RANGE"
	SINGLE   OperatorType = "SINGLE"
)

type OperatorType string

func (cf CronFragment) getPossibleValues() ([]uint8, error) {
	switch cf.kind {
	case WILDCARD:
		return wildcard(cf)
	default:
		return nil, nil
	}
}

func wildcard(cf CronFragment) ([]uint8, error) {
	return []uint8{}, nil
}