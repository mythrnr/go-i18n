package message

import (
	"fmt"
	"strconv"
	"strings"
)

type Replacement map[string]interface{}

type Formatter func(v interface{}) string

type Selector func(n uint) uint

type Localizer struct {
	msg       *M
	formatter Formatter
	selector  Selector
}

func NewLocalizer(msg *M) *Localizer {
	return &Localizer{
		msg:       msg,
		formatter: defaultFormatter,
		selector:  defaultSelector,
	}
}

func (l *Localizer) Formatter(fn Formatter) *Localizer {
	l.formatter = fn

	return l
}

func (l *Localizer) Selector(fn Selector) *Localizer {
	l.selector = fn

	return l
}

func (l *Localizer) Get(key string) string {
	return l.GetNum(key, 0)
}

func (l *Localizer) GetNum(key string, n uint) string {
	return l.pluralize(l.msg.get(key), n)
}

func (l *Localizer) Getf(key string, args ...interface{}) string {
	return l.GetNumf(key, 0, args...)
}

func (l *Localizer) GetNumf(key string, n uint, args ...interface{}) string {
	rep := Replacement{}

	for i, a := range args {
		rep[fmt.Sprintf("%d", i)] = a
	}

	return l.NamedGetNumf(key, n, rep)
}

func (l *Localizer) NamedGetf(key string, rep Replacement) string {
	return l.NamedGetNumf(key, 0, rep)
}

func (l *Localizer) NamedGetNumf(key string, n uint, rep Replacement) string {
	return l.replace(l.pluralize(l.msg.get(key), n), rep)
}

func (l *Localizer) pluralize(msgs []string, n uint) string {
	if len(msgs) == 0 {
		return ""
	}

	idx := l.selector(n)
	if idx < uint(len(msgs)) {
		return msgs[idx]
	}

	return msgs[0]
}

func (l *Localizer) replace(msg string, rep Replacement) string {
	for k, v := range rep {
		msg = strings.ReplaceAll(msg, "{"+k+"}", l.formatter(v))
	}

	return msg
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