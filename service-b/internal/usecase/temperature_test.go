package usecase

import (
	"testing"
)

func Test_celsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "zero celsius",
			celsius:  0,
			expected: 32,
		},
		{
			name:     "positive value",
			celsius:  25,
			expected: 77,
		},
		{
			name:     "negative value",
			celsius:  -10,
			expected: 14,
		},
		{
			name:     "decimal value",
			celsius:  37.5,
			expected: 99.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := celsiusToFahrenheit(tt.celsius)

			if result != tt.expected {
				t.Errorf("celsiusToFahrenheit(%v) = %v, want %v", tt.celsius, result, tt.expected)
			}
		})
	}
}

func Test_celsiusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "zero celsius",
			celsius:  0,
			expected: 273,
		},
		{
			name:     "positive value",
			celsius:  25,
			expected: 298,
		},
		{
			name:     "negative value",
			celsius:  -10,
			expected: 263,
		},
		{
			name:     "decimal value",
			celsius:  37.5,
			expected: 310.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := celsiusToKelvin(tt.celsius)

			if result != tt.expected {
				t.Errorf("celsiusToKelvin(%v) = %v, want %v", tt.celsius, result, tt.expected)
			}
		})
	}
}
