package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
)

//#include <unistd.h>
import "C"

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
		fmt.Println(C.getuid())
	} else {
		for _, arg := range args {
			fmt.Println(arg)
		}
	}

}
