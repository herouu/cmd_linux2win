package main

import (
	"cmd_linux2win/src/common"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

var cmdName = "base64"

func main() {
	helpInfo := common.HelpInfo{
		Name:       os.Args[0],
		UsageLines: []string{"[OPTION]... [FILE]"},
		Description: `base64 encode or decode F LE, or standard input, to standard output.

With no FILE, or when FILE is -, read standard input.

Mandatory arguments to long options are mandatory for short options too.`,
		Options: []common.Option{
			{Verbose: "decode", Short: "d", Description: "decode data"},
			{Verbose: "ignore-garbage", Short: "i", Description: "when decoding, ignore non-alphabet characters"},
			{Verbose: "wrap", Short: "w", Description: "wrap encoded lines after COLS character (default 76).\n Use 0 to disable line wrapping", DefaultP: "COLS", Type: "int", Default: 76},
			{Verbose: "help", Description: "display this help and exit", Func: func() {
				flag.Usage()
				os.Exit(0)
			}},
			{Verbose: "version", Description: "output version information and exit", Func: func() {
				fmt.Print(common.Version(cmdName))
			}},
		},
		Note: `The data are encoded as described for the base64 alphabet in RFC 4648.
When decoding, the input may contain newlines in addition to the bytes of
the formal base64 alphabet.  Use --ignore-garbage to attempt to recover
from any other non-alphabet bytes in the encoded stream.`,
	}
	helpInfo.Parse()

	var input io.Reader
	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		input = os.Stdin
	} else if flag.NArg() == 1 {
		arg0 := flag.Arg(0) + "\n"
		input = &common.MockReader{Data: []byte(arg0)}
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
		decoder := base64.StdEncoding
		if ignoreGarbage {
			decoder = decoder.WithPadding(base64.NoPadding)
		}
		reader := base64.NewDecoder(decoder, input)
		_, err := io.Copy(output, reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v", cmdName, err)
			os.Exit(1)
		}
	} else {
		var writer io.Writer = output
		if wrap > 0 {
			writer = &common.WrappedWriter{W: output, Wrap: wrap}
		}
		encoder := base64.NewEncoder(base64.StdEncoding, writer)
		_, err := io.Copy(encoder, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", cmdName, err)
			os.Exit(1)
		}
		if err := encoder.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", cmdName, err)
			os.Exit(1)
		}
	}

}
