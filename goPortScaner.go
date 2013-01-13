package main

import (
	"fmt"
	"goPortScaner/iptool"
	"os"
)

type PingResult struct {
	Ip string
	IsOpen bool
}

func seachIP(c chan<- PingResult,ip string) {
	err1 := iptool.SearchIP(ip)
	if err1 != nil {
		c <- PingResult{ip,false}
		return
	}
	c <- PingResult{ip,true}
	return
}


func main() {
	ips, err := iptool.ParseIps("36.249.134.219~36.249.134.255")
	if err != nil {
		fmt.Println(err)
		return
	}
	file,err1 := os.OpenFile("ips.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err1 != nil {
		fmt.Println(err)
		return
	}
	o := make(chan PingResult)
	for i, count := 0, len(ips); i < count; i++ {
		go seachIP(o, ips[i])
	}
	okIps := make([]string,0)
	for i, count := 0, len(ips); i < count; i++ {
		result := <-o
		if result.IsOpen {
			okIps = append(okIps,result.Ip)
		}
	}

	for i,count := 0,len(okIps); i < count; i++ {
		file.WriteString(fmt.Sprintln(okIps[i],"....在线"))
	}

}
