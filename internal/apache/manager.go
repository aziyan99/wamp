package apache

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aziyan99/wamp/internal/util"
)

type Manager struct {
	apacheDir string
	tmpDir    string
}

func New(apacheDir, tmpDir string) *Manager {
	return &Manager{
		apacheDir: apacheDir,
		tmpDir:    tmpDir,
	}
}

func (m *Manager) Install(a2Version, defaultPHPPath string) error {
	var err error

	if _, err = os.Stat(path.Join(m.apacheDir, a2Version)); err == nil {
		return errors.New("apache installation already exists")
	}

	util.PrintLog("INFO").Println("Downloading Apache2...")
	err = util.DownloadFile(path.Join(m.tmpDir, a2Version+".zip"), "https://www.apachelounge.com/download/VS17/binaries/"+a2Version+".zip")
	if err != nil {
		return err
	}

	util.PrintLog("INFO").Println("Extracting Apache2...")
	err = util.Unzip(path.Join(m.tmpDir, a2Version+".zip"), path.Join(m.apacheDir, a2Version))
	if err != nil {
		return err
	}

	util.PrintLog("INFO").Println("Setup Apache2...")
	a2SrcPath := path.Join(m.apacheDir, a2Version, "Apache24")
	// for now we only support Apache24
	if _, err := os.Stat(a2SrcPath); err != nil {
		return err
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
		return err
	}

	if err = os.Remove(path.Join(m.tmpDir, a2Version+".zip")); err != nil {
		return err
	}

	fcgidVersion := "mod_fcgid-2.3.10-win64-VS17"
	util.PrintLog("INFO").Println("Downloading mod_fcgid...")
	err = util.DownloadFile(path.Join(m.tmpDir, fcgidVersion+".zip"), "https://www.apachelounge.com/download/VS17/modules/"+fcgidVersion+".zip")
	if err != nil {
		return err
	}

	err = util.Unzip(path.Join(m.tmpDir, fcgidVersion+".zip"), path.Join(m.tmpDir, fcgidVersion))
	if err != nil {
		return err
	}

	err = os.Rename(path.Join(m.tmpDir, fcgidVersion, "mod_fcgid.so"), path.Join(m.apacheDir, a2Version, "modules", "mod_fcgid.so"))
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path.Join(m.apacheDir, a2Version, "conf", "extra", "httpd-fcgid.conf"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(ModFcgidConfStub(defaultPHPPath, util.NormalizePath(defaultPHPPath))))
	if err != nil {
		return err
	}

	f, err = os.OpenFile(path.Join(m.apacheDir, a2Version, "conf", "extra", "httpd-ssl.conf"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(HttpSslConfStub()))
	if err != nil {
		return err
	}

	httpdConfPath := path.Join(m.apacheDir, a2Version, "conf", "httpd.conf")
	if err := UpdateSrvRoot(httpdConfPath, path.Join(util.NormalizePath(m.apacheDir), a2Version)); err != nil {
		return err
	}

	if err := UpdateServerName(httpdConfPath, "127.0.0.1:80"); err != nil {
		return err
	}

	if err := IncludeSslConf(httpdConfPath); err != nil {
		return err
	}

	// setup ssl-cert path
	_, err = os.Stat(path.Join(m.apacheDir, a2Version, "conf", "sites-ssl"))
	if err != nil && !errors.Is(err, fs.ErrExist) {
		if err = os.MkdirAll(path.Join(m.apacheDir, a2Version, "conf", "sites-ssl"), 0755); err != nil {
			return err
		}
	}

	if err := EnableRequiredModules(httpdConfPath); err != nil {
		return err
	}

	if err := SetupFcgidModule(httpdConfPath); err != nil {
		return err
	}

	// setup vhost
	_, err = os.Stat(path.Join(m.apacheDir, a2Version, "conf", "sites-enabled"))
	if err != nil && !errors.Is(err, fs.ErrExist) {
		if err = os.MkdirAll(path.Join(m.apacheDir, a2Version, "conf", "sites-enabled"), 0755); err != nil {
			return err
		}
	}

	if err = SetupIncludeVhost(httpdConfPath, util.NormalizePath(path.Join(m.apacheDir, a2Version, "conf", "sites-enabled"))); err != nil {
		return err
	}

	if err = os.Remove(path.Join(m.tmpDir, fcgidVersion+".zip")); err != nil {
		return err
	}

	if err = os.RemoveAll(path.Join(m.tmpDir, fcgidVersion)); err != nil {
		return err
	}
	return nil
}
