package cron

import (
	"errors"
	"slices"
)

type CronFragmentType string

var (
	WILDCARD OperatorType = "WILDCARD"
	DIVISOR  OperatorType = "DIVISOR"
	LIST     OperatorType = "LIST"
	RANGE    OperatorType = "RANGE"
	SINGLE   OperatorType = "SINGLE"

	cronOutputOrder = []CronFragmentType{
		"MINUTE",
		"HOUR",
		"DAY",
		"MONTH",
		"WEEKDAY",
	}
)

func (c Cron) ExpressionOrder() []CronFragmentType { return cronOutputOrder }

type OperatorFn func() ([]uint8, error)
type Operator interface {
	Fn() ([]uint8, error)
}
type OperatorType string
type CronFragment struct {
	Expr    string
	kind    OperatorType
	factors []uint8
}

type Cron []CronFragment

func (c Cron) Eq(other Cron) bool {
	for idx, cf := range c {
		if cf.Expr != other[idx].Expr {
			return false
		}

		if cf.kind != other[idx].kind {
			return false
		}

		if !slices.Equal(cf.factors, other[idx].factors) {
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
