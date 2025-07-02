package main

import (
	"cmd_linux2win/src/common"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

const cmdName = "false"

func main() {

	helpInfo := common.HelpInfo{
		Name: os.Args[0],
		UsageLines: []string{
			"[ignored command line arguments]",
			"OPTION",
		},
		Description: "Exit with a status code indicating failure.",
		Options: []common.Option{
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
			}},
		},
		Note: fmt.Sprintf(`NOTE: your shell may have its own version of %s, which usually supersedes
the version described here.  Please refer to your shell's documentation
for details about the options it supports.`, cmdName),
		IgnoreShorthand: true,
	}
	helpInfo.Parse()
	os.Exit(1)
}
