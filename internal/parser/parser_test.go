package parser

import (
	"testing"
)

func Test_ItHandlesBasicInputs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected func() Cron
	}{
		{
			"wildcard",
			"*",
			func() Cron {
				cf, _ := NewWildCardFragment("*")
				cf.FragmentType = MINUTE

				return Cron{
					Data: []CronFragment{cf},
				}
			},
		},
		{
			"range",
			"1-5",
			func() Cron {
				cf, _ := NewRangeFragment("1-5", []uint8{1, 5})
				cf.FragmentType = MINUTE

				return Cron{
					Data: []CronFragment{cf},
				}
			},
		},
		{
			"list",
			"1,5",
			func() Cron {
				cf, _ := NewListFragment("1,5", []uint8{1, 5})
				cf.FragmentType = MINUTE

				return Cron{
					Data: []CronFragment{cf},
				}
			},
		},
		{
			"divisor",
			"*/30",
			func() Cron {
				cf, _ := NewDivisorFragment("*/30", []uint8{30})
				cf.FragmentType = MINUTE

				return Cron{
					Data: []CronFragment{cf},
				}
			},
		},
		{
			"single",
			"30",
			func() Cron {
				cf, _ := NewSingleFragment("30", []uint8{30})
				cf.FragmentType = MINUTE

				return Cron{
					Data: []CronFragment{cf},
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			p, _ := NewParser(
				WithInput(testCase.input, false),
			)
			
			out, err := p.Parse()
			expected := testCase.expected()

			if err != nil {
				t.Errorf("Got an unexpected error. Got: %s", err)
				return
			}

			if !expected.Eq(out) {
				t.Errorf("wanted: %v, got: %v", expected, out)
			}
		})
	}
}

func Test_ItHandlesComplexInputs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected func() Cron
	}{
		{
			"all wildcards",
			"* * * * *",
			func() Cron {
				return Cron{
					Data: []CronFragment{
						{
							Expr: "*",
							FragmentType: MINUTE,
							Kind: WILDCARD,
						},
						{
							Expr: "*",
							FragmentType: HOUR,
							Kind: WILDCARD,
						},
						{
							Expr: "*",
							FragmentType: DAY,
							Kind: WILDCARD,
						},
						{
							Expr: "*",
							FragmentType: MONTH,
							Kind: WILDCARD,
						},
						{
							Expr: "*",
							FragmentType: WEEKDAY,
							Kind: WILDCARD,
						},
					},
				}
			},
		},
		{
			"all ranges",
			"1-5 1-5 1-5 1-5 1-5",
			func() Cron {
				return Cron{
					Data: []CronFragment{
						{
							Expr: "1-5",
							FragmentType: MINUTE,
							Kind: RANGE,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1-5",
							FragmentType: HOUR,
							Kind: RANGE,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1-5",
							FragmentType: DAY,
							Kind: RANGE,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1-5",
							FragmentType: MONTH,
							Kind: RANGE,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1-5",
							FragmentType: WEEKDAY,
							Kind: RANGE,
							Factors: []uint8{1, 5},
						},
					},
				}
			},
		},
		{
			"all lists",
			"1,5 1,5 1,5 1,5 1,5",
			func() Cron {
				return Cron{
					Data: []CronFragment{
						{
							Expr: "1,5",
							FragmentType: MINUTE,
							Kind: LIST,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1,5",
							FragmentType: HOUR,
							Kind: LIST,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1,5",
							FragmentType: DAY,
							Kind: LIST,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1,5",
							FragmentType: MONTH,
							Kind: LIST,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1,5",
							FragmentType: WEEKDAY,
							Kind: LIST,
							Factors: []uint8{1, 5},
						},
					},
				}
			},
		},
		{
			"all divisors",
			"*/5 */5 */5 */5 */5",
			func() Cron {
				return Cron{
					Data: []CronFragment{
						{
							Expr: "*/5",
							FragmentType: MINUTE,
							Kind: DIVISOR,
							Factors: []uint8{5},
						},
						{
							Expr: "*/5",
							FragmentType: HOUR,
							Kind: DIVISOR,
							Factors: []uint8{5},
						},
						{
							Expr: "*/5",
							FragmentType: DAY,
							Kind: DIVISOR,
							Factors: []uint8{5},
						},
						{
							Expr: "*/5",
							FragmentType: MONTH,
							Kind: DIVISOR,
							Factors: []uint8{5},
						},
						{
							Expr: "*/5",
							FragmentType: WEEKDAY,
							Kind: DIVISOR,
							Factors: []uint8{5},
						},
					},
				}
			},
		},
		{
			"all singles",
			"5 5 5 5 5",
			func() Cron {
				return Cron{
					Data: []CronFragment{
						{
							Expr: "5",
							FragmentType: MINUTE,
							Kind: SINGLE,
							Factors: []uint8{5},
						},
						{
							Expr: "5",
							FragmentType: HOUR,
							Kind: SINGLE,
							Factors: []uint8{5},
						},
						{
							Expr: "5",
							FragmentType: DAY,
							Kind: SINGLE,
							Factors: []uint8{5},
						},
						{
							Expr: "5",
							FragmentType: MONTH,
							Kind: SINGLE,
							Factors: []uint8{5},
						},
						{
							Expr: "5",
							FragmentType: WEEKDAY,
							Kind: SINGLE,
							Factors: []uint8{5},
						},
					},
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			p, _ := NewParser(
				WithInput(testCase.input, true),
			)
			
			out, err := p.Parse()
			expected := testCase.expected()

			if err != nil {
				t.Errorf("Got an unexpected error. Got: %s", err)
				return
			}

			if !expected.Eq(out) {
				t.Errorf("wanted: %v, got: %v", expected, out)
			}
		})
	}
}
