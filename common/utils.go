package common

import (
	"bytes"
	"log"
	"os/exec"
)
func logError(err error){
    if err != nil{
        log.Fatal(err)
    }
}
func ExecCommand(command string,args ... string) string{
    cmd := exec.Command(command,args...)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    logError(err)
    res := ""
    for{
        line,err := out.ReadString('\n')
        if err != nil{
            break
        }
        res += line + "\n"
    }
    return res
}