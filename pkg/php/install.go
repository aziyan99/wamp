package php

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aziyan99/wamp/pkg/util"
)

func SetPhpIni(p string) error {
	err := util.CopyFile(path.Join(p, "php.ini-development"), path.Join(p, "php.ini"))
	if err != nil {
		return err
	}

	return nil
}

func SetExtDir(p string) error {
	_, err := os.Stat(path.Join(p, "php.ini"))
	if err != nil {
		return err
	}

	file, err := os.Open(path.Join(p, "php.ini"))
	if err != nil {
		return err
	}
	defer file.Close()

	var outputLines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		parts := strings.Fields(trimmedLine)

		if len(parts) >= 3 && parts[0] == ";extension_dir" && parts[1] == "=" && parts[2] == "\"ext\"" {
			identation := line[:strings.Index(line, trimmedLine)]
			newLine := fmt.Sprintf(`%sextension_dir = "%s"`, identation, util.NormalizePath(path.Join(p, "ext")))
			outputLines = append(outputLines, newLine)
			util.PrintLog("INFO").Printf("Found and updated PHP extension_dir.\n--- %s\n+++ %s\n", line, newLine)
		} else {
			outputLines = append(outputLines, line)
		}
	}

	updatedConf := strings.Join(outputLines, "\n")

	if err := scanner.Err(); err != nil {
		return err
	}

	err = os.WriteFile(path.Join(p, "php.ini"), []byte(updatedConf), 0644)
	if err != nil {
		return err
	}

	return nil
}
