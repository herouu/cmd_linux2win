package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"runtime"
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
	a := common.GetBool("a")
	if a {
		fmt.Println(runtime.GOARCH)
		fmt.Println(runtime.GOOS)
	}

}
