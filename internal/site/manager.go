package site

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path"

	"github.com/aziyan99/wamp/internal/hostsrw"
	"github.com/aziyan99/wamp/internal/util"
)

type Manager struct {
	wwwDir          string
	activeApacheDir string
	selectedPHPDir  string
	etcDir          string
}

func New(wwwDir, activeApacheDir, selectedPHPDir, etcDir string) *Manager {
	return &Manager{
		wwwDir:          wwwDir,
		activeApacheDir: activeApacheDir,
		selectedPHPDir:  selectedPHPDir,
		etcDir:          etcDir,
	}
}

func (m *Manager) Add(sitename string, sslEnable bool) error {

	// TODO: Validate sitename must include domain
	// TODO: Handle site with 'public' path
	// TODO: Accept project type (e.g., laravel, wordpress, moodle)

	siteDir := path.Join(m.wwwDir, sitename)
	siteConf := path.Join(m.activeApacheDir, "conf", "sites-enabled", sitename+".conf")

	isDirSiteExists, err := util.DirExists(siteDir)
	if err != nil {
		return err
	}

	_, err = os.Stat(siteConf)
	if err == nil && isDirSiteExists {
		return errors.New("site exists")
	}

	isPHPExists, err := util.DirExists(m.selectedPHPDir)
	if err != nil {
		return err
	}

	if !isPHPExists {
		return errors.New("selected PHP version do not exists")
	}

	if sslEnable {
		_, err = os.Stat(path.Join(m.activeApacheDir, "conf", "sites-ssl", sitename+".pem"))
		if err == nil {
			return errors.New("site ssl .pem exists")
		}

		_, err = os.Stat(path.Join(m.activeApacheDir, "conf", "sites-ssl", sitename+"-key.pem"))
		if err == nil {
			return errors.New("site ssl .pem exists")
		}
	}

	if !isDirSiteExists {
		if err = os.Mkdir(siteDir, 0755); err != nil {
			return errors.New("unable to create site dir")
		}
	}

	confFileValue := []byte(SiteVHostStub(siteDir, sitename, m.selectedPHPDir))

	if sslEnable {
		err = os.Chdir(path.Join(m.activeApacheDir, "conf", "sites-ssl"))
		if err != nil {
			return err
		}

		mkCertCmd := exec.Command(path.Join(m.etcDir, "mkcert.exe"), sitename)
		mkCertCmd.Stdout = os.Stdout
		err = mkCertCmd.Run()
		if err != nil {
			return errors.New("unable to create site ssl conf")
		}

		confFileValue = []byte(SiteVHostSSLStub(siteDir, sitename, m.selectedPHPDir, path.Join(m.activeApacheDir, "conf", "sites-ssl", sitename+".pem"), path.Join(m.activeApacheDir, "conf", "sites-ssl", sitename+"-key.pem")))
	}

	if err = os.WriteFile(siteConf, confFileValue, 0755); err != nil {
		return err
	}

	hostsManager := hostsrw.New(m.etcDir)
	err = hostsManager.Add(sitename)
	if err != nil {
		util.PrintLog("INFO").Printf("Unable to write '%s' into windows hosts file. Please add '127.0.0.1 %s' to your windows hosts file manually. Error: %v\n", sitename, sitename, err)
	} else {
		util.PrintLog("INFO").Printf("Wrote '%s' into windows hosts file.\n", sitename)
	}

	return nil
}

func (m *Manager) Remove(sitename string) error {
	siteDir := path.Join(m.wwwDir, sitename)
	siteConf := path.Join(m.activeApacheDir, "conf", "sites-enabled", sitename+".conf")
	accessLog := path.Join(m.activeApacheDir, "logs", sitename+"-access.log")
	errLog := path.Join(m.activeApacheDir, "logs", sitename+"-error.log")
	siteSSLConf := path.Join(m.activeApacheDir, "conf", "sites-ssl", sitename+".pem")
	siteSSLKeyConf := path.Join(m.activeApacheDir, "conf", "sites-ssl", sitename+"-key.pem")

	if err := os.RemoveAll(siteDir); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if err := os.Remove(siteConf); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if err := os.Remove(accessLog); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if err := os.Remove(errLog); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if _, err := os.Stat(siteSSLConf); err == nil {
		if err := os.Remove(siteSSLConf); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	if _, err := os.Stat(siteSSLKeyConf); err == nil {
		if err := os.Remove(siteSSLKeyConf); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	hostsManager := hostsrw.New(m.etcDir)
	err := hostsManager.Remove(sitename)
	if err != nil {
		util.PrintLog("INFO").Printf("Unable to remove '%s' into windows hosts file. Please remove '127.0.0.1 %s' from your windows hosts file manually. error: %v\n", sitename, sitename, err)
	} else {
		util.PrintLog("INFO").Printf("Remove '%s' from windows hosts file.\n", sitename)
	}

	return nil
}
