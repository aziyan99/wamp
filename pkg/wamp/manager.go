package wamp

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/aziyan99/wamp/pkg/apache"
	"github.com/aziyan99/wamp/pkg/php"
	"github.com/aziyan99/wamp/pkg/util"
)

type Manager struct {
	wampDir   string
	binDir    string
	apacheDir string
	phpDir    string
	mysqlDir  string
	wwwDir    string
	tmpDir    string
}

func New(wampDir string) *Manager {
	binDir := path.Join(wampDir, "bin")
	return &Manager{
		wampDir:   wampDir,
		binDir:    binDir,
		apacheDir: path.Join(binDir, "apache"),
		phpDir:    path.Join(binDir, "php"),
		mysqlDir:  path.Join(binDir, "mysql"),
		wwwDir:    path.Join(wampDir, "www"),
		tmpDir:    path.Join(wampDir, "tmp"),
	}
}

func (m *Manager) Init() error {
	util.PrintLog("INFO").Println("Initializing wamp directory...")

	topDirExist, err := util.DirExists(m.wampDir)
	if err != nil {
		return err
	}

	if !topDirExist {
		err = os.MkdirAll(m.wampDir, 0755)
		if err != nil {
			return err
		}
	}

	binDirExist, err := util.DirExists(m.binDir)

	if err != nil {
		return err
	}

	if !binDirExist {
		err = os.MkdirAll(m.binDir, 0755)
		if err != nil {
			return err
		}
	}

	apacheDirExist, err := util.DirExists(m.apacheDir)
	if err != nil {
		return err
	}

	if !apacheDirExist {
		err = os.MkdirAll(m.apacheDir, 0755)
		if err != nil {
			return err
		}
	}

	phpDirExist, err := util.DirExists(m.phpDir)
	if err != nil {
		return err
	}

	if !phpDirExist {
		err = os.MkdirAll(m.phpDir, 0755)
		if err != nil {
			return err
		}
	}

	mysqlDirExist, err := util.DirExists(m.mysqlDir)
	if err != nil {
		return err
	}

	if !mysqlDirExist {
		err = os.MkdirAll(m.mysqlDir, 0755)
		if err != nil {
			return err
		}
	}

	wwwDirExist, err := util.DirExists(m.wwwDir)
	if err != nil {
		return err
	}

	if !wwwDirExist {
		err = os.MkdirAll(m.wwwDir, 0755)
		if err != nil {
			return err
		}
	}

	tmpDirExist, err := util.DirExists(m.tmpDir)
	if err != nil {
		return err
	}

	if !tmpDirExist {
		err = os.MkdirAll(m.tmpDir, 0755)
		if err != nil {
			return err
		}
	}

	util.PrintLog("INFO").Println("Wamp directory initialized")

	return nil
}

