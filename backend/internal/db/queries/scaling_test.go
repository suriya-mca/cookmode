package queries

import "testing"

func TestShoppingListScaling(t *testing.T) {
	tests := []struct {
		name     string
		qty      float64
		base     int
		target   int
		expected float64
	}{
		{"double", 200, 2, 4, 400},
		{"half", 200, 4, 2, 100},
		{"same", 100, 2, 2, 100},
		{"zero base treated as 1", 100, 0, 4, 400},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			factor := 1.0
			base := tc.base
			if base == 0 {
				base = 1
			}
			factor = float64(tc.target) / float64(base)
			got := tc.qty * factor
			if got != tc.expected {
				t.Errorf("got %v want %v", got, tc.expected)
			}
		})
	}
}
