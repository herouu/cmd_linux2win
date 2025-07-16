package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const cmdName = "pwd"

func main() {
	helpInfo := common.HelpInfo{
		Name:        os.Args[0],
		UsageLines:  []string{"[OPTION]... [FILE]"},
		Description: `Print the full filename of the current working directory.`,
		Options: []common.Option{
			{Verbose: "logical", Short: "L", Description: "use PWD from environment, even if it contains symlinks", Default: false},
			{Verbose: "physical", Short: "P", Description: "avoid all symlinks", Default: true},
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
				os.Exit(0)
			}},
		},
		Note: `If no option is specified, -P is assumed.

NOTE: your shell may have its own version of pwd, which usually supersedes
the version described here.  Please refer to your shell's documentation
for details about the options it supports.`,
	}
	helpInfo.Parse()

	logical := common.GetBool("logical")
	physical := common.GetBool("physical")
	fmt.Println(runtime.GOOS)
	if logical {
		dir := getDir()
		fmt.Println(dir)
	} else if physical {
		dir := getDir()
		symlinks, _ := filepath.EvalSymlinks(dir)
		fmt.Println(symlinks)
	} else {
		dir := getDir()
		fmt.Println(dir)
	}

}

func getDir() string {
	if os.Getenv("PWD") != "" {
		return os.Getenv("PWD")
	}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return dir
}
