package main

import (
	"fmt"
	"os"
)
func main(){
	args := os.Args
	if args == nil || len(args) < 2{
		fmt.Println("无输入")
	}else{
		for _,v := range args{
			fmt.Println(v)
		}
	}
}