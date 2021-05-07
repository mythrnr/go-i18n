package i18n_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/mythrnr/i18n-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const _testJSONContent = `{
	"0-1": "apple",
	"0-2": ["no apples", "one apple", "{count} apples"],
	"0-3": {
		"1-1": "orange",
		"1-2": ["no oranges", "one orange", "{count} oranges"],
		"1-3": "cannot access"
	},
	"0-3.1-3": "direct-value-[0-3.1-3]",
	"0-4": {
		"1-1": {
			"2-1": "banana",
			"2-2": ["no bananas", "one banana", "{count} bananas"],
			"2-3": "cannot access"
		},
		"1-1.2-3": "direct-value-[1-1.2-3]"
	},
	"0-5": "strawberry"
}`

const _tsetYAMLContent = `
"0-1": "apple"
"0-2":
  - "no apples"
  - "one apple"
  - "{count} apples"
"0-3":
  "1-1": "orange"
  "1-2":
    - "no oranges"
    - "one orange"
    - "{count} oranges"
  "1-3": "cannot access"
"0-3.1-3": "direct-value-[0-3.1-3]"
"0-4":
  "1-1":
    "2-1": "banana"
    "2-2":
      - "no bananas"
      - "one banana"
      - "{count} bananas"
    "2-3": "cannot access"
  "1-1.2-3": "direct-value-[1-1.2-3]"
"0-5": "strawberry"
`

func Test_loader_Load_JSON(t *testing.T) {
	t.Parallel()

	t.Run("Load invalid JSON bytes", func(t *testing.T) {
		t.Parallel()

		value := []byte(`{ "0-1": "apple", }`)

		loader := i18n.NewJSONLoader()
		m, err := loader.Load(value)

		assert.NotNil(t, err)
		assert.Nil(t, m)
	})

	t.Run("Load valid JSON bytes, but invalid value", func(t *testing.T) {
		t.Parallel()

		tests := []struct{ v []byte }{{
			v: []byte(`{ "0-1": 1 }`),
		}, {
			v: []byte(`{ "0-2": true }`),
		}, {
			v: []byte(`{ "0-3": 3.3 }`),
		}, {
			v: []byte(`{ "0-4": ["value", 4] }`),
		}, {
			v: []byte(`{ "0-5": { "0-5.1-1": false } }`),
		}}

		loader := i18n.NewJSONLoader()

		for _, tt := range tests {
			m, err := loader.Load(tt.v)

			require.NotNil(t, err)
			assert.ErrorIs(t, err, i18n.ErrMessageValueTypeInvalid)
			assert.Nil(t, m)
		}
	})

	t.Run("Load JSON bytes", func(t *testing.T) {
		t.Parallel()

		value := []byte(_testJSONContent)
		loader := i18n.NewJSONLoader()
		m, err := loader.Load(value)

		assert.Nil(t, err)
		require.NotNil(t, m)
		assert.Equal(t, "apple", (*m)["0-1"])
	})
}

func Test_loader_Load_YAML(t *testing.T) {
	t.Parallel()

	t.Run("Load invalid YAML bytes", func(t *testing.T) {
		t.Parallel()

		value := []byte(`invalid yaml`)

		loader := i18n.NewYAMLLoader()
		m, err := loader.Load(value)

		assert.NotNil(t, err)
		assert.Nil(t, m)
	})

	t.Run("Load valid YAML bytes, but invalid value", func(t *testing.T) {
		t.Parallel()

		tests := []struct{ v []byte }{{
			v: []byte(`"0-1": 1`),
		}, {
			v: []byte(`"0-2": true`),
		}, {
			v: []byte(`"0-3": 3.3`),
		}, {
			v: []byte(`"0-4": ["value", 4]`),
		}, {
			v: []byte(`"0-5": { "0-5.1-1": false }`),
		}}

		loader := i18n.NewYAMLLoader()

		for _, tt := range tests {
			m, err := loader.Load(tt.v)

			require.NotNil(t, err)
			assert.ErrorIs(t, err, i18n.ErrYAMLValueTypeInvalid)
			assert.Nil(t, m)
		}
	})

	t.Run("Load YAML bytes", func(t *testing.T) {
		t.Parallel()

		value := []byte(_tsetYAMLContent)
		loader := i18n.NewYAMLLoader()
		m, err := loader.Load(value)

		assert.Nil(t, err)
		require.NotNil(t, m)
		assert.Equal(t, "apple", (*m)["0-1"])
	})
}

func Test_loader_LoadFile(t *testing.T) {
	t.Parallel()

	t.Run("File not exists", func(t *testing.T) {
		t.Parallel()

		loader := i18n.NewJSONLoader()
		m, err := loader.LoadFile("not_exists_file.json")

		assert.NotNil(t, err)
		assert.Nil(t, m)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		path := fmt.Sprintf(
			"/tmp/Test_loader_LoadFile_%d.json",
			time.Now().UnixNano(),
		)

		defer os.Remove(path)

		require.Nil(t, ioutil.WriteFile(
			path,
			[]byte(_testJSONContent),
			os.ModePerm),
		)

		loader := i18n.NewJSONLoader()
		m, err := loader.LoadFile(path)

		require.Nil(t, err)
		assert.Equal(t, "apple", (*m)["0-1"])
	})
}
