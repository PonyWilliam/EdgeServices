package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SysInfo struct {
	MemAll         uint64
	MemFree        uint64
	MemUsed        uint64
	MemUsedPercent float64
	Days           int64
	Hours          int64
	Minutes        int64
	Seconds        int64

	CpuUsedPercent float64
	OS             string
	Arch           string
	CpuCores       int
}

func GetSysInfo() (allinfo SysInfo) {
	unit := uint64(1024 * 1024) // MB

	v, _ := mem.VirtualMemory()


	cc, _ := cpu.Percent(time.Millisecond*1500, false)

	// 获取开机时间
	boottime, _ := host.BootTime()
	ntime := time.Now().Unix()
	btime := time.Unix(int64(boottime), 0).Unix()
	deltatime := ntime - btime

	allinfo.MemAll = v.Total
	allinfo.MemFree = v.Free
	allinfo.MemUsed = allinfo.MemAll - allinfo.MemFree
	// 注：使用SwapMemory或VirtualMemory，在不同系统中使用率不一样，因此直接计算一次
	allinfo.MemUsedPercent = float64(allinfo.MemUsed) / float64(allinfo.MemAll) * 100.0 // v.UsedPercent
	allinfo.MemAll /= unit
	allinfo.MemUsed /= unit
	allinfo.MemFree /= unit
	allinfo.OS = runtime.GOOS
	allinfo.Arch = runtime.GOARCH
	allinfo.CpuCores = runtime.GOMAXPROCS(0)
	allinfo.CpuUsedPercent = cc[0]
	allinfo.Seconds = int64(deltatime)
	allinfo.Minutes = allinfo.Seconds / 60
	allinfo.Seconds -= allinfo.Minutes * 60
	allinfo.Hours = allinfo.Minutes / 60
	allinfo.Minutes -= allinfo.Hours * 60
	allinfo.Days = allinfo.Hours / 24
	allinfo.Hours -= allinfo.Days * 24
	return
}
func main(){
	var t SysInfo
	t = GetSysInfo()
	fmt.Println(t)
}