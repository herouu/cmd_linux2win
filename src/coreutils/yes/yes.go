package main

import (
	"cmd_linux2win/src/common"
	"flag"
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
			}},
		},
	}
	flag.Usage = func() {
		helpInfo.Print()
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
