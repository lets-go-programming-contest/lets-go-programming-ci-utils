package envs

import (
	"errors"
	"fmt"
	"os"
)

var (
	errEnvNotSet       = errors.New("env value not set")
	errCommonDirIsFile = errors.New("common dir path target is a files")
)

func GetCommonDirFromEnv() (string, error) {
	dir := os.Getenv("COMMON_DIR")
	if dir == "" {
		return "", fmt.Errorf("%w: %q", errEnvNotSet, "COMMON_DIR")
	}

	info, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("get common dir info: %w", err)
	}

	if !info.IsDir() {
		return "", errCommonDirIsFile
	}

	return dir, nil
}

func GetPwd() string {
	return os.Getenv("PWD")
}
