package sysinfo

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	OS  string
	RAM string
	CPU string
	GPU string
}

func GetSystemInfo() SystemInfo {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Info()

	cpuInfo := "Unknown"
	if len(c) > 0 {
		cpuInfo = c[0].ModelName
	}

	return SystemInfo{
		OS:  runtime.GOOS,
		RAM: fmt.Sprintf("%.2f GB", float64(v.Total)/1024/1024/1024),
		CPU: cpuInfo,
		GPU: "GPU info not implemented", // Implementing GPU info is more complex and OS-specific
	}
}
