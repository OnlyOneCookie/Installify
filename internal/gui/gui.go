package gui

import (
	"Installify/internal/installer"
	"Installify/internal/sysinfo"
	"fmt"
	"github.com/AllenDang/giu"
	"strings"
)

var (
	apps         map[string]*bool
	sysInfo      sysinfo.SystemInfo
	isInstalling bool
	logText      string
	progress     float32
	updateChan   chan struct{}
)

func Setup() {
	apps = make(map[string]*bool)
	apps["spotify"] = new(bool)
	apps["discord@canary"] = new(bool)
	apps["brave"] = new(bool)
	sysInfo = sysinfo.GetSystemInfo()
	updateChan = make(chan struct{}, 1)
}

func Loop() {
	select {
	case <-updateChan:
		giu.Update()
	default:
	}

	giu.SingleWindow().Layout(
		giu.Label(fmt.Sprintf("OS: %s | RAM: %s | CPU: %s", sysInfo.OS, sysInfo.RAM, sysInfo.CPU)),
		giu.Separator(),
		giu.Label("Select apps to install:"),
		createAppCheckboxes(),
		giu.Button("Install Selected Apps").OnClick(installApps).Disabled(isInstalling),
		giu.ProgressBar(progress).Size(giu.Auto, 20),
		giu.InputTextMultiline(&logText).Size(giu.Auto, 200).Flags(giu.InputTextFlagsReadOnly),
	)
}

func createAppCheckboxes() *giu.Layout {
	var checkboxes giu.Layout
	for app, selected := range apps {
		app := app // Create a new variable for each iteration
		checkboxes = append(checkboxes, giu.Checkbox(app, selected))
	}
	return &checkboxes
}

func installApps() {
	if isInstalling {
		return
	}
	isInstalling = true
	logText = ""
	progress = 0

	selectedApps := []string{}
	for app, selected := range apps {
		if *selected {
			selectedApps = append(selectedApps, app)
		}
	}

	go func() {
		var log strings.Builder
		log.WriteString("Installation Log:\n")

		callback := func(p float32, s string) {
			progress = p
			log.WriteString(s + "\n")
			logText = log.String()
			select {
			case updateChan <- struct{}{}:
			default:
			}
		}

		err := installer.Install(selectedApps, callback)

		if err != nil {
			log.WriteString(fmt.Sprintf("Installation failed: %v\n", err))
		} else {
			log.WriteString("All selected apps have been installed successfully.\n")
		}

		logText = log.String()
		isInstalling = false
		select {
		case updateChan <- struct{}{}:
		default:
		}
	}()
}
