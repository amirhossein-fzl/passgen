package internal

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureStderr(fn func()) string {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fn()

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNewCommandLineParser(t *testing.T) {
	parser := NewCommandLineParser()

	assert.NotNil(t, parser)
	assert.NotNil(t, parser.flagSet)
	assert.Equal(t, ProgramName, parser.flagSet.Name())
}

func TestCommandLineParserParseDefaultValues(t *testing.T) {
	parser := NewCommandLineParser()

	options, err := parser.Parse([]string{})

	require.NoError(t, err)
	require.NotNil(t, options)

	assert.Equal(t, DefaultPasswordLength, options.length)
	assert.True(t, options.lowercase)
	assert.True(t, options.uppercase)
	assert.True(t, options.numbers)
	assert.False(t, options.symbols)
	assert.Empty(t, options.custom)
	assert.Equal(t, DefaultAvoidRepeats, options.avoidRepeats)
	assert.False(t, options.qrOutput)
}

func TestCommandLineParserParseShortFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *CommandLineOptions
	}{
		{
			name: "length flag",
			args: []string{"-l", "20"},
			expected: &CommandLineOptions{
				length:       20,
				lowercase:    true,
				uppercase:    true,
				numbers:      true,
				symbols:      false,
				custom:       "",
				avoidRepeats: DefaultAvoidRepeats,
				qrOutput:     false,
			},
		},
		{
			name: "boolean flags true",
			args: []string{"-S=true", "-q=true"},
			expected: &CommandLineOptions{
				length:       DefaultPasswordLength,
				lowercase:    true,
				uppercase:    true,
				numbers:      true,
				symbols:      true,
				custom:       "",
				avoidRepeats: DefaultAvoidRepeats,
				qrOutput:     true,
			},
		},
		{
			name: "boolean flags false",
			args: []string{"-L=false", "-U=false", "-N=false"},
			expected: &CommandLineOptions{
				length:       DefaultPasswordLength,
				lowercase:    false,
				uppercase:    false,
				numbers:      false,
				symbols:      false,
				custom:       "",
				avoidRepeats: DefaultAvoidRepeats,
				qrOutput:     false,
			},
		},
		{
			name: "custom and avoid repeats",
			args: []string{"-C", "abc123", "-a", "3"},
			expected: &CommandLineOptions{
				length:       DefaultPasswordLength,
				lowercase:    true,
				uppercase:    true,
				numbers:      true,
				symbols:      false,
				custom:       "abc123",
				avoidRepeats: 3,
				qrOutput:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewCommandLineParser()
			options, err := parser.Parse(tt.args)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, options)
		})
	}
}

func TestCommandLineParserParseLongFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *CommandLineOptions
	}{
		{
			name: "length flag",
			args: []string{"--length", "25"},
			expected: &CommandLineOptions{
				length:       25,
				lowercase:    true,
				uppercase:    true,
				numbers:      true,
				symbols:      false,
				custom:       "",
				avoidRepeats: DefaultAvoidRepeats,
				qrOutput:     false,
			},
		},
		{
			name: "all boolean flags",
			args: []string{"--lowercase=false", "--uppercase=false", "--numbers=false", "--symbols=true", "--qr=true"},
			expected: &CommandLineOptions{
				length:       DefaultPasswordLength,
				lowercase:    false,
				uppercase:    false,
				numbers:      false,
				symbols:      true,
				custom:       "",
				avoidRepeats: DefaultAvoidRepeats,
				qrOutput:     true,
			},
		},
		{
			name: "custom character set",
			args: []string{"--custom", "!@#$%^&*()"},
			expected: &CommandLineOptions{
				length:       DefaultPasswordLength,
				lowercase:    true,
				uppercase:    true,
				numbers:      true,
				symbols:      false,
				custom:       "!@#$%^&*()",
				avoidRepeats: DefaultAvoidRepeats,
				qrOutput:     false,
			},
		},
		{
			name: "avoid repeats",
			args: []string{"--avoid-repeats", "5"},
			expected: &CommandLineOptions{
				length:       DefaultPasswordLength,
				lowercase:    true,
				uppercase:    true,
				numbers:      true,
				symbols:      false,
				custom:       "",
				avoidRepeats: 5,
				qrOutput:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewCommandLineParser()
			options, err := parser.Parse(tt.args)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, options)
		})
	}
}

