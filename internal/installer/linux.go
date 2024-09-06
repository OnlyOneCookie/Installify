package installer

import (
	"fmt"
	"os/exec"
)

type LinuxInstaller struct{}

func (l LinuxInstaller) Install(apps []string, callback ProgressCallback) error {
	if !isAptInstalled() {
		return fmt.Errorf("apt package manager not found")
	}

	callback(0, "Updating package lists...")
	updateCmd := exec.Command("sudo", "apt-get", "update")
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update package lists: %v", err)
	}

	for i, app := range apps {
		callback(float32(i), fmt.Sprintf("Installing %s...", app))
		cmd := exec.Command("sudo", "apt-get", "install", "-y", app)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install %s: %v\n%s", app, err, string(output))
		}
		callback(float32(i+1), fmt.Sprintf("Finished installing %s", app))
	}
	return nil
}

func isAptInstalled() bool {
	cmd := exec.Command("apt-get", "--version")
	return cmd.Run() == nil
}

func init() {
	linuxInstaller = LinuxInstaller{}
}
