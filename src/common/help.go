package common

import (
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"strings"
)

const (
	Copyright = `
Copyright (C) 2025 Herouu
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.`
)

type HelpInfo struct {
	Name            string   // 命令名称
	UsageLines      []string // 使用示例
	Description     string   // 描述信息
	Options         []Option // 选项列表
	Note            string   // 注意事项
	Copyright       string   // 注意事项
	Invalid         func(c string)
	ErrorHandling   *flag.ErrorHandling
	IgnoreShorthand bool
	IgnoreShortH    bool
	CmdFunc         func([]string) string
}

type Option struct {
	Short       string
	Verbose     string
	Description string
	CmdArr      []Cmd
	Func        func()
	BoolVarP    bool
}

type Cmd struct {
	Short string
	Long  string
}

type FlagOption struct {
	Opt  Option
	Flag *bool
}

func (h HelpInfo) Print() {
	// 输出使用示例
	for index, line := range h.UsageLines {
		if index == 0 {
			fmt.Fprintf(os.Stdout, "Usage: %s %s\n", h.Name, line)
		} else {
			fmt.Fprintf(os.Stdout, "  or:  %s %s\n", h.Name, line)
		}
	}
	// 输出描述信息
	if h.Description != "" {
		lines := strings.Split(h.Description, "\n")
		for _, line := range lines {
			fmt.Fprintln(os.Stdout, line)
		}
		fmt.Println()
	}

	// 输出选项列表
	if len(h.Options) > 0 {
		maxFlagLen := 0
		for _, opt := range h.Options {
			if len(opt.Verbose) > maxFlagLen {
				maxFlagLen = len(opt.Verbose)
			}
		}
		// option选项处理
		for _, opt := range h.Options {
			// 对齐选项描述
			spacing := strings.Repeat(" ", maxFlagLen-len(opt.Verbose)+2)
			if opt.Short != "" && opt.Verbose == "" {
				fmt.Fprintf(os.Stdout, "      -%s    %s%s\n", opt.Short, spacing, opt.Description)
			} else if opt.Short != "" && opt.Verbose != "" {
				fmt.Fprintf(os.Stdout, "      -%s, --%s%s%s\n", opt.Short, opt.Verbose, spacing, opt.Description)
			} else if opt.Short == "" && opt.Verbose != "" {
				fmt.Fprintf(os.Stdout, "          --%s%s%s\n", opt.Verbose, spacing, opt.Description)
			}
		}
		fmt.Fprintln(os.Stdout)
	}

	// 输出注意事项
	if h.Note != "" {
		fmt.Fprintln(os.Stdout, h.Note)
	}

	// 输出Copyright
	if h.Copyright != "" {
		fmt.Fprint(os.Stdout, h.Copyright)
	} else {
		fmt.Fprint(os.Stdout, Copyright)
	}
}

func (h HelpInfo) Parse() []FlagOption {
	flag.Usage = func() {
		h.Print()
	}
	if h.Invalid != nil {
		flag.InvalidFlag = h.Invalid
	}

	flag.CommandLine.SetOutput(os.Stdout)
	flag.CommandLine.SetIgnoreShorthand(h.IgnoreShorthand)
	flag.CommandLine.SetIgnoreShortH(h.IgnoreShortH)

	if h.ErrorHandling != nil {
		flag.CommandLine.Init(os.Args[0], *h.ErrorHandling)
	}

	var sliceOption []FlagOption
	for _, opt := range h.Options {
		var verbose bool

		flag.BoolVarP(&verbose, opt.Verbose, opt.Short, opt.BoolVarP, opt.Description)

		if len(opt.CmdArr) != 0 {
			for _, per := range opt.CmdArr {
				flag.BoolVarP(&verbose, per.Long, per.Short, false, opt.Description)
			}
		}
		sliceOption = append(sliceOption, FlagOption{
			Opt:  opt,
			Flag: &verbose,
		})
	}

	flag.Parse()
	for i := range sliceOption {
		if *sliceOption[i].Flag {
			sliceOption[i].Opt.Func()
		}
	}
	return sliceOption
}

func Version(cmdName string) string {
	const version = "%s (Go coreutils) 1.0%s"
	return fmt.Sprintf(version, cmdName, Copyright)
}
