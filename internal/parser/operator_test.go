package parser

import (
	"slices"
	"testing"
)

func Test_WildCard(t *testing.T) {
	t.Parallel()
	cf := CronFragment{
		Expr:         "*",
		FragmentType: WEEKDAY,
		Kind:         WILDCARD,
	}

	expected := []uint8{1, 2, 3, 4, 5, 6, 7}

	nums, _ := wildcard(cf)

	if !slices.Equal(nums, expected) {
		t.Errorf("wanted %v, got %v", expected, nums)
	}
}

func Test_Range(t *testing.T) {
	t.Parallel()
	cf := CronFragment{
		Expr:         "1-5",
		FragmentType: WEEKDAY,
		Kind:         RANGE,
		Factors:      []uint8{5, 1},
	}

	expected := []uint8{1, 2, 3, 4, 5}

	nums, _ := rangeOp(cf)

	if !slices.Equal(nums, expected) {
		t.Errorf("wanted %v, got %v", expected, nums)
	}
}

func Test_List(t *testing.T) {
	t.Parallel()
	cf := CronFragment{
		Expr:         "1,5",
		FragmentType: WEEKDAY,
		Kind:         RANGE,
		Factors:      []uint8{5, 1},
	}

	expected := []uint8{1, 5}

	nums, _ := list(cf)

	if !slices.Equal(nums, expected) {
		t.Errorf("wanted %v, got %v", expected, nums)
	}
}

func Test_Divisor(t *testing.T) {
	t.Parallel()
	cf := CronFragment{
		Expr:         "*/5",
		FragmentType: WEEKDAY,
		Kind:         RANGE,
		Factors:      []uint8{2},
	}

	expected := []uint8{2, 4, 6}

	nums, _ := divisor(cf)

	if !slices.Equal(nums, expected) {
		t.Errorf("wanted %v, got %v", expected, nums)
	}
}

func Test_Single(t *testing.T) {
	t.Parallel()
	cf := CronFragment{
		Expr:         "5",
		FragmentType: WEEKDAY,
		Kind:         RANGE,
		Factors:      []uint8{2},
	}

	expected := []uint8{2}

	nums, _ := single(cf)

	if !slices.Equal(nums, expected) {
		t.Errorf("wanted %v, got %v", expected, nums)
	}
}

func Test_Validation(t *testing.T) {
	t.Parallel()
	cf := CronFragment{
		Expr:         "100",
		FragmentType: WEEKDAY,
		Kind:         RANGE,
		Factors:      []uint8{100},
	}

	expectedErr := ErrFactorsOutsideBounds(WEEKDAY, 100, 1, 7)
	_, err := single(cf)

	if err == nil {
		t.Error("error was expected")
		return
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error: %d, got error: %d", expectedErr, err)
		return
	}
}
