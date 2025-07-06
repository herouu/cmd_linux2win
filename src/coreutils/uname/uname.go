package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"github.com/shirou/gopsutil/host"
	"golang.org/x/sys/windows/registry"
	"os"
	"strings"
	"time"
)

const cmdName = "uname"

func main() {
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

	info, _ := host.Info()
	if kernelName {
		results = append(results, info.OS) // 硬编码为 Windows，因为使用了 Windows API
	}

	if nodename {
		results = append(results, info.Hostname)
	}

	if kernelRelease {
		results = append(results, info.KernelVersion)
	}

	if kernelVersion {
		key, _ := registry.OpenKey(registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
			registry.QUERY_VALUE)
		defer key.Close()
		installTimeStamp, _, _ := key.GetIntegerValue("InstallDate")
		installTime := time.Unix(int64(installTimeStamp), 0)
		results = append(results, installTime.Format(time.DateTime))
	}

	if machine {
		results = append(results, info.KernelArch)
	}

	if processor {
		results = append(results, info.KernelArch)
	}

	if hardwarePlatform {
		results = append(results, info.KernelArch)
	}

	if operatingSystem {
		results = append(results, info.Platform)
	}

	fmt.Println(strings.Join(results, " "))

}
