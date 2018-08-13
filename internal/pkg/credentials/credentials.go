package lastpass

import (
	"bufio"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strings"
	"syscall"
)

const PASSWORD_PROMPT = "Enter password:"
const EMAIL_PROMPT = "Enter email address:"

func ReadPassword(logger *log.Logger) (string, error) {
	logger.Println(PASSWORD_PROMPT)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	password := string(bytePassword)

	return strings.Trim(password, "\n"), nil
}

func ReadEmail(logger *log.Logger) (string, error) {
	var email string

	reader := bufio.NewReader(os.Stdin)
	logger.Println(EMAIL_PROMPT)
	email, _ = reader.ReadString('\n')

	return strings.Trim(email, "\n"), nil
}
