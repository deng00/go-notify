package pagerduty

import "testing"

func TestPagerduty(t *testing.T) {
	opt := &Options{Token: "e6f5ab6a0d434204c05be34562fb1bea", Text: "K TEST"}
	c := New(*opt)
	c.Send("K TEST")
}
