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
	os.Args = []string{
		os.Args[0],
	}
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
	if all {
		mode = NprocAll
	}
	fmt.Println(NumProcessors(mode))
}

// 解析环境变量为 uint64
func parseOmpThreads(env string) *uint64 {
	val := os.Getenv(env)
	if val == "" {
		return nil
	}
	num, err := strconv.ParseUint(val, 10, 64)
	if err != nil || num == 0 {
		return nil
	}
	return &num
}

func NumProcessors(query NprocQuery) uint64 {
	var ompEnvLimit uint64 = 0

	if query == NprocCurrentOverridable {
		ompEnvThreads := parseOmpThreads("OMP_NUM_THREADS")
		ompEnvLimit := parseOmpThreads("OMP_THREAD_LIMIT")
		if ompEnvLimit == nil {
			*ompEnvLimit = ulongMax
		}
		if ompEnvThreads != 0 {
			return min(*ompEnvThreads, *ompEnvLimit)
		}
		query = NprocCurrent
	}

	// 查询实际 CPU 数量
	nprocs := uint64(runtime.NumCPU())

	if *ompEnvLimit == 1 {
		return 1
	}

	if nprocs < *ompEnvLimit {
		return nprocs
	}
	return *ompEnvLimit
}
