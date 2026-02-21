package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Wildcard Tests
func TestWildcardOperator(t *testing.T) {
	tests := []struct {
		name         string
		fragmentType CronFragmentType
		expected     []uint8
	}{
		{
			name:         "minute_bounds",
			fragmentType: MINUTE,
			expected:     []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59},
		},
		{
			name:         "hour_bounds",
			fragmentType: HOUR,
			expected:     []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
		},
		{
			name:         "day_bounds",
			fragmentType: DAY,
			expected:     []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
		},
		{
			name:         "month_bounds",
			fragmentType: MONTH,
			expected:     []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
		{
			name:         "weekday_bounds",
			fragmentType: WEEKDAY,
			expected:     []uint8{1, 2, 3, 4, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cf := CronFragment{
				Expr:         "*",
				FragmentType: tt.fragmentType,
				Kind:         WILDCARD,
			}

			nums, err := wildcard(cf)
			assert.NoError(t, err, "wildcard operation")
			assert.Equal(t, tt.expected, nums, "wildcard result")
		})
	}
}

// Range Tests
func TestRangeOperator(t *testing.T) {
	tests := []struct {
		name         string
		expr         string
		factors      []uint8
		fragmentType CronFragmentType
		expected     []uint8
		wantErr      bool
	}{
		{
			name:         "basic_range",
			expr:         "1-5",
			factors:      []uint8{1, 5},
			fragmentType: WEEKDAY,
			expected:     []uint8{1, 2, 3, 4, 5},
			wantErr:      false,
		},
		{
			name:         "unordered_factors",
			expr:         "5-1",
			factors:      []uint8{5, 1},
			fragmentType: WEEKDAY,
			expected:     []uint8{1, 2, 3, 4, 5},
			wantErr:      false,
		},
		{
			name:         "full_minute_range",
			expr:         "0-59",
			factors:      []uint8{0, 59},
			fragmentType: MINUTE,
			expected:     []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59},
			wantErr:      false,
		},
		{
			name:         "single_value_range",
			expr:         "5-5",
			factors:      []uint8{5, 5},
			fragmentType: HOUR,
			expected:     []uint8{5},
			wantErr:      false,
		},
		{
			name:         "two_factor_range",
			expr:         "1-5",
			factors:      []uint8{1, 5},
			fragmentType: DAY,
			expected:     []uint8{1, 2, 3, 4, 5},
			wantErr:      false,
		},
		{
			name:         "too_many_factors",
			expr:         "1-5-10",
			factors:      []uint8{1, 5, 10},
			fragmentType: WEEKDAY,
			expected:     nil,
			wantErr:      true,
		},
		{
			name:         "too_few_factors",
			expr:         "5",
			factors:      []uint8{5},
			fragmentType: WEEKDAY,
			expected:     nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cf := CronFragment{
				Expr:         tt.expr,
				FragmentType: tt.fragmentType,
				Kind:         RANGE,
				Factors:      tt.factors,
			}

			nums, err := rangeOp(cf)
			if tt.wantErr {
				assert.Error(t, err, "range operation should error")
				return
			}
			assert.NoError(t, err, "range operation")
			assert.Equal(t, tt.expected, nums, "range result")
		})
	}
}

