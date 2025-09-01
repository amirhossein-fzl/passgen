# PassGen 🔐
❗️ This project is still under development and has not been released, and parts of it may still be incomplete or not work properly.

A secure, customizable command-line password generator built with Go. Generate strong passwords with flexible character sets, avoid character repetition, and even output QR codes for easy mobile transfer.

## Features

🔤 **Character Set Control**: Include/exclude lowercase, uppercase, numbers, and symbols  
🎨 **Custom Character Sets**: Define your own character pool  
🔄 **Smart Repetition Avoidance**: Prevent consecutive character repetition  
📱 **QR Code Output**: Generate ANSI UTF-8 QR codes for easy mobile scanning  
🛡️ **Cryptographically Secure**: Uses secure random number generation  

## Installation

### From Source
```bash
git clone https://github.com/amirhossein-fzl/passgen.git
cd passgen
go build -o passgen cmd/main.go
```

## Quick Start

Generate a default 12-character password with lowercase, uppercase, and numbers:
```bash
passgen
```

Generate a 16-character password with all character types:
```bash
passgen -l 16 -S
```

## Usage

```
passgen [options]
```

### Options

| Short | Long              | Description                                     | Default |
| ----- | ----------------- | ----------------------------------------------- | ------- |
| `-l`  | `--length`        | Password length                                 | `12`    |
| `-L`  | `--lowercase`     | Include lowercase letters (a-z)                 | `true`  |
| `-U`  | `--uppercase`     | Include uppercase letters (A-Z)                 | `true`  |
| `-N`  | `--numbers`       | Include numbers (0-9)                           | `true`  |
| `-S`  | `--symbols`       | Include symbols (!@#$%^&* etc.)                 | `false` |
| `-C`  | `--custom`        | Custom character set to use                     | `""`    |
| `-a`  | `--avoid-repeats` | Number of last characters that shouldn't repeat | `1`     |
| `-q`  | `--qr`            | Generate QR code output in ANSI UTF-8 format    | `false` |

### Character Sets

By default, PassGen includes:
- **Lowercase**: `abcdefghijklmnopqrstuvwxyz`
- **Uppercase**: `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
- **Numbers**: `0123456789`
- **Symbols**: `!@#$%^&*()_+-=[]{}|;:,.<>?`

## Examples

### Basic Usage

Generate a simple 12-character password:
```bash
passgen
# Output: aB3kL9mX2wQ1
```

### Custom Length
```bash
passgen -l 20
# Output: aNaNYiQSO62KUcZbpios
```

### Include Symbols
```bash
passgen -l 16 -S
# Output: d8fP.|#<'I;<cZpQ
```

### Only Numbers and Uppercase
```bash
passgen -l 10 --lowercase=false -U -N
# Output: BGXRH7624Y
```

### Custom Character Set
Use only specific characters:
```bash
passgen -l 15 -U=false -L=false -N=false --custom "abcdef123456\!@#"
# Output: 4c6!f@bf4f#c1a3
```

### Avoid Character Repetition
Prevent the last 3 characters from repeating:
```bash
passgen -l 20 -a 3
# Output: OXH7cMOJyagcCvjrcMln
```

### Generate with QR Code
Perfect for transferring passwords to mobile devices:
```bash
passgen -l 16 -S -q
# Output: 
#
# ███████████████████████████████████████
# ███████████████████████████████████████
# █████ ▄▄▄▄▄ █▄▀ ▄▄ ▀ █▄▄ ██ ▄▄▄▄▄ █████
# █████ █   █ █▄▀██▀  █ ▀█ ▄█ █   █ █████
# █████ █▄▄▄█ █▄█▄▀▀ ▄ ██ █▄█ █▄▄▄█ █████
# █████▄▄▄▄▄▄▄█ █▄▀ ▀▄█ █▄█▄█▄▄▄▄▄▄▄█████
# █████▀▀▀██ ▄██▄█▄▄██▀ █▄▄ ▀ █ ▀▄▀▄█████
# ██████▀▀ █▄▄  ▀██ ▄█  ▀▀█▄▀█▄  ██▀█████
# ██████ █ ▄█▄██▄ ▄▀ █ ▀▀▄ ▀▀   ▄ ▀██████
# ██████▀▀▄▀█▄█▄ ▄█▀▄  █  █▄█▄▀ ▀▀███████
# █████ ▄▄▄ ▄▄ ▀▄▄ █ ▄▀▀▄▀█▀█▄▀▄ █  █████
# █████ ▄██▄▀▄ ▄ ▀▀▄ ▄ ▄ █▀ █▄ █▄▄█▄█████
# █████▄████▄▄▄▀▄▀▀▄ ▄ █▀▀▄ ▄▄▄ ▄▄▄▀█████
# █████ ▄▄▄▄▄ █▀█ ▀▄█▀▀▄▄ ▄ █▄█ ▄█▀██████
# █████ █   █ ██    ▄█▀▀▄▄▄  ▄ ▄▄▀ ██████
# █████ █▄▄▄█ ███ ▀▄ █  ▄▀ ▀▄   ▄▀▄▀█████
# █████▄▄▄▄▄▄▄███▄██▄▄▄▄▄█████▄▄█████████
# ███████████████████████████████████████
# ███████████████████████████████████████
# 
# Password: ~xY%tEtf%xUiJ]D9
```

## Security Features

- **Cryptographically Secure**: Uses Go's `crypto/rand` package for secure random number generation
- **No Predictable Patterns**: Avoids consecutive character repetition with configurable history
- **Flexible Character Sets**: Full control over which characters can appear in your passwords
- **Memory Safe**: Passwords are not stored or logged anywhere

## Contributing

We welcome contributions from the community! Whether it's bug fixes, new features, documentation improvements, or testing, your help makes PassGen better for everyone.

### Ways to Contribute

- 🐛 **Bug Reports**: Found an issue? Please open a GitHub issue
- 💡 **Feature Requests**: Have an idea? We'd love to hear it
- 🔧 **Code Contributions**: Submit a pull request
- 📚 **Documentation**: Help improve our docs
- 🧪 **Testing**: Help us test on different platforms

### Getting Started

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

Your contributions, no matter how big or small, genuinely make me happy and help make this tool better for everyone in the community. Thank you for being part of this project! 🙏

## License

This project is licensed under the GPL License - see the [LICENSE](https://github.com/amirhossein-fzl/passgen?tab=GPL-3.0-1-ov-file) file for details.

## Support

If you encounter any issues or have questions:

- 📖 Check this README for usage examples
- 🐛 Open an issue on GitHub
- 💬 Start a discussion in our GitHub Discussions

---

**Made with ❤️ for developers who care about security**