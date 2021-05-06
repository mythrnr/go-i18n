package message

import "golang.org/x/text/language"

type Translator struct {
	fallback *Localizer
	mapping  map[language.Tag]*Localizer
}

func NewTranslator() *Translator {
	return &Translator{
		fallback: NewLocalizer(&M{}),
		mapping:  map[language.Tag]*Localizer{},
	}
}

func (t *Translator) Add(tag language.Tag, l *Localizer) {
	t.mapping[tag] = l
}

func (t *Translator) Fallback(l *Localizer) {
	t.fallback = l
}

func (t *Translator) Get(tag language.Tag) *Localizer {
	if tr, ok := t.mapping[tag]; ok {
		return tr
	}

	return t.fallback
}

func (t *Translator) IsSupported(tag language.Tag) bool {
	_, ok := t.mapping[tag]

	return ok
}
