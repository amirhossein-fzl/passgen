package internal

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ProgramName           = "passgen"
	DefaultPasswordLength = 12
	DefaultAvoidRepeats   = 1
)

type CommandLineOptions struct {
	length       int
	lowercase    bool
	uppercase    bool
	numbers      bool
	symbols      bool
	custom       string
	avoidRepeats int
	qrOutput     bool
}

type CommandLineParser struct {
	flagSet *flag.FlagSet
}

func NewCommandLineParser() *CommandLineParser {
	return &CommandLineParser{
		flagSet: flag.NewFlagSet(ProgramName, flag.ContinueOnError),
	}
}

func (p *CommandLineParser) Parse(args []string) (*CommandLineOptions, error) {
	length := p.flagSet.Int("l", DefaultPasswordLength, "")
	p.flagSet.IntVar(length, "length", DefaultPasswordLength, "")

	lowercase := p.flagSet.Bool("L", true, "")
	p.flagSet.BoolVar(lowercase, "lowercase", true, "")

	uppercase := p.flagSet.Bool("U", true, "")
	p.flagSet.BoolVar(uppercase, "uppercase", true, "")

	numbers := p.flagSet.Bool("N", true, "")
	p.flagSet.BoolVar(numbers, "numbers", true, "")

	symbols := p.flagSet.Bool("S", false, "")
	p.flagSet.BoolVar(symbols, "symbols", false, "")

	custom := p.flagSet.String("C", "", "")
	p.flagSet.StringVar(custom, "custom", "", "")

	avoidRepeats := p.flagSet.Int("a", DefaultAvoidRepeats, "")
	p.flagSet.IntVar(avoidRepeats, "avoid-repeats", DefaultAvoidRepeats, "")

	qrOutput := p.flagSet.Bool("q", false, "")
	p.flagSet.BoolVar(qrOutput, "qr", false, "")

	p.flagSet.Usage = p.printUsage

	p.flagSet.Parse(args)

	options := &CommandLineOptions{
		length:       *length,
		lowercase:    *lowercase,
		uppercase:    *uppercase,
		numbers:      *numbers,
		symbols:      *symbols,
		custom:       *custom,
		avoidRepeats: *avoidRepeats,
		qrOutput:     *qrOutput,
	}

	return options, nil
}

func InitializeCommandLine() (*CommandLineOptions, error) {
	parser := NewCommandLineParser()

	args := os.Args[1:]

	options, err := parser.Parse(args)
	if err != nil {
		return nil, err
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	return options, nil
}

func (p *CommandLineParser) printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "A secure password generator with customizable options.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")

	fmt.Fprintf(os.Stderr, "  -l, --length <length>\t\t\tPassword length (default: 12)\n")
	fmt.Fprintf(os.Stderr, "  -L, --lowercase\t\t\tInclude lowercase letters (a-z) (default: true)\n")
	fmt.Fprintf(os.Stderr, "  -U, --uppercase\t\t\tInclude uppercase letters (A-Z) (default: true)\n")
	fmt.Fprintf(os.Stderr, "  -N, --numbers\t\t\t\tInclude numbers (0-9) (default: true)\n")
	fmt.Fprintf(os.Stderr, "  -S, --symbols\t\t\t\tInclude symbols (!@#$%%^&* etc.)\n")
	fmt.Fprintf(os.Stderr, "  -C, --custom <custom>\t\t\tCustom character set to use\n")
	fmt.Fprintf(os.Stderr, "  -a, --avoid-repeats <avoid-repeats>\tNumber of last characters that shouldn't repeat (default: 1)\n")
	fmt.Fprintf(os.Stderr, "  -q, --qr\t\t\t\tGenerate QR code output in ANSI UTF-8 format\n")

	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s -l 16 --uppercase --numbers --symbols \n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s --length 12 --uppercase --numbers\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s -l 20 --custom \"abcdef123456!@#\"\n", filepath.Base(os.Args[0]))
}

func (c *CommandLineOptions) ToPasswordGeneratorOptions() *PasswordGeneratorOptions {
	return &PasswordGeneratorOptions{
		Length:       c.length,
		Lowercase:    c.lowercase,
		Uppercase:    c.uppercase,
		Numbers:      c.numbers,
		Symbols:      c.symbols,
		Custom:       c.custom,
		AvoidRepeats: c.avoidRepeats,
		QrCode:       c.qrOutput,
	}
}

func (c *CommandLineOptions) Validate() error {
	if c.length <= 0 {
		return ErrLengthMustBeGreaterThanZero
	}

	if c.avoidRepeats < 0 {
		return ErrAvoidRepeatsMustBeEqualOrGreaterThanZero
	}

	return nil
}
