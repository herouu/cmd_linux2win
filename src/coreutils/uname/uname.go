package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"github.com/elastic/go-sysinfo"
	"os"
	"runtime"
	"strings"
)

const cmdName = "uname"

func main() {
	os.Args = []string{
		os.Args[0],
		"--all",
	}
	helpInfo := common.HelpInfo{
		Name:        os.Args[0],
		UsageLines:  []string{"[OPTION]..."},
		Description: "Print certain system information.  With no OPTION, same as -s.",
		Options: []common.Option{
			{Verbose: "all", Short: "a", Description: "print all information, in the following order,\n  except omit -p and -i if unknown:"},
			{Verbose: "kernel-name", Short: "s", Description: "print the kernel name"},
			{Verbose: "nodename", Short: "n", Description: "print the network node hostname"},
			{Verbose: "kernel-release", Short: "r", Description: "print the kernel release"},
			{Verbose: "kernel-version", Short: "v", Description: "print the kernel version"},
			{Verbose: "machine", Short: "m", Description: "print the machine hardware name"},
			{Verbose: "processor", Short: "p", Description: "print the processor type (non-portable)"},
			{Verbose: "hardware-platform", Short: "i", Description: "print the hardware platform (non-portable)"},
			{Verbose: "operating-system", Short: "o", Description: "print the operating system"},
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

	// 初始化标志变量
	all := common.GetBool("all")
	kernelName := common.GetBool("kernel-name")
	nodename := common.GetBool("nodename")
	kernelRelease := common.GetBool("kernel-release")
	kernelVersion := common.GetBool("kernel-version")
	machine := common.GetBool("machine")
	processor := common.GetBool("processor")
	hardwarePlatform := common.GetBool("hardware-platform")
	operatingSystem := common.GetBool("operating-system")

	// 若未指定选项，默认输出 kernel name
	if !all && !kernelName && !nodename && !kernelRelease && !kernelVersion && !machine && !processor && !hardwarePlatform && !operatingSystem {
		kernelName = true
	}

	if all {
		kernelName = true
		nodename = true
		kernelRelease = true
		kernelVersion = true
		machine = true
		processor = true
		hardwarePlatform = true
		operatingSystem = true
	}

	var results []string

	self, _ := sysinfo.Host()
	fmt.Printf("%v", self)
	if kernelName {
		results = append(results, "Windows") // 硬编码为 Windows，因为使用了 Windows API
	}

	if nodename {
		hostname, err := os.Hostname()
		if err == nil {
			results = append(results, hostname)
		}
	}

	if machine {
		results = append(results, runtime.GOARCH)
	}

	if processor {
		// Windows 下难以获取处理器类型，暂留空
		results = append(results, "unknown")
	}

	if hardwarePlatform {
		// Windows 下难以获取硬件平台，暂留空
		results = append(results, "unknown")
	}

	if operatingSystem {
		results = append(results, "Windows")
	}

	fmt.Println(strings.Join(results, " "))

}
