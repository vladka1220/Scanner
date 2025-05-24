package funding

import (
	"fmt"
	"testing"
	"time"
)

func TestFormatFunding(t *testing.T) {
	now := time.Now()
	next := now.Add(2*time.Hour + 30*time.Minute)
	got := FormatFunding(0.01, next.UnixMilli(), 100, 5, 8)
	// expected funding amount
	fundingAmount := 100 * 0.01 * 5
	expected := fmt.Sprintf("%.3f%% (через %dh %dm | %.6f USDT)", 1.0, 2, 30, fundingAmount)
	if got != expected {
		t.Fatalf("unexpected result: %s", got)
	}
}

func TestFormatNextFundingTime(t *testing.T) {
	now := time.Now()
	next := now.Add(3*time.Hour + 15*time.Minute)
	got := FormatNextFundingTime(next.UnixMilli())
	if got != "3h 15m" {
		t.Fatalf("unexpected result: %s", got)
	}
}
