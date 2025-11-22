package parser

import (
	"testing"

	tc "github.com/captainmango/coco-cron-parser/internal/cron"
)

func Test_ItHandlesWildCard(t *testing.T) {
	input := "* */23"
	p := NewParser(input)
	out, err := p.Parse()
	cf, _ := tc.NewWildCardFragment("*")

	expected := tc.Cron{
		cf,
		cf,
	}

	if err != nil {
		t.Errorf("Got an unexpected error. Got: %s", err)
		return
	}

	if !expected.Eq(out) {
		t.Errorf("wanted: %v, got: %v", expected, out)
	}
}