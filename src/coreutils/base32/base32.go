package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"encoding/base32"
	"fmt"
	"io"
	"os"
)

// wrappedWriter 实现行包装功能
type wrappedWriter struct {
	w    io.Writer
	wrap int
	cnt  int
}

var cmdName = "base32"

func main() {
	//os.Args = []string{
	//	os.Args[0], "1",
	//}
	helpInfo := common.HelpInfo{
		Name:       os.Args[0],
		UsageLines: []string{"[OPTION]... [FILE]"},
		Description: `Base32 encode or decode FILE, or standard input, to standard output.

With no FILE, or when FILE is -, read standard input.

Mandatory arguments to long options are mandatory for short options too.`,
		Options: []common.Option{
			{Verbose: "decode", Short: "d", Description: "decode data"},
			{Verbose: "ignore-garbage", Short: "i", Description: "when decoding, ignore non-alphabet characters"},
			{Verbose: "wrap", Short: "w", Description: "wrap encoded lines after COLS character (default 76).\n Use 0 to disable line wrapping", DefaultP: "COLS"},
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
			}},
		},
	}
	helpInfo.Parse()

	var input io.Reader
	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		input = os.Stdin
	} else {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s: %s: No such file or directory", os.Args[0], flag.Arg(0))
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	decode := common.GetBool("decode")
	ignoreGarbage := common.GetBool("ignoreGarbage")
	wrap := common.GetInt("wrap")

	var output io.Writer = os.Stdout
	if decode {
		decoder := base32.StdEncoding
		if ignoreGarbage {
			decoder = decoder.WithPadding(base32.NoPadding)
		}
		reader := base32.NewDecoder(decoder, input)
		_, err := io.Copy(output, reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v", cmdName, err)
			os.Exit(1)
		}
	} else {
		var writer = output
		if wrap > 0 {
			writer = &wrappedWriter{w: output, wrap: wrap}
		}
		encoder := base32.NewEncoder(base32.StdEncoding, writer)
		_, err := io.Copy(encoder, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", cmdName, err)
			os.Exit(1)
		}
		encoder.Close()
	}

}

func (w *wrappedWriter) Write(p []byte) (n int, err error) {
	var total int
	for i := 0; i < len(p); i++ {
		if w.cnt > 0 && w.cnt%w.wrap == 0 {
			if _, err := w.w.Write([]byte{'\n'}); err != nil {
				return total, err
			}
			w.cnt = 0
		}
		if _, err := w.w.Write(p[i : i+1]); err != nil {
			return total, err
		}
		w.cnt++
		total++
	}
	return total, nil
}