// List Tests
func TestListOperator(t *testing.T) {
	tests := []struct {
		name         string
		expr         string
		factors      []uint8
		fragmentType CronFragmentType
		expected     []uint8
		wantErr      bool
	}{
		{
			name:         "two_items",
			expr:         "1,5",
			factors:      []uint8{1, 5},
			fragmentType: WEEKDAY,
			expected:     []uint8{1, 5},
			wantErr:      false,
		},
		{
			name:         "three_items",
			expr:         "1,5,7",
			factors:      []uint8{1, 5, 7},
			fragmentType: WEEKDAY,
			expected:     []uint8{1, 5, 7},
			wantErr:      false,
		},
		{
			name:         "unordered_sorts",
			expr:         "5,1,7",
			factors:      []uint8{5, 1, 7},
			fragmentType: WEEKDAY,
			expected:     []uint8{1, 5, 7},
			wantErr:      false,
		},
		{
			name:         "duplicates_preserved",
			expr:         "1,1,5",
			factors:      []uint8{1, 1, 5},
			fragmentType: WEEKDAY,
			expected:     []uint8{1, 1, 5},
			wantErr:      false,
		},
		{
			name:         "single_item",
			expr:         "5",
			factors:      []uint8{5},
			fragmentType: HOUR,
			expected:     []uint8{5},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cf := CronFragment{
				Expr:         tt.expr,
				FragmentType: tt.fragmentType,
				Kind:         LIST,
				Factors:      tt.factors,
			}

			nums, err := list(cf)
			if tt.wantErr {
				assert.Error(t, err, "list operation should error")
				return
			}
			assert.NoError(t, err, "list operation")
			assert.Equal(t, tt.expected, nums, "list result")
		})
	}
}

// Divisor Tests
func TestDivisorOperator(t *testing.T) {
	tests := []struct {
		name         string
		expr         string
		factors      []uint8
		fragmentType CronFragmentType
		expected     []uint8
		wantErr      bool
	}{
		{
			name:         "divisor_of_2_weekday",
			expr:         "*/2",
			factors:      []uint8{2},
			fragmentType: WEEKDAY,
			expected:     []uint8{2, 4, 6},
			wantErr:      false,
		},
		{
			name:         "divisor_of_5_minute",
			expr:         "*/5",
			factors:      []uint8{5},
			fragmentType: MINUTE,
			expected:     []uint8{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55},
			wantErr:      false,
		},
		{
			name:         "divisor_of_15_minute",
			expr:         "*/15",
			factors:      []uint8{15},
			fragmentType: MINUTE,
			expected:     []uint8{0, 15, 30, 45},
			wantErr:      false,
		},
		{
			name:         "divisor_of_1_minute",
			expr:         "*/1",
			factors:      []uint8{1},
			fragmentType: MINUTE,
			expected:     []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59},
			wantErr:      false,
		},
		{
			name:         "divisor_of_12_hour",
			expr:         "*/12",
			factors:      []uint8{12},
			fragmentType: HOUR,
			expected:     []uint8{0, 12},
			wantErr:      false,
		},
		{
			name:         "divisor_of_7_day",
			expr:         "*/7",
			factors:      []uint8{7},
			fragmentType: DAY,
			expected:     []uint8{7, 14, 21, 28},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cf := CronFragment{
				Expr:         tt.expr,
				FragmentType: tt.fragmentType,
				Kind:         DIVISOR,
				Factors:      tt.factors,
			}

			nums, err := divisor(cf)
			if tt.wantErr {
				assert.Error(t, err, "divisor operation should error")
				return
			}
			assert.NoError(t, err, "divisor operation")
			assert.Equal(t, tt.expected, nums, "divisor result")
		})
	}
}

// Single Tests
func TestSingleOperator(t *testing.T) {
	tests := []struct {
		name         string
		expr         string
		factors      []uint8
		fragmentType CronFragmentType
		expected     []uint8
		wantErr      bool
	}{
		{
			name:         "single_value",
			expr:         "5",
			factors:      []uint8{5},
			fragmentType: WEEKDAY,
			expected:     []uint8{5},
			wantErr:      false,
		},
		{
			name:         "zero_value",
			expr:         "0",
			factors:      []uint8{0},
			fragmentType: MINUTE,
			expected:     []uint8{0},
			wantErr:      false,
		},
		{
			name:         "max_minute",
			expr:         "59",
			factors:      []uint8{59},
			fragmentType: MINUTE,
			expected:     []uint8{59},
			wantErr:      false,
		},
		{
			name:         "max_hour",
			expr:         "23",
			factors:      []uint8{23},
			fragmentType: HOUR,
			expected:     []uint8{23},
			wantErr:      false,
		},
		{
			name:         "out_of_bounds_errors",
			expr:         "100",
			factors:      []uint8{100},
			fragmentType: WEEKDAY,
			expected:     nil,
			wantErr:      true,
		},
		{
			name:         "day_out_of_lower_bounds_errors",
			expr:         "0",
			factors:      []uint8{0},
			fragmentType: DAY,
			expected:     nil,
			wantErr:      true,
		},
		{
			name:         "day_out_of_upper_bounds_errors",
			expr:         "32",
			factors:      []uint8{32},
			fragmentType: DAY,
			expected:     nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cf := CronFragment{
				Expr:         tt.expr,
				FragmentType: tt.fragmentType,
				Kind:         SINGLE,
				Factors:      tt.factors,
			}

			nums, err := single(cf)
			if tt.wantErr {
				assert.Error(t, err, "single operation should error")
				return
			}
			assert.NoError(t, err, "single operation")
			assert.Equal(t, tt.expected, nums, "single result")
		})
	}
}

