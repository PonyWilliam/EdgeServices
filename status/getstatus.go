package status

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var interval int64

var SysDataArr []Process
var sysinfo string
var allinfo AllInfo
type Process struct {
    Pid int `json:"pid"`
    Cpu float64 `json:"cpu"`
	Mem float64 `json:"mem"`
	Command string `json:"command"`
}

type AllInfo struct{
    Status string `json:"status"`
    CpuInfo string `json:"cpuinfo"`
    Cpu float64 `json:"cpu"`
    Cpucores int64 `json:"cpuCores"`
	Mem int64 `json:"allmem"`
    Usemem int64 `json:"usemem"`
    System string `json:"system"`
    Disk string `json:"disk"`
    Arch    string `json:"arch"`
    Days           int64    `json:"days"`
	Hours          int64    `json:"hours"`
	Minutes        int64    `json:"minutes"`
	Seconds        int64    `json:"seconds"`
    Network string `json:"network"`
}
func init(){
    interval = 5000 //默认5s更新一次数据
    allinfo.Status = "在线"    
    n, _ := host.Info()
    allinfo.System = fmt.Sprintf("%v(%v) %v\n", n.Platform, n.PlatformFamily, n.PlatformVersion)
}
func logError(err error){
    if err != nil{
        log.Fatal(err)
    }
}

func getStatus() []Process{
    cmd := exec.Command("ps", "aux")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    processes := make([]Process, 0)
    for {
        line, err := out.ReadString('\n')
        if err!=nil {
            break;
        }
        tokens := strings.Split(line, " ")
        ft := make([]string, 0)
        for _, t := range(tokens) {
            if t!="" && t!="\t" {
                ft = append(ft, t)
            }
        }
        pid, err := strconv.Atoi(ft[1])
        if err!=nil {
            continue
        }
        cpu, err := strconv.ParseFloat(ft[2], 64)
        if err!=nil {
            log.Fatal(err)
        }
		mem,err := strconv.ParseFloat(ft[3],64)
		if err != nil{
			log.Fatal(err)
		}
		command := ft[10]
        processes = append(processes, Process{pid, cpu,mem,command})
    }
    updateInfo()
    return processes
}
func updateInfo(){
    unit := int64(1024 * 1024) // MB

	v, _ := mem.VirtualMemory()


	cc, _ := cpu.Percent(time.Millisecond*1500, false)

	boottime, _ := host.BootTime()
	ntime := time.Now().Unix()
	btime := time.Unix(int64(boottime), 0).Unix()
	deltatime := ntime - btime

	allinfo.Mem = int64(v.Total)
	allinfo.Usemem = allinfo.Mem - int64(v.Free)
	allinfo.Mem /= unit
	allinfo.Usemem /= unit
	allinfo.Arch = runtime.GOARCH
    cores,_ := cpu.Counts(false)
	allinfo.Cpucores = int64(cores)
	allinfo.Cpu = cc[0]
    allinfo.Cpu = Decimal(allinfo.Cpu)
	allinfo.Seconds = int64(deltatime)
	allinfo.Minutes = allinfo.Seconds / 60
	allinfo.Seconds -= allinfo.Minutes * 60
	allinfo.Hours = allinfo.Minutes / 60
	allinfo.Minutes -= allinfo.Hours * 60
	allinfo.Days = allinfo.Hours / 24
	allinfo.Hours -= allinfo.Days * 24
    nv, _ := net.IOCounters(true)
    allinfo.Network = fmt.Sprintf("Network: %v bytes / %v bytes\n", nv[0].BytesRecv, nv[0].BytesSent)
    c,_ := cpu.Info()
    if len(c) > 1 {
		for _, sub_cpu := range c {
            allinfo.CpuInfo = sub_cpu.ModelName
            break
		}
	} else {
		sub_cpu := c[0]
		allinfo.CpuInfo = sub_cpu.ModelName
	}

}
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
func SetInterVal(data int64){
    interval = data
}
func GetSysDataAddr() []Process{
    return SysDataArr
}
func GetSysInfo() string{
    return sysinfo
}
func GetTotal()AllInfo{
    return allinfo
}
func updatenotimportant(){
    d, _ := disk.Usage("/")
    allinfo.Disk = fmt.Sprintf("总大小: %v GB 剩余大小: %v GB Usage:%.2f%%", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
}
func Update(){
    i := 60
    for{
        if(i == 60){
            //更新不重要的信息
            i = 0
            updatenotimportant()
        }
        i += 1
        SysDataArr = getStatus()
        time.Sleep(time.Millisecond * time.Duration(interval))//更新间隔
    }
}