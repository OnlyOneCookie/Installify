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

	callback(0, "Updating package lists...", StatusNormal)
	updateCmd := exec.Command("sudo", "apt-get", "update")
	if err := updateCmd.Run(); err != nil {
		callback(0, fmt.Sprintf("Failed to update package lists: %v", err), StatusError)
		return fmt.Errorf("failed to update package lists: %v", err)
	}

	totalApps := float32(len(apps))
	for i, app := range apps {
		callback(float32(i)/totalApps, fmt.Sprintf("Installing %s...", app), StatusNormal)
		cmd := exec.Command("sudo", "apt-get", "install", "-y", app)
		output, err := cmd.CombinedOutput()
		if err != nil {
			callback((float32(i)+1)/totalApps, fmt.Sprintf("Failed to install %s: %v\n%s", app, err, string(output)), StatusError)
			continue
		}
		callback((float32(i)+1)/totalApps, fmt.Sprintf("Finished installing %s", app), StatusSuccess)
	}
	return nil
}

func (l LinuxInstaller) Uninstall(apps []string, callback ProgressCallback) error {
	if !isAptInstalled() {
		return fmt.Errorf("apt package manager not found")
	}

	totalApps := float32(len(apps))
	for i, app := range apps {
		callback(float32(i)/totalApps, fmt.Sprintf("Uninstalling %s...", app), StatusNormal)
		cmd := exec.Command("sudo", "apt-get", "remove", "-y", app)
		output, err := cmd.CombinedOutput()
		if err != nil {
			callback((float32(i)+1)/totalApps, fmt.Sprintf("Failed to uninstall %s: %v\n%s", app, err, string(output)), StatusError)
			continue
		}
		callback((float32(i)+1)/totalApps, fmt.Sprintf("Finished uninstalling %s", app), StatusSuccess)
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
