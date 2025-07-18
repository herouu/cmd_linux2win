package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
)

const cmdName = "yes"

func main() {
	helpInfo := common.HelpInfo{
		Name:        os.Args[0],
		UsageLines:  []string{"[STRING]...", "OPTION"},
		Description: "Repeatedly output a line with all specified STRING(s), or 'y'.",
		Options: []common.Option{
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
				os.Exit(0)
			}}},
		Invalid: func(f string) {
			var err error
			if len(f) > 1 {
				err = fmt.Errorf(`%s: unrecognized option '--%s'
Try '%s --help' for more information.`, os.Args[0], f, os.Args[0])
			} else {
				err = fmt.Errorf(`%s: invalid option -- '%s'
Try '%s --help' for more information.`, os.Args[0], f, os.Args[0])
			}
			fmt.Println(err)
			os.Exit(1)
		},
		IgnoreShortH: true,
	}
	helpInfo.Parse()
	var opts = flag.Args()
	if len(opts) == 0 {
		opts = []string{"y"}
	}
	for {
		fmt.Println(opts[0])
	}
}
