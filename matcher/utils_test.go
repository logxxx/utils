package matcher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersionCompare(t *testing.T) {
	testCases := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2.3", "1.2.3", 0},
		{"1.2.3.0.0", "1.2.3", 0},
		{"9.9999", "10", -1},
		{"9.9999", "10.0", -1},
		{"9.100", "9.1000", -1},
		{"9.100", "9.100.1", -1},
		{"9.100", "8.999", 1},
		{"9.100", "9.99.99999", 1},

		{"v1.2.3", "1.2.3", 0},
		{"v1.2.3.0.0", "1.2.3", 0},
		{"v9.9999", "10", -1},
		{"v9.9999", "v10.0", -1},
		{"v9.100", "V9.1000", -1},
		{"9.100", "V9.100.1", -1},
		{"9.100", "V8.999", 1},
		{"9.100", "V9.99.99999", 1},

		{"abcd", "0.0.0", 0},
	}

	for _, v := range testCases {
		got := VersionCompare(v.v1, v.v2)
		t.Logf("compare %s to %s, got: %d", v.v1, v.v2, got)
		assert.Equal(t, v.expected, got)
	}
}
