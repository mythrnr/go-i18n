package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_unmarshalYAML(t *testing.T) {
	t.Parallel()

	t.Run("Invalid YAML", func(t *testing.T) {
		t.Parallel()

		b := []byte(`
Name: Platypus
Order: -
`,
		)

		value := map[string]interface{}{}
		err := unmarshalYAML(b, &value)

		assert.NotNil(t, err)
	})

	t.Run("Invalid target Type", func(t *testing.T) {
		t.Parallel()

		b := []byte(`
Name: Platypus
Order:
  Name: Platypus
Asdf.Zxcv: Monotremata
`,
		)

		value := map[bool]int{}
		err := unmarshalYAML(b, &value)

		require.NotNil(t, err)
		assert.ErrorIs(t, err, ErrYAMLCannotBind)
	})

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		b := []byte(`
Name: Platypus
Order:
  Name: Platypus
Asdf.Zxcv: Monotremata
`,
		)

		value := map[string]interface{}{}
		err := unmarshalYAML(b, &value)

		assert.Nil(t, err)
	})
}

func Test_toStringMap(t *testing.T) {
	t.Parallel()

	t.Run("Invalid Key", func(t *testing.T) {
		t.Parallel()

		value := map[interface{}]interface{}{
			0: "invalid key",
		}

		converted, err := toStringMap(value)
		require.NotNil(t, err)
		assert.ErrorIs(t, err, ErrYAMLKeyNotString)
		assert.Nil(t, converted)
	})

	t.Run("Invalid Value", func(t *testing.T) {
		t.Parallel()

		value := map[interface{}]interface{}{
			"key": struct{}{},
		}

		converted, err := toStringMap(value)
		require.NotNil(t, err)
		assert.ErrorIs(t, err, ErrYAMLValueTypeInvalid)
		assert.Nil(t, converted)
	})

	t.Run("Invalid Nested Value", func(t *testing.T) {
		t.Parallel()

		value := map[interface{}]interface{}{
			"key": map[interface{}]interface{}{
				"nested": struct{}{},
			},
		}

		converted, err := toStringMap(value)
		require.NotNil(t, err)
		assert.ErrorIs(t, err, ErrYAMLValueTypeInvalid)
		assert.Nil(t, converted)
	})

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		value := map[interface{}]interface{}{
			"key": map[interface{}]interface{}{
				"nested": "nested value",
			},
			"key2": "value",
		}

		expected := map[string]interface{}{
			"key": map[string]interface{}{
				"nested": "nested value",
			},
			"key2": "value",
		}

		converted, err := toStringMap(value)
		assert.Nil(t, err)
		assert.Equal(t, expected, converted)
	})
}
