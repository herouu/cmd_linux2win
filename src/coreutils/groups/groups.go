package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"syscall"
)

const appName = "groups"

func main() {
	os.Args = []string{
		os.Args[0],
	}
	helpInfo := common.HelpInfo{
		Name:       os.Args[0],
		UsageLines: []string{"[OPTION]... [USERNAME]..."},
		Description: `Print group memberships for each USERNAME or, if no USERNAME is specified, for
the current process (which may differ if the groups database has changed).
`,
		Options: []common.Option{
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(appName))
			}},
		},
	}
	helpInfo.Parse()

	args := flag.Args()

	if len(args) == 0 {
		msys2 := syscall.NewLazyDLL("msys-2.0.dll")
		proc := msys2.NewProc("getuid")
		err := proc.Find()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(proc)
		r1, r3, err := proc.Call()
		fmt.Println(r1)
		fmt.Println(r3)
		if err != nil {
			fmt.Println(err)
		}
	} else {

	}

}
