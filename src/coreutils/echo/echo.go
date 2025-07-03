package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
)

var cmdName = "echo"

func main() {
	helpInfo := common.HelpInfo{
		Name: os.Args[0],
		UsageLines: []string{
			"[SHORT-OPTION]... [STRING]...",
			"LONG-OPTION",
		},
		Description: "Echo the STRING(s) to standard output.",
		Options: []common.Option{
			{Short: "n", Description: "do not output the trailing newline"},
			{Short: "e", Description: "enable interpretation of backslash escapes"},
			{Short: "E", Description: "disable interpretation of backslash escapes (default)"},
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
}
