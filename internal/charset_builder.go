package internal

const (
	UppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	NumberChars    = "0123456789"
	SymbolChars    = "`~!@#$%^&*()_-+=\\|[]'\";:/?<>.,"
)

type CharsetBuilder struct {
	characters string
}

func NewCharsetBuilder() *CharsetBuilder {
	return &CharsetBuilder{}
}

func (c *CharsetBuilder) Characters() string {
	return c.characters
}

func (c *CharsetBuilder) WithUppercase() *CharsetBuilder {
	c.characters += UppercaseChars
	return c
}

func (c *CharsetBuilder) WithLowercase() *CharsetBuilder {
	c.characters += LowercaseChars
	return c
}

func (c *CharsetBuilder) WithNumbers() *CharsetBuilder {
	c.characters += NumberChars
	return c
}

func (c *CharsetBuilder) WithSymbols() *CharsetBuilder {
	c.characters += SymbolChars
	return c
}

func (c *CharsetBuilder) WithCustom(characters string) *CharsetBuilder {
	if characters != "" {
		c.characters += characters
	}
	return c
}

func (c *CharsetBuilder) Reset() *CharsetBuilder {
	c.characters = ""
	return c
}

func (c *CharsetBuilder) Length() int {
	return len(c.characters)
}

func (c *CharsetBuilder) IsEmpty() bool {
	return c.characters == ""
}
