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

				return d.Cron{
					cf,
				}
			},
		},
		{
			"range",
			"1-5",
			func() d.Cron {
				cf, _ := d.NewRangeFragment("1-5", []uint8{1,5})

				return d.Cron{
					cf,
				}
			},
		},
		{
			"list",
			"1,5",
			func() d.Cron {
				cf, _ := d.NewListFragment("1,5", []uint8{1,5})

				return d.Cron{
					cf,
				}
			},
		},
		{
			"divisor",
			"*/30",
			func() d.Cron {
				cf, _ := d.NewDivisorFragment("*/30", []uint8{30})

				return d.Cron{
					cf,
				}
			},
		},
		{
			"single",
			"30",
			func() d.Cron {
				cf, _ := d.NewSingleFragment("30", []uint8{30})

				return d.Cron{
					cf,
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			p := NewParser(testCase.input)
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
