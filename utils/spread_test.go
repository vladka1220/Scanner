package utils

import (
    "math"
    "testing"
)

func TestCalculateSpotSpread(t *testing.T) {
    if v := CalculateSpotSpread(100, 105); v != 5 {
        t.Errorf("expected 5, got %v", v)
    }
    if v := CalculateSpotSpread(0, 100); v != 0 {
        t.Errorf("expected 0 for invalid input, got %v", v)
    }
}

func TestCalculateFuturesSpread(t *testing.T) {
    if v := CalculateFuturesSpread(200, 210); v != 5 {
        t.Errorf("expected 5, got %v", v)
    }
    if v := CalculateFuturesSpread(-1, 0); v != 0 {
        t.Errorf("expected 0 for invalid input, got %v", v)
    }
}

func TestNetSpread(t *testing.T) {
    // spread 2.5%, maker 0.02%, taker 0.07%, funding 0.01%
    expected := 2.5 - 0.02 - 0.07 - 0.01
    if v := NetSpread(2.5, 0.02, 0.07, 0.01); math.Abs(v-expected) > 1e-9 {
        t.Errorf("expected %v, got %v", expected, v)
    }
}