func TestCommandLineParserParseMixedFlags(t *testing.T) {
	parser := NewCommandLineParser()

	args := []string{
		"-l", "16",
		"--lowercase=false",
		"-S=true",
		"--custom", "mixed123",
		"-a", "4",
		"--qr=false",
	}

	options, err := parser.Parse(args)

	require.NoError(t, err)
	assert.Equal(t, 16, options.length)
	assert.False(t, options.lowercase)
	assert.True(t, options.uppercase)
	assert.True(t, options.numbers)
	assert.True(t, options.symbols)
	assert.Equal(t, "mixed123", options.custom)
	assert.Equal(t, 4, options.avoidRepeats)
	assert.False(t, options.qrOutput)
}

func TestCommandLineParserParseEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		validate func(t *testing.T, options *CommandLineOptions)
	}{
		{
			name: "minimum length",
			args: []string{"-l", "1"},
			validate: func(t *testing.T, options *CommandLineOptions) {
				assert.Equal(t, 1, options.length)
			},
		},
		{
			name: "large length",
			args: []string{"--length", "1000"},
			validate: func(t *testing.T, options *CommandLineOptions) {
				assert.Equal(t, 1000, options.length)
			},
		},
		{
			name: "zero avoid repeats",
			args: []string{"-a", "0"},
			validate: func(t *testing.T, options *CommandLineOptions) {
				assert.Equal(t, 0, options.avoidRepeats)
			},
		},
		{
			name: "large avoid repeats",
			args: []string{"--avoid-repeats", "100"},
			validate: func(t *testing.T, options *CommandLineOptions) {
				assert.Equal(t, 100, options.avoidRepeats)
			},
		},
		{
			name: "empty custom string",
			args: []string{"-C", ""},
			validate: func(t *testing.T, options *CommandLineOptions) {
				assert.Empty(t, options.custom)
			},
		},
		{
			name: "complex custom string",
			args: []string{"--custom", "!@#$%^&*()_+-=[]{}|;:'\",.<>?/~`"},
			validate: func(t *testing.T, options *CommandLineOptions) {
				assert.Equal(t, "!@#$%^&*()_+-=[]{}|;:'\",.<>?/~`", options.custom)
			},
		},
		{
			name: "all character types disabled",
			args: []string{"-L=false", "-U=false", "-N=false", "-S=false"},
			validate: func(t *testing.T, options *CommandLineOptions) {
				assert.False(t, options.lowercase)
				assert.False(t, options.uppercase)
				assert.False(t, options.numbers)
				assert.False(t, options.symbols)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewCommandLineParser()
			options, err := parser.Parse(tt.args)

			require.NoError(t, err)
			require.NotNil(t, options)
			tt.validate(t, options)
		})
	}
}

