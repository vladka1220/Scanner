package exchanges

import "testing"

func TestParseFloat(t *testing.T) {
	if v := parseFloat("1.23"); v != 1.23 {
		t.Fatalf("expected 1.23, got %v", v)
	}
	if v := parseFloat("bad"); v != 0 {
		t.Fatalf("expected 0 for invalid input, got %v", v)
	}
}
