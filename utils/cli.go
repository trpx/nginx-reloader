package utils

import (
	"fmt"
	"os"
	"strconv"
)

func ParseDirsOptions() (pollEvery int, dirs []string, options []string) {
	if len(os.Args) == 1 ||
		(len(os.Args) == 2 &&
			(os.Args[1] == "-h" ||
				os.Args[1] == "--help")) {

		exitCode := 1
		if len(os.Args) == 2 {
			exitCode = 0
		}
		usage()
		os.Exit(exitCode)
	}

	pollEvery, err := strconv.Atoi(os.Args[1])
	if err != nil {
		usage()
		os.Exit(1)
	}

	separatorIndex := len(os.Args)
	for idx, el := range os.Args {
		if el == "--" {
			separatorIndex = idx
			if len(os.Args) > idx+1 {
				options = os.Args[separatorIndex+1:]
			}
			break
		}
	}
	dirs = os.Args[1:separatorIndex]

	return pollEvery, dirs, options
}

func usage() {
	fmt.Println("usage: nginx-reloader <POLL_EVERY_SECONDS> [WATCHED_DIR...] -- [NGINX_CLI_OPTION...]")
}
