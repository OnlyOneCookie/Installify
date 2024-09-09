package installer

import (
	"fmt"
	"os/exec"
)

type MacInstaller struct{}

func (m MacInstaller) Install(apps []string, callback ProgressCallback) error {
	if !isBrewInstalled() {
		callback(0, "Homebrew not found. Installing Homebrew...", StatusNormal)
		if err := installBrew(); err != nil {
			callback(0, fmt.Sprintf("Failed to install Homebrew: %v", err), StatusError)
			return fmt.Errorf("failed to install Homebrew: %v", err)
		}
		callback(0, "Homebrew installed successfully.", StatusSuccess)
	}

	totalApps := float32(len(apps))
	for i, app := range apps {
		callback(float32(i)/totalApps, fmt.Sprintf("Installing %s...", app), StatusNormal)
		var cmd *exec.Cmd
		if app == "hyper" {
			cmd = exec.Command("brew", "install", "--cask", app)
		} else {
			cmd = exec.Command("brew", "install", app)
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			callback((float32(i)+1)/totalApps, fmt.Sprintf("Failed to install %s: %v\n%s", app, err, string(output)), StatusError)
			continue
		}
		callback((float32(i)+1)/totalApps, fmt.Sprintf("Finished installing %s", app), StatusSuccess)
	}
	return nil
}

func (m MacInstaller) Uninstall(apps []string, callback ProgressCallback) error {
	if !isBrewInstalled() {
		return fmt.Errorf("homebrew is not installed, cannot uninstall apps")
	}

	totalApps := float32(len(apps))
	for i, app := range apps {
		callback(float32(i)/totalApps, fmt.Sprintf("Uninstalling %s...", app), StatusNormal)
		var cmd *exec.Cmd
		if app == "hyper" {
			cmd = exec.Command("brew", "uninstall", "--cask", app)
		} else {
			cmd = exec.Command("brew", "uninstall", app)
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			callback((float32(i)+1)/totalApps, fmt.Sprintf("Failed to uninstall %s: %v\n%s", app, err, string(output)), StatusError)
			continue
		}
		callback((float32(i)+1)/totalApps, fmt.Sprintf("Finished uninstalling %s", app), StatusSuccess)
	}
	return nil
}

func isBrewInstalled() bool {
	cmd := exec.Command("brew", "--version")
	return cmd.Run() == nil
}

func installBrew() error {
	cmd := exec.Command("/bin/bash", "-c", `"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
	return cmd.Run()
}

func init() {
	macInstaller = MacInstaller{}
}
