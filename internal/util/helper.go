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

func CornConfig() string {
	return `
# This is the configuration file for the corn Daemon.
# Each line defines a job, and has the following structure:
# 
# .---------------- minute (0 - 59)
# |  .------------- hour (0 - 23)
# |  |  .---------- day of month (1 - 31)
# |  |  |  .------- month (1 - 12)
# |  |  |  |  .---- day of week (0 - 6) (Sunday is 0)
# |  |  |  |  |
# *  *  *  *  * command to execute
# 
# - Use '*' for 'every'.
# - Use comma-separated values for specific times (e.g., "0,15,30,45").
# - Does not support ranges (e.g., "8-17") or steps (e.g., "*/15").
# 

# Job 1: Runs every minute. Pings localhost once.
* * * * * ping 127.0.0.1 -n 1

# Job 2: Every 5 minutes, append the current date and a message to a log file. This demonstrates using shell redirection.
0,5,10,15,20,25,30,35,40,45,50,55 * * * * echo "Here is corn" >> corn_output.log

# Job 3: This job is disabled because it is commented out.
# 30 14 * * * echo "This is a disabled job"
	`
}
