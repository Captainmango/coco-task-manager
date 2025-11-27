package data

import (
	"iter"
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

func (c Cron) ExpressionOrder() (func() (CronFragmentType, bool), func()) {
	cronOrderIterNext, cronIterStop := iter.Pull(slices.Values(cronOutputOrder))

	return cronOrderIterNext, cronIterStop
}

type FragmentBounds struct {
	upper uint8
	lower uint8
}

type CronFragment struct {
	Expr         string
	FragmentType CronFragmentType
	Kind         OperatorType
	Factors      []uint8
}

type Cron struct {
	Data []CronFragment
}

func (c Cron) Eq(other Cron) bool {
	for idx, cf := range c.Data {
		if cf.Expr != other.Data[idx].Expr {
			return false
		}

		if cf.FragmentType != other.Data[idx].FragmentType {
			return false
		}

		if cf.Kind != other.Data[idx].Kind {
			return false
		}

		if !slices.Equal(cf.Factors, other.Data[idx].Factors) {
			return false
		}
	}

	return true
}

func NewWildCardFragment(expr string) (CronFragment, error) {
	return CronFragment{
		Expr:    expr,
		Kind:    WILDCARD,
		Factors: nil,
	}, nil
}

func NewRangeFragment(expr string, factors []uint8) (CronFragment, error) {
	if len(factors) != 2 {
		return CronFragment{}, ErrInvalidRangeFragment()
	}

	return CronFragment{
		Expr:    expr,
		Kind:    RANGE,
		Factors: factors,
	}, nil
}

func NewListFragment(expr string, factors []uint8) (CronFragment, error) {
	if len(factors) < 1 {
		return CronFragment{}, ErrInvalidListFragment()
	}

	return CronFragment{
		Expr:    expr,
		Kind:    LIST,
		Factors: factors,
	}, nil
}

func NewDivisorFragment(expr string, factors []uint8) (CronFragment, error) {
	if len(factors) != 1 {
		return CronFragment{}, ErrInvalidDivisorFragment()
	}

	return CronFragment{
		Expr:    expr,
		Kind:    DIVISOR,
		Factors: factors,
	}, nil
}

func NewSingleFragment(expr string, factors []uint8) (CronFragment, error) {
	if len(factors) != 1 {
		return CronFragment{}, ErrInvalidSingleFragment()
	}

	return CronFragment{
		Expr:    expr,
		Kind:    SINGLE,
		Factors: factors,
	}, nil
}
