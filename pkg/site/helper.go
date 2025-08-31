package site

import (
	"fmt"
	"strings"
)

func SiteVHostStub(rootPath, domain, phprcPath string) string {
	normalizeRootPath := strings.ReplaceAll(rootPath, "\\", "/")
	normalizePHPRCPath := strings.ReplaceAll(phprcPath, "\\", "/")
	return fmt.Sprintf(`
define ROOT "%s"
define DOMAIN "%s"
define PHPRC_PATH "%s"

<VirtualHost *:80>
    DocumentRoot "${ROOT}"
    ServerName ${DOMAIN}
    ServerAlias www.${DOMAIN}
    ErrorLog logs/${DOMAIN}-error.log
    CustomLog logs/${DOMAIN}-access.log common

    <Directory "${ROOT}">
        AllowOverride All
        Require all granted

        DirectoryIndex index.php
    </Directory>

    FcgidInitialEnv PHPRC "${PHPRC_PATH}"
    FcgidInitialEnv PHP_FCGI_MAX_REQUESTS 1000

    <Files ~ "\.php$>"
        AddHandler fcgid-script .php
        FcgidWrapper "${PHPRC_PATH}/php-cgi.exe" .php
    </Files>
</VirtualHost>
	`, normalizeRootPath, domain, normalizePHPRCPath)
}
