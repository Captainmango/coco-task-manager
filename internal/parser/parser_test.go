package parser

import (
	"testing"

	tc "github.com/captainmango/coco-cron-parser/internal/cron"
)

func Test_ItHandlesBasicInputs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected func() tc.Cron
	}{
		{
			"wildcard",
			"*",
			func() tc.Cron {
				cf, _ := tc.NewWildCardFragment("*")

				return tc.Cron{
					cf,
				}
			},
		},
		{
			"range",
			"1-5",
			func() tc.Cron {
				cf, _ := tc.NewRangeFragment("1-5", []uint8{1,5})

				return tc.Cron{
					cf,
				}
			},
		},
		{
			"list",
			"1,5",
			func() tc.Cron {
				cf, _ := tc.NewListFragment("1,5", []uint8{1,5})

				return tc.Cron{
					cf,
				}
			},
		},
		{
			"divisor",
			"*/30",
			func() tc.Cron {
				cf, _ := tc.NewDivisorFragment("*/30", []uint8{30})

				return tc.Cron{
					cf,
				}
			},
		},
		{
			"single",
			"30",
			func() tc.Cron {
				cf, _ := tc.NewSingleFragment("30", []uint8{30})

				return tc.Cron{
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
