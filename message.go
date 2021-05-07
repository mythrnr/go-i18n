package i18n

import "strings"

type Message map[string]interface{}

type M = Message

func (m *Message) get(key string) []string {
	if msg := m.lookup(key); 0 < len(msg) {
		return msg
	}

	return []string{key}
}

func (m *Message) lookup(key string) []string {
	v, ok := (*m)[key]
	if ok {
		switch v := v.(type) {
		case string:
			return []string{v}
		case []string:
			return v
		default:
			return nil
		}
	}

	i := strings.IndexByte(key, '.')
	if i == -1 {
		return nil
	}

	v, ok = (*m)[key[:i]]
	if !ok {
		return nil
	}

	child, ok := v.(*Message)
	if !ok {
		return nil
	}

	return child.lookup(key[i+1:])
}
