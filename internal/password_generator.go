package internal

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
)

var (
	ErrEmptyCharset                = errors.New("Character set is empty.")
	ErrCnnotSelectFromEmptyCharset = errors.New("Cannot select from empty charset.")
	ErrMinCnnotGreaterThanMax      = errors.New("min cannot be greater than max.")
)

func GeneratePassword(options PasswordGeneratorOptions) (string, error) {
	charset := NewCharsetBuilderFromPasswordGeneratorOptions(options)
	if charset.IsEmpty() {
		return "", ErrEmptyCharset
	}

	characters := charset.Characters()
	avoidRepeats := normalizeAvoidRepeats(options.AvoidRepeats, charset.Length())

	var password strings.Builder
	password.Grow(options.Length)

	for password.Len() < options.Length {
		character, err := selectValidPasswordChar(characters, password.String(), avoidRepeats)
		if err != nil {
			return "", err
		}
		password.WriteString(character)
	}

	return password.String(), nil
}

func normalizeAvoidRepeats(avoidRepeats, charsetLength int) int {
	if avoidRepeats >= charsetLength {
		return charsetLength - 1
	}

	if avoidRepeats <= 0 {
		return 0
	}

	return avoidRepeats
}

func selectValidPasswordChar(charset, currentPassword string, avoidRepeats int) (string, error) {
	for {
		character, err := pickRandomChar(charset)
		if err != nil {
			return "", err
		}

		if isValidPasswordChar(string(character), currentPassword, avoidRepeats) {
			return string(character), nil
		}
	}
}

func isValidPasswordChar(character, password string, avoidRepeats int) bool {
	if avoidRepeats <= 0 {
		return true
	}

	recentChars := getPasswordSuffix(password, avoidRepeats)
	return !strings.Contains(recentChars, character)
}

func getPasswordSuffix(password string, suffixLength int) string {
	if suffixLength >= len(password) {
		return password
	}

	return password[len(password)-suffixLength:]
}

func pickRandomChar(charset string) (rune, error) {
	if charset == "" {
		return 0, ErrCnnotSelectFromEmptyCharset
	}

	randomIndex, err := secureRandomInt(0, len(charset)-1)
	if err != nil {
		return 0, err
	}

	return []rune(charset[randomIndex : randomIndex+1])[0], nil
}

func secureRandomInt(min, max int) (int, error) {
	if min > max {
		return 0, ErrMinCnnotGreaterThanMax
	}

	if min == max {
		return min, nil
	}

	diff := max - min + 1
	n, err := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + min, nil
}
