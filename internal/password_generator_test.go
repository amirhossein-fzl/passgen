package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePassword(t *testing.T) {
	t.Run("should generate password with correct length", func(t *testing.T) {
		options := PasswordGeneratorOptions{
			Length:    15,
			Lowercase: true,
		}

		password, err := GeneratePassword(options)

		require.NoError(t, err)
		assert.Len(t, password, options.Length)
	})

	t.Run("should return error for empty charset", func(t *testing.T) {
		options := PasswordGeneratorOptions{
			Length: 10,
		}

		password, err := GeneratePassword(options)

		assert.Empty(t, password)
		assert.Equal(t, ErrEmptyCharset, err)
	})
	t.Run("should generate valid password with avoid repeats", func(t *testing.T) {
		options := PasswordGeneratorOptions{
			Length:       10,
			Lowercase:    true,
			Uppercase:    true,
			AvoidRepeats: 2,
		}

		password, err := GeneratePassword(options)

		require.NoError(t, err)
		assert.Len(t, password, 10)
	})

	t.Run("should handle zero length password", func(t *testing.T) {
		options := PasswordGeneratorOptions{Length: 0, Lowercase: true}

		password, err := GeneratePassword(options)

		require.NoError(t, err)
		assert.Empty(t, password)
	})
}