// GetPossibleValues Tests
func TestGetPossibleValues(t *testing.T) {
	tests := []struct {
		name         string
		kind         OperatorType
		factors      []uint8
		fragmentType CronFragmentType
		expected     []uint8
		wantErr      bool
	}{
		{
			name:         "wildcard_minute",
			kind:         WILDCARD,
			factors:      nil,
			fragmentType: MINUTE,
			expected:     []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59},
			wantErr:      false,
		},
		{
			name:         "single_value",
			kind:         SINGLE,
			factors:      []uint8{30},
			fragmentType: MINUTE,
			expected:     []uint8{30},
			wantErr:      false,
		},
		{
			name:         "range_values",
			kind:         RANGE,
			factors:      []uint8{1, 5},
			fragmentType: WEEKDAY,
			expected:     []uint8{1, 2, 3, 4, 5},
			wantErr:      false,
		},
		{
			name:         "list_values",
			kind:         LIST,
			factors:      []uint8{1, 15, 30},
			fragmentType: MINUTE,
			expected:     []uint8{1, 15, 30},
			wantErr:      false,
		},
		{
			name:         "divisor_values",
			kind:         DIVISOR,
			factors:      []uint8{15},
			fragmentType: MINUTE,
			expected:     []uint8{0, 15, 30, 45},
			wantErr:      false,
		},
		{
			name:         "invalid_kind_errors",
			kind:         OperatorType("INVALID"),
			factors:      []uint8{5},
			fragmentType: MINUTE,
			expected:     nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cf := CronFragment{
				Expr:         "test",
				FragmentType: tt.fragmentType,
				Kind:         tt.kind,
				Factors:      tt.factors,
			}

			nums, err := cf.GetPossibleValues()
			if tt.wantErr {
				assert.Error(t, err, "GetPossibleValues should error")
				return
			}
			assert.NoError(t, err, "GetPossibleValues")
			assert.Equal(t, tt.expected, nums, "GetPossibleValues result")
		})
	}
}

