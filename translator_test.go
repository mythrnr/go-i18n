package i18n_test

import (
	"testing"

	"github.com/mythrnr/i18n-go"
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

	assert.Equal(t, "test", tr.Localizer(language.AmericanEnglish).T("key"))
	assert.Equal(t, "テスト", tr.Localizer(language.Japanese).T("key"))
	assert.Equal(t, "key", tr.Localizer(language.English).T("key"))

	tr.Fallback(i18n.NewLocalizer(&i18n.M{"key": "fallback-test"}))

	assert.Equal(t, "fallback-test", tr.Localizer(language.English).T("key"))
}
