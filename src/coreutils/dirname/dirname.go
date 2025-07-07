package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"path"
)

var cmdName = "dirname"

func main() {
	//os.Args = []string{os.Args[0], "/usr/bin/bin1", "dir1/str", "stdio.h"}
	helpInfo := common.HelpInfo{
		Name:       os.Args[0],
		UsageLines: []string{"[OPTION] NAME..."},
		Description: `Output each NAME with its last non-slash component and trailing slashes
removed; if NAME contains no /'s, output '.' (meaning the current directory).`,
		Options: []common.Option{
			{Verbose: "zero", Short: "z", Description: "end each output line with NUL, not newline"},
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
				os.Exit(0)
			}},
		},
		Note: `Examples:
  dirname /usr/bin/          -> "/usr"
  dirname dir1/str dir2/str  -> "dir1" followed by "dir2"
  dirname stdio.h            -> "."`,
	}
	helpInfo.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stdout, "%s: missing operand\nTry '%s --help' for more information.", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	zero := common.GetBool("zero")

	for i, arg := range flag.Args() {
		dirname := getDirname(arg)
		if zero {
			fmt.Print(dirname)
		} else if i < len(flag.Args())-1 {
			fmt.Println(dirname)
		} else {
			fmt.Print(dirname)
		}
	}
}

// getDirname 实现 dirname 逻辑
func getDirname(pathStr string) string {
	// 去除路径末尾的斜杠
	cleanPath := path.Clean(pathStr)
	// 获取路径的目录部分
	dir := path.Dir(cleanPath)

	return dir
}
