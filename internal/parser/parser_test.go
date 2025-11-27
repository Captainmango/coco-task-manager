package parser

import (
	"testing"

	d "github.com/captainmango/coco-cron-parser/internal/data"
)

func Test_ItHandlesBasicInputs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected func() d.Cron
	}{
		{
			"wildcard",
			"*",
			func() d.Cron {
				cf, _ := d.NewWildCardFragment("*")
				cf.FragmentType = d.MINUTE

				return d.Cron{
					Data: []d.CronFragment{cf},
				}
			},
		},
		{
			"range",
			"1-5",
			func() d.Cron {
				cf, _ := d.NewRangeFragment("1-5", []uint8{1, 5})
				cf.FragmentType = d.MINUTE

				return d.Cron{
					Data: []d.CronFragment{cf},
				}
			},
		},
		{
			"list",
			"1,5",
			func() d.Cron {
				cf, _ := d.NewListFragment("1,5", []uint8{1, 5})
				cf.FragmentType = d.MINUTE

				return d.Cron{
					Data: []d.CronFragment{cf},
				}
			},
		},
		{
			"divisor",
			"*/30",
			func() d.Cron {
				cf, _ := d.NewDivisorFragment("*/30", []uint8{30})
				cf.FragmentType = d.MINUTE

				return d.Cron{
					Data: []d.CronFragment{cf},
				}
			},
		},
		{
			"single",
			"30",
			func() d.Cron {
				cf, _ := d.NewSingleFragment("30", []uint8{30})
				cf.FragmentType = d.MINUTE

				return d.Cron{
					Data: []d.CronFragment{cf},
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
		expected func() d.Cron
	}{
		{
			"all wildcards",
			"* * * * *",
			func() d.Cron {
				return d.Cron{
					Data: []d.CronFragment{
						{
							Expr: "*",
							FragmentType: d.MINUTE,
							Kind: d.WILDCARD,
						},
						{
							Expr: "*",
							FragmentType: d.HOUR,
							Kind: d.WILDCARD,
						},
						{
							Expr: "*",
							FragmentType: d.DAY,
							Kind: d.WILDCARD,
						},
						{
							Expr: "*",
							FragmentType: d.MONTH,
							Kind: d.WILDCARD,
						},
						{
							Expr: "*",
							FragmentType: d.WEEKDAY,
							Kind: d.WILDCARD,
						},
					},
				}
			},
		},
		{
			"all ranges",
			"1-5 1-5 1-5 1-5 1-5",
			func() d.Cron {
				return d.Cron{
					Data: []d.CronFragment{
						{
							Expr: "1-5",
							FragmentType: d.MINUTE,
							Kind: d.RANGE,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1-5",
							FragmentType: d.HOUR,
							Kind: d.RANGE,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1-5",
							FragmentType: d.DAY,
							Kind: d.RANGE,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1-5",
							FragmentType: d.MONTH,
							Kind: d.RANGE,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1-5",
							FragmentType: d.WEEKDAY,
							Kind: d.RANGE,
							Factors: []uint8{1, 5},
						},
					},
				}
			},
		},
		{
			"all lists",
			"1,5 1,5 1,5 1,5 1,5",
			func() d.Cron {
				return d.Cron{
					Data: []d.CronFragment{
						{
							Expr: "1,5",
							FragmentType: d.MINUTE,
							Kind: d.LIST,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1,5",
							FragmentType: d.HOUR,
							Kind: d.LIST,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1,5",
							FragmentType: d.DAY,
							Kind: d.LIST,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1,5",
							FragmentType: d.MONTH,
							Kind: d.LIST,
							Factors: []uint8{1, 5},
						},
						{
							Expr: "1,5",
							FragmentType: d.WEEKDAY,
							Kind: d.LIST,
							Factors: []uint8{1, 5},
						},
					},
				}
			},
		},
		{
			"all divisors",
			"*/5 */5 */5 */5 */5",
			func() d.Cron {
				return d.Cron{
					Data: []d.CronFragment{
						{
							Expr: "*/5",
							FragmentType: d.MINUTE,
							Kind: d.DIVISOR,
							Factors: []uint8{5},
						},
						{
							Expr: "*/5",
							FragmentType: d.HOUR,
							Kind: d.DIVISOR,
							Factors: []uint8{5},
						},
						{
							Expr: "*/5",
							FragmentType: d.DAY,
							Kind: d.DIVISOR,
							Factors: []uint8{5},
						},
						{
							Expr: "*/5",
							FragmentType: d.MONTH,
							Kind: d.DIVISOR,
							Factors: []uint8{5},
						},
						{
							Expr: "*/5",
							FragmentType: d.WEEKDAY,
							Kind: d.DIVISOR,
							Factors: []uint8{5},
						},
					},
				}
			},
		},
		{
			"all singles",
			"5 5 5 5 5",
			func() d.Cron {
				return d.Cron{
					Data: []d.CronFragment{
						{
							Expr: "5",
							FragmentType: d.MINUTE,
							Kind: d.SINGLE,
							Factors: []uint8{5},
						},
						{
							Expr: "5",
							FragmentType: d.HOUR,
							Kind: d.SINGLE,
							Factors: []uint8{5},
						},
						{
							Expr: "5",
							FragmentType: d.DAY,
							Kind: d.SINGLE,
							Factors: []uint8{5},
						},
						{
							Expr: "5",
							FragmentType: d.MONTH,
							Kind: d.SINGLE,
							Factors: []uint8{5},
						},
						{
							Expr: "5",
							FragmentType: d.WEEKDAY,
							Kind: d.SINGLE,
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
