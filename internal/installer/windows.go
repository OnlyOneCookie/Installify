package installer

import (
	"fmt"
	"os/exec"
)

type WindowsInstaller struct{}

func (w WindowsInstaller) Install(apps []string, callback ProgressCallback) error {
	if !isChocoInstalled() {
		callback(0, "Chocolatey not found. Installing Chocolatey...")
		if err := installChoco(); err != nil {
			return fmt.Errorf("failed to install Chocolatey: %v", err)
		}
		callback(0, "Chocolatey installed successfully.")
	}

	for i, app := range apps {
		callback(float32(i), fmt.Sprintf("Installing %s...", app))
		cmd := exec.Command("choco", "install", app, "-y")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install %s: %v\n%s", app, err, string(output))
		}
		callback(float32(i+1), fmt.Sprintf("Finished installing %s", app))
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
