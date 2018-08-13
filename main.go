package main

import (
	"fmt"
	"log"
	"os"

	creds "github.com/mbildner/lastpass-loader/internal/pkg/credentials"
	env "github.com/mbildner/lastpass-loader/internal/pkg/environment"
	lpass "github.com/mbildner/lastpass-loader/internal/pkg/lastpass"
)

func main() {
	logger := log.New(os.Stderr, "", 0)

	lastPassNotePath := os.Getenv("LASTPASS_LOADER_NOTE_NAME")
	if len(lastPassNotePath) == 0 {
		logger.Println("please set LASTPASS_LOADER_NOTE_NAME in your environment before running")
		os.Exit(1)
	}

	if lpass.IsLpassLoggedIn() {
		logger.Println("logged into lastpass, retrieving env vars")
		err := lpass.ReadEnvNote(logger, lastPassNotePath)
		if err != nil {
			logger.Println("could not load data from lastpass")
			os.Exit(1)
		}

		logger.Println("success")
		os.Exit(0)
	}

	if !lpass.IsInPathInEnv(os.Environ()) {
		logger.Println("lpass not found in path, exiting")
		os.Exit(1)
	}

	logger.Println("need to log into lastpass")
	email, err := creds.ReadEmail(logger)
	if err != nil {
		logger.Println("could not read email")
		os.Exit(1)
	}

	password, err := creds.ReadPassword(logger)
	if err != nil {
		logger.Println("could not read password")
		os.Exit(1)
	}

	logger.Print("connecting to lastpass... ")
	newEnv, dirForPath, err := env.EnvWithoutPinEntryInPath(logger)
	go func(dir string) {
		defer os.RemoveAll(dir)
	}(dirForPath)

	if err != nil {
		fmt.Println("could not create a valid env")
		os.Exit(1)
	}

	logger.Println("waiting for 2 factor auth (check your phone)")
	err = lpass.Login(logger, newEnv, email, password)
	if err != nil {
		logger.Println("could not log in to lastpass")
		os.Exit(1)
	}

	logger.Println("login successful, reading data")
	err = lpass.ReadEnvNote(logger, lastPassNotePath)
	if err != nil {
		logger.Println("could not load data from lastpass")
		os.Exit(1)
	}

	logger.Println("done")
}
