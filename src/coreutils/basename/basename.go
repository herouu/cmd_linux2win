package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var cmdName = "basename"

//	os.Args = []string{
//		os.Args[0],
//		"-a",
//		"any/str1",
//		"any/str2",
//	}
func main() {
	helpInfo := common.HelpInfo{
		Name:       os.Args[0],
		UsageLines: []string{"NAME [SUFFIX]", "OPTION... NAME..."},
		Description: `Print NAME with any leading directory components removed.
If specified, also remove a trailing SUFFIX.

Mandatory arguments to long options are mandatory for short options too.`,
		Options: []common.Option{
			{Verbose: "multiple", Short: "a", Description: "support multiple arguments and treat each as a NAME"},
			{Verbose: "suffix", Short: "s", Description: "remove a trailing SUFFIX; implies -a", VerboseUsage: "suffix", Type: "string"},
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
  basename /usr/bin/sort          -> "sort"
  basename include/stdio.h .h     -> "stdio"
  basename -s .h include/stdio.h  -> "stdio"
  basename -a any/str1 any/str2   -> "str1" followed by "str2"`,
	}
	helpInfo.Parse()

	multiple := common.GetBool("a")
	suffix := common.GetString("s")
	zero := common.GetBool("z")

	// 获取命令行参数
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "missing operand")
		flag.Usage()
		os.Exit(1)
	}

	// 处理参数
	var names []string
	if multiple || len(args) > 2 {
		names = args
	} else if len(args) == 2 {
		names = []string{args[0]}
		suffix = args[1]
	} else {
		names = args
	}

	// 处理每个名称
	for i, name := range names {
		// 去掉路径部分
		base := filepath.Base(name)

		// 去掉后缀
		if suffix != "" && strings.HasSuffix(base, suffix) {
			base = base[:len(base)-len(suffix)]
		}

		// 输出结果
		if zero {
			fmt.Print(base)
		} else if i < len(names)-1 {
			fmt.Print(base + "\n")
		} else {
			fmt.Print(base)
		}
	}
}
