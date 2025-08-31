package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const customChars = "لورم ایپسوم"

func TestNewCharsetBuilder(t *testing.T) {
	builder := NewCharsetBuilder()
	assert.NotNil(t, builder)
	assert.Empty(t, builder.Characters())
	assert.True(t, builder.IsEmpty())
	assert.Equal(t, 0, builder.Length())
}

func TestCharsetBuilderWithUppercase(t *testing.T) {
	builder := NewCharsetBuilder()
	result := builder.WithUppercase()

	assert.Same(t, builder, result)

	chars := builder.Characters()
	assert.Equal(t, UppercaseChars, chars)
	assert.Len(t, UppercaseChars, builder.Length())
	assert.Subset(t, []byte(chars), []byte(UppercaseChars))
}

func TestCharsetBuilderWithLowercase(t *testing.T) {
	builder := NewCharsetBuilder()
	result := builder.WithLowercase()

	assert.Same(t, builder, result)

	chars := builder.Characters()
	assert.Equal(t, LowercaseChars, chars)
	assert.Len(t, LowercaseChars, builder.Length())
	assert.Subset(t, []byte(chars), []byte(LowercaseChars))
}

func TestCharsetBuilderWithNumbers(t *testing.T) {
	builder := NewCharsetBuilder()
	result := builder.WithNumbers()

	assert.Same(t, builder, result)

	chars := builder.Characters()
	assert.Equal(t, NumberChars, chars)
	assert.Len(t, NumberChars, builder.Length())

	assert.Subset(t, []byte(chars), []byte(NumberChars))
}

func TestCharsetBuilderWithSymbols(t *testing.T) {
	builder := NewCharsetBuilder()
	result := builder.WithSymbols()

	assert.Same(t, builder, result)

	chars := builder.Characters()
	assert.Equal(t, SymbolChars, chars)
	assert.Len(t, SymbolChars, builder.Length())
	assert.Subset(t, []byte(chars), []byte(SymbolChars))
}

func TestCharsetBuilderWithCustom(t *testing.T) {
	t.Run("with valid custom characters", func(t *testing.T) {
		builder := NewCharsetBuilder()
		result := builder.WithCustom(customChars)

		assert.Same(t, builder, result)

		chars := builder.Characters()
		assert.Equal(t, customChars, chars)
		assert.Equal(t, len(customChars), builder.Length())
		assert.Subset(t, []byte(chars), []byte(customChars))
	})

	t.Run("with empty string", func(t *testing.T) {
		builder := NewCharsetBuilder()
		result := builder.WithCustom("")

		assert.Same(t, builder, result)
		assert.Empty(t, builder.Characters())
		assert.True(t, builder.IsEmpty())
	})
}

func TestCharsetBuilderChainedCalls(t *testing.T) {
	builder := NewCharsetBuilder()
	result := builder.
		WithUppercase().
		WithLowercase().
		WithNumbers().
		WithSymbols().
		WithCustom(customChars)

	assert.Same(t, builder, result)

	chars := builder.Characters()
	expectedLength := len(UppercaseChars) + len(LowercaseChars) + len(NumberChars) + len(SymbolChars) + len(customChars)
	assert.Len(t, chars, expectedLength)
	chars_byte := []byte(chars)

	assert.Subset(t, chars_byte, []byte(UppercaseChars))
	assert.Subset(t, chars_byte, []byte(LowercaseChars))
	assert.Subset(t, chars_byte, []byte(NumberChars))
	assert.Subset(t, chars_byte, []byte(SymbolChars))
	assert.Subset(t, chars_byte, []byte(customChars))
}

func TestCharsetBuilderReset(t *testing.T) {
	builder := NewCharsetBuilder()
	builder.WithUppercase().WithLowercase().WithNumbers()

	assert.False(t, builder.IsEmpty())
	assert.Positive(t, builder.Length())

	result := builder.Reset()
	assert.Same(t, builder, result)
	assert.Empty(t, builder.Characters())
	assert.True(t, builder.IsEmpty())
	assert.Equal(t, 0, builder.Length())
}

func TestCharsetBuilderIsEmpty(t *testing.T) {
	builder := NewCharsetBuilder()

	assert.Condition(t, func() (success bool) {
		success = len(builder.characters) == 0 && builder.IsEmpty()
		return
	})
}

func TestCharsetBuilderLength(t *testing.T) {
	builder := NewCharsetBuilder()

	assert.Condition(t, func() (success bool) {
		success = len(builder.characters) == builder.Length()
		return
	})
}

func TestCharsetBuilderEdgeCases(t *testing.T) {
	t.Run("custom with unicode characters", func(t *testing.T) {
		builder := NewCharsetBuilder()
		unicodeChars := "🔒🔑🛡️"
		builder.WithCustom(unicodeChars)

		chars := builder.Characters()
		assert.Len(t, chars, len(unicodeChars))
		for _, char := range unicodeChars {
			assert.Contains(t, chars, string(char))
		}
	})
}

func TestConstants(t *testing.T) {
	t.Run("uppercase constants", func(t *testing.T) {
		uppercase_chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		assert.Len(t, UppercaseChars, len(uppercase_chars))
		assert.Subset(t, []byte(UppercaseChars), []byte(uppercase_chars))
	})

	t.Run("lowercase constants", func(t *testing.T) {
		lowercase_chars := "abcdefghijklmnopqrstuvwxyz"
		assert.Len(t, LowercaseChars, len(lowercase_chars))
		assert.Subset(t, []byte(LowercaseChars), []byte(lowercase_chars))
	})

	t.Run("number constants", func(t *testing.T) {
		number_chars := "0123456789"
		assert.Len(t, NumberChars, len(number_chars))
		assert.Subset(t, []byte(NumberChars), []byte(number_chars))
	})

	t.Run("symbol constants", func(t *testing.T) {
		symbol_chars := "`~!@#$%^&*()_-+=\\|[]'\";:/?<>.,"
		assert.Len(t, SymbolChars, len(symbol_chars))
		assert.Subset(t, []byte(SymbolChars), []byte(symbol_chars))
	})
}

func BenchmarkCharsetBuilderSingleType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := NewCharsetBuilder()
		builder.WithUppercase()
	}
}

func BenchmarkCharsetBuilderAllTypes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := NewCharsetBuilder()
		builder.WithUppercase().WithLowercase().WithNumbers().WithSymbols()
	}
}
