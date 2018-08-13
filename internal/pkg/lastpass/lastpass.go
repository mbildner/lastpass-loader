package lastpass

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func IsLpassLoggedIn() bool {
	cmd := exec.Command("lpass", "status")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(out), "Logged in as")
}

func IsInPathInEnv(env []string) bool {
	cmd := exec.Command("which", "lpass")

	cmd.Env = env
	_, err := cmd.Output()

	return err == nil
}

func ReadEnvNote(logger *log.Logger, notePath string) error {
	readLpassNoteCmd := exec.Command(
		"lpass",
		"show",
		"--name",
		"--notes",
		notePath,
	)
	lpassNoteBytes, err := readLpassNoteCmd.Output()
	if err != nil {
		logger.Println("could not read from lastpass")
		return err
	}

	var envVarsToLoad map[string]string

	err = json.Unmarshal(lpassNoteBytes, &envVarsToLoad)
	if err != nil {
		return errors.New("lastpass note is not valid json")
	}

	var varsToLoad []string

	for k, v := range envVarsToLoad {
		varsToLoad = append(varsToLoad, fmt.Sprintf("export %v=\"%v\"", k, v))
	}

	exportString := strings.Join(varsToLoad, "\n")

	fmt.Println(exportString)

	return nil
}

func Login(logger *log.Logger, environment []string, email string, password string) error {
	lpassCmd := exec.Command("lpass", "login", email)
	lpassCmd.Env = environment
	stdin, err := lpassCmd.StdinPipe()

	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, password+"\n")
	}()

	stdout, err := lpassCmd.StdoutPipe()
	if err != nil {
		return errors.New("could not access stdout from lpass login command")
	}

	go func() {
		defer stdout.Close()
		ioutil.ReadAll(stdout)
	}()

	err = lpassCmd.Run()

	if err != nil {
		return err
	}

	return nil
}
