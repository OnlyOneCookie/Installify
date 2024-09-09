package gui

import (
	"Installify/internal/applist"
	"Installify/internal/installer"
	"Installify/internal/sysinfo"
	"fmt"
	"github.com/AllenDang/giu"
	"image/color"
	"math"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type GUI struct {
	apps           map[string]*bool
	appOrder       []string
	sysInfo        sysinfo.SystemInfo
	isInstalling   bool
	isUninstalling bool
	logText        string
	progress       float32
	window         *giu.MasterWindow
	logMutex       sync.Mutex
	progressMutex  sync.Mutex
	totalApps      int
	installedApps  int
	failedApps     int
}

func New() *GUI {
	g := &GUI{
		apps:    make(map[string]*bool),
		sysInfo: sysinfo.GetSystemInfo(),
	}

	// Initialize apps based on the current OS
	for app := range applist.AppMappings[runtime.GOOS] {
		g.apps[app] = new(bool)
		g.appOrder = append(g.appOrder, app)
	}

	// Sort the app order alphabetically
	sort.Strings(g.appOrder)

	return g
}

func (g *GUI) Setup() {
	g.window = giu.NewMasterWindow("Installify", 800, 600, 0)
	g.window.Run(g.loop)
}

func (g *GUI) loop() {
	giu.SingleWindow().Layout(
		g.staticLayout(),
		g.dynamicLayout(),
	)
}

func (g *GUI) staticLayout() *giu.Layout {
	return &giu.Layout{
		giu.Label(fmt.Sprintf("OS: %s | RAM: %s | CPU: %s", g.sysInfo.OS, g.sysInfo.RAM, g.sysInfo.CPU)),
		giu.Separator(),
		giu.Label("Select apps:"),
		g.createAppCheckboxes(),
		giu.Row(
			giu.Button("Install selected").OnClick(g.installApps).Disabled(g.isInstalling || g.isUninstalling),
			giu.Button("Uninstall selected").OnClick(g.uninstallApps).Disabled(g.isInstalling || g.isUninstalling),
		),
	}
}

func (g *GUI) createAppCheckboxes() *giu.Layout {
	totalApps := len(g.appOrder)
	appsPerColumn := int(math.Ceil(float64(totalApps) / 3.0))

	var columns [3][]giu.Widget

	for i, app := range g.appOrder {
		columnIndex := i / appsPerColumn
		if columnIndex > 2 {
			columnIndex = 2
		}
		checkbox := giu.Checkbox(app, g.apps[app])
		columns[columnIndex] = append(columns[columnIndex], checkbox)
	}

	return &giu.Layout{
		giu.Row(
			giu.Column(columns[0]...),
			giu.Column(columns[1]...),
			giu.Column(columns[2]...),
		),
	}
}

func (g *GUI) dynamicLayout() *giu.Layout {
	progressPercentage := float32(g.installedApps+g.failedApps) / float32(g.totalApps)
	progressText := fmt.Sprintf("%.0f%% (%d/%d)", progressPercentage*100, g.installedApps+g.failedApps, g.totalApps)

	return &giu.Layout{
		giu.Custom(func() {
			giu.ProgressBar(progressPercentage).Overlay(progressText).Size(giu.Auto, 20).Build()
		}),
		giu.Custom(func() {
			g.logMutex.Lock()
			defer g.logMutex.Unlock()
			g.coloredTextWidget(g.logText)
		}),
	}
}

func (g *GUI) coloredTextWidget(text string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[green]") {
			giu.PushStyleColor(giu.StyleColorText, color.RGBA{0, 255, 0, 255})
			giu.Label(strings.TrimPrefix(line, "[green]")).Build()
			giu.PopStyleColor()
		} else if strings.HasPrefix(line, "[yellow]") {
			giu.PushStyleColor(giu.StyleColorText, color.RGBA{255, 255, 0, 255})
			giu.Label(strings.TrimPrefix(line, "[yellow]")).Build()
			giu.PopStyleColor()
		} else if strings.HasPrefix(line, "[red]") {
			giu.PushStyleColor(giu.StyleColorText, color.RGBA{255, 0, 0, 255})
			giu.Label(strings.TrimPrefix(line, "[red]")).Build()
			giu.PopStyleColor()
		} else {
			giu.Label(line).Build()
		}
	}
}

