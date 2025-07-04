package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"strings"
)

var cmdName = "echo"

func main() {
	enableEscapeChars := false
	omitNewline := false

	helpInfo := common.HelpInfo{
		Name: os.Args[0],
		UsageLines: []string{
			"[SHORT-OPTION]... [STRING]...",
			"LONG-OPTION",
		},
		Description: "Echo the STRING(s) to standard output.",
		Options: []common.Option{
			{Verbose: "n", Short: "n", Description: "do not output the trailing newline", Func: func() {
				omitNewline = true
			}},
			{Verbose: "e", Short: "e", Description: "enable interpretation of backslash escapes", Func: func() {
				enableEscapeChars = true
			}},
			{Verbose: "E", Short: "E", Description: "disable interpretation of backslash escapes (default)", Func: func() {
				enableEscapeChars = false
			}},
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
				os.Exit(0)
			}},
		},
		Note: fmt.Sprintf(`If -e is in effect, the following sequences are recognized:

  \\      backslash
  \a      alert (BEL)
  \b      backspace
  \c      produce no further output
  \e      escape
  \f      form feed
  \n      new line
  \r      carriage return
  \t      horizontal tab
  \v      vertical tab
  \0NNN   byte with octal value NNN (1 to 3 digits)
  \xHH    byte with hexadecimal value HH (1 to 2 digits)

NOTE: your shell may have its own version of %s, which usually supersedes
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

	concatenated := strings.Join(flag.Args(), " ")

	a := []rune(concatenated)

	length := len(a)

	ai := 0

	if length != 0 {
		for i := 0; i < length; {
			c := a[i]
			i++
			if (enableEscapeChars == true) && c == '\\' && i < length {
				c = a[i]
				i++
				switch c {
				case 'a':
					c = '\a'
				case 'b':
					c = '\b'
				case 'c':
					os.Exit(0)
				case 'e':
					c = '\x1B'
				case 'f':
					c = '\f'
				case 'n':
					c = '\n'
				case 'r':
					c = '\r'
				case 't':
					c = '\t'
				case 'v':
					c = '\v'
				case '\\':
					c = '\\'
				case 'x':
					c = a[i]
					i++
					if '9' >= c && c >= '0' && i < length {
						hex := (c - '0')
						c = a[i]
						i++
						if '9' >= c && c >= '0' && i <= length {
							c = 16*(c-'0') + hex
						}
					}
				}
			}
			a[ai] = c
			ai++
		}
	}
	os.Stdout.WriteString(string(a[:ai]))
	if omitNewline == false {
		fmt.Print("\n")
	}

}
