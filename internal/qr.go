package internal

import (
	"errors"
	"strings"

	"github.com/skip2/go-qrcode"
)

var (
	ErrQrEmptyContent = errors.New("The content for QR generation should not be empty.")
)

type QrCode struct {
	margin int
	data   *qrcode.QRCode
}

func NewQrCode(content string, margin int) (*QrCode, error) {
	if content == "" {
		return nil, ErrQrEmptyContent
	}

	qr, err := qrcode.New(content, qrcode.Highest)

	if err != nil {
		return nil, err
	}

	return &QrCode{margin: margin, data: qr}, nil
}

func (qr *QrCode) GenerateAnisUtf8i() string {

	var output strings.Builder

	white := "\033[40;37;1m"
	reset := "\033[0m"

	empty := " "
	lowhalf := "\342\226\204"
	uphalf := "\342\226\200"
	full := "\342\226\210"

	bitmap := qr.data.Bitmap()
	qrWidth := len(bitmap)

	realwidth := qrWidth + qr.margin*2

	// Top margin
	qr.writeUTF8Margin(&output, realwidth, white, reset, full)

	// Data rows - process two rows at a time for half-block characters
	for y := 0; y < qrWidth; y += 2 {
		output.WriteString(white)

		// Left margin
		for x := 0; x < qr.margin; x++ {
			output.WriteString(full)
		}

		// QR data
		for x := range qrWidth {
			// Get current row
			row1Dark := bitmap[y][x]
			var row2Dark bool

			// Get next row if it exists
			if y < qrWidth-1 {
				row2Dark = bitmap[y+1][x]
			}

			// Determine which character to output based on the two pixels
			if row1Dark {
				if y < qrWidth-1 && row2Dark {
					output.WriteString(empty) // Both dark
				} else {
					output.WriteString(lowhalf) // Top dark, bottom light
				}
			} else if y < qrWidth-1 && row2Dark {
				output.WriteString(uphalf) // Top light, bottom dark
			} else {
				output.WriteString(full) // Both light (or bottom row)
			}
		}

		// Right margin
		for x := 0; x < qr.margin; x++ {
			output.WriteString(full)
		}

		output.WriteString(reset)
		output.WriteString("\n")
	}

	// Bottom margin
	qr.writeUTF8Margin(&output, realwidth, white, reset, full)

	return output.String()
}

func (qr *QrCode) writeUTF8Margin(qrOutput *strings.Builder, realwidth int, white, reset, full string) {
	for y := 0; y < qr.margin/2; y++ {
		qrOutput.WriteString(white)
		for range realwidth {
			qrOutput.WriteString(full)
		}
		qrOutput.WriteString(reset)
		qrOutput.WriteString("\n")
	}
}
