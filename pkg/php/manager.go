package php

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/aziyan99/wamp/pkg/util"
)

type Manager struct {
	phpDir string
	tmpDir string
}

func New(phpDir, tmpDir string) *Manager {
	return &Manager{
		phpDir: phpDir,
		tmpDir: tmpDir,
	}
}

func (m *Manager) Install(phpVersion string) error {
	var err error

	if _, err = os.Stat(path.Join(m.phpDir, phpVersion)); err == nil {
		return errors.New("php installation already exists")
	}

	util.PrintLog("INFO").Println("Downloading PHP: " + phpVersion)
	err = util.DownloadFile(path.Join(m.tmpDir, phpVersion+".zip"), "https://windows.php.net/downloads/releases/archives/"+phpVersion+".zip")
	if err != nil {
		return err
	}

	util.PrintLog("INFO").Println("Extracting PHP...")
	err = util.Unzip(path.Join(m.tmpDir, phpVersion+".zip"), path.Join(m.phpDir, phpVersion))
	if err != nil {
		return err
	}

	util.PrintLog("INFO").Println("Setup PHP...")
	phpIni := filepath.Join(m.phpDir, phpVersion)
	if err := SetPhpIni(phpIni); err != nil {
		return err
	}

	if err := SetExtDir(phpIni); err != nil {
		return err
	}

	if err = os.Remove(path.Join(m.tmpDir, phpVersion+".zip")); err != nil {
		return err
	}

	return nil
}
