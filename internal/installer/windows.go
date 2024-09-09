package installer

import (
	"fmt"
	"os/exec"
)

type WindowsInstaller struct{}

func (w WindowsInstaller) Install(apps []string, callback ProgressCallback) error {
	if !isChocoInstalled() {
		callback(0, "Chocolatey not found. Installing Chocolatey...", StatusNormal)
		if err := installChoco(); err != nil {
			callback(0, fmt.Sprintf("Failed to install Chocolatey: %v", err), StatusError)
			return fmt.Errorf("failed to install Chocolatey: %v", err)
		}
		callback(0, "Chocolatey installed successfully.", StatusSuccess)
	}

	totalApps := float32(len(apps))
	for i, app := range apps {
		callback(float32(i)/totalApps, fmt.Sprintf("Installing %s...", app), StatusNormal)
		cmd := exec.Command("choco", "install", app, "-y")
		output, err := cmd.CombinedOutput()
		if err != nil {
			callback((float32(i)+1)/totalApps, fmt.Sprintf("Failed to install %s: %v\n%s", app, err, string(output)), StatusError)
			continue
		}
		callback((float32(i)+1)/totalApps, fmt.Sprintf("Finished installing %s", app), StatusSuccess)
	}
	return nil
}

func (w WindowsInstaller) Uninstall(apps []string, callback ProgressCallback) error {
	if !isChocoInstalled() {
		return fmt.Errorf("chocolatey is not installed, cannot uninstall apps")
	}

	totalApps := float32(len(apps))
	for i, app := range apps {
		callback(float32(i)/totalApps, fmt.Sprintf("Uninstalling %s...", app), StatusNormal)
		cmd := exec.Command("choco", "uninstall", app, "-y")
		output, err := cmd.CombinedOutput()
		if err != nil {
			callback((float32(i)+1)/totalApps, fmt.Sprintf("Failed to uninstall %s: %v\n%s", app, err, string(output)), StatusError)
			continue
		}
		callback((float32(i)+1)/totalApps, fmt.Sprintf("Finished uninstalling %s", app), StatusSuccess)
	}
	return nil
}

func isChocoInstalled() bool {
	cmd := exec.Command("choco", "--version")
	return cmd.Run() == nil
}

func installChoco() error {
	cmd := exec.Command("powershell", "-Command", `Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))`)
	return cmd.Run()
}

func init() {
	windowsInstaller = WindowsInstaller{}
}
