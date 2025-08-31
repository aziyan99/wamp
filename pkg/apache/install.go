package apache

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aziyan99/wamp/pkg/util"
)

func UpdateSrvRoot(confPath, newSrvRootValue string) error {
	file, err := os.Open(confPath)
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

		if len(parts) >= 3 && parts[0] == "Define" && parts[1] == "SRVROOT" {
			identation := line[:strings.Index(line, trimmedLine)]
			newLine := fmt.Sprintf(`%sDefine SRVROOT "%s"`, identation, newSrvRootValue)
			outputLines = append(outputLines, newLine)
			util.PrintLog("INFO").Printf("Found and updated SRVROOT directive.\n--- %s\n+++ %s\n", line, newLine)
		} else {
			outputLines = append(outputLines, line)
		}
	}

	updatedConf := strings.Join(outputLines, "\n")

	if err := scanner.Err(); err != nil {
		return err
	}

	err = os.WriteFile(confPath, []byte(updatedConf), 0644)
	if err != nil {
		return err
	}

	return nil
}

func UpdateServerName(confPath, newServerNameValue string) error {
	file, err := os.Open(confPath)
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

		if len(parts) >= 2 && parts[0] == "#ServerName" && parts[1] == "www.example.com:80" {
			identation := line[:strings.Index(line, trimmedLine)]
			newLine := fmt.Sprintf("%sServerName %s", identation, newServerNameValue)
			outputLines = append(outputLines, newLine)
			util.PrintLog("INFO").Printf("Found and updated ServerName directive.\n--- %s\n+++ %s\n", line, newLine)
		} else {
			outputLines = append(outputLines, line)
		}
	}

	updatedConf := strings.Join(outputLines, "\n")

	if err := scanner.Err(); err != nil {
		return err
	}

	err = os.WriteFile(confPath, []byte(updatedConf), 0644)
	if err != nil {
		return err
	}

	return nil
}

func IncludeSslConf(confPath string) error {
	file, err := os.Open(confPath)
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

		if len(parts) >= 2 && parts[0] == "#Include" && parts[1] == "conf/extra/httpd-ssl.conf" {
			identation := line[:strings.Index(line, trimmedLine)]
			newLine := fmt.Sprintf("%sInclude conf/extra/httpd-ssl.conf", identation)
			outputLines = append(outputLines, newLine)
			util.PrintLog("INFO").Printf("Found and enabled httpd-ssl.conf.\n--- %s\n+++ %s\n", line, newLine)
		} else {
			outputLines = append(outputLines, line)
		}
	}

	updatedConf := strings.Join(outputLines, "\n")

	if err := scanner.Err(); err != nil {
		return err
	}

	err = os.WriteFile(confPath, []byte(updatedConf), 0644)
	if err != nil {
		return err
	}

	return nil
}

func EnableRequiredModules(confPath string) error {

	// mod_log_config, mod_setenvif, mod_ssl
	requiredModules := map[string]string{
		"access_compat_module": "mod_access_compat.so",
		"rewrite_module":       "mod_rewrite.so",
		"socache_shmcb_module": "mod_socache_shmcb.so",
		"ssl_module":           "mod_ssl.so",
		"log_config_module":    "mod_log_config.so",
		"setenvif_module":      "mod_setenvif.so",
	}

	for key, value := range requiredModules {
		file, err := os.Open(confPath)
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

			if len(parts) >= 3 && parts[0] == "#LoadModule" && parts[1] == key {
				identation := line[:strings.Index(line, trimmedLine)]
				newLine := fmt.Sprintf("%sLoadModule %s modules/%s", identation, key, value)
				outputLines = append(outputLines, newLine)
				util.PrintLog("INFO").Printf("Found and updated %s module.\n--- %s\n+++ %s\n", key, line, newLine)
			} else {
				outputLines = append(outputLines, line)
			}
		}

		updatedConf := strings.Join(outputLines, "\n")

		if err := scanner.Err(); err != nil {
			return err
		}

		err = os.WriteFile(confPath, []byte(updatedConf), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetupFcgidModule(confPath string) error {
	f, err := os.OpenFile(confPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte("\n\n# wamp\nLoadModule fcgid_module modules/mod_fcgid.so\n<IfModule fcgid_module>\nInclude conf/extra/httpd-fcgid.conf\n</IfModule>"))
	if err != nil {
		return err
	}

	return nil
}

func SetupIncludeVhost(confPath, sitesEnabledPath string) error {
	f, err := os.OpenFile(confPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(fmt.Sprintf("\n\nIncludeOptional %s/*.conf", sitesEnabledPath)))
	if err != nil {
		return err
	}

	return nil
}
