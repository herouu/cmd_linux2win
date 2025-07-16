package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

var cmdName = "hostid"

func main() {
	helpInfo := common.HelpInfo{
		Name:        os.Args[0],
		UsageLines:  []string{"[OPTION]"},
		Description: `Print the numeric identifier (in hexadecimal) for the current host.`,
		Options: []common.Option{
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
				os.Exit(0)
			},
			}},
	}
	helpInfo.Parse()

	id, _ := GetHostID()
	fmt.Printf("%08x\n", id)
}

// GetHostID 模拟 gethostid() 功能
func GetHostID() (uint32, error) {
	// 1. 尝试读取 /etc/hostid 文件（类 Unix 系统）
	id, err := readHostIDFile("/etc/hostid")
	if err == nil {
		return id, nil
	}

	// 2. 读取失败时，基于主机名和IP生成标识
	return generateHostIDFromNetwork()
}

// 从 /etc/hostid 文件读取主机ID
func readHostIDFile(path string) (uint32, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	// /etc/hostid 通常存储 32 位或 64 位二进制数据
	if len(data) >= 4 {
		return binary.BigEndian.Uint32(data[:4]), nil
	}

	return 0, fmt.Errorf("invalid hostid file")
}

// 基于网络信息生成主机ID
func generateHostIDFromNetwork() (uint32, error) {
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return 0, err
	}

	// 解析主机名对应的IP地址
	addrs, err := net.LookupIP(hostname)
	if err != nil || len(addrs) == 0 {
		return 0, fmt.Errorf("failed to resolve hostname: %v", err)
	}

	// 优先使用IPv4地址
	var ipv4 net.IP
	for _, addr := range addrs {
		if ip := addr.To4(); ip != nil {
			ipv4 = ip
		}
	}

	// 如果没有IPv4，使用第一个IPv6地址
	if ipv4 == nil && len(addrs) > 0 {
		ipv4 = addrs[0].To4() // 可能为nil，但会在后续处理
	}

	if ipv4 == nil {
		return 0, fmt.Errorf("no valid IP address found")
	}

	// 将IP地址转换为32位整数并进行位运算（类似glibc实现）
	ipInt := binary.BigEndian.Uint32(ipv4)
	return (ipInt << 16) | (ipInt >> 16), nil
}
