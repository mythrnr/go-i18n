package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v2"
)

var (
	ErrMessageValueTypeInvalid = errors.New("invalid message value")
	ErrYAMLCannotBind          = errors.New("invalid target type")
	ErrYAMLKeyNotString        = errors.New("invalid YAML key")
	ErrYAMLValueTypeInvalid    = errors.New("invalid YAML value")
)

type Loader interface {
	Load(data []byte) (*M, error)
	LoadFile(path string) (*M, error)
}

type loader struct {
	unmarshalFunc func(data []byte, v interface{}) error
}

var _ Loader = (*loader)(nil)

func NewJSONLoader() Loader {
	return &loader{unmarshalFunc: json.Unmarshal}
}

func NewYAMLLoader() Loader {
	return &loader{unmarshalFunc: unmarshalYAML}
}

func (l *loader) Load(data []byte) (*M, error) {
	value := map[string]interface{}{}

	if err := l.unmarshalFunc(data, &value); err != nil {
		return nil, err
	}

	m := &M{}

	if err := l.parse(m, "", value); err != nil {
		return nil, err
	}

	return m, nil
}

func (l *loader) LoadFile(path string) (*M, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return l.Load(b)
}

func (l *loader) parse(
	m *M,
	parentKey string,
	value map[string]interface{},
) error {
	if parentKey != "" {
		parentKey += "."
	}

	for k, in := range value {
		switch v := in.(type) {
		case string, []string:
			(*m)[k] = v
		case []interface{}:
			vv := make([]string, 0, len(v))

			for idx, in := range v {
				v, ok := in.(string)
				if !ok {
					return fmt.Errorf(
						"%w: Key: %s [%d], Value Type: %s",
						ErrMessageValueTypeInvalid,
						parentKey+k,
						idx,
						reflect.TypeOf(in),
					)
				}

				vv = append(vv, v)
			}

			(*m)[k] = vv
		case map[string]interface{}:
			child := &M{}

			if err := l.parse(child, parentKey+k, v); err != nil {
				return err
			}

			(*m)[k] = child
		default:
			return fmt.Errorf(
				"%w: Key: %s, Value Type: %s",
				ErrMessageValueTypeInvalid,
				parentKey+k,
				reflect.TypeOf(v),
			)
		}
	}

	return nil
}

// unmarshalYAML is a custom Unmarshal function.
// yaml.Unmarshal cannot decode into map[string]interface{},
// so we define a function converting map[interface{}]interface{}
// to map[string]interface{}.
//
// unmarshalYAML は独自の Unmarshal 関数.
// yaml.Unmarshalが map[string]interface{} にデコードできないので
// map[interface{}]interface{} -> map[string]interface{} の関数を定義.
func unmarshalYAML(data []byte, v interface{}) error {
	value := map[interface{}]interface{}{}

	if err := yaml.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// to map[string]interface{}
	converted, err := toStringMap(value)
	if err != nil {
		return err
	}

	m, ok := v.(*map[string]interface{})
	if !ok {
		return fmt.Errorf(
			"%w: target type: %s",
			ErrYAMLCannotBind,
			reflect.TypeOf(v),
		)
	}

	for k, val := range converted {
		(*m)[k] = val
	}

	return nil
}

// toStringMap is a recursive process of
// map[interface{}]interface{} to map[string]interface{}.
//
// toStringMap は map[interface{}]interface{} -> map[string]interface{} の再帰処理.
func toStringMap(
	value map[interface{}]interface{},
) (map[string]interface{}, error) {
	converted := map[string]interface{}{}

	for k, in := range value {
		key, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf(
				"%w: Key Type: %s",
				ErrYAMLKeyNotString,
				reflect.TypeOf(k),
			)
		}

		switch v := in.(type) {
		case string, []string:
			converted[key] = v
		case []interface{}:
			vv := make([]string, 0, len(v))

			for idx, in := range v {
				v, ok := in.(string)
				if !ok {
					return nil, fmt.Errorf(
						"%w: Key: %s [%d], Value Type: %s",
						ErrYAMLValueTypeInvalid,
						key,
						idx,
						reflect.TypeOf(in),
					)
				}

				vv = append(vv, v)
			}

			converted[key] = vv
		case map[interface{}]interface{}:
			c, err := toStringMap(v)
			if err != nil {
				return nil, err
			}

			converted[key] = c
		default:
			return nil, fmt.Errorf(
				"%w: Key: %s, Value Type: %s",
				ErrYAMLValueTypeInvalid,
				key,
				reflect.TypeOf(v),
			)
		}
	}

	return converted, nil
}
