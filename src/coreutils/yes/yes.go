package main

import (
	"cmd_linux2win/src/common"
	"flag"
	"fmt"
	"os"
)

func main() {
	helpInfo := common.HelpInfo{
		Name:        os.Args[0],
		UsageLines:  []string{"[STRING]...", "OPTION"},
		Description: "Repeatedly output a line with all specified STRING(s), or 'y'.",
		Options: []common.Option{
			{"help", "display this help and exit", func() {
				flag.Usage()
				os.Exit(0)
			}},
			{"version", "output version information and exit", func() {
				fmt.Print(common.Version(os.Args[0]))
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
