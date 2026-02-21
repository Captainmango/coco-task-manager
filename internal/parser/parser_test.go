package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Fixtures - Helper functions to create expected CronFragments
func makeWildCardFragment(fragmentType CronFragmentType) CronFragment {
	cf, _ := NewWildCardFragment("*")
	cf.FragmentType = fragmentType
	return cf
}

func makeRangeFragment(expr string, factors []uint8, fragmentType CronFragmentType) CronFragment {
	cf, _ := NewRangeFragment(expr, factors)
	cf.FragmentType = fragmentType
	return cf
}

func makeListFragment(expr string, factors []uint8, fragmentType CronFragmentType) CronFragment {
	cf, _ := NewListFragment(expr, factors)
	cf.FragmentType = fragmentType
	return cf
}

func makeDivisorFragment(expr string, factors []uint8, fragmentType CronFragmentType) CronFragment {
	cf, _ := NewDivisorFragment(expr, factors)
	cf.FragmentType = fragmentType
	return cf
}

func makeSingleFragment(expr string, factors []uint8, fragmentType CronFragmentType) CronFragment {
	cf, _ := NewSingleFragment(expr, factors)
	cf.FragmentType = fragmentType
	return cf
}

// Basic Fragment Parsing Tests
func TestParser_BasicFragments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Cron
	}{
		{
			name:  "wildcard",
			input: "*",
			expected: Cron{
				Data: []CronFragment{makeWildCardFragment(MINUTE)},
			},
		},
		{
			name:  "range",
			input: "1-5",
			expected: Cron{
				Data: []CronFragment{makeRangeFragment("1-5", []uint8{1, 5}, MINUTE)},
			},
		},
		{
			name:  "list_two_items",
			input: "1,5",
			expected: Cron{
				Data: []CronFragment{makeListFragment("1,5", []uint8{1, 5}, MINUTE)},
			},
		},
		{
			name:  "list_three_items",
			input: "1,5,10",
			expected: Cron{
				Data: []CronFragment{makeListFragment("1,5,10", []uint8{1, 5, 10}, MINUTE)},
			},
		},
		{
			name:  "divisor",
			input: "*/30",
			expected: Cron{
				Data: []CronFragment{makeDivisorFragment("*/30", []uint8{30}, MINUTE)},
			},
		},
		{
			name:  "single",
			input: "30",
			expected: Cron{
				Data: []CronFragment{makeSingleFragment("30", []uint8{30}, MINUTE)},
			},
		},
		{
			name:  "single_zero",
			input: "0",
			expected: Cron{
				Data: []CronFragment{makeSingleFragment("0", []uint8{0}, MINUTE)},
			},
		},
		{
			name:  "large_number",
			input: "59",
			expected: Cron{
				Data: []CronFragment{makeSingleFragment("59", []uint8{59}, MINUTE)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p, err := NewParser(WithInput(tt.input, false))
			assert.NoError(t, err, "creating parser")

			out, err := p.Parse()
			assert.NoError(t, err, "parsing")
			assert.True(t, tt.expected.Eq(out), "expected %v, got %v", tt.expected, out)
		})
	}
}

