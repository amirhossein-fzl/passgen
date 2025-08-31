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
	assert.Empty(t, options.Custom)
	assert.Equal(t, 1, options.AvoidRepeats)
	assert.False(t, options.QrCode)
}

func TestValidate(t *testing.T) {

	tests := []struct {
		name           string
		options        PasswordGeneratorOptions
		expectedResult bool
		ErrWant        error
	}{
		{
			name: "valid length",
			options: PasswordGeneratorOptions{
				Length: 10,
			},
			expectedResult: true,
			ErrWant:        nil,
		},
		{
			name: "zero length",
			options: PasswordGeneratorOptions{
				Length: 0,
			},
			expectedResult: false,
			ErrWant:        ErrLengthMustBeGreaterThanZero,
		},
		{
			name: "negative length",
			options: PasswordGeneratorOptions{
				Length: -33,
			},
			expectedResult: false,
			ErrWant:        ErrLengthMustBeGreaterThanZero,
		},
		{
			name: "valid avoid repeats",
			options: PasswordGeneratorOptions{
				Length:       10,
				AvoidRepeats: 2,
			},
			expectedResult: true,
			ErrWant:        nil,
		},
		{
			name: "zero avoid repeats",
			options: PasswordGeneratorOptions{
				Length:       8,
				AvoidRepeats: 0,
			},
			expectedResult: true,
			ErrWant:        nil,
		},
		{
			name: "negative avoid repeats",
			options: PasswordGeneratorOptions{
				Length:       8,
				AvoidRepeats: -10,
			},
			expectedResult: false,
			ErrWant:        ErrAvoidRepeatsMustBeEqualOrGreaterThanZero,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := NewPasswordGeneratorOptions()
			options.Length = tt.options.Length
			options.AvoidRepeats = tt.options.AvoidRepeats
			result, err := options.Validate()

			assert.Equal(t, tt.expectedResult, result)
			if tt.ErrWant != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.ErrWant, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
