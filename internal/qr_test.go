package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQrCode(t *testing.T) {
	t.Run("should create QR code with valid content", func(t *testing.T) {
		content := "Hello, World!"
		margin := 2

		qr, err := NewQrCode(content, margin)

		require.NoError(t, err)
		assert.NotNil(t, qr)
		assert.Equal(t, margin, qr.margin)
		assert.NotNil(t, qr.data)
	})

	t.Run("should return error for empty content", func(t *testing.T) {
		margin := 2

		qr, err := NewQrCode("", margin)

		assert.Error(t, err)
		assert.Equal(t, ErrQrEmptyContent, err)
		assert.Nil(t, qr)
	})

	t.Run("should create QR code with zero margin", func(t *testing.T) {
		content := "Test"
		margin := 0

		qr, err := NewQrCode(content, margin)

		require.NoError(t, err)
		assert.NotNil(t, qr)
		assert.Equal(t, margin, qr.margin)
	})

	t.Run("should create QR code with large margin", func(t *testing.T) {
		content := "Test"
		margin := 10

		qr, err := NewQrCode(content, margin)

		require.NoError(t, err)
		assert.NotNil(t, qr)
		assert.Equal(t, margin, qr.margin)
	})

	t.Run("should handle special characters in content", func(t *testing.T) {
		content := "Test with special chars: !@#$%^&*()"
		margin := 2

		qr, err := NewQrCode(content, margin)

		require.NoError(t, err)
		assert.NotNil(t, qr)
		assert.Equal(t, margin, qr.margin)
	})

	t.Run("should handle unicode content", func(t *testing.T) {
		content := "Unicode test: 你好世界"
		margin := 2

		qr, err := NewQrCode(content, margin)

		require.NoError(t, err)
		assert.NotNil(t, qr)
		assert.Equal(t, margin, qr.margin)
	})
}

func TestQrCodeGenerateAnisUtf8i(t *testing.T) {
	t.Run("should generate UTF-8 output for simple content", func(t *testing.T) {
		qr, err := NewQrCode("Test", 2)
		require.NoError(t, err)

		output := qr.GenerateAnisUtf8i()

		assert.NotEmpty(t, output)
		assert.Contains(t, output, "\033[40;37;1m")
		assert.Contains(t, output, "\033[0m")
		assert.Contains(t, output, "\n")
	})

	t.Run("should generate output with zero margin", func(t *testing.T) {
		qr, err := NewQrCode("Test", 0)
		require.NoError(t, err)

		output := qr.GenerateAnisUtf8i()

		assert.NotEmpty(t, output)
		assert.Contains(t, output, "\033[40;37;1m")
		assert.Contains(t, output, "\033[0m")
	})

	t.Run("should generate output with large margin", func(t *testing.T) {
		qr, err := NewQrCode("A", 8)
		require.NoError(t, err)

		output := qr.GenerateAnisUtf8i()

		assert.NotEmpty(t, output)
		lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
		assert.Greater(t, len(lines), 10)
	})

	t.Run("should contain expected UTF-8 block characters", func(t *testing.T) {
		qr, err := NewQrCode("Hello", 1)
		require.NoError(t, err)

		output := qr.GenerateAnisUtf8i()

		assert.Condition(t, func() (success bool) {
			success = strings.Contains(output, "\342\226\204") ||
				strings.Contains(output, "\342\226\200") ||
				strings.Contains(output, "\342\226\210") ||
				strings.Contains(output, " ")
			return
		})
	})

	t.Run("should generate consistent output for same input", func(t *testing.T) {
		qr, err := NewQrCode("Consistent", 2)
		require.NoError(t, err)

		output1 := qr.GenerateAnisUtf8i()
		output2 := qr.GenerateAnisUtf8i()

		assert.Equal(t, output1, output2, "Multiple calls should produce identical output")
	})

	t.Run("should handle odd QR code dimensions", func(t *testing.T) {
		qr, err := NewQrCode("A", 1)
		require.NoError(t, err)

		output := qr.GenerateAnisUtf8i()

		assert.NotEmpty(t, output)
		lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
		assert.NotEmpty(t, lines)
	})

	t.Run("should properly format lines with color codes", func(t *testing.T) {
		qr, err := NewQrCode("Format", 1)
		require.NoError(t, err)

		output := qr.GenerateAnisUtf8i()

		for line := range strings.SplitSeq(output, "\n") {
			if len(line) > 0 {
				assert.Contains(t, line, "\033[40;37;1m", "Each content line should start with white color code")
				assert.Contains(t, line, "\033[0m", "Each content line should end with reset code")
			}
		}
	})
}

func TestQrCode_writeUTF8Margin(t *testing.T) {
	t.Run("should write margin correctly", func(t *testing.T) {
		qr, err := NewQrCode("Test", 4)
		require.NoError(t, err)

		var output strings.Builder
		white := "\033[40;37;1m"
		reset := "\033[0m"
		full := "\342\226\210"
		realwidth := 10

		qr.writeUTF8Margin(&output, realwidth, white, reset, full)

		result := output.String()
		assert.NotEmpty(t, result)
		assert.Contains(t, result, white)
		assert.Contains(t, result, reset)
		assert.Contains(t, result, full)

		lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
		assert.Len(t, lines, 2)
	})

	t.Run("should handle zero margin", func(t *testing.T) {
		qr, err := NewQrCode("Test", 0)
		require.NoError(t, err)

		var output strings.Builder
		white := "\033[40;37;1m"
		reset := "\033[0m"
		full := "\342\226\210"
		realwidth := 10

		qr.writeUTF8Margin(&output, realwidth, white, reset, full)

		result := output.String()
		assert.Empty(t, result)
	})

	t.Run("should handle odd margin", func(t *testing.T) {
		qr, err := NewQrCode("Test", 3)
		require.NoError(t, err)

		var output strings.Builder
		white := "\033[40;37;1m"
		reset := "\033[0m"
		full := "\342\226\210"
		realwidth := 10

		qr.writeUTF8Margin(&output, realwidth, white, reset, full)

		result := output.String()
		lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
		assert.Len(t, lines, 1)
	})
}

func TestErrQrEmptyContent(t *testing.T) {
	t.Run("should have correct error message", func(t *testing.T) {
		expectedMessage := "The content for QR generation should not be empty."
		assert.Equal(t, expectedMessage, ErrQrEmptyContent.Error())
	})
}

func TestQrCodeIntegration(t *testing.T) {
	t.Run("should create and generate QR code successfully", func(t *testing.T) {
		content := "Integration Test"
		margin := 3

		qr, err := NewQrCode(content, margin)
		require.NoError(t, err)
		require.NotNil(t, qr)

		output := qr.GenerateAnisUtf8i()
		assert.NotEmpty(t, output)

		lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
		assert.Greater(t, len(lines), margin/2, "Should have at least margin lines")

		for i, line := range lines {
			if len(line) > 0 {
				assert.Contains(t, line, "\033[40;37;1m", "Line %d should contain white color code", i)
				assert.Contains(t, line, "\033[0m", "Line %d should contain reset code", i)
			}
		}
	})
}

func BenchmarkNewQrCode(b *testing.B) {
	content := "Benchmark Test Content"
	margin := 2

	for b.Loop() {
		qr, err := NewQrCode(content, margin)
		if err != nil {
			b.Fatal(err)
		}
		_ = qr
	}
}

func BenchmarkGenerateAnisUtf8i(b *testing.B) {
	qr, err := NewQrCode("Benchmark Test", 2)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		output := qr.GenerateAnisUtf8i()
		_ = output
	}
}