func (g *GUI) installApps() {
	if g.isInstalling {
		return
	}
	g.isInstalling = true
	g.setLogText("")
	g.setProgress(0)

	selectedApps := g.getSelectedApps()
	g.totalApps = len(selectedApps)
	g.installedApps = 0
	g.failedApps = 0

	go func() {
		callback := func(p float32, s string, status installer.InstallStatus) {
			g.setProgress(p)
			switch status {
			case installer.StatusSuccess:
				g.appendLogText(fmt.Sprintf("[green]%s", s))
				g.installedApps++
			case installer.StatusWarning:
				g.appendLogText(fmt.Sprintf("[yellow]%s", s))
			case installer.StatusError:
				g.appendLogText(fmt.Sprintf("[red]%s", s))
				g.failedApps++
			default:
				g.appendLogText(s)
			}
		}

		for _, app := range selectedApps {
			err := installer.Install([]string{app}, callback)
			if err != nil {
				g.appendLogText(fmt.Sprintf("[red]Failed to install %s: %v", app, err))
				g.failedApps++
			}
		}

		if g.failedApps > 0 {
			g.appendLogText(fmt.Sprintf("[yellow]Installation completed with %d failures.", g.failedApps))
		} else {
			g.appendLogText("[green]All selected apps have been installed successfully.")
		}

		g.isInstalling = false
		giu.Update()
	}()
}

func (g *GUI) uninstallApps() {
	if g.isUninstalling {
		return
	}
	g.isUninstalling = true
	g.setLogText("")
	g.setProgress(0)

	selectedApps := g.getSelectedApps()
	g.totalApps = len(selectedApps)
	g.installedApps = 0
	g.failedApps = 0

	go func() {
		callback := func(p float32, s string, status installer.InstallStatus) {
			g.setProgress(p)
			switch status {
			case installer.StatusSuccess:
				g.appendLogText(fmt.Sprintf("[green]%s", s))
				g.installedApps++
			case installer.StatusWarning:
				g.appendLogText(fmt.Sprintf("[yellow]%s", s))
			case installer.StatusError:
				g.appendLogText(fmt.Sprintf("[red]%s", s))
				g.failedApps++
			default:
				g.appendLogText(s)
			}
		}

		for _, app := range selectedApps {
			err := installer.Uninstall([]string{app}, callback)
			if err != nil {
				g.appendLogText(fmt.Sprintf("[red]Failed to uninstall %s: %v", app, err))
				g.failedApps++
			}
		}

		if g.failedApps > 0 {
			g.appendLogText(fmt.Sprintf("[yellow]Uninstallation completed with %d failures.", g.failedApps))
		} else {
			g.appendLogText("[green]All selected apps have been uninstalled successfully.")
		}

		g.isUninstalling = false
		giu.Update()
	}()
}

func (g *GUI) getSelectedApps() []string {
	selectedApps := []string{}
	osSpecificMappings := applist.AppMappings[runtime.GOOS]
	for app, selected := range g.apps {
		if *selected {
			if packageName, ok := osSpecificMappings[app]; ok {
				selectedApps = append(selectedApps, packageName)
			} else {
				g.appendLogText(fmt.Sprintf("[yellow]Warning: No package mapping found for %s on %s", app, runtime.GOOS))
			}
		}
	}
	return selectedApps
}

func (g *GUI) setProgress(p float32) {
	g.progressMutex.Lock()
	defer g.progressMutex.Unlock()
	g.progress = p
	giu.Update()
}

func (g *GUI) setLogText(s string) {
	g.logMutex.Lock()
	defer g.logMutex.Unlock()
	g.logText = s
	giu.Update()
}

func (g *GUI) appendLogText(s string) {
	g.logMutex.Lock()
	defer g.logMutex.Unlock()
	g.logText += s + "\n"
	giu.Update()
}
