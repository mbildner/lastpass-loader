package environment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const PATH_PREFIX = "PATH="
const PINENTRY_BIN_NAME = "pinentry"

func EnvWithoutPinEntryInPath(logger *log.Logger) ([]string, string, error) {
	var safeToIncludeEnvEntries []string

	var pathEnvEntry string

	for _, envLine := range os.Environ() {
		if strings.HasPrefix(envLine, PATH_PREFIX) {
			pathEnvEntry = envLine
		} else {
			safeToIncludeEnvEntries = append(safeToIncludeEnvEntries, envLine)
		}
	}

	cmd := exec.Command("which", PINENTRY_BIN_NAME)

	var pinentryDir string

	pinentryLocation, err := cmd.Output()
	if err != nil {
		return os.Environ(), "", nil

	}

	pinentryDir, _ = filepath.Split(string(pinentryLocation))
	pinentryDir = strings.TrimRight(pinentryDir, string(os.PathSeparator))

	if !strings.Contains(pathEnvEntry, pinentryDir) {
		errorMessage := "pinentry is not in path despite being found in `which`, this is a problem, I am bailing"
		logger.Println(errorMessage)
		return nil, "", errors.New(errorMessage)
	}

	var pinentryOptions []os.FileInfo
	var safeToInclude []os.FileInfo

	files, err := ioutil.ReadDir(pinentryDir)
	if err != nil {
		logger.Println("couldn't read pinentryDir, there is probably a problem")
		return nil, "", errors.New("couldn't read pinentryDir, there is probably a problem")
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), "pinentry") {
			pinentryOptions = append(pinentryOptions, f)
		} else {
			safeToInclude = append(safeToInclude, f)
		}
	}

	binariesToAdd := []string{}

	replacementPathDir, err := ioutil.TempDir("/tmp", "temp-fake-path-dir")
	if err != nil {
		fmt.Println("could not create a temp dir")
		return nil, "", errors.New("could not create a temp dir")
	}

	for _, b := range safeToInclude {
		binaryPath := filepath.Join(pinentryDir, b.Name())
		binariesToAdd = append(binariesToAdd, binaryPath)

		destinationPath := filepath.Join(replacementPathDir, b.Name())

		os.Symlink(binaryPath, destinationPath)
	}

	replacementBinaryPath := replacementPathDir

	replacementPath := strings.Replace(
		pathEnvEntry,
		pinentryDir,
		replacementBinaryPath,
		-1,
	)

	newEnv := append(safeToIncludeEnvEntries, replacementPath)
	return newEnv, replacementPathDir, nil
}