// Full Expression Parsing Tests
func TestParser_FullExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Cron
	}{
		{
			name:  "all_wildcards",
			input: "* * * * *",
			expected: Cron{
				Data: []CronFragment{
					makeWildCardFragment(MINUTE),
					makeWildCardFragment(HOUR),
					makeWildCardFragment(DAY),
					makeWildCardFragment(MONTH),
					makeWildCardFragment(WEEKDAY),
				},
			},
		},
		{
			name:  "all_ranges",
			input: "1-5 1-5 1-5 1-5 1-5",
			expected: Cron{
				Data: []CronFragment{
					makeRangeFragment("1-5", []uint8{1, 5}, MINUTE),
					makeRangeFragment("1-5", []uint8{1, 5}, HOUR),
					makeRangeFragment("1-5", []uint8{1, 5}, DAY),
					makeRangeFragment("1-5", []uint8{1, 5}, MONTH),
					makeRangeFragment("1-5", []uint8{1, 5}, WEEKDAY),
				},
			},
		},
		{
			name:  "all_lists",
			input: "1,5 1,5 1,5 1,5 1,5",
			expected: Cron{
				Data: []CronFragment{
					makeListFragment("1,5", []uint8{1, 5}, MINUTE),
					makeListFragment("1,5", []uint8{1, 5}, HOUR),
					makeListFragment("1,5", []uint8{1, 5}, DAY),
					makeListFragment("1,5", []uint8{1, 5}, MONTH),
					makeListFragment("1,5", []uint8{1, 5}, WEEKDAY),
				},
			},
		},
		{
			name:  "all_divisors",
			input: "*/5 */5 */5 */5 */5",
			expected: Cron{
				Data: []CronFragment{
					makeDivisorFragment("*/5", []uint8{5}, MINUTE),
					makeDivisorFragment("*/5", []uint8{5}, HOUR),
					makeDivisorFragment("*/5", []uint8{5}, DAY),
					makeDivisorFragment("*/5", []uint8{5}, MONTH),
					makeDivisorFragment("*/5", []uint8{5}, WEEKDAY),
				},
			},
		},
		{
			name:  "all_singles",
			input: "5 5 5 5 5",
			expected: Cron{
				Data: []CronFragment{
					makeSingleFragment("5", []uint8{5}, MINUTE),
					makeSingleFragment("5", []uint8{5}, HOUR),
					makeSingleFragment("5", []uint8{5}, DAY),
					makeSingleFragment("5", []uint8{5}, MONTH),
					makeSingleFragment("5", []uint8{5}, WEEKDAY),
				},
			},
		},
		{
			name:  "mixed_types",
			input: "*/15 0-23 1,15 * 1-5",
			expected: Cron{
				Data: []CronFragment{
					makeDivisorFragment("*/15", []uint8{15}, MINUTE),
					makeRangeFragment("0-23", []uint8{0, 23}, HOUR),
					makeListFragment("1,15", []uint8{1, 15}, DAY),
					makeWildCardFragment(MONTH),
					makeRangeFragment("1-5", []uint8{1, 5}, WEEKDAY),
				},
			},
		},
		{
			name:  "every_minute_every_day",
			input: "* * * * *",
			expected: Cron{
				Data: []CronFragment{
					makeWildCardFragment(MINUTE),
					makeWildCardFragment(HOUR),
					makeWildCardFragment(DAY),
					makeWildCardFragment(MONTH),
					makeWildCardFragment(WEEKDAY),
				},
			},
		},
		{
			name:  "specific_time",
			input: "30 14 15 3 *",
			expected: Cron{
				Data: []CronFragment{
					makeSingleFragment("30", []uint8{30}, MINUTE),
					makeSingleFragment("14", []uint8{14}, HOUR),
					makeSingleFragment("15", []uint8{15}, DAY),
					makeSingleFragment("3", []uint8{3}, MONTH),
					makeWildCardFragment(WEEKDAY),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p, err := NewParser(WithInput(tt.input, true))
			assert.NoError(t, err, "creating parser")

			out, err := p.Parse()
			assert.NoError(t, err, "parsing")
			assert.True(t, tt.expected.Eq(out), "expected %v, got %v", tt.expected, out)
		})
	}
}

