package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"syscall"
	"unsafe"
)

// 加载 Windows DLL
var (
	wtsapi32                        = syscall.NewLazyDLL("wtsapi32.dll")
	procWTSEnumerateSessionsExW     = wtsapi32.NewProc("WTSEnumerateSessionsExW")
	procWTSQuerySessionInformationW = wtsapi32.NewProc("WTSQuerySessionInformationW")
	procWTSFreeMemory               = wtsapi32.NewProc("WTSFreeMemory")
)

// WTS_SESSION_INFO_1 定义 WTS_SESSION_INFO_1 结构体
type WTS_SESSION_INFO_1 struct {
	SessionID    uint32
	State        uint32
	pSessionName *uint16
	pUserName    *uint16
	pDomainName  *uint16
	pFarmName    *uint16
}

// WTS_INFO_CLASS 定义 WTS_INFO_CLASS 枚举
type WTS_INFO_CLASS int

const (
	WTSUserName WTS_INFO_CLASS = iota
	WTSDomainName
	WTSConnectState
	WTSClientName
)

var cmdName = "who"

func main() {
	helpInfo := common.HelpInfo{
		Name:        os.Args[0],
		UsageLines:  []string{"[OPTION]... [ FILE | ARG1 ARG2 ]"},
		Description: "Print information about users who are currently logged in.",
		Options: []common.Option{
			{Verbose: "all", Short: "a", Description: "same as -b -d --login -p -r -t -T -u"},
			{Verbose: "boot", Short: "b", Description: "time of last system boot"},
			{Verbose: "dead", Short: "d", Description: "print dead processes"},
			{Verbose: "heading", Short: "H", Description: "print line of column headings"},
			{Verbose: "ips", Description: "print ips instead of hostnames. with --lookup,\ncanonicalizes based on stored IP, if available,\nrather than stored hostname"},
			{Verbose: "login", Short: "l", Description: "print system login processes"},
			{Verbose: "lookup", Description: "attempt to canonicalize hostnames via DNS"},
			{Short: "m", Description: "only hostname and user associated with stdin"},
			{Verbose: "process", Short: "p", Description: "print active processes spawned by init"},
			{Verbose: "count", Short: "q", Description: "all login names and number of users logged on"},
			{Verbose: "runlevel", Short: "r", Description: "print current runlevel"},
			{Verbose: "short", Short: "s", Description: "print only name, line, and time (default)"},
			{Verbose: "time", Short: "t", Description: "print last system clock change"},
			{Verbose: "mesg", Short: "T", AliasShort: "w", Description: "add user's message status as +, -, or ?"},
			{Verbose: "users", Short: "u", Description: "list users logged in"},
			{Verbose: "message", Description: "same as -T"},
			{Verbose: "writable", Description: "same as -T"},
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
				os.Exit(0)
			}},
		},
		Note: fmt.Sprintf(`If FILE is not specified, use /var/run/utmp.  /var/log/wtmp as FILE is common.
If ARG1 ARG2 given, -m presumed: 'am i' or 'mom likes' are usual.`),
	}
	helpInfo.Parse()

	// 枚举会话
	var ppSessionInfo *WTS_SESSION_INFO_1
	var dwCount uint32
	_, _, err := procWTSEnumerateSessionsExW.Call(
		uintptr(0xFFFFFFFF), // WTS_CURRENT_SERVER_HANDLE
		uintptr(0),          // Reserved
		1,                   // Level
		uintptr(unsafe.Pointer(&ppSessionInfo)),
		uintptr(unsafe.Pointer(&dwCount)),
	)
	if err != nil && err != syscall.Errno(0) {
		fmt.Fprintf(os.Stderr, "Failed to enumerate sessions: %v\n", err)
		os.Exit(1)
	}
	defer procWTSFreeMemory.Call(uintptr(unsafe.Pointer(ppSessionInfo)))

	// 遍历会话信息
	for i := 0; i < int(dwCount); i++ {
		sessionInfo := (*WTS_SESSION_INFO_1)(unsafe.Pointer(uintptr(unsafe.Pointer(ppSessionInfo)) + uintptr(i)*unsafe.Sizeof(WTS_SESSION_INFO_1{})))

		// 获取用户名
		var pBuffer *uint16
		var dwBytesReturned uint32
		_, _, err := procWTSQuerySessionInformationW.Call(
			uintptr(0xFFFFFFFF), // WTS_CURRENT_SERVER_HANDLE
			uintptr(sessionInfo.SessionID),
			uintptr(WTSUserName),
			uintptr(unsafe.Pointer(&pBuffer)),
			uintptr(unsafe.Pointer(&dwBytesReturned)),
		)
		if err != nil && err != syscall.Errno(0) {
			continue
		}
		userName := windows.UTF16PtrToString(pBuffer)
		procWTSFreeMemory.Call(uintptr(unsafe.Pointer(pBuffer)))

		// 获取会话名
		sessionName := windows.UTF16PtrToString(sessionInfo.pSessionName)

		fmt.Printf("%-8s %-12s\n", userName, sessionName)
	}
}
