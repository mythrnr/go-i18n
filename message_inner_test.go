package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Message_get(t *testing.T) {
	t.Parallel()

	// {nest}-{case}
	m := &M{
		"0-1": []string{},
		"0-2": []string{"apple"},
		"0-3": []string{"no apples", "one apple", "{count} apples"},
		"0-4": nil,
	}

	tests := []struct {
		key  string
		want []string
	}{{
		key:  "0-1",
		want: []string{"0-1"},
	}, {
		key:  "0-2",
		want: []string{"apple"},
	}, {
		key:  "0-3",
		want: []string{"no apples", "one apple", "{count} apples"},
	}, {
		key:  "0-4",
		want: []string{"0-4"},
	}}

	for _, tt := range tests {
		t.Log(tt.key, tt.want)
		assert.Equal(t, tt.want, m.get(tt.key))
	}
}

func Test_Message_lookup(t *testing.T) {
	t.Parallel()

	// {nest}-{case}
	m := &M{
		"0-1": "apple",
		"0-2": []string{"no apples", "one apple", "{count} apples"},
		"0-3": struct{ v string }{"invalid-value"},
		"0-4": 0, // invalid
		"0-5": &M{
			"1-1": "orange",
			"1-2": []string{"no oranges", "one orange", "{count} oranges"},
			"1-3": struct{ v string }{"invalid-value"},
			"1-4": 1, // invalid
			"1-5": "cannot access",
		},
		"0-5.1-5": "direct-value-[0-5.1-5]",
		"0-6": &M{
			"1-1": &M{
				"2-1": "banana",
				"2-2": []string{"no bananas", "one banana", "{count} bananas"},
				"2-3": struct{ v string }{"invalid-value"},
				"2-4": 2, // invalid
				"2-5": "cannot access",
			},
			"1-1.2-5": "direct-value-[1-1.2-5]",
		},
		"0-7": "strawberry",
	}

	tests := []struct {
		key  string
		want []string
	}{{
		key:  "0-0",
		want: nil,
	}, {
		key:  "0-1",
		want: []string{"apple"},
	}, {
		key:  "0-2",
		want: []string{"no apples", "one apple", "{count} apples"},
	}, {
		key:  "0-3",
		want: nil,
	}, {
		key:  "0-4",
		want: nil,
	}, {
		key:  "0-5",
		want: nil,
	}, {
		key:  "0-5.1-1",
		want: []string{"orange"},
	}, {
		key:  "0-5.1-2",
		want: []string{"no oranges", "one orange", "{count} oranges"},
	}, {
		key:  "0-5.1-3",
		want: nil,
	}, {
		key:  "0-5.1-4",
		want: nil,
	}, {
		key:  "0-5.1-5",
		want: []string{"direct-value-[0-5.1-5]"},
	}, {
		key:  "0-6",
		want: nil,
	}, {
		key:  "0-6.1-1",
		want: nil,
	}, {
		key:  "0-6.1-1.2-1",
		want: []string{"banana"},
	}, {
		key:  "0-6.1-1.2-2",
		want: []string{"no bananas", "one banana", "{count} bananas"},
	}, {
		key:  "0-6.1-1.2-3",
		want: nil,
	}, {
		key:  "0-6.1-1.2-4",
		want: nil,
	}, {
		key:  "0-6.1-1.2-5",
		want: []string{"direct-value-[1-1.2-5]"},
	}, {
		key:  "0-7.1-1",
		want: nil,
	}, {
		key:  "0-8.1-1",
		want: nil,
	}}

	for _, tt := range tests {
		t.Log(tt.key, tt.want)
		assert.Equal(t, tt.want, m.lookup(tt.key))
	}
}
