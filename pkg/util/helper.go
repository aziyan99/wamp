package util

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

func PrintLog(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s]: ", prefix), log.LstdFlags)
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func DirExists(path string) (bool, error) {
	var err error

	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func NormalizePath(original string) string {
	return strings.ReplaceAll(original, "\\", "/")
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	srcStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}
	err = os.Chmod(dst, srcStat.Mode())
	if err != nil {
		return fmt.Errorf("failed to set destination file permissions: %w", err)
	}

	return nil
}

func CleanDirs(dirs ...string) error {
	for i := range dirs {
		if err := os.RemoveAll(dirs[i]); err != nil {
			return err
		}
	}

	return nil
}
