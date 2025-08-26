// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exstr

import (
	"fmt"
	"math"
	"testing"
)

func TestInt64Float64(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
		value    float64
	}{
		{31457280, "MB", 30.00},
		{2684354560, "GB", 2.50},
		{4398046511104, "TB", 4.00},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Input: %d", test.input), func(t *testing.T) {
			unit, result := Int64Float64(test.input)
			if unit != test.expected {
				t.Errorf("Expected unit %s, but got %s", test.expected, unit)
			}
			if math.Abs(result-test.value) > 0.001 {
				t.Errorf("Expected value %.2f, but got %.2f", test.value, result)
			}
		})
	}
}
