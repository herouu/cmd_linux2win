package main

import flag "cmd_linux2win/src/lib/github.com/spf13/pflag"

func main() {
	flag.BoolPAlias("verbose", "v", "", "x", false, "verbose output")
	flag.Parse()
	flag.PrintDefaults()

}
