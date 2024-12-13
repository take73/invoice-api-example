package validation_test

import (
	"testing"

	"github.com/take73/invoice-api-example/internal/shared/validation"
)

func TestValidRate(t *testing.T) {
	tests := []struct {
		name     string
		rate     float64
		expected bool
	}{
		{
			name:     "Rate within range (0.0)",
			rate:     0.0,
			expected: true,
		},
		{
			name:     "Rate within range (1.0)",
			rate:     1.0,
			expected: true,
		},
		{
			name:     "Rate within range (0.5)",
			rate:     0.5,
			expected: true,
		},
		{
			name:     "Rate below range (-0.1)",
			rate:     -0.1,
			expected: false,
		},
		{
			name:     "Rate above range (1.1)",
			rate:     1.1,
			expected: false,
		},
		{
			name:     "Rate at edge case (0.00000001)",
			rate:     0.00000001,
			expected: true,
		},
		{
			name:     "Rate at edge case (0.99999999)",
			rate:     0.99999999,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validation.ValidRate(tt.rate)
			if got != tt.expected {
				t.Errorf("ValidRate(%f) = %v; want %v", tt.rate, got, tt.expected)
			}
		})
	}
}
