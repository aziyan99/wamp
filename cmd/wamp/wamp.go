package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/aziyan99/wamp/internal/cli"
	"github.com/aziyan99/wamp/internal/manager"
	"github.com/aziyan99/wamp/internal/php"
	"github.com/aziyan99/wamp/internal/site"
	"github.com/aziyan99/wamp/internal/util"
	"github.com/aziyan99/wamp/internal/wamp"
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
	var err error
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

	app := cli.NewCommand(
		os.Args[0],
		"Wamp CLI",
		"Yet another Wamp stack manager",
		func(cmd *cli.Command, args []string) {
			fmt.Println("Use 'help' to see available commands.")
		},
	)

	wampManager := wamp.New(wampDir)
	initCmd := cli.NewCommand("init", "Initializes the application", "", func(cmd *cli.Command, args []string) {
		if err = wampManager.Init(); err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}
	})

	installCmd := cli.NewCommand("install", "Installs the application", "", func(cmd *cli.Command, args []string) {
		if err = wampManager.Install(); err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}
	})

	uninstallCmd := cli.NewCommand("uninstall", "Uninstall WAMP", "", func(cmd *cli.Command, args []string) {
		if err = wampManager.Clean(); err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}
	})
	app.AddCommands(initCmd, installCmd, uninstallCmd)

	apacheCmd := cli.NewCommand("apache", "Manages Apache", "", nil)
	apacheStartCmd := cli.NewCommand("start", "Starts Apache", "", func(cmd *cli.Command, args []string) {
		if err = loadConf(); err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		apacheProcess := manager.New(
			activeApache,
			path.Join(apacheDir, activeApache, "bin")+"\\httpd.exe",
			tmpDir,
		)

		util.PrintLog("INFO").Printf("Use Apache: %s\n", activeApache)
		util.PrintLog("INFO").Println("Apache starting...")

		err = apacheProcess.Start()
		if err != nil {
			util.PrintLog("ERROR").Fatalf("Apache unable to start. Error: %v\n", err)
		}

		util.PrintLog("INFO").Println("Apache started")
	})

	apacheStopCmd := cli.NewCommand("stop", "Stops Apache", "", func(cmd *cli.Command, args []string) {
		if err = loadConf(); err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		apacheProcess := manager.New(
			activeApache,
			path.Join(apacheDir, activeApache, "bin")+"\\httpd.exe",
			tmpDir,
		)

		util.PrintLog("INFO").Printf("Use Apache: %s\n", activeApache)
		util.PrintLog("INFO").Println("Apache stopping...")

		err = apacheProcess.Stop()
		if err != nil {
			util.PrintLog("ERROR").Fatalf("Apache unable to stop. Error: %v\n", err)
		}

		util.PrintLog("INFO").Println("Apache stopped")
	})
	apacheCmd.AddCommands(apacheStartCmd, apacheStopCmd)

	mysqlCmd := cli.NewCommand("mysql", "Manages MySQL", "", nil)
	mysqlStartCmd := cli.NewCommand("start", "Starts MySQL", "", func(cmd *cli.Command, args []string) {
		if err = loadConf(); err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		util.PrintLog("INFO").Printf("Use MySQL: %s\n", activeMysql)

		mysqlProcess := manager.New(
			activeMysql,
			path.Join(mysqlDir, activeMysql, "bin", "mysqld.exe"),
			tmpDir,
			"--console",
		)

		util.PrintLog("INFO").Println("MySQL starting...")

		err = mysqlProcess.Start()
		if err != nil {
			util.PrintLog("ERROR").Fatalf("MySQL unable to start. Error: %v\n", err)
		}

		util.PrintLog("INFO").Println("MySQL started.")
	})

	mysqlStopCmd := cli.NewCommand("stop", "Stops MySQL", "", func(cmd *cli.Command, args []string) {
		if err = loadConf(); err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		util.PrintLog("INFO").Printf("Use MySQL: %s\n", activeMysql)

		mysqlProcess := manager.New(
			activeMysql,
			path.Join(mysqlDir, activeMysql, "bin", "mysqld.exe"),
			tmpDir,
			"--console",
		)

		util.PrintLog("INFO").Println("MySQL stopping...")

		err = mysqlProcess.Stop()
		if err != nil {
			util.PrintLog("ERROR").Fatalf("MySQL unable to stop. Error: %v\n", err)
		}

		util.PrintLog("INFO").Println("MySQL stopped.")
	})
	mysqlCmd.AddCommands(mysqlStartCmd, mysqlStopCmd)

	siteCmd := cli.NewCommand("site", "Manages sites", "", nil)
	siteAddCmd := cli.NewCommand("add", "Adds a site", "", func(cmd *cli.Command, args []string) {
		// TODO: Validate sitename must include domain
		// TODO: Accept project type (e.g., laravel, wordpress, moodle)

		if err = loadConf(); err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		util.PrintLog("INFO").Printf("Use Apache: %s\n", activeApache)
		util.PrintLog("INFO").Printf("Use MySQL: %s\n", activeMysql)

		util.PrintLog("INFO").Println("Creating site...")
		sslEnable, err := strconv.ParseBool(*cmd.Flags["ssl"])
		if err != nil {
			util.PrintLog("ERROR").Fatalf("unable to parse ssl flag. Error: %v\n", err)
		}

		// sslEnable := false
		phpVersion, err := php.Search(phpDir, *cmd.Flags["php"])
		if err != nil {
			util.PrintLog("ERROR").Fatalf("unable to get php. Error: %v\n", err)
		}

		util.PrintLog("INFO").Printf("Use PHP: %s\n", phpVersion)

		sitename := args[0]
		selectedPHPPath := path.Join(phpDir, phpVersion)
		siteManager := site.New(wwwDir, path.Join(apacheDir, activeApache), selectedPHPPath, path.Join(binDir, "etc"))

		if err := siteManager.Add(sitename, sslEnable); err != nil {
			util.PrintLog("ERROR").Fatalf("unable to create site: %s. Error: %v\n", sitename, err)
		}

		util.PrintLog("INFO").Printf("Site '%s' created.\n", sitename)
	})
	siteAddCmd.AddFlag("php", "p", "php-8.3", "The php version")
	siteAddCmd.AddFlag("ssl", "s", "false", "Whether to use SSL")

	siteRmCmd := cli.NewCommand("rm", "Removes a site", "", func(cmd *cli.Command, args []string) {
		if err = loadConf(); err != nil {
			util.PrintLog("ERROR").Fatalf("%v\n", err)
		}

		util.PrintLog("INFO").Printf("Use Apache: %s\n", activeApache)
		util.PrintLog("INFO").Printf("Use MySQL: %s\n", activeMysql)
		util.PrintLog("INFO").Println("Removing site...")

		sitename := args[0]

		siteManager := site.New(wwwDir, path.Join(apacheDir, activeApache), "", path.Join(binDir, "etc"))

		if err := siteManager.Remove(sitename); err != nil {
			util.PrintLog("ERROR").Fatalf("unable to remove site: %s. Error: %v\n", sitename, err)
		}

		util.PrintLog("INFO").Printf("site: '%s' removed.\n", sitename)
	})
	siteCmd.AddCommands(siteAddCmd, siteRmCmd)

	//php-8.4.9-nts-Win32-vs17-x64
	phpCmd := cli.NewCommand("php", "Manages PHP", "Manage PHP instances", nil)
	phpInstallCmd := cli.NewCommand("install", "Install specific PHP version", "", func(cmd *cli.Command, args []string) {

		phpManager := php.New(phpDir, tmpDir)
		err = phpManager.Install(args[0])
		if err != nil {
			util.PrintLog("ERROR").Fatalf("failed to donwload php %s. Error %v\n", args[0], err)
		}

		util.PrintLog("INFO").Printf("PHP %s installed", args[0])
	})
	phpCmd.AddCommands(phpInstallCmd)

	app.AddCommands(apacheCmd, mysqlCmd, siteCmd, phpCmd)
	cli.AddHelpCommands(app, apacheCmd, mysqlCmd, siteCmd, phpCmd)
	app.Execute()
}
