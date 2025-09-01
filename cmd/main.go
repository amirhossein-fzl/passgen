package main

import (
	"amirhossein-fzl/passgen/internal"
	"fmt"
	"os"
)

func main() {
	cmd, err := internal.InitializeCommandLine()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(64)
	}

	err = cmd.Validate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(64)
	}

	password, err := internal.GeneratePassword(*cmd.ToPasswordGeneratorOptions())

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(74)
	}

	fmt.Println(password)
}
