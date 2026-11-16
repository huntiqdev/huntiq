package service

import (
	"math"
	"testing"
)

func TestRoundCents(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{input: 10.555, expected: 10.56},
		{input: 10.554, expected: 10.55},
		{input: 1.000, expected: 1.00},
		{input: 0.05, expected: 0.05},
		{input: 12345.6789, expected: 12345.68},
	}

	for _, tt := range tests {
		got := roundCents(tt.input)
		if math.Abs(got-tt.expected) > 0.001 {
			t.Errorf("roundCents(%f) = %f, want %f", tt.input, got, tt.expected)
		}
	}
}
