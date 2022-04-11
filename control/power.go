package control

import "manage/common"
func ShutDown(){
	common.ExecCommand("poweroff")
}
func Reboot(){
	common.ExecCommand("reboot")
}
