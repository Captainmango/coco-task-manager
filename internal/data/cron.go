package data

import (
	"errors"
	"slices"
)

type CronFragmentType string

var (
	MINUTE  = CronFragmentType("MINUTE")
	HOUR    = CronFragmentType("HOUR")
	DAY     = CronFragmentType("DAY")
	MONTH   = CronFragmentType("MONTH")
	WEEKDAY = CronFragmentType("WEEKDAY")

	cronOutputOrder = []CronFragmentType{
		MINUTE,
		HOUR,
		DAY,
		MONTH,
		WEEKDAY,
	}

	cronOutputBounds = map[CronFragmentType]FragmentBounds{
		MINUTE:  {59, 0},
		HOUR:    {23, 0},
		DAY:     {31, 1},
		MONTH:   {12, 1},
		WEEKDAY: {7, 1},
	}
)

func (c Cron) ExpressionOrder() []CronFragmentType { return cronOutputOrder }

type FragmentBounds struct {
	upper uint8
	lower uint8
}
type CronFragment struct {
	Expr         string
	FragmentType CronFragmentType
	kind         OperatorType
	factors      []uint8
}

type Cron struct {
	Data []CronFragment
}

func (c Cron) Eq(other Cron) bool {
	for idx, cf := range c.Data {
		if cf.Expr != other.Data[idx].Expr {
			return false
		}

		if cf.kind != other.Data[idx].kind {
			return false
		}

		if !slices.Equal(cf.factors, other.Data[idx].factors) {
			return false
		}
	}

	return true
}

func NewWildCardFragment(expr string) (CronFragment, error) {
	return CronFragment{
		Expr:    expr,
		kind:    WILDCARD,
		factors: nil,
	}, nil
}

func NewRangeFragment(expr string, factors []uint8) (CronFragment, error) {
	if len(factors) != 2 {
		return CronFragment{}, errors.New("range only accepts 2 factors")
	}

	return CronFragment{
		Expr:    expr,
		kind:    RANGE,
		factors: factors,
	}, nil
}

func NewListFragment(expr string, factors []uint8) (CronFragment, error) {
	if len(factors) < 1 {
		return CronFragment{}, errors.New("list requires at least 2 factors")
	}

	return CronFragment{
		Expr:    expr,
		kind:    LIST,
		factors: factors,
	}, nil
}

func NewDivisorFragment(expr string, factors []uint8) (CronFragment, error) {
	if len(factors) != 1 {
		return CronFragment{}, errors.New("divisor rule only accepts one factor")
	}

	return CronFragment{
		Expr:    expr,
		kind:    DIVISOR,
		factors: factors,
	}, nil
}

func NewSingleFragment(expr string, factors []uint8) (CronFragment, error) {
	if len(factors) != 1 {
		return CronFragment{}, errors.New("divisor rule only accepts one factor")
	}

	return CronFragment{
		Expr:    expr,
		kind:    SINGLE,
		factors: factors,
	}, nil
}
