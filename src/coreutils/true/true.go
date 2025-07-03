package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
)

const cmdName = "true"

func main() {
	helpInfo := common.HelpInfo{
		Name: os.Args[0],
		UsageLines: []string{
			"[ignored command line arguments]",
			"OPTION",
		},
		Description: "Exit with a status code indicating success.",
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
		Note: fmt.Sprintf(`NOTE: your shell may have its own version of %s, which usually supersedes
the version described here.  Please refer to your shell's documentation
for details about the options it supports.`, cmdName),
		IgnoreShorthand: true,
	}
	helpInfo.Parse()
	os.Exit(0)
}
