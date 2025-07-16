package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"encoding/hex"
	"fmt"
	"github.com/shirou/gopsutil/disk"
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

	id, _ := getWindowsHostID()
	fmt.Println(id)
}

func getWindowsHostID() (string, error) {
	// 使用 WMI 获取主板 UUID（需要管理员权限，或通过注册表）
	// 简化方案：使用系统盘的 UUID
	disks, err := disk.Partitions(false)
	if err != nil {
		return "", err
	}
	for _, d := range disks {
		if d.Mountpoint == "C:" { // 假设系统盘为 C 盘
			info, err := disk.Usage(d.Mountpoint)
			if err == nil && info.Fstype != "" {
				return info.Fstype[:8], nil
			}
		}
	}
	return getMacBasedHostID()
}

// 通用 fallback：基于网卡 MAC 地址生成
func getMacBasedHostID() (string, error) {
	// 获取所有网卡 MAC 地址
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback == 0 && len(iface.HardwareAddr) > 0 {
			// 取 MAC 地址前 4 字节，转换为 8 位 16 进制
			mac := iface.HardwareAddr[:4]
			return hex.EncodeToString(mac), nil
		}
	}
	return "", fmt.Errorf("no valid network interface found")
}
