package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultFormatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		v    interface{}
		want string
	}{{
		v:    "string-value",
		want: "string-value",
	}, {
		v:    int(-1),
		want: "-1",
	}, {
		v:    int(1),
		want: "1",
	}, {
		v:    uint(0),
		want: "0",
	}, {
		v:    uint(1),
		want: "1",
	}, {
		v:    float32(0.123456789),
		want: "0.123457",
	}, {
		v:    float64(9.876543210),
		want: "9.876543",
	}, {
		v:    true,
		want: "true",
	}, {
		v:    false,
		want: "false",
	}, {
		v:    struct{ n int }{1},
		want: "{1}",
	}}

	for _, tt := range tests {
		t.Log(tt.want, tt.v)
		assert.Equal(t, tt.want, defaultFormatter(tt.v))
	}
}
