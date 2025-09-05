package main

import (
	"amirhossein-fzl/passgen/internal"
	"fmt"
	"os"
	"slices"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if slices.Contains(os.Args, "-v") || slices.Contains(os.Args, "--version") {
		fmt.Printf("passgen version %s (%s) released at %s\n", version, commit, date)
		os.Exit(0)
	}

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

	options := *cmd.ToPasswordGeneratorOptions()
	password, err := internal.GeneratePassword(options)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(74)
	}

	if options.QrCode {
		qrcode, err := internal.NewQrCode(password, 1)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(74)
		}

		fmt.Println(qrcode.GenerateAnisUtf8i())
		fmt.Print("Password: ")
	}

	fmt.Println(password)
}
