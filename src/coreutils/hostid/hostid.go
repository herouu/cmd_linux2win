package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"unsafe"
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
	return uint32(getHostIDLikeMsys2()), nil
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

func getHostIDLikeMsys2() int32 {
	// 初始哈希种子，与C代码保持一致
	hostid := int32(0x40291372)

	// 读取注册表中的MachineGuid
	guid, err := getMachineGUID()
	if err != nil {
		// 如果读取失败，使用默认GUID
		guid = "00000000-0000-0000-0000-000000000000"
	}

	// 应用SDBM哈希算法
	for _, r := range guid {
		// 与C代码中的哈希逻辑完全一致：hostid = *wp + (hostid << 6) + (hostid << 16) - hostid
		hostid = int32(r) + (hostid << 6) + (hostid << 16) - hostid
	}

	return hostid
}

// 从Windows注册表读取MachineGuid
func getMachineGUID() (string, error) {
	// 注册表路径：HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Cryptography
	keyPath, err := syscall.UTF16PtrFromString(`SOFTWARE\Microsoft\Cryptography`)
	if err != nil {
		return "", err
	}

	// 打开注册表项
	var hKey syscall.Handle
	err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, keyPath, 0, syscall.KEY_READ, &hKey)
	if err != nil {
		return "", err
	}
	defer syscall.RegCloseKey(hKey)

	// 读取MachineGuid值
	valueName, err := syscall.UTF16PtrFromString("MachineGuid")
	if err != nil {
		return "", err
	}

	// 先查询所需缓冲区大小
	var bufSize uint32
	var bufType uint32
	err = syscall.RegQueryValueEx(hKey, valueName, nil, &bufType, nil, &bufSize)
	if err != nil {
		return "", err
	}

	// 分配缓冲区并读取值
	buf := make([]uint16, bufSize/2)
	err = syscall.RegQueryValueEx(hKey, valueName, nil, &bufType, (*byte)(unsafe.Pointer(&buf[0])), &bufSize)
	if err != nil {
		return "", err
	}

	// 将UTF-16字节转换为Go字符串
	guid := syscall.UTF16ToString(buf)
	return guid, nil
}
