package internal

import "errors"

var (
	ErrLengthMustBeGreaterThanZero              = errors.New("Length must be greater than 0.")
	ErrAvoidRepeatsMustBeEqualOrGreaterThanZero = errors.New("Avoid repeats must be greater than or equal to 0.")
)

type PasswordGeneratorOptions struct {
	Length       int
	Lowercase    bool
	Uppercase    bool
	Numbers      bool
	Symbols      bool
	Custom       string
	AvoidRepeats int
	QrCode       bool
}

func NewPasswordGeneratorOptions() *PasswordGeneratorOptions {
	return &PasswordGeneratorOptions{
		Length:       12,
		Lowercase:    true,
		Uppercase:    true,
		Numbers:      true,
		Symbols:      false,
		Custom:       "",
		AvoidRepeats: 1,
		QrCode:       false,
	}
}

func (p *PasswordGeneratorOptions) Validate() (bool, error) {
	if p.Length <= 0 {
		return false, ErrLengthMustBeGreaterThanZero
	}

	if p.AvoidRepeats < 0 {
		return false, ErrAvoidRepeatsMustBeEqualOrGreaterThanZero
	}

	return true, nil
}
