package i18n

import (
	"fmt"
	"strconv"
)

type Fallback func(key string) string

type Formatter func(v interface{}) string

type Selector func(n uint) uint

func defaultFallback(key string) string {
	return fmt.Sprintf("Undefined key in message: %s", key)
}

func defaultFormatter(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return strconv.FormatBool(v)
	}

	return fmt.Sprintf("%v", v)
}

func defaultSelector(_ uint) uint { return 0 }
