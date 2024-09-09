package installer

import (
	"fmt"
	"runtime"
)

type InstallStatus int

const (
	StatusNormal InstallStatus = iota
	StatusSuccess
	StatusWarning
	StatusError
)

type ProgressCallback func(progress float32, status string, installStatus InstallStatus)

type OSInstaller interface {
	Install(apps []string, callback ProgressCallback) error
	Uninstall(apps []string, callback ProgressCallback) error
}

var (
	windowsInstaller OSInstaller
	macInstaller     OSInstaller
	linuxInstaller   OSInstaller
)

func Install(apps []string, callback ProgressCallback) error {
	switch runtime.GOOS {
	case "windows":
		return windowsInstaller.Install(apps, callback)
	case "darwin":
		return macInstaller.Install(apps, callback)
	case "linux":
		return linuxInstaller.Install(apps, callback)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func Uninstall(apps []string, callback ProgressCallback) error {
	switch runtime.GOOS {
	case "windows":
		return windowsInstaller.Uninstall(apps, callback)
	case "darwin":
		return macInstaller.Uninstall(apps, callback)
	case "linux":
		return linuxInstaller.Uninstall(apps, callback)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
