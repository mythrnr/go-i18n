package i18n

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Localizer_pluralize(t *testing.T) {
	t.Parallel()

	t.Run("Specified empty messages", func(t *testing.T) {
		t.Parallel()

		l := NewLocalizer(&M{})
		m := l.pluralize([]string{}, 1)

		assert.Equal(t, "", m)
	})

	t.Run("With default selector", func(t *testing.T) {
		t.Parallel()

		l := NewLocalizer(&M{})
		m := l.pluralize([]string{
			"no apples",
			"one apple",
			"{count} apples",
		}, 1)

		assert.Equal(t, "no apples", m)
	})

	t.Run("With custom selector", func(t *testing.T) {
		t.Parallel()

		l := NewLocalizer(&M{}).Selector(selector)
		m := l.pluralize([]string{
			"no apples",
			"one apple",
			"{count} apples",
		}, 1)

		assert.Equal(t, "one apple", m)
	})

	t.Run("With custom selector and out of range", func(t *testing.T) {
		t.Parallel()

		l := NewLocalizer(&M{}).Selector(selector)
		m := l.pluralize([]string{
			"no apples",
			"one apple",
		}, 2)

		assert.Equal(t, "no apples", m)
	})
}

func Test_Localizer_replace(t *testing.T) {
	t.Parallel()

	t.Run("With default formatter", func(t *testing.T) {
		t.Parallel()

		l := NewLocalizer(&M{})

		msg := "{stringKey}, {intKey}, {stringKey}, test string, {keyNotSet}"
		rep := R{
			"stringKey": "string-value",
			"intKey":    2,
		}

		assert.Equal(t,
			"string-value, 2, string-value, test string, {keyNotSet}",
			l.replace(msg, rep),
		)
	})

	t.Run("With custom formatter", func(t *testing.T) {
		t.Parallel()

		l := NewLocalizer(&M{}).Formatter(func(v interface{}) string {
			// just print name of type
			return reflect.TypeOf(v).String()
		})

		msg := "{stringKey}, {intKey}, {stringKey}, test string, {keyNotSet}"
		rep := R{
			"stringKey": "string-value",
			"intKey":    2,
		}

		assert.Equal(t,
			"string, int, string, test string, {keyNotSet}",
			l.replace(msg, rep),
		)
	})
}

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

var selector = func(n uint) uint { return uint(math.Min(2, float64(n))) }

func Test_selector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		n    uint
		want uint
	}{{
		n:    0,
		want: 0,
	}, {
		n:    1,
		want: 1,
	}, {
		n:    2,
		want: 2,
	}, {
		n:    3,
		want: 2,
	}}

	for _, tt := range tests {
		t.Log(tt.want, tt.n)
		assert.Equal(t, tt.want, selector(tt.n))
	}
}
