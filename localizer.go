package i18n

import (
	"fmt"
	"strings"
)

type Replace map[string]interface{}

type R = Replace

type Localizer struct {
	msg       *Message
	fallback  Fallback
	formatter Formatter
	selector  Selector
}

func NewLocalizer(msg *M) *Localizer {
	return &Localizer{
		msg:       msg,
		fallback:  defaultFallback,
		formatter: defaultFormatter,
		selector:  defaultSelector,
	}
}

func (l *Localizer) Fallback(fn Fallback) *Localizer {
	l.fallback = fn

	return l
}

func (l *Localizer) Formatter(fn Formatter) *Localizer {
	l.formatter = fn

	return l
}

func (l *Localizer) Selector(fn Selector) *Localizer {
	l.selector = fn

	return l
}

func (l *Localizer) T(key string) string {
	return l.Get(key)
}

func (l *Localizer) TC(key string, n uint) string {
	return l.GetNum(key, n)
}

func (l *Localizer) Tf(key string, args ...interface{}) string {
	return l.Getf(key, args...)
}

func (l *Localizer) TCf(key string, n uint, args ...interface{}) string {
	return l.GetNumf(key, n, args...)
}

func (l *Localizer) NTf(key string, rep R) string {
	return l.NamedGetf(key, rep)
}

func (l *Localizer) NTCf(key string, n uint, rep R) string {
	return l.NamedGetNumf(key, n, rep)
}

func (l *Localizer) Get(key string) string {
	return l.GetNum(key, 0)
}

func (l *Localizer) GetNum(key string, n uint) string {
	if msg := l.lookup(key); msg != nil {
		return l.pluralize(msg, n)
	}

	return l.fallback(key)
}

func (l *Localizer) Getf(key string, args ...interface{}) string {
	return l.GetNumf(key, 0, args...)
}

func (l *Localizer) GetNumf(key string, n uint, args ...interface{}) string {
	rep := R{}

	for i, arg := range args {
		rep[fmt.Sprintf("%d", i)] = arg
	}

	return l.NamedGetNumf(key, n, rep)
}

func (l *Localizer) NamedGetf(key string, rep R) string {
	return l.NamedGetNumf(key, 0, rep)
}

func (l *Localizer) NamedGetNumf(key string, n uint, rep R) string {
	if msg := l.lookup(key); msg != nil {
		return l.replace(l.pluralize(msg, n), rep)
	}

	return l.fallback(key)
}

func (l *Localizer) lookup(key string) []string {
	if msg := l.msg.get(key); 0 < len(msg) {
		return msg
	}

	return nil
}

func (l *Localizer) pluralize(msgs []string, n uint) string {
	if len(msgs) == 0 {
		return ""
	}

	if idx := l.selector(n); idx < uint(len(msgs)) {
		return msgs[idx]
	}

	return msgs[0]
}

func (l *Localizer) replace(msg string, rep R) string {
	for k, v := range rep {
		msg = strings.ReplaceAll(msg, "{"+k+"}", l.formatter(v))
	}

	return msg
}
