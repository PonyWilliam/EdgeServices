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
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
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
}


func init(){
    interval = 5000 //默认5s更新一次数据

    allinfo.Status = "在线"
    
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
    allcpu := 0.0
    alluseMem := 0.0
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
        allcpu += cpu
		mem,err := strconv.ParseFloat(ft[3],64)
		if err != nil{
			log.Fatal(err)
		}
        alluseMem += mem
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

	// 获取开机时间
	boottime, _ := host.BootTime()
	ntime := time.Now().Unix()
	btime := time.Unix(int64(boottime), 0).Unix()
	deltatime := ntime - btime

	allinfo.Mem = int64(v.Total)
	allinfo.Usemem = allinfo.Mem - int64(v.Free)
	// 注：使用SwapMemory或VirtualMemory，在不同系统中使用率不一样，因此直接计算一次
	allinfo.Mem /= unit
	allinfo.Mem /= unit
	allinfo.Usemem /= unit
	allinfo.System = runtime.GOOS
	allinfo.Arch = runtime.GOARCH
	allinfo.Cpucores = int64(runtime.GOMAXPROCS(0))
	allinfo.Cpu = cc[0]
	allinfo.Seconds = int64(deltatime)
	allinfo.Minutes = allinfo.Seconds / 60
	allinfo.Seconds -= allinfo.Minutes * 60
	allinfo.Hours = allinfo.Minutes / 60
	allinfo.Minutes -= allinfo.Hours * 60
	allinfo.Days = allinfo.Hours / 24
	allinfo.Hours -= allinfo.Days * 24
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


func Update(){
    for{
        SysDataArr = getStatus()
        time.Sleep(time.Millisecond * time.Duration(interval))//更新间隔
    }
}