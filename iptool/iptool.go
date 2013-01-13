package iptool

import (
	"io"
	// "io/ioutil"
	// "bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func isCanConnect(c chan<- error, ip string) {
	cmd := exec.Command("ping", ip)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c <- err
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		c <- err
		return
	}

	r := io.MultiReader(stdout, stderr)

	cmd.Start()
	// for {
	buf := make([]byte, 1024)

	count, err := r.Read(buf)
	if err != nil || count == 0 {
		c <- errors.New("no data")
		return

	} else {

		rex, err := regexp.Compile(`ttl`)
		if err != nil {
			c <- err
			return
		}
		if rex.Match(buf) {
			c <- nil
			return
		} else {
			c <- errors.New("return data no match ttl")
			return
		}
	}
	// }
	// return
}

func SearchIP(ip string) error {
	c := make(chan error)
	go isCanConnect(c, ip)
	for {
		select {
		case v := <-c:
			return v
		case <-time.After(1 * time.Second):
			return errors.New("timeout")
		}
	}
	return errors.New("timeout")
}

func isIp(str string) (ips []int, err error) {
	temp := strings.Split(str, ".")
	if len(temp) != 4 {
	} else {
		ips = make([]int, 4, 4)
		for i := 0; i < 4; i++ {
			ips[i], err = strconv.Atoi(temp[i])
			if err != nil {
				return nil, err
			}
		}

		for i := 0; i < 4; i++ {
			if ips[i] < 0 || ips[i] > 255 {
				return nil, errors.New(fmt.Sprintln(str, "err:", ips[i], " ---- .<0 or >255"))
			}
		}
		return ips, nil
	}
	return nil, errors.New("data no mach,192.168.1.1 or 192.168.1.1~192.168.255.255")
}

func ParseIps(str string) (ips []string, err error) {
	temp := strings.Split(str, "~")
	switch len(temp) {
	case 1:
		var intIps []int
		intIps, err = isIp(temp[0])
		if err != nil {
			return nil, err
		} else {
			ips = []string{fmt.Sprintf("%d.%d.%d.%d", intIps[0], intIps[1], intIps[2], intIps[3])}
			return ips, nil
		}
	case 2:
		var intIps []int
		intIps, err = isIp(temp[0])
		if err != nil {
			return nil, err
		}

		var intIps1 []int
		intIps1, err = isIp(temp[1])
		if err != nil {
			return nil, err
		}

		num1 := binary.BigEndian.Uint32([]byte{byte(intIps[0]), byte(intIps[1]), byte(intIps[2]), byte(intIps[3])})
		num2 := binary.BigEndian.Uint32([]byte{byte(intIps1[0]), byte(intIps1[1]), byte(intIps1[2]), byte(intIps1[3])})
		if num1 > num2 {
			return nil, errors.New(fmt.Sprintf("%s must large %s", temp[1], temp[0]))
		} else if num1 == num2 {
			ips = []string{fmt.Sprintf("%d.%d.%d.%d", intIps[0], intIps[1], intIps[2], intIps[3])}
			return ips, nil
		}

		uintIps := make([]uint32, 0)
		for {

			if num1 > num2 {
				break
			}
			uintIps = append(uintIps, num1)
			num1++
		}

		ips = make([]string, 0)
		for i, count := 0, len(uintIps); i < count; i++ {
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, uintIps[i])
			ips = append(ips, fmt.Sprintf("%d.%d.%d.%d", bytes[0], bytes[1], bytes[2], bytes[3]))

		}

		return ips, nil

	}
	return nil, errors.New("data no mach,192.168.1.1 or 192.168.1.1~192.168.255.255")
}

func ConnectIPPort(ip string, port int) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		// fmt.Println(err)
		return err
	}
	conn.Close()
	return nil
}

func ConnectIPPosts(ip string, ports []int) (results []bool) {
	results = make([]bool, 0)
	for i, count := 0, len(ports); i < count; i++ {
		err := ConnectIPPort(ip, ports[i])

		if err != nil {
			results = append(results, false)
		} else {
			results = append(results, true)
		}

	}
	return results
}

func Ports(ports string) (resultPort []int, err error) {
	temp := strings.Split(ports, "-")
	resultPort = make([]int, 0)
	if len(temp) == 1 {
		port, err := strconv.Atoi(temp[0])
		if err != nil {
			return nil, err
		}
		
		if err = isPort(port); err != nil {
			return nil,err
		}
		resultPort = append(resultPort, port)
		return resultPort, nil
	} else if len(temp) == 2 {
		startPort, err := strconv.Atoi(temp[0])
		if err != nil {
			return nil, err
		}
		endPort, err := strconv.Atoi(temp[1])
		if err != nil {
			return nil,err
		}

		if err = isPort(startPort); err != nil {
			return nil,err
		}

		if err = isPort(endPort); err != nil {
			return nil,err
		}

		for i := startPort; i <= endPort; i++ {
			resultPort = append(resultPort,i)
		}
		return resultPort,nil

	}
	return nil,errors.New("ports's format is 21-50")
}

func isPort(port int)error {
	if port < 1 && port > 65535 {
		return errors.New("this not is port")
	} 
	return nil
}
