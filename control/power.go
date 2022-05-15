package control

import "manage/common"
func ShutDown(){
	args := []string{"-s","-t","3"}
	common.ExecCommand("shutdown",args...)
}
func Reboot(){
	args := []string{"-r","-t","3"}
	common.ExecCommand("shutdown",args...)
}
