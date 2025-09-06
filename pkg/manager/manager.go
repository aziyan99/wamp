package manager

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strconv"
)

type Manager struct {
	name   string
	bin    string
	tmpDir string
	args   []string
}

func New(name, bin, tmpDir string, args ...string) *Manager {
	return &Manager{
		name:   name,
		bin:    bin,
		tmpDir: tmpDir,
		args:   args,
	}
}

func (m *Manager) Start() error {
	_, err := os.Stat(path.Join(m.tmpDir, m.name+"_pid"))
	if err == nil {
		return fmt.Errorf("instance %s is running", m.name)
	}

	cmd := exec.Command(m.bin, m.args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	pid := strconv.Itoa(cmd.Process.Pid)

	if _, err := os.Stat(path.Join(m.tmpDir, m.name+"_pid")); err != nil && errors.Is(err, fs.ErrNotExist) {
		currentPid := []byte(pid)
		if err = os.WriteFile(path.Join(m.tmpDir, m.name+"_pid"), currentPid, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) Stop() error {
	_, err := os.Stat(path.Join(m.tmpDir, m.name+"_pid"))
	if err != nil {
		return fmt.Errorf("instance %s is not running", m.name)
	}

	pidValue, err := os.ReadFile(path.Join(m.tmpDir, m.name+"_pid"))
	if err != nil {
		return err
	}

	pid := string(pidValue)
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pidInt)
	if err != nil {
		return err
	}

	if err = process.Kill(); err != nil {
		return err
	}

	if err = os.Remove(path.Join(m.tmpDir, m.name+"_pid")); err != nil {
		return err
	}

	return nil
}
