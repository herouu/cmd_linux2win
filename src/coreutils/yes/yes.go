package main

import (
	"cmd_linux2win/src/common"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

const cmdName = "yes"

func main() {
	helpInfo := common.NewHelpInfo()
	helpInfo.Name = os.Args[0]
	helpInfo.UsageLines = []string{"[STRING]...", "OPTION"}
	helpInfo.Description = "Repeatedly output a line with all specified STRING(s), or 'y'."
	helpInfo.Options = []common.Option{
		{Verbose: "help", Description: "display this help and exit", Func: func() {
			flag.Usage()
			os.Exit(0)
		}},
		{Verbose: "version", Description: "output version information and exit", Func: func() {
			fmt.Print(common.Version(cmdName))
			os.Exit(0)
		}}}
	helpInfo.Invalid = func(c uint8) {
		err := fmt.Errorf(`%s: invalid option -- %q
Try '%s --help' for more information.`, os.Args[0], c, os.Args[0])
		fmt.Println(err)
		os.Exit(1)
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
