package envs

import (
	"fmt"
	"os"
)

func GetCommonDirFromEnv() (string, error) {
	dir := os.Getenv("COMMON_DIR")
	if dir == "" {
		return "", fmt.Errorf("common dir path not set")
	}

	info, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("get common dir info: %w", err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("common dir is not a dir")
	}

	return dir, nil
}

func GetPwd() string {
	return os.Getenv("PWD")
}
