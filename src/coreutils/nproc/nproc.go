package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

const appName = "nproc"
const ulongMax = ^uint64(0)

type NprocQuery int

const (
	NprocAll NprocQuery = iota
	NprocCurrent
	NprocCurrentOverridable
)

func main() {
	helpInfo := common.HelpInfo{
		Name:        os.Args[0],
		UsageLines:  []string{"[OPTION]..."},
		Description: "Print the number of processing units available to the current process,\nwhich may be less than the number of online processors",
		Options: []common.Option{
			{Verbose: "all", Description: "print the number of processors"},
			{Verbose: "ignore", Description: "if possible, exclude N processing units", VerboseUsage: "N", Type: "string"},
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(appName))
			}},
		},
	}
	helpInfo.Parse()

	mode := NprocCurrentOverridable
	all := common.GetBool("all")
	ignore := common.GetInt("ignore")
	if all {
		mode = NprocAll
	}

	nproc := NumProcessors(mode)
	uignore := uint64(ignore)
	if uignore < nproc {
		nproc -= uignore
	} else {
		nproc = 1
	}
	fmt.Println(nproc)
}

// 解析环境变量为 uint64
func parseOmpThreads(env string) uint64 {
	val := os.Getenv(env)
	if val == "" {
		return 0
	}
	num, err := strconv.ParseUint(val, 10, 64)
	if err != nil || num == 0 {
		return 0
	}
	return num
}

func NumProcessors(query NprocQuery) uint64 {
	var ompEnvLimit = ulongMax
	if query == NprocCurrentOverridable {
		ompEnvThreads := parseOmpThreads("OMP_NUM_THREADS")
		ompEnvLimit := parseOmpThreads("OMP_THREAD_LIMIT")
		if ompEnvLimit == 0 {
			ompEnvLimit = ulongMax
		}
		if ompEnvThreads > 0 {
			return min(ompEnvThreads, ompEnvLimit)
		}
		query = NprocCurrent
	}

	if ompEnvLimit == 1 {
		return 1
	}
	// 查询实际 CPU 数量
	nprocs := uint64(runtime.NumCPU())
	return min(nprocs, ompEnvLimit)
}
