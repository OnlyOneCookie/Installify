package installer

import (
	"fmt"
	"runtime"
)

type ProgressCallback func(progress float32, status string)

type OSInstaller interface {
	Install(apps []string, callback ProgressCallback) error
}

var (
	windowsInstaller OSInstaller
	macInstaller     OSInstaller
	linuxInstaller   OSInstaller
)

func Install(apps []string, callback ProgressCallback) error {
	totalApps := float32(len(apps))
	wrappedCallback := func(progress float32, status string) {
		callback(progress/totalApps, status)
	}

	switch runtime.GOOS {
	case "windows":
		return windowsInstaller.Install(apps, wrappedCallback)
	case "darwin":
		return macInstaller.Install(apps, wrappedCallback)
	case "linux":
		return linuxInstaller.Install(apps, wrappedCallback)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
