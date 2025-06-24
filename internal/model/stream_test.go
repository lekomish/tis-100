package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lekomish/tis-100/internal/model"
)

/* TESTS */

// --- Len ---
func TestStreamLen(t *testing.T) {
	tests := []struct {
		name     string
		stream   model.Stream
		expected int
	}{
		{
			name:     "empty stream",
			stream:   model.Stream{Values: []int16{}},
			expected: 0,
		},
		{
			name:     "single value",
			stream:   model.Stream{Values: []int16{42}},
			expected: 1,
		},
		{
			name:     "multiple values",
			stream:   model.Stream{Values: []int16{1, 2, 3, 4, 5}},
			expected: 5,
		},
		{
			name:     "nil slice",
			stream:   model.Stream{Values: nil},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.stream.Len()
			require.Equal(t, tt.expected, actual)
		})
	}
}
