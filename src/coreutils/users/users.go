package main

import (
	"cmd_linux2win/src/common"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var cmdName = "users"

func main() {
	helpInfo := common.HelpInfo{
		Name:        os.Args[0],
		UsageLines:  []string{"[OPTION]... [FILE]"},
		Description: "Output who is currently logged in according to FILE.\nIf FILE is not specified, use /var/run/utmp.  /var/log/wtmp as FILE is common.",
		Options: []common.Option{
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				common.HelpInfo{}.Print()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
			}},
		},
	}
	helpInfo.Parse()

	users, err := getLoggedOnUsers()
	if err != nil {
		fmt.Printf("获取登录用户失败: %v\n", err)
		return
	}

	for i, user := range users {
		fmt.Printf("%s", user)
		if i != len(users)-1 {
			fmt.Print(" ")
		}
	}

}

// 定义Windows API所需的常量和结构体
const (
	WTS_CURRENT_SERVER_HANDLE = 0
	WTS_SESSION_INFO_LEVEL    = 1
	WTSUserName               = 5
)

type WTS_SESSION_INFO struct {
	SessionID       uint32
	pWinStationName *uint16
	State           uint32
}

var (
	wtsapi32                       = syscall.NewLazyDLL("wtsapi32.dll")
	procWTSEnumerateSessions       = wtsapi32.NewProc("WTSEnumerateSessionsW")
	procWTSQuerySessionInformation = wtsapi32.NewProc("WTSQuerySessionInformationW")
	procWTSFreeMemory              = wtsapi32.NewProc("WTSFreeMemory")
)

// 获取所有登录的用户
func getLoggedOnUsers() ([]string, error) {
	var pSessions unsafe.Pointer
	var count uint32

	// 枚举所有会话
	ret, _, err := procWTSEnumerateSessions.Call(
		uintptr(WTS_CURRENT_SERVER_HANDLE),
		0,
		uintptr(WTS_SESSION_INFO_LEVEL),
		uintptr(unsafe.Pointer(&pSessions)),
		uintptr(unsafe.Pointer(&count)),
	)

	if ret == 0 {
		return nil, fmt.Errorf("WTSEnumerateSessions failed: %v", err)
	}
	defer procWTSFreeMemory.Call(uintptr(pSessions))

	sessions := (*[1 << 20]WTS_SESSION_INFO)(pSessions)[:count]
	users := make([]string, 0)

	// 遍历每个会话获取用户名
	for _, session := range sessions {
		var pBuffer unsafe.Pointer
		var bytesReturned uint32

		// 查询会话的用户名信息
		ret, _, _ := procWTSQuerySessionInformation.Call(
			uintptr(WTS_CURRENT_SERVER_HANDLE),
			uintptr(session.SessionID),
			uintptr(WTSUserName),
			uintptr(unsafe.Pointer(&pBuffer)),
			uintptr(unsafe.Pointer(&bytesReturned)),
		)

		if ret == 0 {
			continue // 有些会话可能没有用户名
		}
		defer procWTSFreeMemory.Call(uintptr(pBuffer))

		// 转换为字符串
		userName := syscall.UTF16ToString((*[1 << 20]uint16)(pBuffer)[:])
		if userName != "" {
			users = append(users, userName)
		}
	}
	return users, nil
}