// Error Cases
func TestParser_ErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		validateLen   bool
		expectedError string
	}{
		{
			name:          "too_few_parts",
			input:         "* * * *",
			validateLen:   true,
			expectedError: "is not a valid input",
		},
		{
			name:          "too_many_parts",
			input:         "* * * * * *",
			validateLen:   true,
			expectedError: "is not a valid input",
		},
		{
			name:          "double_spaces_caught_as_invalid_input",
			input:         "*  * * * *",
			validateLen:   true,
			expectedError: "is not a valid input",
		},
		{
			name:          "invalid_character",
			input:         "a * * * *",
			validateLen:   true,
			expectedError: "malformed cron expression",
		},
		{
			name:          "invalid_after_asterisk",
			input:         "*a * * * *",
			validateLen:   false,
			expectedError: "malformed cron expression",
		},
		{
			name:          "empty_input",
			input:         "",
			validateLen:   true,
			expectedError: "is not a valid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p, err := NewParser(WithInput(tt.input, tt.validateLen))
			if err != nil {
				// Error during parser creation
				assert.Contains(t, err.Error(), tt.expectedError, "parser creation error")
				return
			}

			_, err = p.Parse()
			assert.Error(t, err, "parsing should error")
			assert.Contains(t, err.Error(), tt.expectedError, "parse error")
		})
	}
}

// Parser Creation Tests
func TestParser_NewParser(t *testing.T) {
	t.Run("valid_string_input", func(t *testing.T) {
		t.Parallel()
		p, err := NewParser(WithInput("* * * * *", true))
		assert.NoError(t, err, "creating parser with string input")
		assert.NotNil(t, p, "expected parser to be created")
	})

	t.Run("valid_slice_input", func(t *testing.T) {
		t.Parallel()
		p, err := NewParser(WithInput([]string{"*", "*", "*", "*", "*"}, true))
		assert.NoError(t, err, "creating parser with slice input")
		assert.NotNil(t, p, "expected parser to be created")
	})

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p, err := NewParser()
		assert.NoError(t, err, "creating parser without options")
		assert.NotNil(t, p, "expected parser to be created")
	})
}

// Edge Cases
func TestParser_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Cron
	}{
		{
			name:  "range_with_zero",
			input: "0-59",
			expected: Cron{
				Data: []CronFragment{
					makeRangeFragment("0-59", []uint8{0, 59}, MINUTE),
				},
			},
		},
		{
			name:  "large_list",
			input: "0,15,30,45",
			expected: Cron{
				Data: []CronFragment{
					makeListFragment("0,15,30,45", []uint8{0, 15, 30, 45}, MINUTE),
				},
			},
		},
		{
			name:  "divisor_by_one",
			input: "*/1",
			expected: Cron{
				Data: []CronFragment{
					makeDivisorFragment("*/1", []uint8{1}, MINUTE),
				},
			},
		},
		{
			name:  "single_zero_minute",
			input: "0 0 1 1 1",
			expected: Cron{
				Data: []CronFragment{
					makeSingleFragment("0", []uint8{0}, MINUTE),
					makeSingleFragment("0", []uint8{0}, HOUR),
					makeSingleFragment("1", []uint8{1}, DAY),
					makeSingleFragment("1", []uint8{1}, MONTH),
					makeSingleFragment("1", []uint8{1}, WEEKDAY),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p, err := NewParser(WithInput(tt.input, false))
			assert.NoError(t, err, "creating parser")

			out, err := p.Parse()
			assert.NoError(t, err, "parsing")
			assert.True(t, tt.expected.Eq(out), "expected %v, got %v", tt.expected, out)
		})
	}
}

// String Array Input Tests
func TestParser_StringArrayInput(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		validate bool
		wantErr  bool
	}{
		{
			name:     "valid_array",
			input:    []string{"*", "*", "*", "*", "*"},
			validate: true,
			wantErr:  false,
		},
		{
			name:     "too_short_array",
			input:    []string{"*", "*", "*", "*"},
			validate: true,
			wantErr:  true,
		},
		{
			name:     "too_long_array",
			input:    []string{"*", "*", "*", "*", "*", "*"},
			validate: true,
			wantErr:  true,
		},
		{
			name:     "no_validation_short",
			input:    []string{"*"},
			validate: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p, err := NewParser(WithInput(tt.input, tt.validate))
			if tt.wantErr {
				assert.Error(t, err, "creating parser should error")
				return
			}
			assert.NoError(t, err, "creating parser")
			assert.NotNil(t, p, "expected parser to be created")
		})
	}
}
