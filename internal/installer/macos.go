package installer

import (
	"fmt"
	"os/exec"
)

type MacInstaller struct{}

func (m MacInstaller) Install(apps []string, callback ProgressCallback) error {
	if !isBrewInstalled() {
		callback(0, "Homebrew not found. Installing Homebrew...")
		if err := installBrew(); err != nil {
			return fmt.Errorf("failed to install Homebrew: %v", err)
		}
		callback(0, "Homebrew installed successfully.")
	}

	for i, app := range apps {
		callback(float32(i), fmt.Sprintf("Installing %s...", app))
		cmd := exec.Command("brew", "install", app)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install %s: %v\n%s", app, err, string(output))
		}
		callback(float32(i+1), fmt.Sprintf("Finished installing %s", app))
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
