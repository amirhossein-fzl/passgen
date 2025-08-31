package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordGeneratorOptions(t *testing.T) {
	options := NewPasswordGeneratorOptions()

	assert.Equal(t, 12, options.Length)
	assert.True(t, options.Lowercase)
	assert.True(t, options.Uppercase)
	assert.True(t, options.Numbers)
	assert.False(t, options.Symbols)
	assert.Empty(t, options.Exclude)
	assert.False(t, options.ExcludeSimilar)
	assert.Empty(t, options.Custom)
	assert.Equal(t, 1, options.AvoidRepeats)
	assert.False(t, options.QrCode)
}
