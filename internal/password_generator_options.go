package internal

type PasswordGeneratorOptions struct {
	Length         int
	Lowercase      bool
	Uppercase      bool
	Numbers        bool
	Symbols        bool
	Exclude        string
	ExcludeSimilar bool
	Custom         string
	AvoidRepeats   int
	QrCode         bool
}

func NewPasswordGeneratorOptions() *PasswordGeneratorOptions {
	return &PasswordGeneratorOptions{
		Length:         12,
		Lowercase:      true,
		Uppercase:      true,
		Numbers:        true,
		Symbols:        false,
		Exclude:        "",
		ExcludeSimilar: false,
		Custom:         "",
		AvoidRepeats:   1,
		QrCode:         false,
	}
}
