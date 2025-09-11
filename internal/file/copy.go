package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyDir(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("get info about src dir: %w", err)
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("create destination folder: %w", err)
	}

	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("get relative path: %w", err)
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			if err := os.MkdirAll(destPath, info.Mode()); err != nil {
				return fmt.Errorf("create dst dir: %w", err)
			}

			return nil
		} else {
			return CopyFile(path, destPath, info.Mode())
		}
	}); err != nil {
		return fmt.Errorf("copy files from %q to %q: %w", src, dst, err)
	}

	return nil
}

func CopyFile(srcFile, dstFile string, perm os.FileMode) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("open src file: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return fmt.Errorf("open dst file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy data from src to dst file: %w", err)
	}

	return nil
}
