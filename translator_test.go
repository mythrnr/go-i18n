package message_test

import (
	"testing"

	"github.com/mythrnr/go-message"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_Translator(t *testing.T) {
	t.Parallel()

	tr := message.NewTranslator()

	tr.Add(
		language.AmericanEnglish,
		message.NewLocalizer(&message.M{"key": "test"}),
	)

	tr.Add(
		language.Japanese,
		message.NewLocalizer(&message.M{"key": "テスト"}),
	)

	assert.True(t, tr.IsSupported(language.AmericanEnglish))
	assert.True(t, tr.IsSupported(language.Japanese))
	assert.False(t, tr.IsSupported(language.English))

	assert.Equal(t, "test", tr.Get(language.AmericanEnglish).Getf("key"))
	assert.Equal(t, "テスト", tr.Get(language.Japanese).Getf("key"))
	assert.Equal(t, "key", tr.Get(language.English).Getf("key"))

	tr.Fallback(message.NewLocalizer(&message.M{"key": "fallback-test"}))

	assert.Equal(t, "fallback-test", tr.Get(language.English).Getf("key"))
}
