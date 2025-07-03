package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"log"
	"os"
	"os/user"
)

var cmdName = "whoami"

func main() {
	helpInfo := common.HelpInfo{
		Name: os.Args[0],
		UsageLines: []string{
			"[OPTION]...",
		},
		Description: "Print the user name associated with the current effective user ID.\nSame as id -un.",
		Options: []common.Option{
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
				os.Exit(0)
			}},
		},
		Invalid: func(f string) {
			var err error
			if len(f) > 1 {
				err = fmt.Errorf(`%s: unrecognized option '--%s'
Try '%s --help' for more information.`, os.Args[0], f, os.Args[0])
			} else {
				err = fmt.Errorf(`%s: unknown option -- '%s'
Try '%s --help' for more information.`, os.Args[0], f, os.Args[0])
			}
			fmt.Println(err)
			os.Exit(1)
		},
		IgnoreShortH: true,
	}
	helpInfo.Parse()
	fmt.Println(whoami())
}

func whoami() string {
	u, err := user.Current()
	if err != nil {
		log.Fatalln("cannot find name for current user")
	}
	return u.Name
}