// Validation Tests
func TestValidation(t *testing.T) {
	tests := []struct {
		name         string
		fragmentType CronFragmentType
		factors      []uint8
		wantErr      bool
		errContains  string
	}{
		{
			name:         "minute_valid_upper",
			fragmentType: MINUTE,
			factors:      []uint8{59},
			wantErr:      false,
		},
		{
			name:         "minute_valid_zero",
			fragmentType: MINUTE,
			factors:      []uint8{0},
			wantErr:      false,
		},
		{
			name:         "minute_too_high",
			fragmentType: MINUTE,
			factors:      []uint8{60},
			wantErr:      true,
			errContains:  "not within",
		},
		{
			name:         "hour_valid_upper",
			fragmentType: HOUR,
			factors:      []uint8{23},
			wantErr:      false,
		},
		{
			name:         "hour_valid_zero",
			fragmentType: HOUR,
			factors:      []uint8{0},
			wantErr:      false,
		},
		{
			name:         "hour_too_high",
			fragmentType: HOUR,
			factors:      []uint8{24},
			wantErr:      true,
			errContains:  "not within",
		},
		{
			name:         "day_valid_upper",
			fragmentType: DAY,
			factors:      []uint8{31},
			wantErr:      false,
		},
		{
			name:         "day_valid_lower",
			fragmentType: DAY,
			factors:      []uint8{1},
			wantErr:      false,
		},
		{
			name:         "day_too_high",
			fragmentType: DAY,
			factors:      []uint8{32},
			wantErr:      true,
			errContains:  "not within",
		},
		{
			name:         "day_zero_invalid",
			fragmentType: DAY,
			factors:      []uint8{0},
			wantErr:      true,
			errContains:  "not within",
		},
		{
			name:         "month_valid_upper",
			fragmentType: MONTH,
			factors:      []uint8{12},
			wantErr:      false,
		},
		{
			name:         "month_valid_lower",
			fragmentType: MONTH,
			factors:      []uint8{1},
			wantErr:      false,
		},
		{
			name:         "month_too_high",
			fragmentType: MONTH,
			factors:      []uint8{13},
			wantErr:      true,
			errContains:  "not within",
		},
		{
			name:         "month_zero_invalid",
			fragmentType: MONTH,
			factors:      []uint8{0},
			wantErr:      true,
			errContains:  "not within",
		},
		{
			name:         "weekday_valid_upper",
			fragmentType: WEEKDAY,
			factors:      []uint8{7},
			wantErr:      false,
		},
		{
			name:         "weekday_valid_lower",
			fragmentType: WEEKDAY,
			factors:      []uint8{1},
			wantErr:      false,
		},
		{
			name:         "weekday_too_high",
			fragmentType: WEEKDAY,
			factors:      []uint8{8},
			wantErr:      true,
			errContains:  "not within",
		},
		{
			name:         "weekday_zero_invalid",
			fragmentType: WEEKDAY,
			factors:      []uint8{0},
			wantErr:      true,
			errContains:  "not within",
		},
		{
			name:         "unknown_fragment_type",
			fragmentType: CronFragmentType("UNKNOWN"),
			factors:      []uint8{5},
			wantErr:      true,
			errContains:  "unknown bounds type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cf := CronFragment{
				Expr:         "test",
				FragmentType: tt.fragmentType,
				Kind:         SINGLE,
				Factors:      tt.factors,
			}

			err := cf.validate()
			if tt.wantErr {
				assert.Error(t, err, "validation should error")
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			assert.NoError(t, err, "validation")
		})
	}
}

// GetBounds Tests
func TestGetBounds(t *testing.T) {
	tests := []struct {
		name         string
		fragmentType CronFragmentType
		wantUpper    uint8
		wantLower    uint8
		wantErr      bool
	}{
		{
			name:         "minute_bounds",
			fragmentType: MINUTE,
			wantUpper:    59,
			wantLower:    0,
			wantErr:      false,
		},
		{
			name:         "hour_bounds",
			fragmentType: HOUR,
			wantUpper:    23,
			wantLower:    0,
			wantErr:      false,
		},
		{
			name:         "day_bounds",
			fragmentType: DAY,
			wantUpper:    31,
			wantLower:    1,
			wantErr:      false,
		},
		{
			name:         "month_bounds",
			fragmentType: MONTH,
			wantUpper:    12,
			wantLower:    1,
			wantErr:      false,
		},
		{
			name:         "weekday_bounds",
			fragmentType: WEEKDAY,
			wantUpper:    7,
			wantLower:    1,
			wantErr:      false,
		},
		{
			name:         "unknown_type_errors",
			fragmentType: CronFragmentType("UNKNOWN"),
			wantUpper:    0,
			wantLower:    0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bounds, err := getBounds(tt.fragmentType)
			if tt.wantErr {
				assert.Error(t, err, "getBounds should error")
				return
			}
			assert.NoError(t, err, "getBounds")
			assert.Equal(t, tt.wantUpper, bounds.upper, "upper bound")
			assert.Equal(t, tt.wantLower, bounds.lower, "lower bound")
		})
	}
}
