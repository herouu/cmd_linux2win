package common

import (
	"flag"
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
	Name        string   // 命令名称
	UsageLines  []string // 使用示例
	Description string   // 描述信息
	Options     []Option // 选项列表
	Note        string   // 注意事项
	Copyright   string   // 注意事项
}

type Option struct {
	Flag        string
	Description string
	Func        func()
}

type flagOption struct {
	opt  Option
	flag *bool
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
			if len(opt.Flag) > maxFlagLen {
				maxFlagLen = len(opt.Flag)
			}
		}
		// option选项处理
		for _, opt := range h.Options {
			// 对齐选项描述
			spacing := strings.Repeat(" ", maxFlagLen-len(opt.Flag)+2)
			fmt.Fprintf(os.Stdout, "      --%s%s%s\n", opt.Flag, spacing, opt.Description)
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

func (h HelpInfo) Parse() {
	var sliceOption []flagOption
	for _, opt := range h.Options {
		b := flag.Bool(opt.Flag, false, opt.Description)
		sliceOption = append(sliceOption, flagOption{
			opt:  opt,
			flag: b,
		})
	}
	flag.Parse()

	for i := range sliceOption {
		if *sliceOption[i].flag {
			sliceOption[i].opt.Func()
		}
	}
}

func Version(cmdName string) string {
	const version = "%s (Go coreutils) 1.0"
	return fmt.Sprintf(version, cmdName)
}