func TestCommandLineOptionsValidate(t *testing.T) {
	tests := []struct {
		name        string
		options     *CommandLineOptions
		expectedErr error
	}{
		{
			name: "valid options - default values",
			options: &CommandLineOptions{
				length:       DefaultPasswordLength,
				avoidRepeats: DefaultAvoidRepeats,
			},
			expectedErr: nil,
		},
		{
			name: "valid options - custom values",
			options: &CommandLineOptions{
				length:       20,
				avoidRepeats: 5,
			},
			expectedErr: nil,
		},
		{
			name: "valid options - minimum values",
			options: &CommandLineOptions{
				length:       1,
				avoidRepeats: 0,
			},
			expectedErr: nil,
		},
		{
			name: "invalid length - zero",
			options: &CommandLineOptions{
				length:       0,
				avoidRepeats: DefaultAvoidRepeats,
			},
			expectedErr: ErrLengthMustBeGreaterThanZero,
		},
		{
			name: "invalid length - negative",
			options: &CommandLineOptions{
				length:       -5,
				avoidRepeats: DefaultAvoidRepeats,
			},
			expectedErr: ErrLengthMustBeGreaterThanZero,
		},
		{
			name: "invalid avoid repeats - negative",
			options: &CommandLineOptions{
				length:       DefaultPasswordLength,
				avoidRepeats: -1,
			},
			expectedErr: ErrAvoidRepeatsMustBeEqualOrGreaterThanZero,
		},
		{
			name: "multiple invalid arguments - length checked first",
			options: &CommandLineOptions{
				length:       -1,
				avoidRepeats: -1,
			},
			expectedErr: ErrLengthMustBeGreaterThanZero,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommandLineOptionsToPasswordGeneratorOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  *CommandLineOptions
		expected *PasswordGeneratorOptions
	}{
		{
			name: "default values conversion",
			options: &CommandLineOptions{
				length:       DefaultPasswordLength,
				lowercase:    true,
				uppercase:    true,
				numbers:      true,
				symbols:      false,
				custom:       "",
				avoidRepeats: DefaultAvoidRepeats,
				qrOutput:     false,
			},
			expected: &PasswordGeneratorOptions{
				Length:       DefaultPasswordLength,
				Lowercase:    true,
				Uppercase:    true,
				Numbers:      true,
				Symbols:      false,
				Custom:       "",
				AvoidRepeats: DefaultAvoidRepeats,
				QrCode:       false,
			},
		},
		{
			name: "custom values conversion",
			options: &CommandLineOptions{
				length:       20,
				lowercase:    false,
				uppercase:    false,
				numbers:      false,
				symbols:      true,
				custom:       "abc123!@#",
				avoidRepeats: 5,
				qrOutput:     true,
			},
			expected: &PasswordGeneratorOptions{
				Length:       20,
				Lowercase:    false,
				Uppercase:    false,
				Numbers:      false,
				Symbols:      true,
				Custom:       "abc123!@#",
				AvoidRepeats: 5,
				QrCode:       true,
			},
		},
		{
			name: "edge values conversion",
			options: &CommandLineOptions{
				length:       1,
				lowercase:    true,
				uppercase:    false,
				numbers:      true,
				symbols:      false,
				custom:       "",
				avoidRepeats: 0,
				qrOutput:     false,
			},
			expected: &PasswordGeneratorOptions{
				Length:       1,
				Lowercase:    true,
				Uppercase:    false,
				Numbers:      true,
				Symbols:      false,
				Custom:       "",
				AvoidRepeats: 0,
				QrCode:       false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.ToPasswordGeneratorOptions()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCommandLineOptionsToPasswordGeneratorOptionsFieldMapping(t *testing.T) {
	options := &CommandLineOptions{qrOutput: true}
	result := options.ToPasswordGeneratorOptions()
	assert.True(t, result.QrCode)

	options = &CommandLineOptions{qrOutput: false}
	result = options.ToPasswordGeneratorOptions()
	assert.False(t, result.QrCode)
}

func TestInitializeCommandLineSuccess(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name     string
		args     []string
		validate func(t *testing.T, options *CommandLineOptions)
	}{
		{
			name: "default arguments",
			args: []string{"testprogram"},
			validate: func(t *testing.T, options *CommandLineOptions) {
				assert.Equal(t, DefaultPasswordLength, options.length)
				assert.True(t, options.lowercase)
				assert.True(t, options.uppercase)
				assert.True(t, options.numbers)
				assert.False(t, options.symbols)
			},
		},
		{
			name: "valid custom arguments",
			args: []string{"testprogram", "-l", "16", "--symbols=true", "-a", "2"},
			validate: func(t *testing.T, options *CommandLineOptions) {
				assert.Equal(t, 16, options.length)
				assert.True(t, options.symbols)
				assert.Equal(t, 2, options.avoidRepeats)
			},
		},
		{
			name: "complex valid arguments",
			args: []string{
				"testprogram",
				"--length", "30",
				"--lowercase=false",
				"--uppercase=true",
				"--numbers=false",
				"--symbols=true",
				"--custom", "specialChars!@#",
				"--avoid-repeats", "10",
				"--qr=true",
			},
			validate: func(t *testing.T, options *CommandLineOptions) {
				assert.Equal(t, 30, options.length)
				assert.False(t, options.lowercase)
				assert.True(t, options.uppercase)
				assert.False(t, options.numbers)
				assert.True(t, options.symbols)
				assert.Equal(t, "specialChars!@#", options.custom)
				assert.Equal(t, 10, options.avoidRepeats)
				assert.True(t, options.qrOutput)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			options, err := InitializeCommandLine()

			require.NoError(t, err)
			require.NotNil(t, options)
			tt.validate(t, options)
		})
	}
}

func TestInitializeCommandLineValidationErrors(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name        string
		args        []string
		expectedErr error
	}{
		{
			name:        "invalid length - zero",
			args:        []string{"testprogram", "-l", "0"},
			expectedErr: ErrLengthMustBeGreaterThanZero,
		},
		{
			name:        "invalid length - negative",
			args:        []string{"testprogram", "--length", "-5"},
			expectedErr: ErrLengthMustBeGreaterThanZero,
		},
		{
			name:        "invalid avoid repeats - negative",
			args:        []string{"testprogram", "-a", "-1"},
			expectedErr: ErrAvoidRepeatsMustBeEqualOrGreaterThanZero,
		},
		{
			name:        "invalid avoid repeats - negative long flag",
			args:        []string{"testprogram", "--avoid-repeats", "-3"},
			expectedErr: ErrAvoidRepeatsMustBeEqualOrGreaterThanZero,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			options, err := InitializeCommandLine()

			assert.Error(t, err)
			assert.Nil(t, options)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestCommandLineParserParseInvalidFlags(t *testing.T) {
	parser := NewCommandLineParser()

	output := captureStderr(func() {
		parser.Parse([]string{"--invalid-flag"})
	})

	assert.NotEmpty(t, output)
	assert.Contains(t, output, "flag provided but not defined")
}

func TestIntegrationParseValidateConvert(t *testing.T) {
	parser := NewCommandLineParser()

	args := []string{
		"--length", "16",
		"--symbols=true",
		"--custom", "mycharset",
		"--avoid-repeats", "2",
		"--qr=true",
	}

	options, err := parser.Parse(args)
	require.NoError(t, err)
	require.NotNil(t, options)

	err = options.Validate()
	assert.NoError(t, err)

	genOptions := options.ToPasswordGeneratorOptions()
	require.NotNil(t, genOptions)

	assert.Equal(t, 16, genOptions.Length)
	assert.True(t, genOptions.Lowercase)
	assert.True(t, genOptions.Uppercase)
	assert.True(t, genOptions.Numbers)
	assert.True(t, genOptions.Symbols)
	assert.Equal(t, "mycharset", genOptions.Custom)
	assert.Equal(t, 2, genOptions.AvoidRepeats)
	assert.True(t, genOptions.QrCode)
}

func TestCommandLineConstants(t *testing.T) {
	assert.Equal(t, 12, DefaultPasswordLength)
	assert.Equal(t, 1, DefaultAvoidRepeats)
}

func BenchmarkCommandLineParserParse(b *testing.B) {
	args := []string{"-l", "12", "--symbols=true"}

	for b.Loop() {
		parser := NewCommandLineParser()
		parser.Parse(args)
	}
}

func BenchmarkGeneratePasswordWitAllCharset(b *testing.B) {
	os.Args = []string{"passgen", "-l", "50", "-a", "3", "-L", "-U", "-N", "-S"}
	cmd, _ := InitializeCommandLine()

	b.ResetTimer()
	for b.Loop() {
		GeneratePassword(*cmd.ToPasswordGeneratorOptions())
	}
}