func TestNormalizeAvoidRepeats(t *testing.T) {
	tests := []struct {
		name          string
		avoidRepeats  int
		charsetLength int
		expected      int
	}{
		{
			name:          "avoid repeats less than charset length",
			avoidRepeats:  3,
			charsetLength: 10,
			expected:      3,
		},
		{
			name:          "avoid repeats equal to charset length",
			avoidRepeats:  10,
			charsetLength: 10,
			expected:      9,
		},
		{
			name:          "avoid repeats greater than charset length",
			avoidRepeats:  15,
			charsetLength: 10,
			expected:      9,
		},
		{
			name:          "zero avoid repeats",
			avoidRepeats:  0,
			charsetLength: 10,
			expected:      0,
		},
		{
			name:          "negative avoid repeats",
			avoidRepeats:  -5,
			charsetLength: 10,
			expected:      0,
		},
		{
			name:          "charset length 1",
			avoidRepeats:  1,
			charsetLength: 1,
			expected:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeAvoidRepeats(tt.avoidRepeats, tt.charsetLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSelectValidPasswordChar(t *testing.T) {
	t.Run("should select char from charset", func(t *testing.T) {
		charset := "abc"
		currentPassword := ""
		avoidRepeats := 0

		char, err := selectValidPasswordChar(charset, currentPassword, avoidRepeats)

		require.NoError(t, err)
		assert.Len(t, char, 1)
		assert.Contains(t, charset, char)
	})

	t.Run("should respect avoid repeats constraint", func(t *testing.T) {
		charset := "abcdef"
		currentPassword := "ab"
		avoidRepeats := 2

		char, err := selectValidPasswordChar(charset, currentPassword, avoidRepeats)

		require.NoError(t, err)
		assert.Len(t, char, 1)
		assert.Contains(t, charset, char)
		assert.NotContains(t, "ab", char)
	})

	t.Run("should return error for empty charset", func(t *testing.T) {
		charset := ""
		currentPassword := ""
		avoidRepeats := 0

		char, err := selectValidPasswordChar(charset, currentPassword, avoidRepeats)

		assert.Empty(t, char)
		assert.Error(t, err)
		assert.Equal(t, ErrCnnotSelectFromEmptyCharset, err)
	})

	t.Run("should work with no avoid repeats", func(t *testing.T) {
		charset := "a"
		currentPassword := "aaaa"
		avoidRepeats := 0

		char, err := selectValidPasswordChar(charset, currentPassword, avoidRepeats)

		require.NoError(t, err)
		assert.Equal(t, "a", char)
	})
}

func TestIsValidPasswordChar(t *testing.T) {
	tests := []struct {
		name         string
		character    string
		password     string
		avoidRepeats int
		expected     bool
	}{
		{
			name:         "valid char with no restrictions",
			character:    "a",
			password:     "bcdef",
			avoidRepeats: 0,
			expected:     true,
		},
		{
			name:         "valid char not in recent chars",
			character:    "d",
			password:     "abc",
			avoidRepeats: 2,
			expected:     true,
		},
		{
			name:         "invalid char in recent chars",
			character:    "c",
			password:     "abc",
			avoidRepeats: 2,
			expected:     false,
		},
		{
			name:         "valid char with avoid repeats larger than password",
			character:    "d",
			password:     "abc",
			avoidRepeats: 5,
			expected:     true,
		},
		{
			name:         "invalid char with avoid repeats larger than password",
			character:    "a",
			password:     "abc",
			avoidRepeats: 5,
			expected:     false,
		},
		{
			name:         "zero avoid repeats should allow any char",
			character:    "a",
			password:     "aaa",
			avoidRepeats: 0,
			expected:     true,
		},
		{
			name:         "negative avoid repeats should allow any char",
			character:    "a",
			password:     "aaa",
			avoidRepeats: -1,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidPasswordChar(tt.character, tt.password, tt.avoidRepeats)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPickRandomPasswordChar(t *testing.T) {
	t.Run("should pick char from single character charset", func(t *testing.T) {
		charset := "a"

		char, err := pickRandomChar(charset)

		require.NoError(t, err)
		assert.Equal(t, 'a', char)
	})

	t.Run("should pick char from multi character charset", func(t *testing.T) {
		charset := "abcdef"

		char, err := pickRandomChar(charset)

		require.NoError(t, err)
		assert.Contains(t, charset, string(char))
	})

	t.Run("should return error for empty charset", func(t *testing.T) {
		charset := ""

		char, err := pickRandomChar(charset)

		assert.Empty(t, char)
		assert.Error(t, err)
		assert.Equal(t, ErrCnnotSelectFromEmptyCharset, err)
	})
}

func TestSecureRandomInt(t *testing.T) {
	t.Run("should return value in range", func(t *testing.T) {
		min, max := 5, 10

		result, err := secureRandomInt(min, max)

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result, min)
		assert.LessOrEqual(t, result, max)
	})

	t.Run("should handle single value range", func(t *testing.T) {
		min, max := 5, 5

		result, err := secureRandomInt(min, max)

		require.NoError(t, err)
		assert.Equal(t, 5, result)
	})

	t.Run("should return error when min > max", func(t *testing.T) {
		min, max := 10, 5

		result, err := secureRandomInt(min, max)

		assert.Zero(t, result)
		assert.Error(t, err)
		assert.Equal(t, ErrMinCnnotGreaterThanMax, err)
	})

	t.Run("should handle negative values", func(t *testing.T) {
		min, max := -10, -5

		result, err := secureRandomInt(min, max)

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result, min)
		assert.LessOrEqual(t, result, max)
	})
}

func TestGetPasswordSuffix(t *testing.T) {
	tests := []struct {
		name         string
		password     string
		suffixLength int
		expected     string
	}{
		{
			name:         "suffix shorter than password",
			password:     "abcdef",
			suffixLength: 3,
			expected:     "def",
		},
		{
			name:         "suffix equal to password length",
			password:     "abc",
			suffixLength: 3,
			expected:     "abc",
		},
		{
			name:         "suffix longer than password",
			password:     "abc",
			suffixLength: 5,
			expected:     "abc",
		},
		{
			name:         "empty password",
			password:     "",
			suffixLength: 3,
			expected:     "",
		},
		{
			name:         "zero suffix length",
			password:     "abcdef",
			suffixLength: 0,
			expected:     "",
		},
		{
			name:         "single character password",
			password:     "a",
			suffixLength: 1,
			expected:     "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPasswordSuffix(tt.password, tt.suffixLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}
