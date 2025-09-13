package apache

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aziyan99/wamp/internal/util"
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

func ModFcgidConfStub(phpPath, normalizePhpPath string) string {
	newPhpPath := strings.ReplaceAll(phpPath, "/", "\\")
	newPhpPath = strings.ReplaceAll(newPhpPath, "\\", "\\\\")
	return fmt.Sprintf(`
FcgidInitialEnv PATH "%s;C:\\WINDOWS\\system32;C:\\WINDOWS;C:\\WINDOWS\\System32\\Wbem;"
FcgidInitialEnv SystemRoot "C:\\Windows"
FcgidInitialEnv SystemDrive "C:"
FcgidInitialEnv TEMP "C:\\WINDOWS\\TEMP"
FcgidInitialEnv TMP "C:\\WINDOWS\\TEMP"
FcgidInitialEnv windir "C:\\WINDOWS"
FcgidInitialEnv PHPRC "%s"

# ---------------------------
# Process Management
# ---------------------------
FcgidMaxProcesses 100
FcgidMaxRequestsPerProcess 1000
FcgidMinProcessesPerClass 0
FcgidProcessLifeTime 900

# ---------------------------
# Request and Buffer Limits
# ---------------------------
FcgidMaxRequestLen 268435456
FcgidOutputBufferSize 65536

# ---------------------------
# Timeouts
# ---------------------------
FcgidIOTimeout 210
FcgidConnectTimeout 120
FcgidBusyTimeout 120
FcgidIdleTimeout 300
FcgidBusyScanInterval 120
FcgidErrorScanInterval 120
FcgidIdleScanInterval 120
FcgidZombieScanInterval 120

# ---------------------------
# PHP FastCGI Settings
# ---------------------------
FcgidInitialEnv PHP_FCGI_CHILDREN 0
FcgidInitialEnv PHP_FCGI_MAX_REQUESTS 0

# ---------------------------
# PHP File Handler
# ---------------------------
<Files ~ "\.php$">
  Options ExecCGI SymLinksIfOwnerMatch
  AddHandler fcgid-script .php
  FcgidWrapper "%s/php-cgi.exe" .php
</Files>
	`, newPhpPath, normalizePhpPath, normalizePhpPath)
}

func HttpSslConfStub() string {
	return `
#
# This is the Apache server configuration file providing SSL support.
# It contains the configuration directives to instruct the server how to
# serve pages over an https connection. For detailed information about these 
# directives see <URL:http://httpd.apache.org/docs/2.4/mod/mod_ssl.html>
# 
# Do NOT simply read the instructions in here without understanding
# what they do.  They're here only as hints or reminders.  If you are unsure
# consult the online docs. You have been warned.  
#
# Required modules: mod_log_config, mod_setenvif, mod_ssl,
#          socache_shmcb_module (for default value of SSLSessionCache)

#
# Pseudo Random Number Generator (PRNG):
# Configure one or more sources to seed the PRNG of the SSL library.
# The seed data should be of good random quality.
# WARNING! On some platforms /dev/random blocks if not enough entropy
# is available. This means you then cannot use the /dev/random device
# because it would lead to very long connection times (as long as
# it requires to make more entropy available). But usually those
# platforms additionally provide a /dev/urandom device which doesn't
# block. So, if available, use this one instead. Read the mod_ssl User
# Manual for more details.
#
#SSLRandomSeed startup file:/dev/random  512
#SSLRandomSeed startup file:/dev/urandom 512
#SSLRandomSeed connect file:/dev/random  512
#SSLRandomSeed connect file:/dev/urandom 512


#
# When we also provide SSL we have to listen to the 
# standard HTTP port (see above) and to the HTTPS port
#
Listen 443

##
##  SSL Global Context
##
##  All SSL configuration in this context applies both to
##  the main server and all SSL-enabled virtual hosts.
##

#   SSL Cipher Suite:
#   List the ciphers that the client is permitted to negotiate,
#   and that httpd will negotiate as the client of a proxied server.
#   See the OpenSSL documentation for a complete list of ciphers, and
#   ensure these follow appropriate best practices for this deployment.
#   httpd 2.2.30, 2.4.13 and later force-disable aNULL, eNULL and EXP ciphers,
#   while OpenSSL disabled these by default in 0.9.8zf/1.0.0r/1.0.1m/1.0.2a.
SSLCipherSuite HIGH:MEDIUM:!MD5:!RC4:!3DES
SSLProxyCipherSuite HIGH:MEDIUM:!MD5:!RC4:!3DES

#  By the end of 2016, only TLSv1.2 ciphers should remain in use.
#  Older ciphers should be disallowed as soon as possible, while the
#  kRSA ciphers do not offer forward secrecy.  These changes inhibit
#  older clients (such as IE6 SP2 or IE8 on Windows XP, or other legacy
#  non-browser tooling) from successfully connecting.  
#
#  To restrict mod_ssl to use only TLSv1.2 ciphers, and disable
#  those protocols which do not support forward secrecy, replace
#  the SSLCipherSuite and SSLProxyCipherSuite directives above with
#  the following two directives, as soon as practical.
# SSLCipherSuite HIGH:MEDIUM:!SSLv3:!kRSA
# SSLProxyCipherSuite HIGH:MEDIUM:!SSLv3:!kRSA

#   User agents such as web browsers are not configured for the user's
#   own preference of either security or performance, therefore this
#   must be the prerogative of the web server administrator who manages
#   cpu load versus confidentiality, so enforce the server's cipher order.
SSLHonorCipherOrder on 

#   SSL Protocol support:
#   List the protocol versions which clients are allowed to connect with.
#   Disable SSLv3 by default (cf. RFC 7525 3.1.1).  TLSv1 (1.0) should be
#   disabled as quickly as practical.  By the end of 2016, only the TLSv1.2
#   protocol or later should remain in use.
SSLProtocol all -SSLv3
SSLProxyProtocol all -SSLv3

#   Pass Phrase Dialog:
#   Configure the pass phrase gathering process.
#   The filtering dialog program ('builtin' is an internal
#   terminal dialog) has to provide the pass phrase on stdout.
SSLPassPhraseDialog  builtin

#   Inter-Process Session Cache:
#   Configure the SSL Session Cache: First the mechanism 
#   to use and second the expiring timeout (in seconds).
#SSLSessionCache         "dbm:${SRVROOT}/logs/ssl_scache"
SSLSessionCache        "shmcb:${SRVROOT}/logs/ssl_scache(512000)"
SSLSessionCacheTimeout  300

#   OCSP Stapling (requires OpenSSL 0.9.8h or later)
#
#   This feature is disabled by default and requires at least
#   the two directives SSLUseStapling and SSLStaplingCache.
#   Refer to the documentation on OCSP Stapling in the SSL/TLS
#   How-To for more information.
#
#   Enable stapling for all SSL-enabled servers:
#SSLUseStapling On

#   Define a relatively small cache for OCSP Stapling using
#   the same mechanism that is used for the SSL session cache
#   above.  If stapling is used with more than a few certificates,
#   the size may need to be increased.  (AH01929 will be logged.)
#SSLStaplingCache "shmcb:${SRVROOT}/logs/ssl_stapling(32768)"

#   Seconds before valid OCSP responses are expired from the cache
#SSLStaplingStandardCacheTimeout 3600

#   Seconds before invalid OCSP responses are expired from the cache
#SSLStaplingErrorCacheTimeout 600
	`
}
