package iptool

import (
	"fmt"
	"regexp"
	"testing"
)

func TestSearchIPFail(t *testing.T) {
	err := SearchIP("192.168.2.2")
	if err == nil {
		t.Log(err)
		t.Fail()
	}
}

func TestSearchIPSucess(t *testing.T) {
	err := SearchIP("www.suncco.com")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestParseIpsOne(t *testing.T) {
	ips, err := ParseIps("192.1.1.1")
	fmt.Println(len(ips))
	if err != nil {
		t.Log(err)
		t.Fail()
	}

}

func TestParseIpsTwo(t *testing.T) {
	ips, err := ParseIps("192.1.1.a")

	if err == nil {
		t.Log(ips)
		t.Fail()
	}
	fmt.Println(len(ips))
}

func TestParseIpsThree(t *testing.T) {
	ips, err := ParseIps("192.1.1.266")
	fmt.Println(len(ips))
	if err == nil {
		t.Log(ips)
		t.Fail()
	}
}

func TestParseIpsFour(t *testing.T) {
	ips, err := ParseIps("192.1.1.111~192.1.1.111")
	fmt.Println(len(ips))
	if err != nil || ips[0] != "192.1.1.111" {
		// t.Log(ips)
		t.Fail()
	}
	// fmt.Println(ips)
}

func TestParseIps5(t *testing.T) {
	ips, err := ParseIps("192.1.1.111~192.1.1.211")
	fmt.Println(len(ips))
	if err != nil {
		// t.Log(ips)
		t.Fail()
	}
}

func TestParseIps6(t *testing.T) {
	ips, err := ParseIps("192.1.1.111~192.1.2.211")
	fmt.Println(len(ips))
	if err != nil {
		// t.Log(ips)
		t.Fail()
	}
}

func TestParseIps7(t *testing.T) {
	ips, err := ParseIps("192.1.1.111~192.2.2.211")
	fmt.Println(len(ips))
	if err != nil {
		// t.Log(ips)
		t.Fail()
	}
}

func TestParseIps8(t *testing.T) {
	ips, err := ParseIps("192.1.1.111~193.2.2.211")
	fmt.Println(len(ips))
	if err != nil {
		// t.Log(ips)
		t.Fail()
	}
}

func TestConnectPort(t *testing.T) {
	err := ConnectIPPort("172.16.10.20", 21)
	if err != nil {
		reg,_ := regexp.Compile(`refused`)
		if err != nil && reg.MatchString(err.Error()) {
			fmt.Println("没打开21端口", err)
			return
		} else {
			t.Log(err)
			t.Fail()
			return
		}
	}

	fmt.Println("打开了21端口")

}

func TestConnectPorts(t *testing.T) {
	result:= ConnectIPPosts("172.16.10.20", []int{21,23,80})
	if len(result) == 0 {
		t.Fail()
	}
	fmt.Println(result)
}

func TestPort(t *testing.T) {
	result,err := Ports("20-20")
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	fmt.Println(result)
}

func TestPort1(t *testing.T) {
	result,err := Ports("20-50")
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	fmt.Println(result)
}
