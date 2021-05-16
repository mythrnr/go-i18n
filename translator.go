package i18n

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

func (t *Translator) IsSupported(tag language.Tag) bool {
	_, supported := t.mapping[tag]

	return supported
}

func (t *Translator) L(tag language.Tag) *Localizer {
	return t.Localizer(tag)
}

func (t *Translator) Localizer(tag language.Tag) *Localizer {
	if tr, ok := t.mapping[tag]; ok {
		return tr
	}

	return t.fallback
}
