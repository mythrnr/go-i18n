package i18n_test

import (
	"testing"

	"github.com/mythrnr/go-i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_Translator(t *testing.T) {
	t.Parallel()

	tr := i18n.NewTranslator()

	tr.Add(
		language.AmericanEnglish,
		i18n.NewLocalizer(&i18n.M{"key": "test"}),
	)

	tr.Add(
		language.Japanese,
		i18n.NewLocalizer(&i18n.M{"key": "テスト"}),
	)

	assert.True(t, tr.IsSupported(language.AmericanEnglish))
	assert.True(t, tr.IsSupported(language.Japanese))
	assert.False(t, tr.IsSupported(language.English))

	assert.Equal(t, "test", tr.L(language.AmericanEnglish).T("key"))
	assert.Equal(t, "テスト", tr.L(language.Japanese).T("key"))
	assert.Equal(t,
		"Undefined key in message: key",
		tr.L(language.English).T("key"),
	)

	tr.Fallback(i18n.NewLocalizer(&i18n.M{"key": "fallback-test"}))

	assert.Equal(t, "fallback-test", tr.Localizer(language.English).T("key"))
}
