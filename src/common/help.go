package common

import (
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"os"
	"strconv"
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
	Short        string
	AliasShort   string
	Verbose      string
	VerboseUsage string // 用于描述的用法
	AliseVerbose string
	Description  string
	CmdArr       []Cmd
	Func         func()
	DefaultP     string // 默认参数值
	Type         string // 参数类型，默认为bool
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
		// 计算所有选项的最大长度
		maxFlagLen := 0
		for _, opt := range h.Options {
			flagStr := getFlagString(opt)
			if len(flagStr) > maxFlagLen {
				maxFlagLen = len(flagStr)
			}
		}

		// 输出每个选项
		for _, opt := range h.Options {
			flagStr := getFlagString(opt)
			spacing := strings.Repeat(" ", maxFlagLen-len(flagStr)+2)
			// 拆分多行描述
			descLines := strings.Split(opt.Description, "\n")
			for i, line := range descLines {
				if i == 0 {
					fmt.Fprintf(os.Stdout, "  %s%s%s\n", flagStr, spacing, line)
				} else {
					fmt.Fprintf(os.Stdout, "  %s%s%s\n", strings.Repeat(" ", len(flagStr)), spacing, line)
				}
			}
		}
		fmt.Println()
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

// getFlagString 生成选项的字符串表示
func getFlagString(opt Option) string {
	var flags []string
	if opt.Short != "" {
		flags = append(flags, "-"+opt.Short)
	}
	if opt.AliasShort != "" {
		flags = append(flags, "-"+opt.AliasShort)
	}
	if opt.Verbose != "" {
		if opt.Type == "string" && opt.VerboseUsage != "" {
			flags = append(flags, "--"+opt.Verbose+"="+strings.ToUpper(opt.VerboseUsage))
		} else {
			flags = append(flags, "--"+opt.Verbose)
		}
	}
	if opt.AliseVerbose != "" {
		flags = append(flags, "--"+opt.AliseVerbose)
	}

	needPrefix := len(flags) > 0 && strings.HasPrefix(flags[0], "--")
	var builder strings.Builder
	if needPrefix {
		builder.WriteString("    ")
	}
	for i, flag := range flags {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(flag)
	}
	join := builder.String()
	return join
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

		if opt.Type == "" {
			opt.Type = "bool"
		}

		if opt.Type == "bool" {
			var verbose bool
			var defaultBoolP bool
			if opt.DefaultP != "" {
				defaultBoolP, _ = strconv.ParseBool(opt.DefaultP)
			}
			flag.BoolVarPAlias(&verbose, opt.Verbose, opt.Short, opt.AliseVerbose, opt.AliasShort, defaultBoolP, opt.Description)
			if len(opt.CmdArr) != 0 {
				for _, per := range opt.CmdArr {
					flag.BoolVarP(&verbose, per.Long, per.Short, false, opt.Description)
				}
			}
			sliceOption = append(sliceOption, FlagOption{
				Opt:  opt,
				Flag: &verbose,
			})
		} else if opt.Type == "string" {
			var verbose string
			flag.StringVarP(&verbose, opt.Verbose, opt.Short, opt.DefaultP, opt.Description)
		}
	}

	flag.Parse()
	for i := range sliceOption {
		if *sliceOption[i].Flag && sliceOption[i].Opt.Func != nil {
			sliceOption[i].Opt.Func()
		}
	}
	return sliceOption
}

func Version(cmdName string) string {
	const version = "%s (Go coreutils) 1.0%s"
	return fmt.Sprintf(version, cmdName, Copyright)
}

func GetBool(nameOrShort string) bool {
	var f2 *flag.Flag
	if len(nameOrShort) == 1 {
		// 如果是单字符短选项，查找对应的短选项
		f2 = flag.ShorthandLookup(nameOrShort)
	} else {
		// 如果是长选项，查找对应的长选项
		f2 = flag.Lookup(nameOrShort)
	}
	if f2 != nil && f2.Value.String() == "true" {
		return true
	}
	return false
}

func GetString(nameOrShort string) string {
	var f2 *flag.Flag
	if len(nameOrShort) == 1 {
		// 如果是单字符短选项，查找对应的短选项
		f2 = flag.ShorthandLookup(nameOrShort)
	} else {
		// 如果是长选项，查找对应的长选项
		f2 = flag.Lookup(nameOrShort)
	}
	return f2.Value.String()

}
