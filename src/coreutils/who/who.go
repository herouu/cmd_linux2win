package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"time"
)

var cmdName = "who"

func main() {
	helpInfo := common.HelpInfo{
		Name:        os.Args[0],
		UsageLines:  []string{"[OPTION]... [ FILE | ARG1 ARG2 ]"},
		Description: "Print information about users who are currently logged in.",
		Options: []common.Option{
			{Verbose: "all", Short: "a", Description: "same as -b -d --login -p -r -t -T -u"},
			{Verbose: "boot", Short: "b", Description: "time of last system boot"},
			{Verbose: "dead", Short: "d", Description: "print dead processes"},
			{Verbose: "heading", Short: "H", Description: "print line of column headings"},
			{Verbose: "ips", Description: "print ips instead of hostnames. with --lookup,\ncanonicalizes based on stored IP, if available,\nrather than stored hostname"},
			{Verbose: "login", Short: "l", Description: "print system login processes"},
			{Verbose: "lookup", Description: "attempt to canonicalize hostnames via DNS"},
			{Verbose: "m", Short: "m", Description: "only hostname and user associated with stdin"},
			{Verbose: "process", Short: "p", Description: "print active processes spawned by init"},
			{Verbose: "count", Short: "q", Description: "all login names and number of users logged on"},
			{Verbose: "runlevel", Short: "r", Description: "print current runlevel"},
			{Verbose: "short", Short: "s", Description: "print only name, line, and time (default)"},
			{Verbose: "time", Short: "t", Description: "print last system clock change"},
			{Verbose: "mesg", Short: "T", Description: "add user's message status as +, -, or ?"},
			{Verbose: "users", Short: "u", Description: "list users logged in"},
			{Verbose: "message", Description: "same as -T"},
			{Verbose: "writable", Description: "same as -T"},
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
				os.Exit(0)
			}},
		},
		Note: fmt.Sprintf(`If FILE is not specified, use /var/run/utmp.  /var/log/wtmp as FILE is common.
If ARG1 ARG2 given, -m presumed: 'am i' or 'mom likes' are usual.`),
	}
	helpInfo.Parse()
}

func timeOfDay() int64 {
	now := time.Now().Unix()
	return now
}
