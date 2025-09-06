package hostsrw

import (
	"errors"
	"os/exec"
	"path"
)

type Manager struct {
	etcDir string
}

func New(etcDir string) *Manager {
	return &Manager{
		etcDir: etcDir,
	}
}

func (m *Manager) Add(sitename string) error {
	hostsrwExistsCmd := exec.Command(path.Join(m.etcDir, "hostsrw.exe"), "exists", sitename)
	output, err := hostsrwExistsCmd.Output()
	if err != nil {
		return err
	}

	str := string(output)

	if len(str) > 0 {
		return errors.New("sitename: " + sitename + " already registered on hosts file")
	}

	hostsrwAddCmd := exec.Command(path.Join(m.etcDir, "hostsrw.exe"), "add", sitename)
	_, err = hostsrwAddCmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Remove(sitename string) error {
	hostsrwRemoveCmd := exec.Command(path.Join(m.etcDir, "hostsrw.exe"), "rm", sitename)
	output, err := hostsrwRemoveCmd.Output()
	if err != nil {
		return err
	}

	str := string(output)

	if len(str) > 0 {
		return errors.New("sitename: " + sitename + ". error: " + str)
	}

	return nil
}
