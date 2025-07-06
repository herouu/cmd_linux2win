package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var cmdName = "sleep"

//	os.Args = []string{
//		os.Args[0],
//		"1s", "20s",
//	}
func main() {

	helpInfo := common.HelpInfo{
		Name: os.Args[0],
		UsageLines: []string{
			"NUMBER[SUFFIX]...", "OPTION",
		},
		Description: `Pause for NUMBER seconds.  SUFFIX may be 's' for seconds (the default),
'm' for minutes, 'h' for hours or 'd' for days.  Unlike most implementations
that require NUMBER be an integer, here NUMBER may be an arbitrary floating
point number.  Given two or more arguments, pause for the amount of time
specified by the sum of their values.`,
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
	helpInfo.Parse()

	if flag.NArg() == 0 {
		fmt.Print("sleep: missing operand\ntry 'sleep --help' for more information")
		os.Exit(1)
	}

	var totalDuration float64
	for i := 0; i < flag.NArg(); i++ {
		arg := flag.Arg(i)
		var duration float64
		var err error

		// 处理无后缀的情况，默认单位为秒
		if len(arg) > 0 && (arg[len(arg)-1] < '0' || arg[len(arg)-1] > '9') {
			duration, err = strconv.ParseFloat(arg[:len(arg)-1], 64)
			if err != nil {
				fmt.Printf("sleep: invalid time interval '%s'\nTry 'sleep --help' for more information.", arg)
				os.Exit(1)
			}
			switch arg[len(arg)-1] {
			case 's':
				duration = duration

				time.Sleep(time.Duration(totalDuration * float64(time.Second)))
			case 'm':
				duration = duration * 60
			case 'h':
				duration = duration * 3600
			case 'd':
				duration = duration * 86400
			default:
				fmt.Printf("sleep: invalid time interval '%s'\nTry 'sleep --help' for more information.", arg)
				os.Exit(1)
			}
		} else {
			duration, err = strconv.ParseFloat(arg, 64)
			if err != nil {
				fmt.Printf("sleep: invalid time interval '%s'\nTry 'sleep --help' for more information.", arg)
				os.Exit(1)
			}
		}
		totalDuration += duration
	}
	time.Sleep(time.Duration(totalDuration) * time.Second)
}