func (m *Manager) Install() error {
	var err error

	phpVersion := "php-8.3.8-nts-Win32-vs16-x64"
	phpManager := php.New(m.phpDir, m.tmpDir)
	err = phpManager.Install(phpVersion)
	if err != nil {
		return err
	}

	a2Version := "httpd-2.4.65-250724-Win64-VS17"
	apacheManager := apache.New(m.apacheDir, m.tmpDir)
	err = apacheManager.Install(a2Version, path.Join(m.phpDir, phpVersion))
	if err != nil {
		return err
	}

	mariadbVersion := "mariadb-11.8.3-winx64"
	if _, err = os.Stat(path.Join(m.mysqlDir, mariadbVersion)); err == nil {
		return errors.New("mysql installation already exists")
	}
	util.PrintLog("INFO").Println("Downloading MariaDB...")
	err = util.DownloadFile(path.Join(m.tmpDir, mariadbVersion+".zip"), "http://downloads.mariadb.org/rest-api/mariadb/11.8.3/mariadb-11.8.3-winx64.zip")
	if err != nil {
		return err
	}

	util.PrintLog("INFO").Println("Extracting MariaDB...")
	err = util.Unzip(path.Join(m.tmpDir, mariadbVersion+".zip"), path.Join(m.mysqlDir, mariadbVersion))
	if err != nil {
		return err
	}

	mariaDbSrcPath := path.Join(m.mysqlDir, mariadbVersion)
	err = filepath.Walk(mariaDbSrcPath, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			newSubdir, err := filepath.Abs(strings.Replace(p, mariadbVersion, "", 1))
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return err
			}

			if err = os.MkdirAll(newSubdir, 0755); err != nil {
				return err
			}
			return nil
		}

		newPath := filepath.Join(strings.Replace(p, mariadbVersion, "", 1))
		if err := os.Rename(p, newPath); err != nil {
			util.PrintLog("ERROR").Printf("unable to move file %s: %v\n", p, err)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// setup MariaDB
	util.PrintLog("INFO").Println("Setup MariaDB...")
	mysqlDataDirExists, err := util.DirExists(path.Join(m.mysqlDir, mariadbVersion, "data"))
	if err != nil {
		return err
	}

	if !mysqlDataDirExists {
		if err = os.Mkdir(path.Join(m.mysqlDir, mariadbVersion, "data"), 0755); err != nil {
			return err
		}
	}

	mariaDbConf := []byte("[mysqld]\nbasedir=" + path.Join(m.mysqlDir, mariadbVersion) + "\ndatadir=" + path.Join(m.mysqlDir, mariadbVersion, "data"))
	if err = os.WriteFile(path.Join(m.mysqlDir, mariadbVersion, "my.ini"), mariaDbConf, 0755); err != nil {
		return err
	}

	mariaDdInstallDbCmd := exec.Command(path.Join(m.mysqlDir, mariadbVersion, "bin", "mariadb-install-db.exe"))
	mariaDdInstallDbCmd.Stdout = os.Stdout
	err = mariaDdInstallDbCmd.Run()
	if err != nil {
		return err
	}

	if err = os.Remove(path.Join(m.tmpDir, mariadbVersion+".zip")); err != nil {
		return err
	}

	// setup wamp conf
	if _, err := os.Stat(path.Join(m.wampDir, "wamp.ini")); err != nil && errors.Is(err, fs.ErrNotExist) {
		util.PrintLog("INFO").Println("Setup wamp.conf...")
		wampConf := []byte("[apache]\nactive=" + a2Version + "\n[mysql]\nactive=" + mariadbVersion)
		if err = os.WriteFile(path.Join(m.wampDir, "wamp.ini"), wampConf, 0755); err != nil {
			return err
		}
	}

	// setup etc to download 3rd party bin
	_, err = os.Stat(path.Join(m.binDir, "etc"))
	if err != nil && !errors.Is(err, fs.ErrExist) {
		if err = os.MkdirAll(path.Join(m.binDir, "etc"), 0755); err != nil {
			return err
		}
	}

	// setup mkcert

	util.PrintLog("INFO").Println("Setup mkcert...")

	if _, err = os.Stat(path.Join(m.binDir, "etc", "mkcert.exe")); err == nil {
		return errors.New("mkcert already exists")
	}

	util.PrintLog("INFO").Println("Downloading mkcert...")
	err = util.DownloadFile(path.Join(m.binDir, "etc", "mkcert.exe"), "https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-windows-amd64.exe")
	if err != nil {
		return err
	}

	mkcertInstallCmd := exec.Command(path.Join(m.binDir, "etc", "mkcert.exe"), "-install")
	mkcertInstallCmd.Stdout = os.Stdout
	err = mkcertInstallCmd.Run()
	if err != nil {
		return err
	}

	util.PrintLog("INFO").Println("mkcert installed")

	util.PrintLog("INFO").Println("Downloading hostsrw...")
	err = util.DownloadFile(path.Join(m.binDir, "etc", "hostsrw.exe"), "https://github.com/aziyan99/hostsrw/releases/download/v2.3.2/hostsrw.exe")
	if err != nil {
		return err
	}

	util.PrintLog("INFO").Println("hostsrw installed")

	// install composer: https://getcomposer.org/download/latest-stable/composer.phar
	util.PrintLog("INFO").Println("Installing composer...")
	err = util.DownloadFile(path.Join(m.binDir, "etc", "composer.phar"), "https://getcomposer.org/download/latest-stable/composer.phar")
	if err != nil {
		return err
	}

	util.PrintLog("INFO").Println("composer installed")

	return nil
}

func (m *Manager) Clean() error {
	var err error
	util.PrintLog("INFO").Println("Cleaning all wamp directories...")

	mkcertUninstallCmd := exec.Command(path.Join(m.binDir, "etc", "mkcert.exe"), "-uninstall")
	mkcertUninstallCmd.Stdout = os.Stdout
	err = mkcertUninstallCmd.Run()
	if err != nil {
		return err
	}

	if err = util.CleanDirs(m.binDir, m.tmpDir, m.wwwDir); err != nil {
		return err
	}

	if _, err = os.Stat(path.Join(m.wampDir, "wamp.ini")); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if err = os.Remove(path.Join(m.wampDir, "wamp.ini")); err != nil {
		return err
	}

	util.PrintLog("INFO").Println("Done")

	return nil
}
