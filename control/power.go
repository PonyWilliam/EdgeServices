package control

import "manage/common"
func ShutDown(){
	common.ExecCommand("bash","-c","poweroff")
}
func Reboot(){
	common.ExecCommand("bash","-c","reboot")
}
