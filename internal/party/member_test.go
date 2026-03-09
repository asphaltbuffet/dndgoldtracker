package party

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateLevel(t *testing.T) {
	tests := []struct {
		name string
		xp   int
		want int
	}{
		{"no xp", 0, 1},
		{"very high", 999_999, 20},
		{"exact", 300, 2},
		{"under by 1", 2699, 3},
		{"mid", 30_000, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Member{XP: tt.xp}

			m.UpdateLevel()

			assert.Equal(t, tt.want, m.Level)
		})
	}
}
