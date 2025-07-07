package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var cmdName = "tty"

func main() {
	//os.Args = []string{
	//	os.Args[0], "-s",
	//}
	helpInfo := common.HelpInfo{
		Name:        os.Args[0],
		UsageLines:  []string{"[OPTION]..."},
		Description: "Print the file name of the terminal connected to standard input.",
		Options: []common.Option{
			{Verbose: "silent", AliseVerbose: "quiet", Short: "s", Description: "print nothing, only return an exit status"},
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
				os.Exit(0)
			}},
		},
	}
	helpInfo.Parse()

	silent := common.GetBool("silent")
	if !isTerminal(os.Stdin.Fd()) {
		if !silent {
			fmt.Print("not a tty")
		}
		os.Exit(1)
	}
	if !silent {
		title, _ := getConsoleTitle()
		fmt.Println(title)
	}
	os.Exit(0)

}

// isTerminal 判断给定文件描述符是否关联到终端设备
func isTerminal(fd uintptr) bool {
	var mode uint32
	// 尝试获取控制台模式，若成功则为终端设备
	err := syscall.GetConsoleMode(syscall.Handle(fd), &mode)
	return err == nil
}

func getConsoleTitle() (string, error) {
	var kernel32 = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleTitle := kernel32.NewProc("GetConsoleTitleW")

	// 分配足够大的缓冲区
	buf := make([]uint16, 256)
	r0, _, e1 := syscall.Syscall(procGetConsoleTitle.Addr(), 2,
		uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), 0)
	if r0 == 0 {
		return "", error(e1)
	}
	return syscall.UTF16ToString(buf), nil
}
