package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aziyan99/wamp/pkg/apache"
	"github.com/aziyan99/wamp/pkg/php"
	"github.com/aziyan99/wamp/pkg/util"
)

var wampDir string
var binDir string
var apacheDir string
var mysqlDir string
var phpDir string
var wwwDir string
var tmpDir string

var activeApache string
var activeMysql string

func loadConf() error {
	conf, err := util.LoadConf(path.Join(wampDir, "wamp.ini"))
	if err != nil {
		return err
	}

	var found bool

	activeApache, found = conf.GetConf("apache", "active")
	if !found {
		return errors.New("apache conf unavailable")
	}

	activeMysql, found = conf.GetConf("mysql", "active")
	if !found {
		return errors.New("mySQL conf unavailable")
	}

	return nil
}

func main() {
	args := os.Args[1:]
	var err error

	if len(args) < 1 {
		util.PrintLog("ERROR").Fatalln("No argument specified")
	}

	wampDir, err = os.Executable()
	if err != nil {
		util.PrintLog("ERROR").Panic(err)
	}

	wampDir = filepath.Dir(wampDir)

	util.PrintLog("INFO").Printf("Wamp dir: %s\n", wampDir)

	binDir = path.Join(wampDir, "bin")
	apacheDir = path.Join(binDir, "apache")
	mysqlDir = path.Join(binDir, "mysql")
	phpDir = path.Join(binDir, "php")
	wwwDir = path.Join(wampDir, "www")
	tmpDir = path.Join(wampDir, "tmp")

	switch args[0] {
	case "init":

		util.PrintLog("INFO").Println("Initializing wamp directory...")

		topDirExist, err := util.DirExists(wampDir)
		if err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		if !topDirExist {
			err = os.MkdirAll(wampDir, 0755)
			if err != nil {
				util.PrintLog("ERROR").Fatalf("%v\n", err)
			}
		}

		binDir = fmt.Sprintf("%s\\bin", wampDir)
		binDirExist, err := util.DirExists(binDir)

		if err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		if !binDirExist {
			err = os.MkdirAll(binDir, 0755)
			if err != nil {
				util.PrintLog("ERROR").Fatalf("%v\n", err)
			}
		}

		apacheDir = fmt.Sprintf("%s\\apache", binDir)
		apacheDirExist, err := util.DirExists(apacheDir)
		if err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		if !apacheDirExist {
			err = os.MkdirAll(apacheDir, 0755)
			if err != nil {
				util.PrintLog("ERROR").Fatalf("%v\n", err)
			}
		}

		phpDir = fmt.Sprintf("%s\\php", binDir)
		phpDirExist, err := util.DirExists(phpDir)
		if err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		if !phpDirExist {
			err = os.MkdirAll(phpDir, 0755)
			if err != nil {
				util.PrintLog("ERROR").Fatalf("%v\n", err)
			}
		}

		mysqlDir = fmt.Sprintf("%s\\mysql", binDir)
		mysqlDirExist, err := util.DirExists(mysqlDir)
		if err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		if !mysqlDirExist {
			err = os.MkdirAll(mysqlDir, 0755)
			if err != nil {
				util.PrintLog("ERROR").Fatalf("%v\n", err)
			}
		}

		wwwDir = fmt.Sprintf("%s\\www", wampDir)
		wwwDirExist, err := util.DirExists(wwwDir)
		if err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		if !wwwDirExist {
			err = os.MkdirAll(wwwDir, 0755)
			if err != nil {
				util.PrintLog("ERROR").Fatalf("%v\n", err)
			}
		}

		tmpDir = fmt.Sprintf("%s\\tmp", wampDir)
		tmpDirExist, err := util.DirExists(tmpDir)
		if err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		if !tmpDirExist {
			err = os.MkdirAll(tmpDir, 0755)
			if err != nil {
				util.PrintLog("ERROR").Fatalf("%v\n", err)
			}
		}

		util.PrintLog("INFO").Println("Wamp directory initialized")

	case "install":

		phpVersion := "php-8.3.8-nts-Win32-vs16-x64"
		if _, err = os.Stat(path.Join(phpDir, phpVersion)); err == nil {
			util.PrintLog("ERROR").Fatalln("php installation already exists")
		}

		util.PrintLog("INFO").Println("Downloading default PHP...")
		err = util.DownloadFile(path.Join(tmpDir, phpVersion+".zip"), "https://windows.php.net/downloads/releases/archives/php-8.3.8-nts-Win32-vs16-x64.zip")
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		util.PrintLog("INFO").Println("Extracting PHP...")
		err = util.Unzip(path.Join(tmpDir, phpVersion+".zip"), path.Join(phpDir, phpVersion))
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		util.PrintLog("INFO").Println("Setup PHP...")
		phpIni := filepath.Join(phpDir, phpVersion)
		if err := php.SetPhpIni(phpIni); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		if err := php.SetExtDir(phpIni); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		if err = os.Remove(path.Join(tmpDir, phpVersion+".zip")); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		a2Version := "httpd-2.4.65-250724-Win64-VS17"
		if _, err = os.Stat(path.Join(apacheDir, a2Version)); err == nil {
			util.PrintLog("ERROR").Fatalln("apache installation already exists")
		}

		util.PrintLog("INFO").Println("Downloading Apache2...")
		err = util.DownloadFile(path.Join(tmpDir, a2Version+".zip"), "https://www.apachelounge.com/download/VS17/binaries/httpd-2.4.65-250724-Win64-VS17.zip")
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		util.PrintLog("INFO").Println("Extracting Apache2...")
		err = util.Unzip(path.Join(tmpDir, a2Version+".zip"), path.Join(apacheDir, a2Version))
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		util.PrintLog("INFO").Println("Setup Apache2...")
		a2SrcPath := path.Join(apacheDir, a2Version, "Apache24")
		// for now we only support Apache24
		if _, err := os.Stat(a2SrcPath); err != nil {
			util.PrintLog("ERROR").Fatalln("can not find Apache24 dir")
		}

		err = filepath.Walk(a2SrcPath, func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				newSubdir, err := filepath.Abs(strings.Replace(p, "Apache24"+string(os.PathSeparator), "", 1))
				if err != nil && !errors.Is(err, fs.ErrNotExist) {
					return err
				}

				if err = os.MkdirAll(newSubdir, 0755); err != nil {
					return err
				}
				return nil
			}

			newPath := filepath.Join(strings.Replace(p, "Apache24"+string(os.PathSeparator), "", 1))
			if err := os.Rename(p, newPath); err != nil {
				util.PrintLog("ERROR").Printf("unable to move file %s: %v\n", p, err)
				return err
			}

			return nil
		})

		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		if err = os.Remove(path.Join(tmpDir, a2Version+".zip")); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		fcgidVersion := "mod_fcgid-2.3.10-win64-VS17"
		util.PrintLog("INFO").Println("Downloading mod_fcgid...")
		err = util.DownloadFile(path.Join(tmpDir, fcgidVersion+".zip"), "https://www.apachelounge.com/download/VS17/modules/"+fcgidVersion+".zip")
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		err = util.Unzip(path.Join(tmpDir, fcgidVersion+".zip"), path.Join(tmpDir, fcgidVersion))
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		err = os.Rename(path.Join(tmpDir, fcgidVersion, "mod_fcgid.so"), path.Join(apacheDir, a2Version, "modules", "mod_fcgid.so"))
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		f, err := os.OpenFile(path.Join(apacheDir, a2Version, "conf", "extra", "httpd-fcgid.conf"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}
		defer f.Close()

		_, err = f.Write([]byte(util.ModFcgidConfStub(path.Join(phpDir, phpVersion), util.NormalizePath(path.Join(phpDir, phpVersion)))))
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		f, err = os.OpenFile(path.Join(apacheDir, a2Version, "conf", "extra", "httpd-ssl.conf"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}
		defer f.Close()
		_, err = f.Write([]byte(util.HttpSslConfStub()))
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		httpdConfPath := path.Join(apacheDir, a2Version, "conf", "httpd.conf")
		if err := apache.UpdateSrvRoot(httpdConfPath, path.Join(util.NormalizePath(apacheDir), a2Version)); err != nil {
			util.PrintLog("ERROR").Fatalln(err)
		}

		if err := apache.UpdateServerName(httpdConfPath, "127.0.0.1:80"); err != nil {
			util.PrintLog("ERROR").Fatalln(err)
		}

		if err := apache.IncludeSslConf(httpdConfPath); err != nil {
			util.PrintLog("ERROR").Fatalln(err)
		}

		if err := apache.EnableRequiredModules(httpdConfPath); err != nil {
			util.PrintLog("ERROR").Fatalln(err)
		}

		if err := apache.SetupFcgidModule(httpdConfPath); err != nil {
			util.PrintLog("ERROR").Fatalln(err)
		}

		// setup vhost
		_, err = os.Stat(path.Join(apacheDir, a2Version, "conf", "sites-enabled"))
		if err != nil && !errors.Is(err, fs.ErrExist) {
			if err = os.MkdirAll(path.Join(apacheDir, a2Version, "conf", "sites-enabled"), 0755); err != nil {
				util.PrintLog("ERROR").Panic(err)
			}
		}

		if err = apache.SetupIncludeVhost(httpdConfPath, util.NormalizePath(path.Join(apacheDir, a2Version, "conf", "sites-enabled"))); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		if err = os.Remove(path.Join(tmpDir, fcgidVersion+".zip")); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		if err = os.RemoveAll(path.Join(tmpDir, fcgidVersion)); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		mariadbVersion := "mariadb-11.8.3-winx64"
		if _, err = os.Stat(path.Join(mysqlDir, mariadbVersion)); err == nil {
			util.PrintLog("ERROR").Fatalln("mariadb installation already exists")
		}
		util.PrintLog("INFO").Println("Downloading MariaDB...")
		err = util.DownloadFile(path.Join(tmpDir, mariadbVersion+".zip"), "http://downloads.mariadb.org/rest-api/mariadb/11.8.3/mariadb-11.8.3-winx64.zip")
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		util.PrintLog("INFO").Println("Extracting MariaDB...")
		err = util.Unzip(path.Join(tmpDir, mariadbVersion+".zip"), path.Join(mysqlDir, mariadbVersion))
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		mariaDbSrcPath := path.Join(mysqlDir, mariadbVersion)
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
			util.PrintLog("ERROR").Panic(err)
		}

		// setup MariaDB
		util.PrintLog("INFO").Println("Setup MariaDB...")
		mysqlDataDirExists, err := util.DirExists(path.Join(mysqlDir, mariadbVersion, "data"))
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		if !mysqlDataDirExists {
			if err = os.Mkdir(path.Join(mysqlDir, mariadbVersion, "data"), 0755); err != nil {
				util.PrintLog("ERROR").Panic(err)
			}
		}

		mariaDbConf := []byte("[mysqld]\nbasedir=" + path.Join(mysqlDir, mariadbVersion) + "\ndatadir=" + path.Join(mysqlDir, mariadbVersion, "data"))
		if err = os.WriteFile(path.Join(mysqlDir, mariadbVersion, "my.ini"), mariaDbConf, 0755); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		mariaDdInstallDbCmd := exec.Command(path.Join(mysqlDir, mariadbVersion, "bin", "mariadb-install-db.exe"))
		mariaDdInstallDbCmd.Stdout = os.Stdout
		err = mariaDdInstallDbCmd.Run()
		if err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		if err = os.Remove(path.Join(tmpDir, mariadbVersion+".zip")); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		// setup wamp conf
		if _, err := os.Stat(path.Join(wampDir, "wamp.ini")); err != nil && errors.Is(err, fs.ErrNotExist) {
			util.PrintLog("INFO").Println("Setup wamp.conf...")
			wampConf := []byte("[apache]\nactive=" + a2Version + "\n[mysql]\nactive=" + mariadbVersion)
			if err = os.WriteFile(path.Join(wampDir, "wamp.ini"), wampConf, 0755); err != nil {
				util.PrintLog("ERROR").Panic(err)
			}
		}

	case "clean":

		util.PrintLog("INFO").Println("Cleaning all wamp directories...")

		if err = os.RemoveAll(binDir); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		if err = os.RemoveAll(tmpDir); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		if err = os.RemoveAll(wwwDir); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

	case "apache":
		if err = loadConf(); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		util.PrintLog("INFO").Printf("Use Apache: %s\n", activeApache)

		switch args[1] {
		case "start":

			_, err := os.Stat(path.Join(tmpDir, "apache_pid"))
			if err == nil {
				util.PrintLog("INFO").Println("There is apache instance running, please stop it before start")
				os.Exit(0)
			}

			util.PrintLog("INFO").Println("Starting Apache")

			cmd := exec.Command(path.Join(apacheDir, activeApache, "bin") + "\\httpd.exe")
			if err = cmd.Start(); err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			apachePid := strconv.Itoa(cmd.Process.Pid)

			// write PID for tracking the process
			if _, err := os.Stat(path.Join(tmpDir, "apache_pid")); err != nil && errors.Is(err, fs.ErrNotExist) {
				currentPid := []byte(apachePid)
				if err = os.WriteFile(path.Join(tmpDir, "apache_pid"), currentPid, 0755); err != nil {
					util.PrintLog("ERROR").Panic(err)
				}
			}

			util.PrintLog("INFO").Println("Apache started")

		case "stop":
			_, err := os.Stat(path.Join(tmpDir, "apache_pid"))
			if err != nil {
				util.PrintLog("INFO").Println("There is no apache instance running")
				os.Exit(0)
			}

			util.PrintLog("INFO").Println("Stopping Apache...")

			apachePidValue, err := os.ReadFile(path.Join(tmpDir, "apache_pid"))
			if err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			apachePid := string(apachePidValue)
			apachePidInt, err := strconv.Atoi(apachePid)
			if err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			apacheProcess, err := os.FindProcess(apachePidInt)
			if err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			if err = apacheProcess.Kill(); err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			if err = os.Remove(path.Join(tmpDir, "apache_pid")); err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			util.PrintLog("INFO").Println("Apache stopped")

		default:
			fmt.Printf("Unknown command: %s\n", args[0])
			os.Exit(1)
		}

	case "mysql":
		if err = loadConf(); err != nil {
			util.PrintLog("ERROR").Panic(err)
		}

		util.PrintLog("INFO").Printf("Use MySQL: %s\n", activeMysql)

		switch args[1] {
		case "start":

			_, err := os.Stat(path.Join(tmpDir, "mysql_pid"))
			if err == nil {
				util.PrintLog("INFO").Println("There is mysql instance running, please stop it before start")
				os.Exit(0)
			}

			util.PrintLog("INFO").Println("Starting MySQL")

			// TODO: Handle both mysql and mariadb
			cmd := exec.Command(path.Join(mysqlDir, activeMysql, "bin")+"\\mariadbd.exe", "--console")
			if err = cmd.Start(); err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			apachePid := strconv.Itoa(cmd.Process.Pid)

			// write PID for tracking the process
			if _, err := os.Stat(path.Join(tmpDir, "mysql_pid")); err != nil && errors.Is(err, fs.ErrNotExist) {
				currentPid := []byte(apachePid)
				if err = os.WriteFile(path.Join(tmpDir, "mysql_pid"), currentPid, 0755); err != nil {
					util.PrintLog("ERROR").Panic(err)
				}
			}

			util.PrintLog("INFO").Println("MySQL started")
		case "stop":
			_, err := os.Stat(path.Join(tmpDir, "mysql_pid"))
			if err != nil {
				util.PrintLog("INFO").Println("There is no mysql instance running")
				os.Exit(0)
			}

			util.PrintLog("INFO").Println("Stopping MySQL...")

			mysqlPidValue, err := os.ReadFile(path.Join(tmpDir, "mysql_pid"))
			if err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			mysqlPid := string(mysqlPidValue)
			mysqlPidInt, err := strconv.Atoi(mysqlPid)
			if err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			mysqlProcess, err := os.FindProcess(mysqlPidInt)
			if err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			if err = mysqlProcess.Kill(); err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			if err = os.Remove(path.Join(tmpDir, "mysql_pid")); err != nil {
				util.PrintLog("ERROR").Panic(err)
			}

			util.PrintLog("INFO").Println("MySQL stopped")

		default:
			fmt.Printf("Unknown command: %s\n", args[0])
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		os.Exit(1)
	}
}
